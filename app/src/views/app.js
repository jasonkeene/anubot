
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
                "cmd": "subscribe",
            }));
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
            return <ChatTab messages={this.state.messages}
                            listeners={this.props.listeners}
                            connection={this.props.connection} />;
        default:
            return <div className="tab">Content for {this.state.tab} tab!</div>;
        }
    },
    render: function () {
        if (!this.state.authenticated) {
            return (
                <div>
                    <AuthOverlay parent={this}
                                 listeners={this.props.listeners}
                                 connection={this.props.connection} />
                    <span>Some Content</span>
                </div>
            );
        }
        return (
            <div id="app">
                <Menu parent={this} selected={this.state.tab} />
                {this.renderTab()}
            </div>
        );
    },
});

module.exports = App;
