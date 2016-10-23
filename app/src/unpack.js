
function unpack(raw) {
    var result = JSON.parse(raw);
    return [result.cmd, result.request_id, result.payload, result.error];
}

module.exports = unpack;
