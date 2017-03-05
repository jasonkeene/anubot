
const React = require('react');

const _defaultLogo = "data:image/png;base64,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";

const Menu = React.createClass({
    getInitialState: function () {
        return {
            items: [
                ["Chat", "chat", <i className="icon mdi mdi-message-outline"></i>],
                ["Playlist", "playlist", <i className="icon material-icons" style={{padding: "3px 0 0 0"}}>queue_music</i>],
                ["Currency", "currency", <i className="icon mdi mdi-coin"></i>],
                ["Commands", "commands", <i className="icon mdi mdi-code-tags"></i>],
                ["Mini-games", "mini-games", <i className="icon mdi mdi-google-controller" style={{marginTop: "-1px"}}></i>],
                ["Stats", "stats", <i className="icon mdi mdi-chart-bar"></i>],
            ],
            userItems: [
                ["Connections", "connections", <i className="icon mdi mdi-power-plug"></i>],
                ["Settings", "settings", <i className="icon mdi mdi-settings"></i>],
                ["Logout", "logout", <i className="icon mdi mdi-logout"></i>],
            ],
            showUserMenu: false,
        };
    },

    // events from children
    setActiveTab: function (id) {
        if (id === "logout") {
            this.logout();
        }
        this.props.parent.setState({tab: id});
    },

    toggleUserMenu: function () {
        this.setState({
            showUserMenu: !this.state.showUserMenu,
        })
    },

    logout: function () {
        this.props.parent.logout();
    },

    renderMenuItems: function () {
        var nodes = [];
        for (let i = 0; i < this.state.items.length; i++) {
            nodes.push(<MenuItem parent={this}
                                 text={this.state.items[i][0]}
                                 id={this.state.items[i][1]}
                                 icon={this.state.items[i][2]}
                                 selected={this.props.selected == this.state.items[i][1]}
                                 key={i+"-"+this.state.items[i][1]} />);
        }
        return nodes;
    },
    renderUserMenuItems: function () {
        var nodes = [];
        for (let i = 0; i < this.state.userItems.length; i++) {
            nodes.push(<MenuItem parent={this}
                                 text={this.state.userItems[i][0]}
                                 id={this.state.userItems[i][1]}
                                 icon={this.state.userItems[i][2]}
                                 selected={this.props.selected == this.state.userItems[i][1]}
                                 key={i+"-"+this.state.userItems[i][1]} />);
        }
        return nodes;
    },
    render: function () {
        return <div id="menu">
            <img id="logo" src="../assets/images/anubot-logo.svg" />
            <ul>{this.renderMenuItems()}</ul>
            <ul id="user-menu" className={this.state.showUserMenu ? "show" : ""}>
                {this.renderUserMenuItems()}
            </ul>
            <div id="user-info" onClick={this.toggleUserMenu}>
                <img id="streamer-logo" src={this.props.streamerLogo || _defaultLogo } />
                <div id="display-name">{this.props.streamerDisplayName}</div>
            </div>
        </div>;
    },
});

const MenuItem = React.createClass({
    getInitialState: function () {
        return {};
    },

    // ui events
    handleClick: function () {
        this.props.parent.setActiveTab(this.props.id);
    },

    render: function () {
        return (
            <li onClick={this.handleClick}
                className={this.props.selected ? "selected" : ""}>
                {this.props.icon}
                {this.props.text}
            </li>
        );
    },
});

module.exports = Menu;
