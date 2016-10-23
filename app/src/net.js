/* global window: false */

const uuid = require('./uuid.js');

const Net = function (listeners, connection) {
    this.listeners = listeners;
    this.connection =  connection;
}

Net.prototype.request = function (cmd, payload) {
    var that = this;
    var request_id = uuid();

    var p = new Promise((resolve, reject) => {
        var l = (result, error) => {
            if (error !== null) {
                reject(error);
                return;
            }
            resolve(result);
        };
        that.listeners.request(request_id, l);
        that.send({
            cmd,
            request_id,
            payload,
        });
        window.setTimeout(() => {
            reject("timeout");
        }, 4000);
    });

    var deregister = () => {
        that.listeners.remove(request_id);
    };
    p.then(deregister, deregister);

    return p;
};

Net.prototype.send = function (obj) {
    this.connection.sendUTF(JSON.stringify(obj));
    console.log("[app] Sent:", obj);
}

module.exports = Net;
