/* global window: false */

const websocket = require('websocket'),
      Listeners = require('../lib/listeners.js'),
      Net = require('../lib/net.js'),
      unpack = require('../lib/unpack.js'),
      views = require('../lib/views/main.js'),
      electron = require('electron'),
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
        console.log("[app] Connection Errored: " + error.toString());
        app.disconnect();
    });
    connection.on('close', function() {
        console.log('[app] Connection Closed');
        app.disconnect();
    });
});

client.on('connectFailed', function(error) {
    console.log('[app] Connect Failed: ' + error.toString());
    app.disconnect();
});

var win = electron.remote.getCurrentWindow();
function collect() {
    var b = win.getBounds();
    if (window.localStorage.getItem("window_bounds_collection") !== null) {
        window.localStorage.setItem("window_bounds_width", b.width);
        window.localStorage.setItem("window_bounds_height", b.height);
        window.localStorage.setItem("window_bounds_x", b.x);
        window.localStorage.setItem("window_bounds_y", b.y);
    }
}
win.on("resize", collect);
win.on("move", collect);

const app = views.render(window.localStorage, connect);
context_menu.register(window);

function connect() {
    client.connect(settings.getSync('api'), '', settings.getSync('origin'));
}
connect();
