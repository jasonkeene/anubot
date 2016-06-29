
const React = require("react");

function render(message) {
    var nodes = [message.body];
    if (message.tags === undefined) {
        return nodes;
    }
    if (message.tags.emotes === undefined) {
        return nodes;
    }
    var splices = _splices(message.tags.emotes);

    var insertions = 0,
        processed = 0;
    for (var i = 0; i < splices.length; i++) {
        var splice = splices[i],
            prefix = message.body.substring(processed, splice[0]),
            imgNode = <img className="emoji" src={_emoteURL(splice[2])} />,
            rest = message.body.substring(splice[1]);
        nodes.splice(insertions, 1, prefix, imgNode, rest);
        insertions += 2;
        processed = splice[1];
    }
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

module.exports = {
    render: render,
};
