
const React = require('react'),
      ReactDOM = require('react-dom'),
      AuthOverlay = require('./auth_overlay.js');

const App = React.createClass({
    getInitialState: function () {
        this.checkIfAuthenticated();
        return {
            initCredentialsCheck: 0,
            userCredentialsSet: false,
            botCredentialsSet: false,
        };
    },

    handleCredentialsSetEvent: function (payload) {
        var setValues = {};
        setValues[payload.kind + "CredentialsSet"] = payload.result;
        setValues.initCredentialsCheck = this.state.initCredentialsCheck + 1;
        if (setValues.initCredentialsCheck >= 2) {
            setValues.initCredentialsCheck = 0;
            listeners.remove("has-credentials-set", this.handleCredentialsSetEvent);
        }
        this.setState(setValues);
    },
    handleAuth: function (credentials) {
        conn.sendUTF(JSON.stringify({
            "cmd": "set-credentials",
            "payload": {
                "kind": "bot",
                "username": credentials.botUsername,
                "password": credentials.botPassword,
            },
        }));
        conn.sendUTF(JSON.stringify({
            "cmd": "set-credentials",
            "payload": {
                "kind": "user",
                "username": credentials.channelUsername,
                "password": credentials.channelPassword,
            },
        }));
        this.checkIfAuthenticated();
    },
    checkIfAuthenticated: function () {
        listeners.add("has-credentials-set", this.handleCredentialsSetEvent);
        conn.sendUTF(JSON.stringify({
            "cmd": "has-credentials-set",
            "payload": "user",
        }));
        conn.sendUTF(JSON.stringify({
            "cmd": "has-credentials-set",
            "payload": "bot",
        }));
    },
    authenticated: function () {
        return this.state.userCredentialsSet && this.state.botCredentialsSet;
    },

    render: function () {
        if (!this.authenticated()) {
            return (
                <div>
                    <AuthOverlay parent={this} />
                    <span>Some Content</span>
                </div>
            );
        }
        return (
            <div>
                <span>Some Content</span>
            </div>
        );
    },
});

function render() {
    ReactDOM.render(<App />, document.querySelector('#app'));
}

exports.render = render;
