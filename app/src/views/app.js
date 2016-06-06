
const React = require('react'),
      AuthOverlay = require('./auth_overlay.js'),
      Menu = require('./menu.js'),
      ChatTab = require('./chat_tab.js');

const App = React.createClass({
    getInitialState: function () {
        return {
            authenticated: false,
            tab: "chat",
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
        }
    },

    renderTab: function () {
        return <ChatTab />;
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
            <div>
                <Menu parent={this} />
                {this.renderTab()}
            </div>
        );
    },
});

module.exports = App;
