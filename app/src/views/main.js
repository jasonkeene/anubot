
const React = require('react'),
      ReactDOM = require('react-dom'),
      App = require('./app.js');

function render(connection, listeners, localStorage) {
    ReactDOM.render(
        <App connection={connection} listeners={listeners}
            localStorage={localStorage} />,
        document.querySelector('#react-anchor')
    );
}

module.exports.render = render;
