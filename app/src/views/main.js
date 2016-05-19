
const React = require('react'),
      ReactDOM = require('react-dom'),
      AuthOverlay = require('./auth_overlay.js');

function render() {
    ReactDOM.render(<AuthOverlay />, document.querySelector('#app'));
}

exports.render = render;
