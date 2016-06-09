
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
                <div className="footer">
                    <div className="selection">
                        <select>
                            <option>postcrypt</option>
                            <option>pc_anubot</option>
                        </select>
                    </div>
                    <div className="input">
                        <input type="text" placeholder="Enter a message here" />
                        <span></span>
                    </div>
                </div>
            </div>
        );
    },
});

module.exports = ChatTab;
