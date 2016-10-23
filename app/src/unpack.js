
function unpack(raw) {
    var result = JSON.parse(raw);
    return [result.cmd, result.payload, result.error];
}

module.exports = unpack;
