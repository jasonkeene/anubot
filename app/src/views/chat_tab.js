
const React = require('react');

const ChatTab = React.createClass({
    getInitialState: function () {
        return {
        };
    },

    renderMessage: function (message) {
        return (
            <div className="message" key={message.id}>
                {message.nick}:&nbsp;
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
