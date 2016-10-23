
const React = require('react'),
      Setup = require('./setup.js'),
      Menu = require('./menu.js'),
      ChatTab = require('./chat_tab.js');

const App = React.createClass({
    getInitialState: function () {
        return {
            loaded: false,
            authenticated: false,

            tab: "chat",
            messages: [],

            streamer_username: "",
            bot_username: "",
            status: "",
            game: "",
        };
    },
    componentWillMount: function () {
        var credentials = this.getLocalCredentials();
        if (credentials !== null) {
            this.props.listeners.add("authenticate", this.handleAuthenticate);
            console.log("authenticating");
            this.props.connection.sendUTF(JSON.stringify({
                cmd: "authenticate",
                payload: credentials,
            }));
            return;
        }
        this.setState({
            loaded: true,
        });
    },

    // network events
    handleAuthenticate: function (payload, error) {
        if (error === null) {
            this.setState({
                loaded: true,
                authenticated: true,
            });
            this.props.listeners.add("twitch-user-details", this.handleUserDetails);
            this.props.connection.sendUTF(JSON.stringify({
                cmd: "twitch-user-details",
            }));
            this.props.listeners.add("chat-message", this.handleChatMessage);
            this.props.connection.sendUTF(JSON.stringify({
                cmd: "twitch-stream-messages",
            }));
            return
        }
        this.setState({
            loaded: true,
        })
    },
    handleUserDetails: function (payload, error) {
        this.setState({
            streamer_username: payload.streamer_username,
            bot_username: payload.bot_username,
            status: payload.streamer_status,
            game: payload.streamer_game,
        });
    },
    handleChatMessage: function (payload, error) {
        var messages = this.state.messages;
        messages.push(payload);
        this.setState({messages: messages});
    },

    getLocalCredentials: function () {
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

    renderTab: function () {
        switch (this.state.tab) {
        case "chat":
            return <ChatTab streamer_username={this.state.streamer_username}
                            bot_username={this.state.bot_username}
                            status={this.state.status}
                            game={this.state.game}
                            messages={this.state.messages}
                            listeners={this.props.listeners}
                            connection={this.props.connection}
                            key="chat-tab" />;
        default:
            return <div className="tab">Content for {this.state.tab} tab!</div>;
        }
    },
    renderLoading: function () {
        return <div id="loading">Loading</div>;
    },
    renderSetup: function () {
        return <Setup parent={this}
                      listeners={this.props.listeners}
                      connection={this.props.connection} />;
    },
    renderNormal: function () {
        return [
            <Menu parent={this} selected={this.state.tab} key="menu" />,
            this.renderTab()
        ];
    },
    renderApp: function () {
        if (!this.state.loaded) {
            return this.renderLoading();
        }
        if (!this.state.authenticated) {
            return this.renderSetup();
        }
        return this.renderNormal();
    },
    render: function () {
        return <div id="app">
            {this.renderApp()}
        </div>;
    },
});

module.exports = App;
