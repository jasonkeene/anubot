
const websocket = require('websocket'),
      Listeners = require('./lib/listeners.js'),
      unpack = require('./lib/unpack.js'),
      views = require('./lib/views/main.js');

const client = new websocket.client();

// save off connection for development/debugging
var conn;

client.on('connect', function(connection) {
    conn = connection;
    const listeners = new Listeners();
    views.render(connection, listeners);
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
