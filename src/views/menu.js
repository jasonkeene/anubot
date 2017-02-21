
const React = require('react');

const Menu = React.createClass({
    getInitialState: function () {
        return {
            items: [
                ["Chat", "chat"],
                ["Playlist", "playlist"],
                //["Currency", "currency"],
                //["Commands", "commands"],
                //["Mini-games", "mini-games"],
                //["Stats", "stats"],
                //["Settings", "settings"],
            ],
        };
    },

    // events from children
    setActiveTab: function (id) {
        this.props.parent.setState({tab: id});
    },

    render: function () {
        return <ul id="menu">{this.renderMenuItems()}</ul>;
    },
    renderMenuItems: function () {
        var nodes = [];
        for (let i = 0; i < this.state.items.length; i++) {
            nodes.push(<MenuItem parent={this}
                                 text={this.state.items[i][0]}
                                 id={this.state.items[i][1]}
                                 selected={this.props.selected == this.state.items[i][1]}
                                 key={i+"-"+this.state.items[i][1]} />);
        }
        return nodes;
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
                {this.props.text}
            </li>
        );
    },
});

module.exports = Menu;