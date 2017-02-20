/* global window: false */

const websocket = require('websocket'),
      Listeners = require('../lib/listeners.js'),
      Net = require('../lib/net.js'),
      unpack = require('../lib/unpack.js'),
      views = require('../lib/views/main.js'),
      settings = require('electron-settings'),
      context_menu = require('../lib/context_menu.js');

// save off net for development/debugging
var net;

const client = new websocket.client();

client.on('connect', function(connection) {
    console.log('[app] WebSocket Client Connected');

    const listeners = new Listeners();
    net = new Net(listeners, connection);
    app.connectionReady(net);

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
});

client.on('connectFailed', function(error) {
    console.log('[app] Connect Error: ' + error.toString());
});

var app = views.render(window.localStorage);
context_menu.register(window);
client.connect(settings.getSync('api'), '', settings.getSync('origin'));
