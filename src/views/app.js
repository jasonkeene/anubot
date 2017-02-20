
const React = require('react'),
      electron = require('electron'),
      Setup = require('./setup.js'),
      Menu = require('./menu.js'),
      ChatTab = require('./chat_tab.js'),
      emoji = require('./emoji.js');

const App = React.createClass({
    getInitialState: function () {
        return {
            loaded: false,
            connected: false,
            authenticated: false,

            tab: "chat",
            messages: [],

            streamer_username: "",
            bot_username: "",
            status: "",
            game: "",

            net: null,
        };
    },

    connectionReady: function (net) {
        this.setState({
            net,
            connected: true,
        });
        var creds = this.localCredentials();
        if (creds !== null) {
            net.request("authenticate", creds).then(
                this.handleAuthenticateSuccess(creds),
                this.handleAuthenticateFailure,
            );
            return;
        }
    },

    authenticated: function (creds) {
        this.setLocalCredentials(creds);
        this.setState({
            authenticated: true,
        });
        this.finishLoading();
    },

    // network events
    handleAuthenticateSuccess: function (creds) {
        var that = this;
        return function (payload) {
            that.authenticated(creds);
        };
    },
    handleAuthenticateFailure: function (error) {
        // TODO: handle failure
        console.log("got error while authenticating:", error);
    },
    handleUserDetailsSuccess: function (payload) {
        var win = electron.remote.getCurrentWindow(),
            bounds = win.getBounds(),
            goalWidth = 1024,
            goalHeight = 768,
            widthDelta = goalWidth-bounds.width,
            heightDelta = goalHeight-bounds.height;
        win.setBounds({
            width: goalWidth,
            height: goalHeight,
            x: bounds.x - widthDelta/2,
            y: bounds.y - heightDelta/2,
        }, true);
        this.setState({
            streamer_username: payload.streamer_username,
            bot_username: payload.bot_username,
            status: payload.streamer_status,
            game: payload.streamer_game,
            loaded: true,
        });
    },
    handleUserDetailsFailure: function (error) {
        // TODO: handle failure
        console.log("got error while getting user details:", error);
    },
    handleChatMessage: function (payload, error) {
        var messages = this.state.messages;
        this.setState({
            messages: messages.concat([payload]),
        });
    },

    setLocalCredentials: function (creds) {
        this.props.localStorage.setItem("username", creds.username),
        this.props.localStorage.setItem("password", creds.password);
    },
    localCredentials: function () {
        var username = this.props.localStorage.getItem("username"),
            password = this.props.localStorage.getItem("password");
        if (!username || !password) {
            return null;
        }
        return {
            username: username,
            password: password,
        };
    },
    finishLoading: function () {
        this.state.net.request("twitch-user-details", null).then(
            this.handleUserDetailsSuccess,
            this.handleUserDetailsFailure,
        );
        this.state.net.request("bttv-emoji").then((payload) => {
            emoji.initBTTV(payload);
        }, (error) => {
            console.log("got error while requesting BTTV emoji:", error);
        })

        this.state.net.listeners.cmd("chat-message", this.handleChatMessage);
        this.state.net.send({
            cmd: "twitch-stream-messages",
        });
    },

    renderTab: function () {
        switch (this.state.tab) {
        case "chat":
            return <ChatTab streamer_username={this.state.streamer_username}
                            bot_username={this.state.bot_username}
                            status={this.state.status}
                            game={this.state.game}
                            messages={this.state.messages}
                            net={this.state.net}
                            key="chat-tab" />;
        default:
            return <div className="tab">Content for {this.state.tab} tab!</div>;
        }
    },
    renderLoading: function () {
        return <div id="loading">Loading</div>;
    },
    renderSetup: function () {
        return <Setup parent={this} net={this.state.net} />;
    },
    renderNormal: function () {
        return [
            <Menu parent={this} selected={this.state.tab} key="menu" />,
            this.renderTab()
        ];
    },
    renderApp: function () {
        if (!this.state.connected) {
            return this.renderLoading();
        }
        if (!this.state.authenticated) {
            return this.renderSetup();
        }
        if (!this.state.loaded) {
            return this.renderLoading();
        }
        return this.renderNormal();
    },
    render: function () {
        return <div id="app">
            <div id="drag-area" />
            {this.renderApp()}
        </div>;
    },
});

module.exports = App;
