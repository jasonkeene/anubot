/* global window: false */

const React = require('react'),
      ReactDOM = require('react-dom'),
      badges = require('./badges.js'),
      mentions = require('./mentions.js'),
      emoji = require('./emoji.js');

const _defaultNickColors = [
    "#FF0000",
    "#0000FF",
    "#008000",
    "#B22222",
    "#FF7F50",
    "#9ACD32",
    "#FF4500",
    "#2E8B57",
    "#DAA520",
    "#D2691E",
    "#5F9EA0",
    "#1E90FF",
    "#FF69B4",
    "#8A2BE2",
    "#00FF7F",
];

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

    defaultNickColor: function (username) {
        return _defaultNickColors[_hash(username) % _defaultNickColors.length];
    },
    nickStyle: function (message) {
        if (message.tags === undefined ||
            message.tags.color === undefined ||
            message.tags.color === "") {
            var color = this.defaultNickColor(message.tags['display-name']);
            return {color};
        }
        return {color: message.tags.color};
    },
    renderPrivmsg: function (message) {
        return (
            <div className="message" key={message.twitch.tags.id}>
                <span className="badges">{badges.render(message.twitch)}</span>
                <span className="nick" style={this.nickStyle(message.twitch)}>{message.twitch.tags['display-name']}</span>:&nbsp;
                {mentions.render(this.props.streamer_username, emoji.render(message.twitch))}
            </div>
        );
    },
    renderAction: function (message) {
        return (
            <div className="message action" key={message.twitch.tags.id}>
                <span className="badges">{badges.render(message.twitch)}</span>
                <span className="nick" style={this.nickStyle(message.twitch)}>{message.twitch.tags['display-name']}</span>&nbsp;
                <span className="message" style={this.nickStyle(message.twitch)}>{mentions.render(this.props.streamer_username, emoji.render(message.twitch))}</span>
            </div>
        );
    },
    renderWhisper: function (message) {
        return (
            <div className="message whisper" key={message.twitch.tags.id}>
                <span className="badges">{badges.render(message.twitch)}</span>
                <span className="nick" style={this.nickStyle(message.twitch)}>{message.twitch.tags['display-name']}</span>&nbsp;
                <span className="arrow">&#x25B8;</span>&nbsp;
                <span className="target" style={{color: this.defaultNickColor(message.twitch.target)}}>{message.twitch.target}</span>:&nbsp;
                {mentions.render(this.props.streamer_username, emoji.render(message.twitch))}
            </div>
        );
    },
    renderMessage: function (message) {
        switch (message.twitch.cmd) {
        case "PRIVMSG":
            return this.renderPrivmsg(message);
        case "ACTION":
            return this.renderAction(message);
        case "WHISPER":
            return this.renderWhisper(message);
        }
        return null;
    },
    render: function () {
        return (
            <div id="chat-tab" className="tab">
                <ChatHeader channel={"#" + this.props.streamer_username}
                            status={this.props.status}
                            game={this.props.game}
                            net={this.props.net} />
                <div className="body">
                    <div className="spacer" />
                    {this.props.messages.map(this.renderMessage)}
                </div>
                <ChatFooter streamer_username={this.props.streamer_username}
                            bot_username={this.props.bot_username}
                            net={this.props.net} />
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
    componentWillReceiveProps: function (nextProps) {
        this.setState({
            status: nextProps.status,
            game: nextProps.game,
        });
    },

    handleEditClick: function () {
        this.setState({
            editing: true,
        });
    },

    renderEditing: function () {
        return (
            <div className="header">
                <div className="channel">{this.props.channel}</div>
                <div className="description">
                    <ChatHeaderInput parent={this}
                                     net={this.props.net}
                                     inputGame={this.state.game}
                                     inputStatus={this.state.status} />
                </div>
            </div>
        );
    },
    renderDisplay: function () {
        return (
            <div className="header">
                <div className="channel">{this.props.channel}</div>
                <div className="description" onClick={this.handleEditClick}>
                    <b>{this.state.game}:</b>&nbsp;
                    {this.state.status}
                    &nbsp;<i className="material-icons edit">edit</i>
                </div>
            </div>
        );
    },
    render: function () {
        if (this.state.editing) {
            return this.renderEditing();
        }
        return this.renderDisplay();
    },
});

