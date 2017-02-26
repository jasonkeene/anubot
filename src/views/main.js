
const React = require('react'),
      ReactDOM = require('react-dom'),
      App = require('./app.js');

function render(localStorage, connect) {
    return ReactDOM.render(
        <App localStorage={localStorage} connect={connect} />,
        document.querySelector('#react-anchor'),
    );
}

module.exports.render = render;
