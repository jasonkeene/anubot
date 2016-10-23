
const uuid = require('./uuid.js');

const Listeners = function () {
    this.cmdSet = {};
    this.cmdListeners = {};
    this.requestListeners = {};
}

Listeners.prototype.cmd = function (cmd, listener) {
    var id = uuid();

    this.cmdSet[id] = cmd;

    if (this.cmdListeners[cmd] === undefined) {
        this.cmdListeners[cmd] = [];
    }
    this.cmdListeners[cmd].push({
        listener,
        id,
    });
    return id;
};

Listeners.prototype.request = function (request_id, listener) {
    this.requestListeners[request_id] = listener;
    return request_id;
};

Listeners.prototype.remove = function (id) {
    if (this.requestListeners[id] !== undefined) {
        delete this.requestListeners[id];
        return;
    }

    var cmd = this.cmdSet[id];
    if (cmd !== undefined) {
        delete this.cmdSet[id];
        var listeners = this.cmdListeners[cmd];

        if (listeners === undefined) {
            return;
        }

        for (var i = 0; i < listeners.length; i++) {
            var l = listeners[i];
            if (l.id === id) {
                listeners.splice(i, 1);
                return;
            }
        }
    }
};

Listeners.prototype.dispatch = function (cmd, request_id, payload, error) {
    if (request_id) {
        var l = this.requestListeners[request_id];
        if (l !== undefined) {
            l(payload, error);
        }
    }

    var listeners = this.cmdListeners[cmd];
    if (listeners === undefined) {
        return;
    }
    for (var i = 0; i < listeners.length; i++) {
        listeners[i].listener(payload, error);
    }
};

module.exports = Listeners;
