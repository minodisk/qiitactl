/// <reference path="../../typings/tsd.d.ts" />

export class Socket {
  s = new WebSocket("ws://localhost:9000/socket")

  constructor() {
    this.s.addEventListener('error', this.onError, false);
    this.s.addEventListener('open', this.onOpen, false);
    this.s.addEventListener('message', this.onMessage, false);
    this.s.addEventListener('close', this.onClose, false);
  }

  call(method: string, data?: any, callback?: (data?: any) => void) {
    let req = new Message(method, data);
    if (callback != null) {
      let cb = (e: MessageEvent) => {
        let res = JSON.parse(e.data) as Message
        if (req.id === res.id) {
          callback(res.data)
          this.s.removeEventListener('message', cb, false)
        }
      }
      this.s.addEventListener('message', cb, false)
    }
    this.s.send(JSON.stringify(req));
  }

  onError = (e) => {
    console.log('socket error:', e)
  }

  onOpen = (e) => {
    this.call('GetAllFiles', null, this.onGetAllFiles)
  }

  onGetAllFiles = (paths: string[]) => {
    console.log(paths)
  }

  onMessage = (e) => {
    console.log('socket messaged:', e)
  }

  onClose = (e) => {
    console.log('socket closed:', e)
  }
}

class Message {
  id: string
  method: string
  data: any
  constructor(method: string, data?: any) {
    this.id = String(Date.now())
    this.method = method
    this.data = data
  }
}