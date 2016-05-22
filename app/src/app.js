
const websocket = require('websocket'),
      Listeners = require('./lib/listeners.js'),
      unpack = require('./lib/unpack.js'),
      views = require('./lib/views/main.js');

// TODO: need to consider mutation of this by other connects
var conn;

const client = new websocket.client(),
      listeners = new Listeners();

client.on('connect', function(connection) {
    conn = connection;
    views.render();
    console.log('WebSocket Client Connected');

    connection.on('message', function(message) {
        console.log("Received: '" + message.utf8Data + "'");
        listeners.dispatch(...unpack(message.utf8Data));
    });
    connection.on('error', function(error) {
        console.log("Connection Error: " + error.toString());
    });
    connection.on('close', function() {
        console.log('echo-protocol Connection Closed');
    });
});

client.on('connectFailed', function(error) {
    console.log('Connect Error: ' + error.toString());
});

client.connect('ws://localhost:12345/api', '', 'http://localhost:12345');
