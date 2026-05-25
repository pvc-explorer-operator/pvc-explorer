import { ref } from 'vue';
import { useExplorerStore } from '../stores/explorerStore';

interface WsFrame {
  id: string;
  type: string;
  serverTime: string;
  payload: any;
}

export interface IdleTickPayload {
  namespace: string;
  name: string;
  remainingSeconds: number;
  idleTimeout: string;
}

export interface IdleWarningPayload {
  namespace: string;
  name: string;
  remainingSeconds: number;
  warningThreshold: number;
}

export interface IdleExpiredPayload {
  namespace: string;
  name: string;
  expiredAt: string;
}

export interface AgentReadyPayload {
  namespace: string;
  name: string;
}

export interface ConsumerEventPayload {
  namespace: string;
  pvcName: string;
  podName: string;
  nodeName: string;
  ownerKind: string;
  ownerName: string;
  readOnly: boolean;
}

export interface WebSocketCallbacks {
  onIdleTick?: (payload: IdleTickPayload) => void;
  onIdleWarning?: (payload: IdleWarningPayload) => void;
  onIdleExpired?: (payload: IdleExpiredPayload) => void;
  onAgentReady?: (payload: AgentReadyPayload) => void;
  onConsumerAttached?: (payload: ConsumerEventPayload) => void;
  onConsumerDetached?: (payload: ConsumerEventPayload) => void;
}

export function useWebSocket(callbacks?: WebSocketCallbacks) {
  const connected = ref(false);
  const lastEventId = ref('');
  let ws: WebSocket | null = null;
  let reconnectTimeout: ReturnType<typeof setTimeout> | null = null;
  let manualDisconnect = false;
  let backoff = 1000;
  const maxBackoff = 30000;
  const explorerStore = useExplorerStore();

  function getWsUrl(): string {
    const loc = window.location;
    const proto = loc.protocol === 'https:' ? 'wss:' : 'ws:';
    let url = proto + '//' + loc.host + '/ws/v1/events';
    if (lastEventId.value) {
      url += '?since=' + encodeURIComponent(lastEventId.value);
    }
    return url;
  }

  function connect(): void {
    manualDisconnect = false;
    if (ws) ws.close();
    ws = new WebSocket(getWsUrl());
    ws.onopen = () => {
      connected.value = true;
      backoff = 1000;
    };
    ws.onclose = () => {
      connected.value = false;
      ws = null;
      if (!manualDisconnect) {
        reconnectTimeout = setTimeout(connect, backoff);
        backoff = Math.min(backoff * 2, maxBackoff);
      }
    };
    ws.onerror = () => {
      ws?.close();
    };
    ws.onmessage = (ev) => {
      let frame: WsFrame;
      try {
        frame = JSON.parse(ev.data);
      } catch {
        return;
      }
      if (frame.id) lastEventId.value = frame.id;
      switch (frame.type) {
        case 'snapshot':
          explorerStore.setSnapshot(frame.payload.explorers, frame.payload.scopes);
          break;
        case 'explorer.updated':
          explorerStore.upsertExplorer(frame.payload);
          break;
        case 'explorer.deleted':
          explorerStore.removeExplorer(frame.payload.namespace, frame.payload.name);
          break;
        case 'scope.updated':
          explorerStore.upsertScope(frame.payload);
          break;
        case 'scope.deleted':
          explorerStore.removeScope(frame.payload.name);
          break;
        case 'ping':
          ws?.send(JSON.stringify({ type: 'pong' }));
          break;
        case 'idle.tick':
          callbacks?.onIdleTick?.(frame.payload as IdleTickPayload);
          break;
        case 'idle.warning':
          callbacks?.onIdleWarning?.(frame.payload as IdleWarningPayload);
          break;
        case 'idle.expired':
          callbacks?.onIdleExpired?.(frame.payload as IdleExpiredPayload);
          break;
        case 'agent.ready':
          callbacks?.onAgentReady?.(frame.payload as AgentReadyPayload);
          break;
        case 'consumer.attached':
          callbacks?.onConsumerAttached?.(frame.payload as ConsumerEventPayload);
          break;
        case 'consumer.detached':
          callbacks?.onConsumerDetached?.(frame.payload as ConsumerEventPayload);
          break;
      }
    };
  }

  function disconnect(): void {
    manualDisconnect = true;
    if (reconnectTimeout) {
      clearTimeout(reconnectTimeout);
      reconnectTimeout = null;
    }
    if (ws) {
      ws.close();
      ws = null;
    }
    connected.value = false;
  }

  return { connect, disconnect, connected, lastEventId };
}
