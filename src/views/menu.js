
const React = require('react');

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
                <img id="streamer-logo" src={this.props.streamerLogo} />
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
