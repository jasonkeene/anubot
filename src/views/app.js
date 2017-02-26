
const React = require('react'),
      electron = require('electron'),
      Setup = require('./setup.js'),
      TwitchStreamerSetup = require('./twitch_streamer_setup.js'),
      TwitchBotSetup = require('./twitch_bot_setup.js'),
      Menu = require('./menu.js'),
      ChatTab = require('./chat_tab.js'),
      emoji = require('./emoji.js');

const App = React.createClass({
    getInitialState: function () {
        return {
            loading: true,
            connected: false,
            authenticated: false,
            twitchStreamerAuthenticated: false,
            twitchBotAuthenticated: false,
            disconnected: false,

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
                (payload) => {
                    this.authenticated(creds);
                },
                (error) => {
                    console.log("got error while authenticating:", error);
                    this.setState({
                        loading: false,
                    });
                },
            );
            return;
        }
        this.setState({
            loading: false,
        });
    },
    authenticated: function (creds) {
        this.setLocalCredentials(creds);
        this.setState({
            authenticated: true,
            loading: true,
        });
        this.queryUserDetails();
    },
    queryUserDetails: function () {
        this.state.net.request("twitch-user-details", null).then(
            (payload) => {
                this.setState({
                    loading: false,
                });
                if (payload.streamer_authenticated) {
                    this.twitchStreamerAuthenticated(
                        payload.streamer_username,
                        payload.streamer_status,
                        payload.streamer_game,
                    );
                }
                if (payload.bot_authenticated) {
                    this.twitchBotAuthenticated(
                        payload.bot_username,
                    );
                }
            },
            (error) => {
                console.log("got error while getting user details:", error);
                this.setState({
                    loading: false,
                });
            },
        );
    },
    twitchStreamerAuthenticated: function (streamer_username, status, game) {
        this.setState({
            streamer_username,
            status,
            game,
            twitchStreamerAuthenticated: true,
        })
    },
    twitchBotAuthenticated: function (bot_username) {
        this.setState({
            bot_username,
            twitchBotAuthenticated: true,
        })
        this.finish();
    },
    setLocalCredentials: function (creds) {
        this.props.localStorage.setItem("username", creds.username);
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
    finish: function () {
        this.state.net.request("bttv-emoji").then(
            (payload) => {
                emoji.initBTTV(payload);
            },
            (error) => {
                console.log("got error while requesting BTTV emoji:", error);
            },
        )
        this.state.net.listeners.cmd("chat-message", (payload, error) => {
            var messages = this.state.messages;
            this.setState({
                messages: messages.concat([payload]),
            });
        });
        this.state.net.send({
            cmd: "twitch-stream-messages",
        });

        window.localStorage.setItem("window_bounds_collection", "true");
        var win = electron.remote.getCurrentWindow();
        win.setBounds(this.targetBounds(win.getBounds()), true);
        win.setResizable(true);
    },
    targetBounds: function (currentBounds) {
        var previousWidth = this.props.localStorage.getItem("window_bounds_width"),
            previousHeight = this.props.localStorage.getItem("window_bounds_height"),
            previousX = this.props.localStorage.getItem("window_bounds_x"),
            previousY = this.props.localStorage.getItem("window_bounds_y");
        if (
            previousWidth !== null &&
            previousHeight !== null &&
            previousX !== null &&
            previousY !== null
        ) {
            var foo = {
                width: parseInt(previousWidth),
                height: parseInt(previousHeight),
                x: parseInt(previousX),
                y: parseInt(previousY),
            };
            return foo;
        }
        var goalWidth = 1024,
            goalHeight = 768,
            widthDelta = goalWidth-currentBounds.width,
            heightDelta = goalHeight-currentBounds.height;
        var bar = {
            width: goalWidth,
            height: goalHeight,
            x: currentBounds.x - Math.floor(widthDelta/2),
            y: currentBounds.y - Math.floor(heightDelta/2),
        };
        return bar;
    },
    disconnect: function () {
        this.setState({
            disconnected: true,
        });
    },
    logout: function () {
        this.state.net.request("logout").then(
            () => {
                window.localStorage.removeItem("window_bounds_collection");
                var win = electron.remote.getCurrentWindow(),
                    bounds = win.getBounds(),
                    goalWidth = 560,
                    goalHeight = 620,
                    widthDelta = goalWidth-bounds.width,
                    heightDelta = goalHeight-bounds.height;
                win.setBounds({
                    width: goalWidth,
                    height: goalHeight,
                    x: bounds.x - Math.floor(widthDelta/2),
                    y: bounds.y - Math.floor(heightDelta/2),
                }, true);
                win.setResizable(false);
                var net = this.state.net;
                this.setLocalCredentials("", "");
                this.setState(this.getInitialState());
                this.connectionReady(net);
            },
        );
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
    renderDisconnected: function () {
        return <div id="disconnected">Disconnected</div>;
    },
    renderLoading: function () {
        return <div id="loading">Loading</div>;
    },
    renderSetup: function () {
        return <Setup parent={this} net={this.state.net} />;
    },
    renderTwitchStreamerSetup: function () {
        return <TwitchStreamerSetup parent={this} net={this.state.net} />;
    },
    renderTwitchBotSetup: function () {
        return <TwitchBotSetup parent={this} net={this.state.net} />;
    },
    renderNormal: function () {
        return [
            <Menu parent={this} selected={this.state.tab} key="menu" />,
            this.renderTab()
        ];
    },
    renderApp: function () {
        if (this.state.disconnected) {
            return this.renderDisconnected();
        }
        if (this.state.loading) {
            return this.renderLoading();
        }
        if (!this.state.connected) {
            return this.renderLoading();
        }
        if (!this.state.authenticated) {
            return this.renderSetup();
        }
        if (!this.state.twitchStreamerAuthenticated) {
            return this.renderTwitchStreamerSetup();
        }
        if (!this.state.twitchBotAuthenticated) {
            return this.renderTwitchBotSetup();
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
