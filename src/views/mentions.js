
const React = require("react");

function render(username, nodes) {
    renderAtMentions(username, nodes);
    renderPlainMentions(username, nodes);
    return nodes;
}

function renderAtMentions(username, nodes) {
    for (var i = 0; i < nodes.length; i++) {
        var node = nodes[i],
            searchIndex = 0;
        if (typeof node === "string") {
            searchIndex = node.indexOf("@" + username);
            if (searchIndex !== -1) {
                var prefix = node.substring(0, searchIndex),
                    spanNode = <span className="mention"
                                     key={"at-mention-" + username + "-" + i}>{"@" + username}</span>,
                    rest = node.substring(searchIndex + 1 + username.length);
                nodes.splice(i, 1, prefix, spanNode, rest);
                continue;
            }
        }
    }
}

function renderPlainMentions(username, nodes) {
    for (var i = 0; i < nodes.length; i++) {
        var node = nodes[i],
            searchIndex = 0;
        if (typeof node === "string") {
            searchIndex = node.indexOf(username);
            if (searchIndex !== -1) {
                var prefix = node.substring(0, searchIndex),
                    spanNode = <span className="mention"
                                     key={"plain-mention-" + username + "-" + i}>{username}</span>,
                    rest = node.substring(searchIndex + username.length);
                nodes.splice(i, 1, prefix, spanNode, rest);
                continue;
            }
        }
    }
}

module.exports = {
    render: render,
};
