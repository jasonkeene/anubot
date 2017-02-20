
const React = require('react'),
      ReactDOM = require('react-dom');

const Setup = React.createClass({
    getInitialState: function () {
        return {
            step: "login",
            thinking: false,
        };
    },

    renderStep: function () {
        switch (this.state.step) {
        case "login":
            return <Login root={this.props.parent}
                          parent={this}
                          net={this.props.net} />;
        case "register":
            return <Register root={this.props.parent}
                             parent={this}
                             net={this.props.net} />;
        default:
            return null;
        }
    },
    renderThinking: function () {
        if (this.state.thinking) {
            return <div id="modal-spinner"></div>;
        }
        return null;
    },
    render: function () {
        return <div id="setup">
            {this.renderThinking()}
            {this.renderStep()}
        </div>;
    },
});

const Login = React.createClass({
    getInitialState: function () {
        return {
            username: "",
            password: "",
            badCreds: false,
        };
    },

    // network events
    handleAuthenticateSuccess: function (payload) {
        this.setState({
            badCreds: false,
        });
        this.props.parent.setState({
            thinking: false,
        });
        this.props.root.authenticated({
            username: this.state.username,
            password: this.state.password,
        });
    },
    handleAuthenticateFailure: function (error) {
        this.setState({
            badCreds: true,
        });
        this.props.parent.setState({
            thinking: false,
        });
    },

    // ui events
    submitLogin: function (e) {
        e.preventDefault();
        document.activeElement.blur();
        this.props.parent.setState({
            thinking: true,
        });
        net.request("authenticate", {
            username: this.state.username,
            password: this.state.password,
        }).then(
            this.handleAuthenticateSuccess,
            this.handleAuthenticateFailure,
        );
    },
    handleUsernameChange: function (e) {
        this.setState({
            username: e.target.value,
        });
    },
    handlePasswordChange: function (e) {
        this.setState({
            password: e.target.value,
        });
    },
    handleSignup: function (e) {
        e.preventDefault();
        this.props.parent.setState({
            step: "register",
        });
    },

    renderBadCreds: function () {
        if (this.state.badCreds) {
            return <div id="bad-creds">The username and password you entered are not correct.</div>;
        }
        return null;
    },
    render: function () {
        return <div id="login-form">
            {this.renderBadCreds()}
            <form onSubmit={this.submitLogin}>
                <input type="text" placeholder="username" onChange={this.handleUsernameChange} /><br />
                <input type="password" placeholder="password" onChange={this.handlePasswordChange} />
                <input type="submit" style={{display: "none"}} />
            </form>
            Don't have an account? <a href="#" onClick={this.handleSignup}>Sign up</a>
        </div>;
    },
});

const Register = React.createClass({
    getInitialState: function () {
        return {
            username: "",
            password: "",
            referral: "",
            error: null,
        };
    },

    // network events
    handleRegisterSuccess: function (payload) {
        this.setState({
            error: null,
        });
        this.props.parent.setState({
            thinking: false,
        });
        this.props.root.authenticated({
            username: this.state.username,
            password: this.state.password,
        });
    },
    handleRegisterFailure: function (error) {
        this.setState({error});
        this.props.parent.setState({
            thinking: false,
        });
    },

    // ui events
    submitRegister: function (e) {
        e.preventDefault();
        document.activeElement.blur();
        this.props.parent.setState({
            thinking: true,
        });

        net.request("register", {
            username: this.state.username,
            password: this.state.password,
        }).then(
            this.handleRegisterSuccess,
            this.handleRegisterFailure,
        );
    },
    handleUsernameChange: function (e) {
        this.setState({
            username: e.target.value,
        });
    },
    handlePasswordChange: function (e) {
        this.setState({
            password: e.target.value,
        });
    },

    renderError: function () {
        // todo: add wrapping nodes for display
        if (this.state.error !== null) {
            return this.state.error;
        }
        return null;
    },
    render: function () {
        return <div id="registration-form">
            {this.renderError()}
            <form onSubmit={this.submitRegister}>
                <input type="text" placeholder="username" onChange={this.handleUsernameChange} /><br />
                <input type="password" placeholder="password" onChange={this.handlePasswordChange} />
                <input type="submit" style={{display: "none"}} />
            </form>
        </div>
    },
});

module.exports = Setup;
