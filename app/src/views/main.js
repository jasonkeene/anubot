
const React = require('react'),
      ReactDOM = require('react-dom'),
      AuthOverlay = require('./auth_overlay.js');

const App = React.createClass({
    getInitialState: function () {
        return {
            authenticated: false,
        };
    },

    handleAuth: function (credentials) {
        this.setState({authenticated: true});
    },

    render: function () {
        if (!this.state.authenticated) {
            return (
                <div>
                    <AuthOverlay parent={this} />
                    <span>Some Content</span>
                </div>
            );
        }
        return (
            <div>
                <span>Some Content</span>
            </div>
        );
    },
});

function render() {
    ReactDOM.render(<App />, document.querySelector('#app'));
}

exports.render = render;
