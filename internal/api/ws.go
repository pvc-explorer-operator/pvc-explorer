/*
Copyright 2026.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"nhooyr.io/websocket"        //nolint:staticcheck
	"nhooyr.io/websocket/wsjson" //nolint:staticcheck
)

const (
	ringBufferSize = 500
	clientBufSize  = 64
)

var frameCounter atomic.Uint64

func newFrameID() string {
	seq := frameCounter.Add(1)
	return fmt.Sprintf("%d-%d", time.Now().UnixMilli(), seq)
}

func newFrame(t string, payload any) (WSFrame, error) {
	raw, err := json.Marshal(payload)
	if err != nil {
		return WSFrame{}, err
	}
	return WSFrame{
		ID:         newFrameID(),
		Type:       EventType(t),
		ServerTime: time.Now().UTC(),
		Payload:    raw,
	}, nil
}

type Broadcaster struct {
	mu      sync.RWMutex
	clients map[chan WSFrame]struct{}
	ring    []WSFrame
	ringPos int
}

func NewBroadcaster() *Broadcaster {
	return &Broadcaster{
		clients: make(map[chan WSFrame]struct{}),
		ring:    make([]WSFrame, 0, ringBufferSize),
	}
}

func (b *Broadcaster) Publish(t string, payload any) error {
	frame, err := newFrame(t, payload)
	if err != nil {
		return err
	}

	b.mu.Lock()
	if len(b.ring) < ringBufferSize {
		b.ring = append(b.ring, frame)
	} else {
		b.ring[b.ringPos] = frame
		b.ringPos = (b.ringPos + 1) % ringBufferSize
	}
	for ch := range b.clients {
		select {
		case ch <- frame:
		default:
		}
	}
	b.mu.Unlock()
	return nil
}

func (b *Broadcaster) since(id string) []WSFrame {
	if id == "" {
		return nil
	}
	b.mu.RLock()
	defer b.mu.RUnlock()
	var out []WSFrame
	found := false
	for _, f := range b.ring {
		if found {
			out = append(out, f)
		}
		if f.ID == id {
			found = true
		}
	}
	return out
}

func (b *Broadcaster) ServeWS(snapshotFn func() (any, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{ //nolint:staticcheck
			InsecureSkipVerify: true,
		})
		if err != nil {
			return
		}
		defer conn.CloseNow() //nolint:errcheck,staticcheck

		ctx := r.Context()
		sinceID := r.URL.Query().Get("since")

		connFrame, _ := newFrame(string(EventConnected), map[string]string{"serverTime": time.Now().UTC().Format(time.RFC3339Nano)})
		_ = wsjson.Write(ctx, conn, connFrame)

		if sinceID != "" {
			for _, f := range b.since(sinceID) {
				if err := wsjson.Write(ctx, conn, f); err != nil {
					return
				}
			}
		} else if snapshotFn != nil {
			if snap, err := snapshotFn(); err == nil {
				snapFrame, _ := newFrame(string(EventSnapshot), snap)
				_ = wsjson.Write(ctx, conn, snapFrame)
			}
		}

		ch := make(chan WSFrame, clientBufSize)
		b.mu.Lock()
		b.clients[ch] = struct{}{}
		b.mu.Unlock()

		defer func() {
			b.mu.Lock()
			delete(b.clients, ch)
			b.mu.Unlock()
		}()

		go func() {
			for {
				var msg map[string]string
				if err := wsjson.Read(ctx, conn, &msg); err != nil {
					return
				}
				if msg["type"] == "ping" {
					pong, _ := newFrame(string(EventPong), map[string]string{})
					_ = wsjson.Write(ctx, conn, pong)
				}
			}
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case frame, ok := <-ch:
				if !ok {
					return
				}
				if err := wsjson.Write(ctx, conn, frame); err != nil {
					return
				}
			}
		}
	}
}

func (b *Broadcaster) ServeHTTP(snapshotFn func() (any, error)) func(http.ResponseWriter, *http.Request) {
	return b.ServeWS(snapshotFn)
}

func PublishFrame(ctx context.Context, b *Broadcaster, t EventType, payload any) {
	_ = b.Publish(string(t), payload)
}
