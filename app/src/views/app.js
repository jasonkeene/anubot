
const React = require('react'),
      AuthOverlay = require('./auth_overlay.js'),
      Menu = require('./menu.js'),
      ChatTab = require('./chat_tab.js');

const App = React.createClass({
    getInitialState: function () {
        return {
            authenticated: false,
            tab: "chat",
            messages: [],
            streamer: "",
            bot: "",
            status: "",
            game: "",
        };
    },

    // network events
    connect: function () {
        this.props.listeners.add("connect", this.handleConnect);
        this.props.connection.sendUTF(JSON.stringify({
            "cmd": "connect",
        }));
    },
    handleConnect: function (payload) {
        this.props.listeners.remove("connect", this.handleConnect);
        if (payload) {
            this.props.connection.sendUTF(JSON.stringify({
                // TODO: this event name changed
                "cmd": "subscribe",
            }));
            // TODO: this event name changed
            this.props.listeners.add("chat-message", this.handleChatMessage);
        }
    },
    handleChatMessage: function (payload) {
        var messages = this.state.messages;
        messages.push(payload);
        this.setState({messages: messages});
    },

    renderTab: function () {
        switch (this.state.tab) {
        case "chat":
            return <ChatTab streamer={this.state.streamer}
                            bot={this.state.bot}
                            status={this.state.status}
                            game={this.state.game}
                            messages={this.state.messages}
                            listeners={this.props.listeners}
                            connection={this.props.connection} />;
        default:
            return <div className="tab">Content for {this.state.tab} tab!</div>;
        }
    },
    renderAuthOverlay: function () {
        if (!this.state.authenticated) {
            return <AuthOverlay parent={this}
                                listeners={this.props.listeners}
                                connection={this.props.connection} />;
        }
        return null;
    },
    render: function () {
        return (
            <div id="app">
                {this.renderAuthOverlay()}
                <Menu parent={this} selected={this.state.tab} />
                {this.renderTab()}
            </div>
        );
    },
});

module.exports = App;
