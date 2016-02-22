module Store {
  export class Watcher {
    test():string {
      return 'foo';
    }
    // constructor() {
    //   var s = new WebSocket("ws://localhost:9000/watcher")
    //   s.addEventListener('error', function(e) {
    //     console.log('error:', e)
    //   }, false);
    //   s.addEventListener('open', function(e) {
    //     console.log('open:', e)
    //   }, false);
    //   s.addEventListener('message', function(e) {
    //     console.log('message:', e)
    //   }, false);
    //   s.addEventListener('close', function(e) {
    //     console.log('close:', e)
    //   }, false);
    // }
  }
}
