
const React = require('react'),
      WebView = require('react-electron-web-view');

const TwitchBotSetup = React.createClass({
    getInitialState: function () {
        return {
            thinking: true,
            url: null,
        };
    },
    componentWillMount: function () {
        var id = this.props.net.listeners.cmd("twitch-oauth-complete", (payload, error) => {
            this.props.parent.queryUserDetails();
            this.props.net.listeners.remove(id);
        });
        this.props.net.request("twitch-oauth-start", "bot").then(
            (payload) => {
                this.setState({
                    thinking: false,
                    url: payload,
                });
            },
            (error) => {
                // TODO: consider error path
                console.log(error);
            },
        );
    },

    renderThinking: function () {
        if (this.state.thinking) {
            return <div id="modal-spinner"></div>;
        }
        return null;
    },
    renderWebview: function () {
        if (this.state.url !== null) {
            return <WebView src={this.state.url} partition="twitch-bot-setup" />;
        }
        return null;
    },
    render: function () {
        return <div id="bot-setup">
            <div className="help-text">Login as your Bot user:</div>
            {this.renderThinking()}
            {this.renderWebview()}
        </div>;
    },
});

module.exports = TwitchBotSetup;
