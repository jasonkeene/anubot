
const React = require("react");

function render(message) {
    var nodes = [message.body];
    if (message.tags !== undefined && message.tags.emotes !== undefined) {
        var splices = _splices(message.tags.emotes),
            insertions = 0,
            processed = 0;
        for (var i = 0; i < splices.length; i++) {
            var splice = splices[i],
                prefix = message.body.substring(processed, splice[0]),
                imgNode = <img className="emoji"
                                key={i + "-" + splice[2]}
                                src={_emoteURL(splice[2])} />,
                rest = message.body.substring(splice[1]);
            nodes.splice(insertions, 1, prefix, imgNode, rest);
            insertions += 2;
            processed = splice[1];
        }
    }
    _renderBTTV(nodes);
    return nodes;
}

function _splices(emotes) {
    var emoteSplices = emotes.split('/'),
        splices = [];
    for (var i = 0; i < emoteSplices.length; i++) {
        var splice = emoteSplices[i].split(':');
        if (splice.length !== 2) {
            continue
        }
        var emoteID = parseInt(splice[0]),
            indexes = splice[1].split(',');
        for (var j = 0; j < indexes.length; j++) {
            // TODO: figure out better names for these or find ways to not
            // have named values
            var Indexes = indexes[j].split('-');
            if (Indexes.length !== 2) {
                continue
            }
            var start = parseInt(Indexes[0]),
                end = parseInt(Indexes[1]) + 1;
            splices.push([start, end, emoteID]);
        }
    }
    splices.sort((a, b) => {
        if (a[0] < b[0]) {
            return -1;
        }
        if (a[0] > b[0]) {
            return 1;
        }
        return 0;
    });
    return splices;
}

function _emoteURL(emoteID) {
    return "https://static-cdn.jtvnw.net/emoticons/v1/" + emoteID + "/1.0";
}

function _renderBTTV(nodes) {
    for (var key in _bttv) {
        _processNodes(nodes, key);
    }
}

function _processNodes(nodes, key) {
    var nodeIndex = 0;
    while (nodeIndex < nodes.length) {
        var node = nodes[nodeIndex];
        if (typeof node === "string") {
            _processStringNode(node, nodes, nodeIndex, key);
        }
        nodeIndex++;
    }
}

function _processStringNode(node, nodes, nodeIndex, key) {
    var searchIndex = 0,
        offset = 0;

    while ((searchIndex = node.indexOf(key, offset)) !== -1) {
        // update offset for next iteration
        offset = searchIndex + key.length;

        // check if match has whitespace around it
        if (!_isolated(node, key, searchIndex)) {
            continue;
        }

        // splice in react nodes
        var prefix = node.substring(0, searchIndex),
            imgNode = <img className="emoji"
                            key={nodeIndex + "-" + key}
                            src={_bttv[key]} />,
            rest = node.substring(searchIndex + key.length);
        nodes.splice(nodeIndex, 1, prefix, imgNode, rest);
        break;
    }
}

function _isolated(node, key, i) {
    return _leftIsolated(node, key, i) && _rightIsolated(node, key, i);
}

function _leftIsolated(node, key, i) {
    if (i === 0) {
        return true;
    }
    var before = node[i-1];
    if (before === " " ||
        before === "\t" ||
        before === "\n") {
        return true;
    }
    return false;
}

function _rightIsolated(node, key, i) {
    if (i+key.length === node.length) {
        return true;
    }
    var after = node[i+key.length];
    if (after === " " ||
        after === "\t" ||
        after === "\n") {
        return true;
    }
    return false;
}

function initBTTV(emoji) {
    _bttv = emoji;
}

var _bttv = {};

module.exports = {
    render: render,
    initBTTV: initBTTV,
    _bttv: _bttv,
};
