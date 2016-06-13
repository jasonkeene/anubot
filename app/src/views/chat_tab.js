
const React = require('react'),
      ReactDOM = require('react-dom');

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

    renderMessage: function (message) {
        return (
            <div className="message" key={message.id}>
                <span className="nick">{message.nick}:</span>&nbsp;
                {message.body}
            </div>
        );
    },
    render: function () {
        return (
            <div id="chat-tab" className="tab">
                <div className="header">
                    <div className="channel">#postcrypt</div>
                    <div className="description">Working on my bot today - https://github.com/jasonkeene/anubot</div>
                </div>
                <div className="body">
                    <div className="spacer" />
                    {this.props.messages.map(this.renderMessage)}
                </div>
                <ChatFooter listeners={this.props.listeners}
                            connection={this.props.connection} />
            </div>
        );
    },
});

const ChatFooter = React.createClass({
    getInitialState: function () {
        return {
            user: "streamer",
            message: "",
            streamerUsername: "streamer",
            botUsername: "bot",
        };
    },
    componentWillMount: function () {
        this.props.listeners.add("usernames", this.handleUsernames);
        this.props.connection.sendUTF(JSON.stringify({
            "cmd": "usernames",
        }));
    },

    // network events
    handleUsernames: function (payload) {
        if (payload.streamer === "" || payload.bot === "") {
            this.props.connection.sendUTF(JSON.stringify({
                "cmd": "usernames",
            }));
            return
        }
        this.props.listeners.remove("usernames", this.handleUsernames);
        this.setState({
            streamerUsername: payload.streamer,
            botUsername: payload.bot,
        })
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
                        <option value="streamer">{this.state.streamerUsername}</option>
                        <option value="bot">{this.state.botUsername}</option>
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
