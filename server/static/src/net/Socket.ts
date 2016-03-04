import { Promise } from 'es6-promise';

import {
  openSocket,
  didCloseSocket,
  changeFile,
} from '../actions/socket'
// import { Stack } from '../models/stack'

export class Socket {

  private static _instance: Socket
  private static _internal: boolean
  private static _url:string = "ws://localhost:9000/socket"

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
  // private _stacks:Stack[]

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
      this._ws = new WebSocket(Socket._url)
      this.addEventListeners()
      const onOpened = (e) => {
        removeEventListeners()
        resolve()
      }
      const onClosed = (e) => {
        removeEventListeners()
      }
      const removeEventListeners = () => {
        this._ws.removeEventListener('open', onOpened, false)
        this._ws.removeEventListener('close', onClosed, false)
      }
      this._ws.addEventListener('open', onOpened, false)
      this._ws.addEventListener('close', onClosed, false)
    })
  }

  close() {
    this.removeEventListeners()
    this._ws.close()
  }

  addEventListeners() {
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

  // stack(s:Stack) {
  //   this._stacks.push(s)
  // }
  //
  // doStack() {
  //   while (this._stacks.length > 0) {
  //     const stack = this._stacks.shift()
  //     this[stack.method].apply(this, stack.args)
  //   }
  // }

  send(req: Req) {
    // if (this._ws.readyState !== WebSocket.OPEN) {
    //   this.stack({
    //     method: 'send',
    //     args: [req],
    //   })
    //   return
    // }
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

  private _timeoutID:number;

  reopen() {
    console.log("websocket: reopen")
    this.removeEventListeners()

    clearTimeout(this._timeoutID)
    this._timeoutID = setTimeout(() => {
      this._dispatch(openSocket(this))
    }, 1000)
  }

  onError = (e) => {
    console.log("websocket: error")
  }

  onOpen = (e) => {
    clearTimeout(this._timeoutID)
    console.log("websocket: opened")
  }

  onClose = (e) => {
    console.log("websocket: closed")
    this._dispatch(didCloseSocket())
    this.reopen()
  }

  onMessage = (e) => {
    // console.log("websocket: messaged")
    let res = JSON.parse(e.data) as Res
    switch (res.method) {
      case "ChangeFile":
        console.log(res.data)
        this._dispatch(changeFile(res.data))
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