const ChatHeaderInput = React.createClass({
    getInitialState: function () {
        this.props.net.request("twitch-games").then((payload) => {
            this.setState({
                games: payload,
            });
        });
        return {
            inputStatus: this.props.inputStatus,
            inputGame: this.props.inputGame,
            // tracks if the game field has changed, in which case display the
            // games menu.
            editedGame: false,
            // tracks if the games menu should be displayed.
            displayGamesMenu: false,
            // id of the selected game in the games menu.
            selectedGame: -1,
            // available games.
            games: [],
            // displays an error.
            displayError: false,
        };
    },
    componentDidMount(){
        // select and focus games text input
        var node = ReactDOM.findDOMNode(this.refs.gameInput);
        node.focus();
        node.setSelectionRange(0, node.value.length);
    },

    handleStatusChange: function (e) {
        this.setState({
            inputStatus: e.target.value,
        });
    },
    handleStatusKeyDown: function (e) {
        switch (e.keyCode) {
        case 27: // escape
            this.abort();
            break;
        case 13: // return
            this.complete();
            break;
        }
    },

    handleGameFocus: function (e) {
        if (this.state.editedGame) {
            this.setState({
                displayGamesMenu: true,
            });
        }
    },
    handleGameBlur: function (e) {
        window.setTimeout(() => {
            this.setState({
                displayGamesMenu: false,
            });
        }, 300);
    },
    handleGameChange: function (e) {
        this.setState({
            editedGame: true,
            displayError: false,
            displayGamesMenu: true,
            inputGame: e.target.value,
            selectedGame: -1,
        });
    },
    handleGameKeyDown: function (e) {
        switch (e.keyCode) {
        case 27: // escape
            this.abort();
            break;
        case 13: // return
            if (this.gameSelected()) {
                this.captureGame();
                break;
            }
            this.complete();
            break;
        case 40: // down
            e.preventDefault();
            this.moveMenu("down");
            break;
        case 38: // up
            e.preventDefault();
            this.moveMenu("up");
            break;
        }
    },

    handleGameMenuClick: function (e) {
        var node = e.target,
            count = 0;
        while (node.attributes["data-name"] === undefined && count < 4) {
            node = node.parentNode;
            count++;
        }
        if (node.attributes["data-name"] === undefined) {
            return;
        }
        this.setState({
            inputGame: node.attributes["data-name"].value,
            displayGamesMenu: false,
        });
        ReactDOM.findDOMNode(this.refs.statusInput).focus();
    },
    handleGameMenuMouseOver: function (e) {
        var node = e.target,
            count = 0;
        while (node.attributes["data-id"] === undefined && count < 4) {
            node = node.parentNode;
            count++;
        }
        if (node.attributes["data-id"] === undefined) {
            return;
        }
        this.setState({
            selectedGame: parseInt(node.attributes["data-id"].value),
        });
    },

    moveMenu: function (direction) {
        var games = this.filterGames(),
            game = null;

        if (this.state.selectedGame === -1) {
            if (games.length > 0) {
                game = direction === "up" ? games[games.length-1] : games[0];
                this.setState({
                    selectedGame: game.id,
                });
            }
            return;
        }

        var i = -1;
        for (var j = 0; j < games.length; j++) {
            game = games[j];
            if (game.id === this.state.selectedGame) {
                var delta = direction === "up" ? -1 : 1;
                i = j + delta;
                break;
            }
        }
        if (i > -1 && i < games.length) {
            game = games[i];
            this.setState({
                selectedGame: game.id,
            });
        }
    },
    captureGame: function () {
        var gameName = "";
        for (var i = 0; i < this.state.games.length; i++) {
            var game = this.state.games[i];
            if (game.id === this.state.selectedGame) {
                gameName = game.name;
                break;
            }
        }
        this.setState({
            selectedGame: -1,
            editedGame: false,
            displayGamesMenu: false,
            inputGame: gameName,
        })
        var node = ReactDOM.findDOMNode(this.refs.statusInput);
        node.focus();
        node.setSelectionRange(0, node.value.length);
    },
    gameSelected: function () {
        return this.state.selectedGame !== -1
    },
    abort: function () {
        this.props.parent.setState({
            editing: false,
        });
    },
    complete: function () {
        var game = this.validateGame(this.state.inputGame);
        if (game !== null) {
            this.props.net.send({
                "cmd": "twitch-update-chat-description",
                "payload": {
                    "status": this.state.inputStatus,
                    "game": game,
                },
            });
            this.props.parent.setState({
                editing: false,
                status: this.state.inputStatus,
                game: game,
            });
            return;
        }

        var node = ReactDOM.findDOMNode(this.refs.gameInput);
        node.focus();
        node.setSelectionRange(0, node.value.length);
        this.setState({
            displayError: true,
        });
    },
    validateGame: function (input) {
        var games = this.state.games;
        input = input.toLowerCase();

        for (var i = 0; i < games.length; i++) {
            var game = games[i],
                match = game.name.toLowerCase();
            if (match === input) {
                return game.name;
            }
        }
        return null;
    },
    filterGames: function () {
        var input = this.state.inputGame.toLowerCase(),
            games = this.state.games,
            result = [];

        for (var i = 0; i < games.length; i++) {
            var game = games[i],
                match = game.name.toLowerCase();
            if (match.indexOf(input) !== -1) {
                result.push(game);
            }
        }
        return result.slice(0, 5);
    },

    renderError: function () {
        return <div className="error">
            {'"' + this.state.inputGame + '" is not a valid game.'}
        </div>;
    },
    renderGamesMenuItem: function (game) {
        return <li key={"game-"+game.id}
                   className={game.id === this.state.selectedGame ? "selected" : ""}
                   data-id={game.id}
                   data-name={game.name}
                   onClick={this.handleGameMenuClick}
                   onMouseOver={this.handleGameMenuMouseOver}>
                   <img src={game.image} width="30" height="40" />
                   <span className="name">{game.name}</span>
                </li>;
    },
    renderGamesMenu: function () {
        return <ol>
            {this.filterGames().map(this.renderGamesMenuItem)}
        </ol>;
    },
    render: function () {
        return <div className="input">
            {this.state.displayGamesMenu ? this.renderGamesMenu() : ""}
            {this.state.displayError ? this.renderError() : ""}
            <input type="text" ref="gameInput" className="game-input text-input"
                onChange={this.handleGameChange}
                onKeyDown={this.handleGameKeyDown}
                onFocus={this.handleGameFocus}
                onBlur={this.handleGameBlur}
                value={this.state.inputGame} />&nbsp;
            <input type="text" ref="statusInput" className="status-input text-input"
                onChange={this.handleStatusChange}
                onKeyDown={this.handleStatusKeyDown}
                value={this.state.inputStatus} />
        </div>;
    },
});

