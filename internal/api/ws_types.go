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
	"encoding/json"
	"time"
)

type EventType string

const (
	EventConnected        EventType = "connected"
	EventSnapshot         EventType = "snapshot"
	EventExplorerUpdated  EventType = "explorer.updated"
	EventExplorerDeleted  EventType = "explorer.deleted"
	EventScopeUpdated     EventType = "scope.updated"
	EventScopeDeleted     EventType = "scope.deleted"
	EventConsumerAttached EventType = "consumer.attached"
	EventConsumerDetached EventType = "consumer.detached"
	EventAgentWaking      EventType = "agent.waking"
	EventAgentReady       EventType = "agent.ready"
	EventAgentError       EventType = "agent.error"
	EventPing             EventType = "ping"
	EventPong             EventType = "pong"
	EventIdleTick         EventType = "idle.tick"
	EventIdleWarning      EventType = "idle.warning"
	EventIdleExpired      EventType = "idle.expired"
)

type WSFrame struct {
	ID         string          `json:"id"`
	Type       EventType       `json:"type"`
	ServerTime time.Time       `json:"serverTime"`
	Payload    json.RawMessage `json:"payload"`
}
