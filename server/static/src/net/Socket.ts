import { Promise } from 'es6-promise';

export class Socket {

  private static _instance: Socket
  private static _internal: boolean

  static getInstance(): Socket {
    if (Socket._instance == null) {
      Socket._internal = true
      Socket._instance = new Socket()
      Socket._internal = false
    }
    return Socket._instance
  }

  private _ws: WebSocket
  private _dispatch: Function

  set dispatch(dispatch: Function) {
    this._dispatch = dispatch
  }

  constructor() {
    if (!Socket._internal) {
      throw new Error('Socket is implemented with singleton pattern, use Socket.getInstance()')
    }
  }

  open() {
    return new Promise((resolve, reject) => {
      this._ws = new WebSocket("ws://localhost:9000/socket")
      this.addEventListener()
      const cb = (e) => {
        this._ws.removeEventListener('open', cb, false)
        resolve()
      }
      this._ws.addEventListener('open', cb, false)
    })
  }

  close() {
    this.removeEventListeners()
    this._ws.close()
  }

  addEventListener() {
    this.removeEventListeners()
    this._ws.addEventListener('error', this.onError, false);
    this._ws.addEventListener('open', this.onOpen, false);
    this._ws.addEventListener('close', this.onClose, false);
    this._ws.addEventListener('message', this.onMessage, false);
  }

  removeEventListeners() {
    this._ws.removeEventListener('error', this.onError, false);
    this._ws.removeEventListener('open', this.onOpen, false);
    this._ws.removeEventListener('close', this.onClose, false);
    this._ws.removeEventListener('message', this.onMessage, false);
  }

  send(req: Req) {
    this._ws.send(JSON.stringify(req))
  }

  call(method: string, data?: any) {
    return new Promise((resolve, reject) => {
      let req = new Req(method, data);
      let cb = (e: MessageEvent) => {
        let res = JSON.parse(e.data) as Res
        if (req.id === res.id) {
          this._ws.removeEventListener('message', cb, false)
          if (res.error != null) {
            reject(res.error)
            return
          }
          resolve(res.data)
        }
      }
      this._ws.addEventListener('message', cb, false)
      this.send(req);
    })
  }

  reopen() {
    console.log("websocket: reopen")
    this.removeEventListeners()
    this.open()
  }

  onError = (e) => {
    console.log("websocket: error")
  }

  onOpen = (e) => {
    console.log("websocket: opened")
  }

  onClose = (e) => {
    console.log("websocket: closed")
    this.reopen()
  }

  onMessage = (e) => {
    console.log("websocket: messaged")
    let res = JSON.parse(e.data) as Res
    console.log(res);
    switch (res.method) {
    }
  }
}

class Req {
  id: string
  method: string
  data: any
  constructor(method: string, data?: any) {
    this.id = String(Date.now())
    this.method = method
    this.data = data
  }
}

class Res {
  id: string;
  method: string;
  data: any;
  error: any;
}
