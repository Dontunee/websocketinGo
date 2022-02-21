# WebsocketinGo
# Implement websocket in Golang using gorilla/mux


## How to test from browser
 - open console in browser
 - enter the below commands 
 
 1. `const client = new WebSocket("ws://localhost:3000/v1/ws");`
// This is to initiate connection to websocket 

2.  `client.addEventListener('message', function (event) { 
  console.log('Message from server ', event.data); 
}); `
// This is to display messages from server 
3. `client.send(JSON.stringify({
    "userName": "Tunde",
    "firstName": "Afolabi",
    "lastName": "tunde"
}));` // This is to send json request


