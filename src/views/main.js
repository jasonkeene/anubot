
const React = require('react'),
      ReactDOM = require('react-dom'),
      App = require('./app.js');

function render(localStorage) {
    return ReactDOM.render(
        <App localStorage={localStorage} />,
        document.querySelector('#react-anchor'),
    );
}

module.exports.render = render;
