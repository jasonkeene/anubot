
const React = require('react'),
      ReactDOM = require('react-dom'),
      App = require('./app.js');

function render(connection, listeners) {
    ReactDOM.render(<App connection={connection} listeners={listeners} />,
                    document.querySelector('#app'));
}

module.exports.render = render;
