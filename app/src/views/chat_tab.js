
const React = require('react'),
      ReactDOM = require('react-dom'),
      badges = require('./badges.js'),
      emoji = require('./emoji.js');

const ChatTab = React.createClass({
    getInitialState: function () {
        return {
            scroll: true,
        };
    },
    componentDidMount: function () {
        var domNode = ReactDOM.findDOMNode(this).querySelector('.body');
        this.domNode = domNode;
        domNode.addEventListener('scroll', this.handleScroll);
        this.scrollToBottom();
    },
    componentWillUnmount: function() {
        this.domNode.removeEventListener('scroll', this.handleScroll);
    },
    componentDidUpdate: function () {
        this.scrollToBottom();
    },

    scrollToBottom: function () {
        if (this.state.scroll) {
            this.domNode.scrollTop = this.domNode.scrollHeight;
        }
    },
    scrollIsAtBottom: function () {
        return this.domNode.scrollTop + this.domNode.clientHeight == this.domNode.scrollHeight;
    },

    // event handlers
    handleScroll: function (event) {
        if (this.scrollIsAtBottom()) {
            this.setState({scroll: true});
            return
        }
        this.setState({scroll: false});
    },

    nickStyle: (message) => {
        if (message.tags === undefined) {
            return {};
        }
        if (message.tags.color === undefined) {
            return {};
        }
        return {color: message.tags.color};
    },
    renderMessage: function (message) {
        return (
            <div className="message" key={message.tags.id}>
                <span className="badges">{badges.render(message)}</span>
                <span className="nick" style={this.nickStyle(message)}>{message.tags['display-name']}:</span>&nbsp;
                {emoji.render(message)}
            </div>
        );
    },
    render: function () {
        return (
            <div id="chat-tab" className="tab">
                <ChatHeader channel={"#" + this.props.streamer}
                            status={this.props.status}
                            game={this.props.game}
                            connection={this.props.connection} />
                <div className="body">
                    <div className="spacer" />
                    {this.props.messages.map(this.renderMessage)}
                </div>
                <ChatFooter streamer={this.props.streamer}
                            bot={this.props.bot}
                            listeners={this.props.listeners}
                            connection={this.props.connection} />
            </div>
        );
    },
});

const ChatHeader = React.createClass({
    getInitialState: function () {
        return {
            editing: false,
            status: this.props.status,
            game: this.props.game,
        };
    },

    // event handlers
    handleEditClick: function () {
        this.setState({
            editing: true,
        });
    },
    handleStatusChange: function (e) {
        this.setState({
            status: e.target.value,
        });
    },
    handleGameChange: function (e) {
        this.setState({
            game: e.target.value,
        });
    },
    handleSubmit: function (e) {
        e.preventDefault();
        this.props.connection.sendUTF(JSON.stringify({
            "cmd": "update-description",
            "payload": {
                "status": this.state.status,
                "game": this.state.game,
            },
        }));
        this.setState({
            editing: false,
        });
    },

    render: function () {
        if (this.state.editing) {
            return (
                <div className="header">
                    <div className="channel">{this.props.channel}</div>
                    <div className="description">
                        <form onSubmit={this.handleSubmit}>
                            <input type="text" className="game-input" onChange={this.handleGameChange} value={this.state.game} />&nbsp;
                            <input type="text" className="status-input" onChange={this.handleStatusChange} value={this.state.status} />&nbsp;
                            <input type="submit" value="Done" />
                        </form>
                    </div>
                </div>
            );
        }
        return (
            <div className="header">
                <div className="channel">{this.props.channel}</div>
                <div className="description" onClick={this.handleEditClick}>
                    <b>{this.state.game}:</b>&nbsp;
                    {this.state.status}
                    &nbsp;<i className="material-icons">edit</i>
                </div>
            </div>
        );
    },
});

const ChatFooter = React.createClass({
    getInitialState: function () {
        return {
            user: "streamer",
            message: "",
        };
    },

    // event handlers
    handleSubmit: function (e) {
        e.preventDefault();
        this.props.connection.sendUTF(JSON.stringify({
            "cmd": "send-message",
            "payload": {
                "user": this.state.user,
                "message": this.state.message,
            },
        }));
        this.setState({message: ""});
    },
    handleMessageChange: function (e) {
        this.setState({message: e.target.value});
    },
    handleUserChange: function (e) {
        this.setState({user: e.target.value});
    },

    render: function () {
        return (
            <div className="footer">
                <div className="selection">
                    <select onChange={this.handleUserChange}>
                        <option value="streamer">{this.props.streamer}</option>
                        <option value="bot">{this.props.bot}</option>
                    </select>
                </div>
                <div className="input">
                    <div className="form">
                        <form onSubmit={this.handleSubmit}>
                            <input onChange={this.handleMessageChange} type="text" placeholder="Enter a message here" value={this.state.message} />
                        </form>
                    </div>
                    <div className="spacer"></div>
                </div>
            </div>
        );
    },
});

module.exports = ChatTab;
