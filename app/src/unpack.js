
function unpack(raw) {
    var result = JSON.parse(raw);
    return [result.cmd, result.payload];
}

module.exports = unpack;
