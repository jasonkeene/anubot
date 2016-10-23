/* global window: false */

const websocket = require('websocket'),
      Listeners = require('./lib/listeners.js'),
      Net = require('./lib/net.js'),
      unpack = require('./lib/unpack.js'),
      views = require('./lib/views/main.js');

const client = new websocket.client();

// save off net for development/debugging
var net;

client.on('connect', function(connection) {
    console.log('[app] WebSocket Client Connected');

    const listeners = new Listeners();
    net = new Net(listeners, connection);

    connection.on('message', function(message) {
        console.log("[app] Received:", JSON.parse(message.utf8Data));
        listeners.dispatch(...unpack(message.utf8Data));
    });
    connection.on('error', function(error) {
        console.log("[app] Connection Error: " + error.toString());
    });
    connection.on('close', function() {
        console.log('[app] echo-protocol Connection Closed');
    });

    views.render(net, window.localStorage);

});

client.on('connectFailed', function(error) {
    console.log('[app] Connect Error: ' + error.toString());
});

client.connect('wss://anubot.io:4443/api', '', 'https://anubot.io:4443');
