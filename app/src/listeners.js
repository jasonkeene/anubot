
const Listeners = function () {
    this.listenersMap = {};
}

Listeners.prototype.add = function (cmd, listener) {
    if (this.listenersMap[cmd] === undefined) {
        this.listenersMap[cmd] = [];
    }
    this.listenersMap[cmd].push(listener);
};

Listeners.prototype.remove = function (cmd, listener) {
    var l = this.listenersMap[cmd];
    if (l === undefined) {
        // no listeners for this cmd
        return;
    }
    for (var i = 0; i < l.length; i++) {
        var j = -1;
        while ((j = l.indexOf(listener)) !== -1) {
            l.splice(j, 1);
        }
    }
};

Listeners.prototype.dispatch = function (cmd, payload, error) {
    var l = this.listenersMap[cmd];
    if (l === undefined) {
        // no listeners for this cmd
        return;
    }
    for (var i = 0; i < l.length; i++) {
        l[i](payload, error);
    }
};

module.exports = Listeners;
