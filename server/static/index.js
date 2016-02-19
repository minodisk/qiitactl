var output = document.querySelector(".markdown-body");
var websocket = new WebSocket("ws://localhost:9000/watcher");

websocket.addEventListener('error', function(e) {
  console.log('error:', e)
}, false);

websocket.addEventListener('open', function(e) {
  console.log('open:', e)
}, false);

websocket.addEventListener('message', function(e) {
  console.log('message:', e)
}, false);

websocket.addEventListener('close', function(e) {
  console.log('close:', e)
}, false);
