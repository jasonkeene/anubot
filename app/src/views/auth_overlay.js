
const React = require('react');

const AuthOverlay = React.createClass({
    // lifecycle
    getInitialState: function () {
        return {
            // ui data
            mode: "choice",

            // auth data
            initCredentialsCheck: 0,
            userCredentialsSet: false,
            botCredentialsSet: false,
        };
    },
    componentWillMount: function () {
        this.checkIfAuthenticated();
    },

    // network events
    checkIfAuthenticated: function () {
        this.props.listeners.add("has-credentials-set", this.handleHasCredentialsSetEvent);
        this.props.connection.sendUTF(JSON.stringify({
            "cmd": "has-credentials-set",
            "payload": "user",
        }));
        this.props.connection.sendUTF(JSON.stringify({
            "cmd": "has-credentials-set",
            "payload": "bot",
        }));
    },
    handleHasCredentialsSetEvent: function (payload) {
        var setValues = {};
        setValues[payload.kind + "CredentialsSet"] = payload.result;
        setValues.initCredentialsCheck = this.state.initCredentialsCheck + 1;

        if (setValues.initCredentialsCheck >= 2) {
            // at this point we got two response events from the server
            setValues.initCredentialsCheck = 0;
            this.props.listeners.remove("has-credentials-set", this.handleHasCredentialsSetEvent);
        }

        // write state out
        // TODO: react docs say this isn't guarenteed to happen syncronously
        // so we can't rely on this.state values afterward
        this.setState(setValues);

        if (this.state.userCredentialsSet && this.state.botCredentialsSet) {
            // tell the parent we are authenticated
            this.props.parent.setState({authenticated: true});

            // connect to the remote irc server
            this.props.parent.connect();
        }
    },

    // events from clildren
    handleOptionClick: function (mode) {
        this.setState({mode: mode});
    },
    handleAuth: function (credentials) {
        this.props.connection.sendUTF(JSON.stringify({
            "cmd": "set-credentials",
            "payload": {
                "kind": "bot",
                "username": credentials.botUsername,
                "password": credentials.botPassword,
            },
        }));
        this.props.connection.sendUTF(JSON.stringify({
            "cmd": "set-credentials",
            "payload": {
                "kind": "user",
                "username": credentials.channelUsername,
                "password": credentials.channelPassword,
            },
        }));
        this.checkIfAuthenticated();
    },

    // rendering
    render: function () {
        switch (this.state.mode) {
        case "choice":
            return this.renderChoice();
        case "auto":
            return this.renderAuto();
        case "manual":
            return <ManualMode parent={this} />;
        }
    },
    renderChoice: function () {
        return <div id="auth-overlay">
            <ul>
                <AuthOption parent={this} value="auto" text="Login via Twitch" />
                <AuthOption parent={this} value="manual" text="Manually Enter Oauth Tokens" />
            </ul>
        </div>;
    },
    renderAuto: function () {
        return <div id="auth-overlay"> Auto </div>;
    },
});

const AuthOption = React.createClass({
    handleClick: function () {
        this.props.parent.handleOptionClick(this.props.value);
    },

    render: function () {
        return <li onClick={this.handleClick}>{this.props.text}</li>
    }
})

const ManualMode = React.createClass({
    getInitialState: function () {
        return {
            channelUsername: "",
            channelPassword: "",
            botUsername: "",
            botPassword: "",
            errors: [],
        };
    },
    handleUpdate: function (key, e) {
        var update = {};
        update[key] = e.target.value;
        this.setState(update);
    },

    handleSubmit: function (e) {
        e.preventDefault();
        var credentials = {
                channelUsername: this.state.channelUsername,
                channelPassword: this.state.channelPassword,
                botUsername: this.state.botUsername,
                botPassword: this.state.botPassword,
            },
            errors = [];
        if (credentials.channelUsername.length === 0) {
            errors.push({
                id: "bad-channel-username",
                text: "You didn't provide a channel username",
            });
        }
        if (!credentials.channelPassword.startsWith("oauth:") ||
            credentials.channelPassword.length < 7) {
            errors.push({
                id: "bad-channel-oauth",
                text: "Make sure your channel oauth token starts with 'oauth:'",
            });
        }
        if (credentials.botUsername.length === 0) {
            errors.push({
                id: "bad-bot-username",
                text: "You didn't provide a bot username",
            });
        }
        if (!credentials.botPassword.startsWith("oauth:") ||
            credentials.botPassword.length < 7) {
            errors.push({
                id: "bad-bot-oauth",
                text: "Make sure your bot oauth token starts with 'oauth:'",
            });
        }
        this.setState({errors: errors});
        if (errors.length === 0) {
            this.props.parent.handleAuth(credentials);
        }
    },

    renderError: function (error) {
        return <div key={error.id}>{error.text}</div>
    },
    render: function () {
        return (
            <div id="auth-overlay">
                <form onSubmit={this.handleSubmit}>
                    Channel Username: <input type="text" name="channelUsername" value={this.state.channelUsername} onChange={this.handleUpdate.bind(this, "channelUsername")} /><br />
                    Channel Oauth Token: <input type="password" name="channelPassword" value={this.state.channelPassword} onChange={this.handleUpdate.bind(this, "channelPassword")} /><br />
                    Bot Username: <input type="text" name="botUsername" value={this.state.botUsername} onChange={this.handleUpdate.bind(this, "botUsername")} /><br />
                    Bot Oauth Token: <input type="password" name="botPassword" value={this.state.botPassword} onChange={this.handleUpdate.bind(this, "botPassword")} /><br />
                    <input type="submit" value="Submit" />
                    <div class="errors">
                        {this.state.errors.map(this.renderError)}
                    </div>
                </form>
            </div>
        );
    }
})

module.exports = AuthOverlay;