const ChatFooter = React.createClass({
    getInitialState: function () {
        return {
            user_type: "streamer",
            message: "",

            selected: -1,
            previous_messages: [],
        };
    },

    // event handlers
    handleSubmit: function (e) {
        e.preventDefault();
        this.props.net.send({
            cmd: "twitch-send-message",
            payload: {
                user_type: this.state.user_type,
                message: this.state.message,
            },
        });
        this.setState({
            message: "",
            previous_messages: this.state.previous_messages.concat([this.state.message]),
            selected: -1,
        });
    },
    handleMessageChange: function (e) {
        this.setState({message: e.target.value});
    },
    handleUserChange: function (e) {
        this.setState({user_type: e.target.value});
    },
    handleKeyDown: function (e) {
        if (this.state.previous_messages.length === 0) {
            return;
        }
        var selected;
        switch (e.keyCode) {
        case 38: // up
            if (this.state.selected < 1) {
                selected = this.state.previous_messages.length - 1;
                break;
            }
            selected = (this.state.selected - 1) % this.state.previous_messages.length;
            break;
        case 40: // down
            selected = (this.state.selected + 1) % this.state.previous_messages.length;
            break;
        default:
            return
        }
        this.setState({
            selected,
            message: this.state.previous_messages[selected],
        });
    },

    render: function () {
        return (
            <div className="footer">
                <div className="selection">
                    <select onChange={this.handleUserChange}
                            className="select-input">
                        <option value="streamer">{this.props.streamer_username}</option>
                        <option value="bot">{this.props.bot_username}</option>
                    </select>
                </div>
                <div className="input">
                    <div className="form">
                        <form onSubmit={this.handleSubmit}>
                            <input className="text-input"
                                   onChange={this.handleMessageChange}
                                   onKeyDown={this.handleKeyDown}
                                   type="text"
                                   placeholder="Send a message"
                                   value={this.state.message} />
                        </form>
                    </div>
                    <div className="spacer"></div>
                </div>
            </div>
        );
    },
});

function _hash(str) {
    var hash = 0, i, chr, len;
    if (str.length === 0) {
        return hash;
    }
    for (i = 0; i < str.length; i++) {
        chr = str.charCodeAt(i);
        hash = ((hash << 5) - hash) + chr;
        hash |= 0;
    }
    return hash;
}

module.exports = ChatTab;
