/* global window: false */

const websocket = require('websocket'),
      Listeners = require('./lib/listeners.js'),
      unpack = require('./lib/unpack.js'),
      views = require('./lib/views/main.js');

const client = new websocket.client();

// save off connection for development/debugging
var conn;

client.on('connect', function(connection) {
    console.log('[app] WebSocket Client Connected');

    conn = connection;
    const listeners = new Listeners();

    connection.on('message', function(message) {
        console.log("[app] Received: '" + message.utf8Data + "'");
        listeners.dispatch(...unpack(message.utf8Data));
    });
    connection.on('error', function(error) {
        console.log("[app] Connection Error: " + error.toString());
    });
    connection.on('close', function() {
        console.log('[app] echo-protocol Connection Closed');
    });

    views.render(connection, listeners, window.localStorage);
});

client.on('connectFailed', function(error) {
    console.log('[app] Connect Error: ' + error.toString());
});

client.connect('wss://anubot.io:4443/api', '', 'https://anubot.io:4443');
