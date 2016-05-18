
const websocket = require('websocket');

var client = new websocket.client();

client.on('connectFailed', function(error) {
    console.log('Connect Error: ' + error.toString());
});

var conn; // used for debugging for the time being

client.on('connect', function(connection) {
    conn = connection;
    console.log('WebSocket Client Connected');

    connection.on('error', function(error) {
        console.log("Connection Error: " + error.toString());
    });
    connection.on('close', function() {
        console.log('echo-protocol Connection Closed');
    });
    connection.on('message', function(message) {
        console.log("Received: '" + message.utf8Data + "'");
    });

    function requestIsLoggedIn() {
        connection.sendUTF(JSON.stringify({
            "cmd": "is_logged_in",
        }));
    }
    requestIsLoggedIn();
});

client.connect('ws://localhost:12345/api', '', 'http://localhost:12345');
