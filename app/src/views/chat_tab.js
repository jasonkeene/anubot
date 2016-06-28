
const React = require('react'),
      ReactDOM = require('react-dom'),
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

    renderBadges: function (message) {
        var tags = message.tags;
        if (tags === undefined) {
            return;
        }
        var badges = tags.badges;
        if (badges === undefined) {
            return;
        }
        badges = badges.split(",").map((x) => {return x.split('/')[0]});
        var nodes = [];
        for (var i = 0; i < badges.length; i++) {
            var imageURL = BadgeImages[badges[i]];
            if (imageURL !== undefined) {
                nodes.push(<img src={imageURL} />);
            }
        }
        return nodes;
    },

    renderMessage: function (message) {
        return (
            <div className="message" key={message.id}>
                <span className="badges">{this.renderBadges(message)}</span>
                <span className="nick">{message.nick}:</span>&nbsp;
                {emoji.render(message.body)}
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

const BadgeImages = {
  "admin": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABIAAAASCAYAAABWzo5XAAABfklEQVQ4Ea1Tv0vDQBh9SYyLIDhaEHQRf4CCq6uLglDwH3BWXMXFycXV/8BVoVgHnURQEBR0aUHdKjUtCApW/AFNcud9Z3LeNT+69IM073vvfa93l8RqHw1z9KDsHmTIiD49yC029bYr9ssF5em6Iv7xCN56UANZoGsQqx+C1Q+y5hWfH8QCsOcSmFcCmK+G0oBxRrEhvN8Fq+0DnAHBp6T902nAsmGPrcKZ2oqt6p66ImdyE3ZhWYVItwgkjrS0MoN49ErRP49vJPySE5qs2Bu5jCD+XlHDrHkisTU4AbqoYo6w7qXeOCMyOkOzxEujM7cHe2RF9vLQX84lph89lHrL+EQGRuEuXAnWEs42YPeT579iTmzLP5sHvp6UZmyNhOB2DTz8SYbQiAgmLbhbN0KkpCIjwBvHCC+L4N9epyS4xp/mlROauaI4rFVFcLEE9nqtBtjbjeAWxedSVZwOzDPSFcKWC2dmR7JhZVs8gey3Oz+oMzinT91ajj9T+gXD8I6HuPQ1DgAAAABJRU5ErkJggg==",
  "broadcaster": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABIAAAASCAYAAABWzo5XAAAAhklEQVQ4EWN8LiHxn4EKgIkKZoCNGMYGscDCSOzCBRgTL/3KwAAiz8zMIHb2LAOMT1YYMfHzgw1BtpFkg5gVFRlEDh5ENgPMJtkgwSVLMAwBCZBsEMPfv9Qx6H1MDHUM+vvoEcNrW1sMw0j3GtCI/58/M7wyNkYxjHE006KEBzYOWYGNzSAA6TUbUUpeebAAAAAASUVORK5CYII=",
  "global_mod": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABIAAAASCAYAAABWzo5XAAAAAXNSR0IArs4c6QAAAT1JREFUOBFjZMhX+M9ABcBEBTPARtDHoJXxUxj0pTTBNgpzCzBsz1jAcK/2EFZPsGAVhQraKpsybFOaz8DIyMhw4+VdBg1xZYaLz65j1YLXa1MOLwIbAtIJMuT1l3cMMYsLSDeobfdUFE2iPEIMkUZ+YLEc2zgGEIYBRnzRf6Z4E4MUvzhMLVZaqs4cLI7Ta/uyl4ENefbxJVYD0AUxDALFDsgQUJiAAtik148BFFaEAIpBIENWJ0yDG+I0NQqsHz2sYIYiWwCOfhNZPYZ06ygGT00HBiZoVIcuyIKpRwlUmCDIEGQLwC468/gSg7eWI9iQf///M4AMefv1A0wPAx8HL5wNYqAbAhJDibVnTScZtl8/wJC8vBwkRxJASdmgqDSV0yPJAJhilMAGCZ5+dAkmRxKNYRBJupEUU80gANncWuEJw+OJAAAAAElFTkSuQmCC",
  "moderator": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABIAAAASCAYAAABWzo5XAAABP0lEQVQ4EWM0Wcf1n4EKgIkKZoCNYCHVoECFJIY0zWq4Ns/tymA2SS5ykPRFMQRuGpBBtEGGIjYM5QYTkPWisIkySIVPm6HNdCGKRhBn9o02uBjBMJLkkmOYbL0JrgHG6LpYxLD/2UYYlwHDICZGZoZ///+CFfCzCTHMs98PVwxjVJ+OZzj35giMC6YxvJat3cigwKvOwMnMxbDC+TSKYhAn95g/hiEgcUb0BLnd8y5IHCtIPujE8OzbQ6xyGC4CKcYGIvea4TQEpB4ljBgZGBne/XyFYU7Qbj2G73++YogjC6C4KEu7gWG92xVkeTBbnFMGQwxdAG4QKGx85GLg8iBXwLzpLx/PAIpNfADFazCFWx8tZfjx5xvDsz8PGWB5CSaHi4a7CFnDlKt1DP+BkBSAEf2kaEZWC3cRsiA5bAA3hU+ysIFpygAAAABJRU5ErkJggg==",
  "staff": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABIAAAASCAYAAABWzo5XAAAA1ElEQVQ4EWNU4Df+z0AFwEQFM8BGEDQoLSeW4eazY3CMy2K8BoEMKa7KhOudNWUxnI3OwGkQNkN626aj64fzGXEFNsg7MAByCT5DQOpwughmCIjGZggs7GDqiDJIQlIMph5Mo3sbJEiUQYvWTmGAGYbNEJBBLCACHYAUIwN5BRmGg2c3IAuB2cixiOEiXDaim4IeASgGkWsIyBK413AZgm4zustgfHg6Qk43MEliDQGpR/EazAAQTYohIPVwg0AaYYBUQ0D64F6DGUIuDXcRuQbA9AEAtHhR9n5yEkgAAAAASUVORK5CYII=",
  // TODO: fetch subscriber badge for user dynamically
  //"subscriber": "",
  "turbo": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABIAAAASCAYAAABWzo5XAAAAyElEQVQ4EWNMcVz6n4EKgIkKZoCNGHwGsSB7rWd1IDIXL7skdD2KPFFe+/vnHwO6RhRTgByCBn368IOhPHIjuj4MPorX0GVBruAT5AALF/c6o0jDggHmUpwuAilgZmFi+PT+B9iA3uK9eL2H0yCQblDYIIOQdENkLgobp9dATn/17DNDV/4esIbcVnsGeTUhFM3IHJwGgRTBDKmb7cnAJwAJK5hmWNjA+HgNAinqXO4PDiuYBlw0XoNgMYNLM7I442juRw4OrGwA84AxLKWQUDUAAAAASUVORK5CYII=",
};

module.exports = ChatTab;
