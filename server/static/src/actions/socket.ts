import { createAction } from 'redux-actions';
import { WILL_OPEN_SOCKET, DID_OPEN_SOCKET } from '../constants/ActionTypes';
import { Promise } from 'es6-promise';
import * as EventEmitter from 'eventemitter3';

export class Socket extends EventEmitter {

  s: WebSocket

  constructor() {
    super()
  }

  open() {
    return new Promise((resolve, reject) => {
      this.s = new WebSocket("ws://localhost:9000/socket")
      this.addEventListener()
      const cb = (e) => {
        this.s.removeEventListener('open', cb, false)
        resolve()
      }
      this.s.addEventListener('open', cb, false)
    })
  }

  close() {
    this.removeEventListeners()
    // this.s.close()
  }

  addEventListener() {
    this.removeEventListeners()
    this.s.addEventListener('error', this.onError, false);
    this.s.addEventListener('open', this.onOpen, false);
    this.s.addEventListener('close', this.onClose, false);
    this.s.addEventListener('message', this.onMessage, false);
  }

  removeEventListeners() {
    this.s.removeEventListener('error', this.onError, false);
    this.s.removeEventListener('open', this.onOpen, false);
    this.s.removeEventListener('close', this.onClose, false);
    this.s.removeEventListener('message', this.onMessage, false);
  }

  send(req: Req) {
    this.s.send(JSON.stringify(req))
  }

  call(method: string, data?: any) {
    return new Promise((resolve, reject) => {
      let req = new Req(method, data);
      let cb = (e: MessageEvent) => {
        let res = JSON.parse(e.data) as Res
        if (req.id === res.id) {
          this.s.removeEventListener('message', cb, false)
          if (res.error != null) {
            reject(res.error)
            return
          }
          resolve(res.data)
        }
      }
      this.s.addEventListener('message', cb, false)
      this.send(req);
    })
  }

  reopen() {
    console.log("websocket: reopen")
    this.close()
    this.open()
  }

  onError = (e) => {
    console.log("websocket: error")
    // this.reopen()
  }

  onOpen = (e) => {
    console.log("websocket: opened")
  }

  onClose = (e) => {
    console.log("websocket: closed")
    // this.reopen()
  }

  onMessage = (e) => {
    console.log("websocket: messaged")
    let res = JSON.parse(e.data) as Res
    this.emit(res.method, res.data)
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

export const openingSocket = createAction<void>(
  WILL_OPEN_SOCKET
)

export const openedSocket = createAction<void>(
  DID_OPEN_SOCKET
)

export let socket

export const openSocket = () => {
  return (dispatch, getState) => {
    dispatch(openingSocket())
    socket = new Socket()
    return socket.open()
      .then(() => {
        dispatch(openedSocket())
      })
  }
}
