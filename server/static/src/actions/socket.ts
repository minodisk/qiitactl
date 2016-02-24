import { createAction } from 'redux-actions';
import * as types from '../constants/ActionTypes';
import { Promise } from 'es6-promise';

export class Socket {

  s: WebSocket

  constructor() {}

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
    this.s.close()
    this.removeEventListeners()
  }

  addEventListener() {
    this.removeEventListeners()
    this.s.addEventListener('error', this.onError, false);
    this.s.addEventListener('open', this.onOpen, false);
    this.s.addEventListener('message', this.onMessage, false);
    this.s.addEventListener('close', this.onClose, false);
  }

  removeEventListeners() {
    this.s.removeEventListener('error', this.onError, false);
    this.s.removeEventListener('open', this.onOpen, false);
    this.s.removeEventListener('message', this.onMessage, false);
    this.s.removeEventListener('close', this.onClose, false);
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
      // console.log('send', req)
      this.s.send(JSON.stringify(req));
    })
  }

  onError = (e) => {
    // console.log('socket error:', e)
  }

  onOpen = (e) => {
    // console.log('socket opened:', e)
  }

  onMessage = (e) => {
    // console.log('socket messaged:', e)
  }

  onClose = (e) => {
    // console.log('socket closed:', e)
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
  types.OPEN_SOCKET
)

export const openedSocket = createAction<void>(
  types.OPENED_SOCKET
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
