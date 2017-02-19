
const React = require('react'),
      ReactDOM = require('react-dom'),
      App = require('./app.js');

function render(net, localStorage) {
    ReactDOM.render(
        <App net={net} localStorage={localStorage} />,
        document.querySelector('#react-anchor')
    );
}

module.exports.render = render;
