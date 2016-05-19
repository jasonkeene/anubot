
const React = require('react');

const AuthOption = React.createClass({
    render: function () {
        return <li onClick={this.props.click}>{this.props.text}</li>
    }
})

const ManualMode = React.createClass({
    render: function () {
        return <div id="auth-overlay">
                Channel Username: <input type="text" /><br />
                Channel Oauth Token: <input type="password" /><br />
                Bot Username: <input type="text" /><br />
                Bot Oauth Token: <input type="password" /><br />
                <input type="button" value="Submit" />
            </div>;
    }
})

const AuthOverlay = React.createClass({
    getInitialState: function () {
        return {
            mode: "choice",
        };
    },
    handleOptionClick: function (mode) {
        this.setState({mode: mode});
    },
    renderChoice: function () {
        return <div id="auth-overlay">
            <ul>
                <AuthOption click={this.handleOptionClick.bind(this, "auto")} text="Login via Twitch" />
                <AuthOption click={this.handleOptionClick.bind(this, "manual")} text="Manually Enter Oauth Tokens" />
            </ul>
        </div>;
    },
    renderAuto: function () {
        return <div id="auth-overlay"> Auto </div>;
    },
    render: function () {
        switch (this.state.mode) {
        case "choice":
            return this.renderChoice();
        case "auto":
            return this.renderAuto();
        case "manual":
            return <ManualMode />;
        }
    },
});

module.exports = AuthOverlay;
