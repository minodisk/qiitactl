declare module "websocket" {
  export class WebSocket {
    binaryType: string
    bufferedAmount: number
    extensions: string
    onclose: (event: CloseEvent) => void;
    onerror: (err: Error) => void;
    onmessage: (event: MessageEvent) => void;
    onopen: (event: Event) => void;
    protocol: string
    readyState: number
    url: string

    close(code?: number, reason?: string): void;
    send(data: string): void;

    addEventListener(method: string, listener?: () => void): void;
    addEventListener(method: 'close', cb?: (event: CloseEvent) => void): void;
    addEventListener(method: 'error', cb?: (err: Error) => void): void;
    addEventListener(method: 'message', cb?: (event: MessageEvent) => void): void;
    addEventListener(method: 'open', cb?: (event: Event) => void): void;
  }

  class Event {
    bubbles: boolean;
    cancelable: boolean;
    currentTarget: any;
    defaultPrevented: any;
    eventPhase: any;
    explicitOriginalTarget: any;
    originalTarget: any;
    srcElement: any;
    target: any;
    timeStamp: number;
    type: string;
    isTrusted: boolean;

    initEvent(type: string, bubbles: boolean, cancelable: boolean): void;
    preventDefault(): void;
    stopImmediatePropagation(): void;
    stopPropagation(): void;
  }

  class CloseEvent extends Event {
    code: number;
    reason: string;
    wasClean: boolean;
  }

  class MessageEvent extends Event {
    data: any;
    origin: string;
    ports: any;
    source: any;
  }
}
