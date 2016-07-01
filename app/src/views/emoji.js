
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
                imgNode = <img className="emoji" src={_emoteURL(splice[2])} />,
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
            nodeIndex = _processStringNode(node, nodes, nodeIndex, key);
        }
        nodeIndex++;
    }
}

function _processStringNode(node, nodes, nodeIndex, key) {
    var searchIndex, offset;

    while ((searchIndex = node.indexOf(key, offset)) !== -1) {
        // update offset for next iteration
        offset = searchIndex + key.length;

        // check if match has whitespace around it
        if (!_isolated(node, key, searchIndex)) {
            continue;
        }

        // splice in react nodes
        var prefix = node.substring(0, searchIndex),
            imgNode = <img className="emoji" src={_bttv[key]} />,
            rest = node.substring(searchIndex + key.length);
        nodes.splice(nodeIndex, 1, prefix, imgNode, rest);
        return nodeIndex + 2;
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

const _bttv = {
    "OhMyGoodness": "https://cdn.betterttv.net/emote/54fa925e01e468494b85b54d/1x",
    "PancakeMix": "https://cdn.betterttv.net/emote/54fa927801e468494b85b54e/1x",
    "PedoBear": "https://cdn.betterttv.net/emote/54fa928f01e468494b85b54f/1x",
    "PokerFace": "https://cdn.betterttv.net/emote/54fa92a701e468494b85b550/1x",
    "RageFace": "https://cdn.betterttv.net/emote/54fa92d701e468494b85b552/1x",
    "RebeccaBlack": "https://cdn.betterttv.net/emote/54fa92ee01e468494b85b553/1x",
    "trollface": "https://cdn.betterttv.net/emote/54fa8f1401e468494b85b537/1x",
    "aPliS": "https://cdn.betterttv.net/emote/54fa8f4201e468494b85b538/1x",
    "CiGrip": "https://cdn.betterttv.net/emote/54fa8fce01e468494b85b53c/1x",
    "CHAccepted": "https://cdn.betterttv.net/emote/54fa8fb201e468494b85b53b/1x",
    "FuckYea": "https://cdn.betterttv.net/emote/54fa90d601e468494b85b544/1x",
    "DatSauce": "https://cdn.betterttv.net/emote/54fa903b01e468494b85b53f/1x",
    "ForeverAlone": "https://cdn.betterttv.net/emote/54fa909b01e468494b85b542/1x",
    "GabeN": "https://cdn.betterttv.net/emote/54fa90ba01e468494b85b543/1x",
    "HailHelix": "https://cdn.betterttv.net/emote/54fa90f201e468494b85b545/1x",
    "HerbPerve": "https://cdn.betterttv.net/emote/54fa913701e468494b85b546/1x",
    "iDog": "https://cdn.betterttv.net/emote/54fa919901e468494b85b548/1x",
    "rStrike": "https://cdn.betterttv.net/emote/54fa930801e468494b85b554/1x",
    "ShoopDaWhoop": "https://cdn.betterttv.net/emote/54fa932201e468494b85b555/1x",
    "SwedSwag": "https://cdn.betterttv.net/emote/54fa9cc901e468494b85b565/1x",
    "M&Mjc": "https://cdn.betterttv.net/emote/54fab45f633595ca4c713abc/1x",
    "bttvNice": "https://cdn.betterttv.net/emote/54fab7d2633595ca4c713abf/1x",
    "TopHam": "https://cdn.betterttv.net/emote/54fa934001e468494b85b556/1x",
    "TwaT": "https://cdn.betterttv.net/emote/54fa935601e468494b85b557/1x",
    "WhatAYolk": "https://cdn.betterttv.net/emote/54fa93d001e468494b85b559/1x",
    "WatChuSay": "https://cdn.betterttv.net/emote/54fa99b601e468494b85b55d/1x",
    "aplis!": "https://cdn.betterttv.net/emote/54faa4d801e468494b85b576/1x",
    "Blackappa": "https://cdn.betterttv.net/emote/54faa50d01e468494b85b578/1x",
    "DogeWitIt": "https://cdn.betterttv.net/emote/54faa52f01e468494b85b579/1x",
    "BadAss": "https://cdn.betterttv.net/emote/54faa4f101e468494b85b577/1x",
    "SavageJerky": "https://cdn.betterttv.net/emote/54fb603201abde735115ddb5/1x",
    "Zappa": "https://cdn.betterttv.net/emote/5622aaef3286c42e57d8e4ab/1x",
    "tehPoleCat": "https://cdn.betterttv.net/emote/566ca11a65dbbdab32ec0558/1x",
    "AngelThump": "https://cdn.betterttv.net/emote/566ca1a365dbbdab32ec055b/1x",
    "Kaged": "https://cdn.betterttv.net/emote/54fbf11001abde735115de66/1x",
    "HHydro": "https://cdn.betterttv.net/emote/54fbef6601abde735115de57/1x",
    "TaxiBro": "https://cdn.betterttv.net/emote/54fbefeb01abde735115de5b/1x",
    "BroBalt": "https://cdn.betterttv.net/emote/54fbf00a01abde735115de5c/1x",
    "ButterSauce": "https://cdn.betterttv.net/emote/54fbf02f01abde735115de5d/1x",
    "BaconEffect": "https://cdn.betterttv.net/emote/54fbf05a01abde735115de5e/1x",
    "SuchFraud": "https://cdn.betterttv.net/emote/54fbf07e01abde735115de5f/1x",
    "CandianRage": "https://cdn.betterttv.net/emote/54fbf09c01abde735115de61/1x",
    "She'llBeRight": "https://cdn.betterttv.net/emote/54fbefc901abde735115de5a/1x",
    "OhhhKee": "https://cdn.betterttv.net/emote/54fbefa901abde735115de59/1x",
    "D:": "https://cdn.betterttv.net/emote/55028cd2135896936880fdd7/1x",
    "SexPanda": "https://cdn.betterttv.net/emote/5502874d135896936880fdd2/1x",
    "poolparty": "https://cdn.betterttv.net/emote/5502883d135896936880fdd3/1x",
    ":'(": "https://cdn.betterttv.net/emote/55028923135896936880fdd5/1x",
    "puke": "https://cdn.betterttv.net/emote/550288fe135896936880fdd4/1x",
    "bttvWink": "https://cdn.betterttv.net/emote/550292c0135896936880fdef/1x",
    "bttvAngry": "https://cdn.betterttv.net/emote/550291a3135896936880fde3/1x",
    "bttvConfused": "https://cdn.betterttv.net/emote/550291be135896936880fde4/1x",
    "bttvCool": "https://cdn.betterttv.net/emote/550291d4135896936880fde5/1x",
    "bttvHappy": "https://cdn.betterttv.net/emote/55029200135896936880fde7/1x",
    "bttvSad": "https://cdn.betterttv.net/emote/5502925d135896936880fdea/1x",
    "bttvSleep": "https://cdn.betterttv.net/emote/55029272135896936880fdeb/1x",
    "bttvSurprised": "https://cdn.betterttv.net/emote/55029288135896936880fdec/1x",
    "bttvTongue": "https://cdn.betterttv.net/emote/5502929b135896936880fded/1x",
    "bttvUnsure": "https://cdn.betterttv.net/emote/550292ad135896936880fdee/1x",
    "bttvGrin": "https://cdn.betterttv.net/emote/550291ea135896936880fde6/1x",
    "bttvHeart": "https://cdn.betterttv.net/emote/55029215135896936880fde8/1x",
    "bttvTwink": "https://cdn.betterttv.net/emote/55029247135896936880fde9/1x",
    "VisLaud": "https://cdn.betterttv.net/emote/550352766f86a5b26c281ba2/1x",
    "chompy": "https://cdn.betterttv.net/emote/550b225fff8ecee922d2a3b2/1x",
    "SoSerious": "https://cdn.betterttv.net/emote/5514afe362e6bd0027aede8a/1x",
    "BatKappa": "https://cdn.betterttv.net/emote/550b6b07ff8ecee922d2a3e7/1x",
    "KaRappa": "https://cdn.betterttv.net/emote/550b344bff8ecee922d2a3c1/1x",
    "YetiZ": "https://cdn.betterttv.net/emote/55189a5062e6bd0027aee082/1x",
    "miniJulia": "https://cdn.betterttv.net/emote/552d2fc2236a1aa17a996c5b/1x",
    "FishMoley": "https://cdn.betterttv.net/emote/566ca00f65dbbdab32ec0544/1x",
    "Hhhehehe": "https://cdn.betterttv.net/emote/566ca02865dbbdab32ec0547/1x",
    "KKona": "https://cdn.betterttv.net/emote/566ca04265dbbdab32ec054a/1x",
    "OhGod": "https://cdn.betterttv.net/emote/566ca07965dbbdab32ec0552/1x",
    "PoleDoge": "https://cdn.betterttv.net/emote/566ca09365dbbdab32ec0555/1x",
    "motnahP": "https://cdn.betterttv.net/emote/55288e390fa35376704a4c7a/1x",
    "sosGame": "https://cdn.betterttv.net/emote/553b48a21f145f087fc15ca6/1x",
    "CruW": "https://cdn.betterttv.net/emote/55471c2789d53f2d12781713/1x",
    "RarePepe": "https://cdn.betterttv.net/emote/555015b77676617e17dd2e8e/1x",
    "iamsocal": "https://cdn.betterttv.net/emote/54fbef8701abde735115de58/1x",
    "haHAA": "https://cdn.betterttv.net/emote/555981336ba1901877765555/1x",
    "FeelsBirthdayMan": "https://cdn.betterttv.net/emote/55b6524154eefd53777b2580/1x",
    "RonSmug": "https://cdn.betterttv.net/emote/55f324c47f08be9f0a63cce0/1x",
    "KappaCool": "https://cdn.betterttv.net/emote/560577560874de34757d2dc0/1x",
    "SqShy": "https://cdn.betterttv.net/emote/5622ab523286c42e57d8e4b2/1x",
    "FeelsBadMan": "https://cdn.betterttv.net/emote/566c9fc265dbbdab32ec053b/1x",
    "BasedGod": "https://cdn.betterttv.net/emote/566c9eeb65dbbdab32ec052b/1x",
    "bUrself": "https://cdn.betterttv.net/emote/566c9f3b65dbbdab32ec052e/1x",
    "ConcernDoge": "https://cdn.betterttv.net/emote/566c9f6365dbbdab32ec0532/1x",
    "FapFapFap": "https://cdn.betterttv.net/emote/566c9f9265dbbdab32ec0538/1x",
    "FeelsGoodMan": "https://cdn.betterttv.net/emote/566c9fde65dbbdab32ec053e/1x",
    "FireSpeed": "https://cdn.betterttv.net/emote/566c9ff365dbbdab32ec0541/1x",
    "NaM": "https://cdn.betterttv.net/emote/566ca06065dbbdab32ec054e/1x",
    "SourPls": "https://cdn.betterttv.net/emote/566ca38765dbbdab32ec0560/1x",
    "LUL": "https://cdn.betterttv.net/emote/567b00c61ddbe1786688a633/1x",
    "SaltyCorn": "https://cdn.betterttv.net/emote/56901914991f200c34ffa656/1x",
    "FCreep": "https://cdn.betterttv.net/emote/56d937f7216793c63ec140cb/1x",
    "VapeNation": "https://cdn.betterttv.net/emote/56f5be00d48006ba34f530a4/1x",
    "ariW": "https://cdn.betterttv.net/emote/56fa09f18eff3b595e93ac26/1x",
    "notsquishY": "https://cdn.betterttv.net/emote/5709ab688eff3b595e93c595/1x",
    "DuckerZ": "https://cdn.betterttv.net/emote/573d38b50ffbf6cc5cc38dc9/1x",
};

module.exports = {
    render: render,
    _bttv: _bttv,
};
