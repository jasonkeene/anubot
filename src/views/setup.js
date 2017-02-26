
const React = require('react');

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
        this.props.net.request("authenticate", {
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
            return <div id="bad-creds">Wrong username or password.</div>;
        }
        return null;
    },
    render: function () {
        return <div id="login-form">
            {this.renderBadCreds()}
            <form onSubmit={this.submitLogin}>
                <i className="material-icons">person_outline</i>
                <input type="text" placeholder="username" onChange={this.handleUsernameChange} /><br />
                <i className="material-icons">lock_open</i>
                <input type="password" placeholder="password" onChange={this.handlePasswordChange} /><br />
                <input type="submit" className="submit" value="Sign in" />
            </form>
            Don't have an account? <a href="#" onClick={this.handleSignup}>Sign up</a><br />
        </div>;
    },
});

const Register = React.createClass({
    getInitialState: function () {
        return {
            username: "",
            email: "",
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

        this.props.net.request("register", {
            username: this.state.username,
            email: this.state.email,
            referral: this.state.referral,
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
    handleEmailChange: function (e) {
        this.setState({
            email: e.target.value,
        });
    },
    handleReferralChange: function (e) {
        this.setState({
            referral: e.target.value,
        });
    },
    handlePasswordChange: function (e) {
        this.setState({
            password: e.target.value,
        });
    },
    handleLogin: function (e) {
        e.preventDefault();
        this.props.parent.setState({
            step: "login",
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
                <i className="material-icons">person_outline</i>
                <input type="text" placeholder="username" onChange={this.handleUsernameChange} /><br />
                <i className="material-icons">email</i>
                <input type="text" placeholder="email" onChange={this.handleEmailChange} /><br />
                <i className="material-icons">code</i>
                <input type="text" placeholder="invitation code" onChange={this.handleReferralChange} /><br />
                <i className="material-icons">lock_open</i>
                <input type="password" placeholder="password" onChange={this.handlePasswordChange} /><br />
                <input type="submit" className="submit" />
            </form>
            Already have an account? <a href="#" onClick={this.handleLogin}>Login</a><br />
        </div>
    },
});

module.exports = Setup;
