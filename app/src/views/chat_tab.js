
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

    processEmoji: function (key, nodes) {
        for (var i = 0; i < nodes.length; i++) {
            var el = nodes[i];
            if (typeof el === "string") {
                var index = el.indexOf(key);
                if (index === -1) {
                    continue
                } else if (index > 0) {
                    var prefix = el.substring(0, index),
                        rest = el.substring(index);
                    nodes.splice(i, 1, prefix, rest);
                    continue
                }

                var imgNode = <img className="emoji" src={emotes[key]} />,
                    rest = el.substring(index + key.length);
                nodes.splice(i, 1, imgNode, rest);
            }
        }
        return nodes;
    },
    renderEmoji: function (message) {
        var nodes = [message, <div />];
        for (var key in emotes) {
            this.processEmoji(key, nodes);
        }
        return nodes;
    },
    renderMessage: function (message) {
        return (
            <div className="message" key={message.id}>
                <span className="badges">{this.renderBadges(message)}</span>
                <span className="nick">{message.nick}:</span>&nbsp;
                {this.renderEmoji(message.body)}
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

const robot_emoji = {
    ":)": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-ebf60cd72f7aa600-24x18.png",
    ":(": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-d570c4b3b8d8fc4d-24x18.png",
    ":o": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-ae4e17f5b9624e2f-24x18.png",
    ":z": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-b9cbb6884788aa62-24x18.png",
    "B)": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-2cde79cfe74c6169-24x18.png",
    ":\\": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-374120835234cb29-24x18.png",
    ";)": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-cfaf6eac72fe4de6-24x18.png",
    ";p": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-3407bf911ad2fd4a-24x18.png",
    ":p": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-e838e5e34d9f240c-24x18.png",
    "R)": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-0536d670860bf733-24x18.png",
    "o_O": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-8e128fa8dc1de29c-24x18.png",
    ":D": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-9f2ac5d4b53913d7-24x18.png",
    ">(": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-d31223e81104544a-24x18.png",
    "<3": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-577ade91d46d7edc-24x18.png",
};
const turbo_emoji = {
    ":)": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-64f279c77d6f621d-21x18.png",
    ":(": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-c41c5c6c88f481cd-21x18.png",
    ":o": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-a43f189a61cbddbe-21x18.png",
    ":z": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-ff8b4b697171a170-21x18.png",
    "B)": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-9ad04ce3cf69ffd6-21x18.png",
    ":\\": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-7cd0191276363a02-21x18.png",
    ";)": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-54ab3f91053d8b97-21x18.png",
    ";p": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-a66f1856f37d0f48-21x18.png",
    ":p": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-a3ceb91b93f5082b-21x18.png",
    "R)": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-ffe61c02bd7cd500-21x18.png",
    "o_O": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-38b510fc1dd50022-21x18.png",
    ":D": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-1c8ec529616b79e0-21x18.png",
    ">(": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-91a9cf0c00b30760-21x18.png",
    "<3": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-934a78aa6d805cd7-21x18.png",
};
const monkey_emoji = {
    "<3": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-3f5d7d20df6ee956-20x18.png",
    "R)": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-7791c28e2e965fdf-20x22.png",
    ":>": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-665aec4773011f44-27x42.png",
    "<]": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-fd30ca5440d03927-20x42.png",
    ":7": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-206849962fa002dd-29x24.png",
    ":(": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-e4acdcf1ff2b4cef-20x18.png",
    ":P": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-e7b4f5211a173ff1-20x18.png",
    ";P": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-a63a460b5e1f74fc-20x18.png",
    ":O": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-30f5d7516b695012-20x18.png",
    ":\\": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-5067d1fe40f8e607-20x18.png",
    ":|": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-d6e8b4f562b8f46f-20x18.png",
    ":s": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-f5428b0c125bf4a5-20x18.png",
    ":D": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-89679577f86caf4e-20x18.png",
    "o_O": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-8429b9f83d424cb4-20x18.png",
    ">(": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-a6848b1076547d6f-20x18.png",
    ":)": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-ae6e77b75597c3d6-20x18.png",
    "B)": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-f381447031502180-20x18.png",
    ";)": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-095b4874cbf49881-20x18.png",
    "#/": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-39f51e122c6b2d60-27x18.png",
};
const emotes = {
    "KappaHD": "https://static-cdn.jtvnw.net/jtv_user_pictures/emoticon-2867-src-f02f9d40f66f0840-28x28.png",
    "MiniK": "https://static-cdn.jtvnw.net/jtv_user_pictures/emoticon-2868-src-5a7a81bb829e1a4c-28x28.png",
    "imGlitch": "https://static-cdn.jtvnw.net/emoticons/v1/62837/1.0",
    "copyThis": "https://static-cdn.jtvnw.net/emoticons/v1/62841/1.0",
    "pastaThat": "https://static-cdn.jtvnw.net/emoticons/v1/62842/1.0",
    "4Head": "https://static-cdn.jtvnw.net/emoticons/v1/354/1.0",
    "AMPEnergy": "https://static-cdn.jtvnw.net/emoticons/v1/99263/1.0",
    "AMPEnergyCherry": "https://static-cdn.jtvnw.net/emoticons/v1/99265/1.0",
    "ANELE": "https://static-cdn.jtvnw.net/emoticons/v1/3792/1.0",
    "ArgieB8": "https://static-cdn.jtvnw.net/emoticons/v1/51838/1.0",
    "ArsonNoSexy": "https://static-cdn.jtvnw.net/emoticons/v1/50/1.0",
    "AsianGlow": "https://static-cdn.jtvnw.net/emoticons/v1/74/1.0",
    "AthenaPMS": "https://static-cdn.jtvnw.net/emoticons/v1/32035/1.0",
    "BabyRage": "https://static-cdn.jtvnw.net/emoticons/v1/22639/1.0",
    "BatChest": "https://static-cdn.jtvnw.net/emoticons/v1/1905/1.0",
    "BCouch": "https://static-cdn.jtvnw.net/emoticons/v1/83536/1.0",
    "BCWarrior": "https://static-cdn.jtvnw.net/emoticons/v1/30/1.0",
    "BibleThump": "https://static-cdn.jtvnw.net/emoticons/v1/86/1.0",
    "BiersDerp": "https://static-cdn.jtvnw.net/emoticons/v1/96824/1.0",
    "BigBrother": "https://static-cdn.jtvnw.net/emoticons/v1/1904/1.0",
    "BionicBunion": "https://static-cdn.jtvnw.net/emoticons/v1/24/1.0",
    "BlargNaut": "https://static-cdn.jtvnw.net/emoticons/v1/38/1.0",
    "bleedPurple": "https://static-cdn.jtvnw.net/emoticons/v1/62835/1.0",
    "BloodTrail": "https://static-cdn.jtvnw.net/emoticons/v1/69/1.0",
    "BORT": "https://static-cdn.jtvnw.net/emoticons/v1/243/1.0",
    "BrainSlug": "https://static-cdn.jtvnw.net/emoticons/v1/881/1.0",
    "BrokeBack": "https://static-cdn.jtvnw.net/emoticons/v1/4057/1.0",
    "BudBlast": "https://static-cdn.jtvnw.net/emoticons/v1/97855/1.0",
    "BuddhaBar": "https://static-cdn.jtvnw.net/emoticons/v1/27602/1.0",
    "BudStar": "https://static-cdn.jtvnw.net/emoticons/v1/97856/1.0",
    "ChefFrank": "https://static-cdn.jtvnw.net/emoticons/v1/90129/1.0",
    "cmonBruh": "https://static-cdn.jtvnw.net/emoticons/v1/84608/1.0",
    "CoolCat": "https://static-cdn.jtvnw.net/emoticons/v1/58127/1.0",
    "CorgiDerp": "https://static-cdn.jtvnw.net/emoticons/v1/49106/1.0",
    "CougarHunt": "https://static-cdn.jtvnw.net/emoticons/v1/21/1.0",
    "DAESuppy": "https://static-cdn.jtvnw.net/emoticons/v1/973/1.0",
    "DalLOVE": "https://static-cdn.jtvnw.net/emoticons/v1/96848/1.0",
    "DansGame": "https://static-cdn.jtvnw.net/emoticons/v1/33/1.0",
    "DatSheffy": "https://static-cdn.jtvnw.net/emoticons/v1/170/1.0",
    "DBstyle": "https://static-cdn.jtvnw.net/emoticons/v1/73/1.0",
    "deExcite": "https://static-cdn.jtvnw.net/emoticons/v1/46249/1.0",
    "deIlluminati": "https://static-cdn.jtvnw.net/emoticons/v1/46248/1.0",
    "DendiFace": "https://static-cdn.jtvnw.net/emoticons/v1/58135/1.0",
    "DogFace": "https://static-cdn.jtvnw.net/emoticons/v1/1903/1.0",
    "DOOMGuy": "https://static-cdn.jtvnw.net/emoticons/v1/54089/1.0",
    "DoritosChip": "https://static-cdn.jtvnw.net/emoticons/v1/102242/1.0",
    "duDudu": "https://static-cdn.jtvnw.net/emoticons/v1/62834/1.0",
    "EagleEye": "https://static-cdn.jtvnw.net/emoticons/v1/20/1.0",
    "EleGiggle": "https://static-cdn.jtvnw.net/emoticons/v1/4339/1.0",
    "FailFish": "https://static-cdn.jtvnw.net/emoticons/v1/360/1.0",
    "FPSMarksman": "https://static-cdn.jtvnw.net/emoticons/v1/42/1.0",
    "FrankerZ": "https://static-cdn.jtvnw.net/emoticons/v1/65/1.0",
    "FreakinStinkin": "https://static-cdn.jtvnw.net/emoticons/v1/39/1.0",
    "FUNgineer": "https://static-cdn.jtvnw.net/emoticons/v1/244/1.0",
    "FunRun": "https://static-cdn.jtvnw.net/emoticons/v1/48/1.0",
    "FutureMan": "https://static-cdn.jtvnw.net/emoticons/v1/98562/1.0",
    "FuzzyOtterOO": "https://static-cdn.jtvnw.net/emoticons/v1/168/1.0",
    "GingerPower": "https://static-cdn.jtvnw.net/emoticons/v1/32/1.0",
    "GrammarKing": "https://static-cdn.jtvnw.net/emoticons/v1/3632/1.0",
    "HassaanChop": "https://static-cdn.jtvnw.net/emoticons/v1/20225/1.0",
    "HassanChop": "https://static-cdn.jtvnw.net/emoticons/v1/68/1.0",
    "HeyGuys": "https://static-cdn.jtvnw.net/emoticons/v1/30259/1.0",
    "HotPokket": "https://static-cdn.jtvnw.net/emoticons/v1/357/1.0",
    "HumbleLife": "https://static-cdn.jtvnw.net/emoticons/v1/46881/1.0",
    "ItsBoshyTime": "https://static-cdn.jtvnw.net/emoticons/v1/169/1.0",
    "Jebaited": "https://static-cdn.jtvnw.net/emoticons/v1/90/1.0",
    "JKanStyle": "https://static-cdn.jtvnw.net/emoticons/v1/15/1.0",
    "JonCarnage": "https://static-cdn.jtvnw.net/emoticons/v1/26/1.0",
    "KAPOW": "https://static-cdn.jtvnw.net/emoticons/v1/9803/1.0",
    "Kappa": "https://static-cdn.jtvnw.net/emoticons/v1/25/1.0",
    "KappaClaus": "https://static-cdn.jtvnw.net/emoticons/v1/74510/1.0",
    "KappaPride": "https://static-cdn.jtvnw.net/emoticons/v1/55338/1.0",
    "KappaRoss": "https://static-cdn.jtvnw.net/emoticons/v1/70433/1.0",
    "KappaWealth": "https://static-cdn.jtvnw.net/emoticons/v1/81997/1.0",
    "Keepo": "https://static-cdn.jtvnw.net/emoticons/v1/1902/1.0",
    "KevinTurtle": "https://static-cdn.jtvnw.net/emoticons/v1/40/1.0",
    "Kippa": "https://static-cdn.jtvnw.net/emoticons/v1/1901/1.0",
    "Kreygasm": "https://static-cdn.jtvnw.net/emoticons/v1/41/1.0",
    "Mau5": "https://static-cdn.jtvnw.net/emoticons/v1/30134/1.0",
    "mcaT": "https://static-cdn.jtvnw.net/emoticons/v1/35063/1.0",
    "MikeHogu": "https://static-cdn.jtvnw.net/emoticons/v1/81636/1.0",
    "MingLee": "https://static-cdn.jtvnw.net/emoticons/v1/68856/1.0",
    "MKXRaiden": "https://static-cdn.jtvnw.net/emoticons/v1/102324/1.0",
    "MKXScorpion": "https://static-cdn.jtvnw.net/emoticons/v1/102325/1.0",
    "MrDestructoid": "https://static-cdn.jtvnw.net/emoticons/v1/28/1.0",
    "MVGame": "https://static-cdn.jtvnw.net/emoticons/v1/29/1.0",
    "NinjaTroll": "https://static-cdn.jtvnw.net/emoticons/v1/45/1.0",
    "NomNom": "https://static-cdn.jtvnw.net/emoticons/v1/90075/1.0",
    "NoNoSpot": "https://static-cdn.jtvnw.net/emoticons/v1/44/1.0",
    "NotATK": "https://static-cdn.jtvnw.net/emoticons/v1/34875/1.0",
    "NotLikeThis": "https://static-cdn.jtvnw.net/emoticons/v1/58765/1.0",
    "OhMyDog": "https://static-cdn.jtvnw.net/emoticons/v1/81103/1.0",
    "OMGScoots": "https://static-cdn.jtvnw.net/emoticons/v1/91/1.0",
    "OneHand": "https://static-cdn.jtvnw.net/emoticons/v1/66/1.0",
    "OpieOP": "https://static-cdn.jtvnw.net/emoticons/v1/100590/1.0",
    "OptimizePrime": "https://static-cdn.jtvnw.net/emoticons/v1/16/1.0",
    "OSfrog": "https://static-cdn.jtvnw.net/emoticons/v1/81248/1.0",
    "OSkomodo": "https://static-cdn.jtvnw.net/emoticons/v1/81273/1.0",
    "OSsloth": "https://static-cdn.jtvnw.net/emoticons/v1/81249/1.0",
    "panicBasket": "https://static-cdn.jtvnw.net/emoticons/v1/22998/1.0",
    "PanicVis": "https://static-cdn.jtvnw.net/emoticons/v1/3668/1.0",
    "PartyTime": "https://static-cdn.jtvnw.net/emoticons/v1/76171/1.0",
    "PazPazowitz": "https://static-cdn.jtvnw.net/emoticons/v1/19/1.0",
    "PeoplesChamp": "https://static-cdn.jtvnw.net/emoticons/v1/3412/1.0",
    "PermaSmug": "https://static-cdn.jtvnw.net/emoticons/v1/27509/1.0",
    "PeteZaroll": "https://static-cdn.jtvnw.net/emoticons/v1/81243/1.0",
    "PeteZarollTie": "https://static-cdn.jtvnw.net/emoticons/v1/81244/1.0",
    "PicoMause": "https://static-cdn.jtvnw.net/emoticons/v1/27/1.0",
    "PipeHype": "https://static-cdn.jtvnw.net/emoticons/v1/4240/1.0",
    "PJSalt": "https://static-cdn.jtvnw.net/emoticons/v1/36/1.0",
    "PJSugar": "https://static-cdn.jtvnw.net/emoticons/v1/102556/1.0",
    "PMSTwin": "https://static-cdn.jtvnw.net/emoticons/v1/92/1.0",
    "PogChamp": "https://static-cdn.jtvnw.net/emoticons/v1/88/1.0",
    "Poooound": "https://static-cdn.jtvnw.net/emoticons/v1/358/1.0",
    "PraiseIt": "https://static-cdn.jtvnw.net/emoticons/v1/38586/1.0",
    "PRChase": "https://static-cdn.jtvnw.net/emoticons/v1/28328/1.0",
    "PunchTrees": "https://static-cdn.jtvnw.net/emoticons/v1/47/1.0",
    "PuppeyFace": "https://static-cdn.jtvnw.net/emoticons/v1/58136/1.0",
    "RaccAttack": "https://static-cdn.jtvnw.net/emoticons/v1/27679/1.0",
    "RalpherZ": "https://static-cdn.jtvnw.net/emoticons/v1/1900/1.0",
    "RedCoat": "https://static-cdn.jtvnw.net/emoticons/v1/22/1.0",
    "ResidentSleeper": "https://static-cdn.jtvnw.net/emoticons/v1/245/1.0",
    "riPepperonis": "https://static-cdn.jtvnw.net/emoticons/v1/62833/1.0",
    "RitzMitz": "https://static-cdn.jtvnw.net/emoticons/v1/4338/1.0",
    "RuleFive": "https://static-cdn.jtvnw.net/emoticons/v1/361/1.0",
    "SeemsGood": "https://static-cdn.jtvnw.net/emoticons/v1/64138/1.0",
    "ShadyLulu": "https://static-cdn.jtvnw.net/emoticons/v1/52492/1.0",
    "ShazBotstix": "https://static-cdn.jtvnw.net/emoticons/v1/87/1.0",
    "ShibeZ": "https://static-cdn.jtvnw.net/emoticons/v1/27903/1.0",
    "SmoocherZ": "https://static-cdn.jtvnw.net/emoticons/v1/89945/1.0",
    "SMOrc": "https://static-cdn.jtvnw.net/emoticons/v1/52/1.0",
    "SMSkull": "https://static-cdn.jtvnw.net/emoticons/v1/51/1.0",
    "SoBayed": "https://static-cdn.jtvnw.net/emoticons/v1/1906/1.0",
    "SoonerLater": "https://static-cdn.jtvnw.net/emoticons/v1/355/1.0",
    "SriHead": "https://static-cdn.jtvnw.net/emoticons/v1/14706/1.0",
    "SSSsss": "https://static-cdn.jtvnw.net/emoticons/v1/46/1.0",
    "StinkyCheese": "https://static-cdn.jtvnw.net/emoticons/v1/90076/1.0",
    "StoneLightning": "https://static-cdn.jtvnw.net/emoticons/v1/17/1.0",
    "StrawBeary": "https://static-cdn.jtvnw.net/emoticons/v1/37/1.0",
    "SuperVinlin": "https://static-cdn.jtvnw.net/emoticons/v1/31/1.0",
    "SwiftRage": "https://static-cdn.jtvnw.net/emoticons/v1/34/1.0",
    "TBCheesePull": "https://static-cdn.jtvnw.net/emoticons/v1/94039/1.0",
    "TBTacoLeft": "https://static-cdn.jtvnw.net/emoticons/v1/94038/1.0",
    "TBTacoRight": "https://static-cdn.jtvnw.net/emoticons/v1/94040/1.0",
    "TF2John": "https://static-cdn.jtvnw.net/emoticons/v1/1899/1.0",
    "TheRinger": "https://static-cdn.jtvnw.net/emoticons/v1/18/1.0",
    "TheTarFu": "https://static-cdn.jtvnw.net/emoticons/v1/70/1.0",
    "TheThing": "https://static-cdn.jtvnw.net/emoticons/v1/7427/1.0",
    "ThunBeast": "https://static-cdn.jtvnw.net/emoticons/v1/1898/1.0",
    "TinyFace": "https://static-cdn.jtvnw.net/emoticons/v1/67/1.0",
    "TooSpicy": "https://static-cdn.jtvnw.net/emoticons/v1/359/1.0",
    "TriHard": "https://static-cdn.jtvnw.net/emoticons/v1/171/1.0",
    "TTours": "https://static-cdn.jtvnw.net/emoticons/v1/38436/1.0",
    "twitchRaid": "https://static-cdn.jtvnw.net/emoticons/v1/62836/1.0",
    "TwitchRPG": "https://static-cdn.jtvnw.net/emoticons/v1/102157/1.0",
    "UleetBackup": "https://static-cdn.jtvnw.net/emoticons/v1/49/1.0",
    "UncleNox": "https://static-cdn.jtvnw.net/emoticons/v1/3666/1.0",
    "UnSane": "https://static-cdn.jtvnw.net/emoticons/v1/71/1.0",
    "VaultBoy": "https://static-cdn.jtvnw.net/emoticons/v1/54090/1.0",
    "VoHiYo": "https://static-cdn.jtvnw.net/emoticons/v1/81274/1.0",
    "Volcania": "https://static-cdn.jtvnw.net/emoticons/v1/166/1.0",
    "WholeWheat": "https://static-cdn.jtvnw.net/emoticons/v1/1896/1.0",
    "WinWaker": "https://static-cdn.jtvnw.net/emoticons/v1/167/1.0",
    "WTRuck": "https://static-cdn.jtvnw.net/emoticons/v1/1897/1.0",
    "WutFace": "https://static-cdn.jtvnw.net/emoticons/v1/28087/1.0",
    "YouWHY": "https://static-cdn.jtvnw.net/emoticons/v1/4337/1.0",
};

module.exports = ChatTab;
