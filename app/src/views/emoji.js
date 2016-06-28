
const React = require("react");

function render(message) {
    var nodes = [message];
    for (var key in global_emoji) {
        _processNodes(nodes, key);
    }
    return nodes;
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
            imgNode = <img className="emoji" src={global_emoji[key]} />,
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
        console.log("left most");
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

const robot_emoji = {
    ":)": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-ebf60cd72f7aa600-24x18.png",
    ":(": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-d570c4b3b8d8fc4d-24x18.png",
    ":o": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-ae4e17f5b9624e2f-24x18.png",
    ":z": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-b9cbb6884788aa62-24x18.png",
    "B)": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-2cde79cfe74c6169-24x18.png",
    ":\\": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-374120835234cb29-24x18.png",
    ";)": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-cfaf6eac72fe4de6-24x18.png",
    ";p": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-3407bf911ad2fd4a-24x18.png",
    ":p": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-e838e5e34d9f240c-24x18.png",
    "R)": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-0536d670860bf733-24x18.png",
    "o_O": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-8e128fa8dc1de29c-24x18.png",
    ":D": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-9f2ac5d4b53913d7-24x18.png",
    ">(": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-d31223e81104544a-24x18.png",
    "<3": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-577ade91d46d7edc-24x18.png",
};
const turbo_emoji = {
    ":)": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-64f279c77d6f621d-21x18.png",
    ":(": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-c41c5c6c88f481cd-21x18.png",
    ":o": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-a43f189a61cbddbe-21x18.png",
    ":z": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-ff8b4b697171a170-21x18.png",
    "B)": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-9ad04ce3cf69ffd6-21x18.png",
    ":\\": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-7cd0191276363a02-21x18.png",
    ";)": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-54ab3f91053d8b97-21x18.png",
    ";p": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-a66f1856f37d0f48-21x18.png",
    ":p": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-a3ceb91b93f5082b-21x18.png",
    "R)": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-ffe61c02bd7cd500-21x18.png",
    "o_O": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-38b510fc1dd50022-21x18.png",
    ":D": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-1c8ec529616b79e0-21x18.png",
    ">(": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-91a9cf0c00b30760-21x18.png",
    "<3": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-934a78aa6d805cd7-21x18.png",
};
const monkey_emoji = {
    "<3": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-3f5d7d20df6ee956-20x18.png",
    "R)": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-7791c28e2e965fdf-20x22.png",
    ":>": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-665aec4773011f44-27x42.png",
    "<]": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-fd30ca5440d03927-20x42.png",
    ":7": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-206849962fa002dd-29x24.png",
    ":(": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-e4acdcf1ff2b4cef-20x18.png",
    ":P": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-e7b4f5211a173ff1-20x18.png",
    ";P": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-a63a460b5e1f74fc-20x18.png",
    ":O": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-30f5d7516b695012-20x18.png",
    ":\\": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-5067d1fe40f8e607-20x18.png",
    ":|": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-d6e8b4f562b8f46f-20x18.png",
    ":s": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-f5428b0c125bf4a5-20x18.png",
    ":D": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-89679577f86caf4e-20x18.png",
    "o_O": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-8429b9f83d424cb4-20x18.png",
    ">(": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-a6848b1076547d6f-20x18.png",
    ":)": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-ae6e77b75597c3d6-20x18.png",
    "B)": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-f381447031502180-20x18.png",
    ";)": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-095b4874cbf49881-20x18.png",
    "#/": "https://static-cdn.jtvnw.net/jtv_user_pictures/chansub-global-emoticon-39f51e122c6b2d60-27x18.png",
};
const global_emoji = {
    "KappaHD": "https://static-cdn.jtvnw.net/jtv_user_pictures/emoticon-2867-src-f02f9d40f66f0840-28x28.png",
    "MiniK": "https://static-cdn.jtvnw.net/jtv_user_pictures/emoticon-2868-src-5a7a81bb829e1a4c-28x28.png",
    "imGlitch": "https://static-cdn.jtvnw.net/emoticons/v1/62837/1.0",
    "copyThis": "https://static-cdn.jtvnw.net/emoticons/v1/62841/1.0",
    "pastaThat": "https://static-cdn.jtvnw.net/emoticons/v1/62842/1.0",
    "4Head": "https://static-cdn.jtvnw.net/emoticons/v1/354/1.0",
    "AMPEnergy": "https://static-cdn.jtvnw.net/emoticons/v1/99263/1.0",
    "AMPEnergyCherry": "https://static-cdn.jtvnw.net/emoticons/v1/99265/1.0",
    "ANELE": "https://static-cdn.jtvnw.net/emoticons/v1/3792/1.0",
    "ArgieB8": "https://static-cdn.jtvnw.net/emoticons/v1/51838/1.0",
    "ArsonNoSexy": "https://static-cdn.jtvnw.net/emoticons/v1/50/1.0",
    "AsianGlow": "https://static-cdn.jtvnw.net/emoticons/v1/74/1.0",
    "AthenaPMS": "https://static-cdn.jtvnw.net/emoticons/v1/32035/1.0",
    "BabyRage": "https://static-cdn.jtvnw.net/emoticons/v1/22639/1.0",
    "BatChest": "https://static-cdn.jtvnw.net/emoticons/v1/1905/1.0",
    "BCouch": "https://static-cdn.jtvnw.net/emoticons/v1/83536/1.0",
    "BCWarrior": "https://static-cdn.jtvnw.net/emoticons/v1/30/1.0",
    "BibleThump": "https://static-cdn.jtvnw.net/emoticons/v1/86/1.0",
    "BiersDerp": "https://static-cdn.jtvnw.net/emoticons/v1/96824/1.0",
    "BigBrother": "https://static-cdn.jtvnw.net/emoticons/v1/1904/1.0",
    "BionicBunion": "https://static-cdn.jtvnw.net/emoticons/v1/24/1.0",
    "BlargNaut": "https://static-cdn.jtvnw.net/emoticons/v1/38/1.0",
    "bleedPurple": "https://static-cdn.jtvnw.net/emoticons/v1/62835/1.0",
    "BloodTrail": "https://static-cdn.jtvnw.net/emoticons/v1/69/1.0",
    "BORT": "https://static-cdn.jtvnw.net/emoticons/v1/243/1.0",
    "BrainSlug": "https://static-cdn.jtvnw.net/emoticons/v1/881/1.0",
    "BrokeBack": "https://static-cdn.jtvnw.net/emoticons/v1/4057/1.0",
    "BudBlast": "https://static-cdn.jtvnw.net/emoticons/v1/97855/1.0",
    "BuddhaBar": "https://static-cdn.jtvnw.net/emoticons/v1/27602/1.0",
    "BudStar": "https://static-cdn.jtvnw.net/emoticons/v1/97856/1.0",
    "ChefFrank": "https://static-cdn.jtvnw.net/emoticons/v1/90129/1.0",
    "cmonBruh": "https://static-cdn.jtvnw.net/emoticons/v1/84608/1.0",
    "CoolCat": "https://static-cdn.jtvnw.net/emoticons/v1/58127/1.0",
    "CorgiDerp": "https://static-cdn.jtvnw.net/emoticons/v1/49106/1.0",
    "CougarHunt": "https://static-cdn.jtvnw.net/emoticons/v1/21/1.0",
    "DAESuppy": "https://static-cdn.jtvnw.net/emoticons/v1/973/1.0",
    "DalLOVE": "https://static-cdn.jtvnw.net/emoticons/v1/96848/1.0",
    "DansGame": "https://static-cdn.jtvnw.net/emoticons/v1/33/1.0",
    "DatSheffy": "https://static-cdn.jtvnw.net/emoticons/v1/170/1.0",
    "DBstyle": "https://static-cdn.jtvnw.net/emoticons/v1/73/1.0",
    "deExcite": "https://static-cdn.jtvnw.net/emoticons/v1/46249/1.0",
    "deIlluminati": "https://static-cdn.jtvnw.net/emoticons/v1/46248/1.0",
    "DendiFace": "https://static-cdn.jtvnw.net/emoticons/v1/58135/1.0",
    "DogFace": "https://static-cdn.jtvnw.net/emoticons/v1/1903/1.0",
    "DOOMGuy": "https://static-cdn.jtvnw.net/emoticons/v1/54089/1.0",
    "DoritosChip": "https://static-cdn.jtvnw.net/emoticons/v1/102242/1.0",
    "duDudu": "https://static-cdn.jtvnw.net/emoticons/v1/62834/1.0",
    "EagleEye": "https://static-cdn.jtvnw.net/emoticons/v1/20/1.0",
    "EleGiggle": "https://static-cdn.jtvnw.net/emoticons/v1/4339/1.0",
    "FailFish": "https://static-cdn.jtvnw.net/emoticons/v1/360/1.0",
    "FPSMarksman": "https://static-cdn.jtvnw.net/emoticons/v1/42/1.0",
    "FrankerZ": "https://static-cdn.jtvnw.net/emoticons/v1/65/1.0",
    "FreakinStinkin": "https://static-cdn.jtvnw.net/emoticons/v1/39/1.0",
    "FUNgineer": "https://static-cdn.jtvnw.net/emoticons/v1/244/1.0",
    "FunRun": "https://static-cdn.jtvnw.net/emoticons/v1/48/1.0",
    "FutureMan": "https://static-cdn.jtvnw.net/emoticons/v1/98562/1.0",
    "FuzzyOtterOO": "https://static-cdn.jtvnw.net/emoticons/v1/168/1.0",
    "GingerPower": "https://static-cdn.jtvnw.net/emoticons/v1/32/1.0",
    "GrammarKing": "https://static-cdn.jtvnw.net/emoticons/v1/3632/1.0",
    "HassaanChop": "https://static-cdn.jtvnw.net/emoticons/v1/20225/1.0",
    "HassanChop": "https://static-cdn.jtvnw.net/emoticons/v1/68/1.0",
    "HeyGuys": "https://static-cdn.jtvnw.net/emoticons/v1/30259/1.0",
    "HotPokket": "https://static-cdn.jtvnw.net/emoticons/v1/357/1.0",
    "HumbleLife": "https://static-cdn.jtvnw.net/emoticons/v1/46881/1.0",
    "ItsBoshyTime": "https://static-cdn.jtvnw.net/emoticons/v1/169/1.0",
    "Jebaited": "https://static-cdn.jtvnw.net/emoticons/v1/90/1.0",
    "JKanStyle": "https://static-cdn.jtvnw.net/emoticons/v1/15/1.0",
    "JonCarnage": "https://static-cdn.jtvnw.net/emoticons/v1/26/1.0",
    "KAPOW": "https://static-cdn.jtvnw.net/emoticons/v1/9803/1.0",
    "Kappa": "https://static-cdn.jtvnw.net/emoticons/v1/25/1.0",
    "KappaClaus": "https://static-cdn.jtvnw.net/emoticons/v1/74510/1.0",
    "KappaPride": "https://static-cdn.jtvnw.net/emoticons/v1/55338/1.0",
    "KappaRoss": "https://static-cdn.jtvnw.net/emoticons/v1/70433/1.0",
    "KappaWealth": "https://static-cdn.jtvnw.net/emoticons/v1/81997/1.0",
    "Keepo": "https://static-cdn.jtvnw.net/emoticons/v1/1902/1.0",
    "KevinTurtle": "https://static-cdn.jtvnw.net/emoticons/v1/40/1.0",
    "Kippa": "https://static-cdn.jtvnw.net/emoticons/v1/1901/1.0",
    "Kreygasm": "https://static-cdn.jtvnw.net/emoticons/v1/41/1.0",
    "Mau5": "https://static-cdn.jtvnw.net/emoticons/v1/30134/1.0",
    "mcaT": "https://static-cdn.jtvnw.net/emoticons/v1/35063/1.0",
    "MikeHogu": "https://static-cdn.jtvnw.net/emoticons/v1/81636/1.0",
    "MingLee": "https://static-cdn.jtvnw.net/emoticons/v1/68856/1.0",
    "MKXRaiden": "https://static-cdn.jtvnw.net/emoticons/v1/102324/1.0",
    "MKXScorpion": "https://static-cdn.jtvnw.net/emoticons/v1/102325/1.0",
    "MrDestructoid": "https://static-cdn.jtvnw.net/emoticons/v1/28/1.0",
    "MVGame": "https://static-cdn.jtvnw.net/emoticons/v1/29/1.0",
    "NinjaTroll": "https://static-cdn.jtvnw.net/emoticons/v1/45/1.0",
    "NomNom": "https://static-cdn.jtvnw.net/emoticons/v1/90075/1.0",
    "NoNoSpot": "https://static-cdn.jtvnw.net/emoticons/v1/44/1.0",
    "NotATK": "https://static-cdn.jtvnw.net/emoticons/v1/34875/1.0",
    "NotLikeThis": "https://static-cdn.jtvnw.net/emoticons/v1/58765/1.0",
    "OhMyDog": "https://static-cdn.jtvnw.net/emoticons/v1/81103/1.0",
    "OMGScoots": "https://static-cdn.jtvnw.net/emoticons/v1/91/1.0",
    "OneHand": "https://static-cdn.jtvnw.net/emoticons/v1/66/1.0",
    "OpieOP": "https://static-cdn.jtvnw.net/emoticons/v1/100590/1.0",
    "OptimizePrime": "https://static-cdn.jtvnw.net/emoticons/v1/16/1.0",
    "OSfrog": "https://static-cdn.jtvnw.net/emoticons/v1/81248/1.0",
    "OSkomodo": "https://static-cdn.jtvnw.net/emoticons/v1/81273/1.0",
    "OSsloth": "https://static-cdn.jtvnw.net/emoticons/v1/81249/1.0",
    "panicBasket": "https://static-cdn.jtvnw.net/emoticons/v1/22998/1.0",
    "PanicVis": "https://static-cdn.jtvnw.net/emoticons/v1/3668/1.0",
    "PartyTime": "https://static-cdn.jtvnw.net/emoticons/v1/76171/1.0",
    "PazPazowitz": "https://static-cdn.jtvnw.net/emoticons/v1/19/1.0",
    "PeoplesChamp": "https://static-cdn.jtvnw.net/emoticons/v1/3412/1.0",
    "PermaSmug": "https://static-cdn.jtvnw.net/emoticons/v1/27509/1.0",
    "PeteZaroll": "https://static-cdn.jtvnw.net/emoticons/v1/81243/1.0",
    "PeteZarollTie": "https://static-cdn.jtvnw.net/emoticons/v1/81244/1.0",
    "PicoMause": "https://static-cdn.jtvnw.net/emoticons/v1/27/1.0",
    "PipeHype": "https://static-cdn.jtvnw.net/emoticons/v1/4240/1.0",
    "PJSalt": "https://static-cdn.jtvnw.net/emoticons/v1/36/1.0",
    "PJSugar": "https://static-cdn.jtvnw.net/emoticons/v1/102556/1.0",
    "PMSTwin": "https://static-cdn.jtvnw.net/emoticons/v1/92/1.0",
    "PogChamp": "https://static-cdn.jtvnw.net/emoticons/v1/88/1.0",
    "Poooound": "https://static-cdn.jtvnw.net/emoticons/v1/358/1.0",
    "PraiseIt": "https://static-cdn.jtvnw.net/emoticons/v1/38586/1.0",
    "PRChase": "https://static-cdn.jtvnw.net/emoticons/v1/28328/1.0",
    "PunchTrees": "https://static-cdn.jtvnw.net/emoticons/v1/47/1.0",
    "PuppeyFace": "https://static-cdn.jtvnw.net/emoticons/v1/58136/1.0",
    "RaccAttack": "https://static-cdn.jtvnw.net/emoticons/v1/27679/1.0",
    "RalpherZ": "https://static-cdn.jtvnw.net/emoticons/v1/1900/1.0",
    "RedCoat": "https://static-cdn.jtvnw.net/emoticons/v1/22/1.0",
    "ResidentSleeper": "https://static-cdn.jtvnw.net/emoticons/v1/245/1.0",
    "riPepperonis": "https://static-cdn.jtvnw.net/emoticons/v1/62833/1.0",
    "RitzMitz": "https://static-cdn.jtvnw.net/emoticons/v1/4338/1.0",
    "RuleFive": "https://static-cdn.jtvnw.net/emoticons/v1/361/1.0",
    "SeemsGood": "https://static-cdn.jtvnw.net/emoticons/v1/64138/1.0",
    "ShadyLulu": "https://static-cdn.jtvnw.net/emoticons/v1/52492/1.0",
    "ShazBotstix": "https://static-cdn.jtvnw.net/emoticons/v1/87/1.0",
    "ShibeZ": "https://static-cdn.jtvnw.net/emoticons/v1/27903/1.0",
    "SmoocherZ": "https://static-cdn.jtvnw.net/emoticons/v1/89945/1.0",
    "SMOrc": "https://static-cdn.jtvnw.net/emoticons/v1/52/1.0",
    "SMSkull": "https://static-cdn.jtvnw.net/emoticons/v1/51/1.0",
    "SoBayed": "https://static-cdn.jtvnw.net/emoticons/v1/1906/1.0",
    "SoonerLater": "https://static-cdn.jtvnw.net/emoticons/v1/355/1.0",
    "SriHead": "https://static-cdn.jtvnw.net/emoticons/v1/14706/1.0",
    "SSSsss": "https://static-cdn.jtvnw.net/emoticons/v1/46/1.0",
    "StinkyCheese": "https://static-cdn.jtvnw.net/emoticons/v1/90076/1.0",
    "StoneLightning": "https://static-cdn.jtvnw.net/emoticons/v1/17/1.0",
    "StrawBeary": "https://static-cdn.jtvnw.net/emoticons/v1/37/1.0",
    "SuperVinlin": "https://static-cdn.jtvnw.net/emoticons/v1/31/1.0",
    "SwiftRage": "https://static-cdn.jtvnw.net/emoticons/v1/34/1.0",
    "TBCheesePull": "https://static-cdn.jtvnw.net/emoticons/v1/94039/1.0",
    "TBTacoLeft": "https://static-cdn.jtvnw.net/emoticons/v1/94038/1.0",
    "TBTacoRight": "https://static-cdn.jtvnw.net/emoticons/v1/94040/1.0",
    "TF2John": "https://static-cdn.jtvnw.net/emoticons/v1/1899/1.0",
    "TheRinger": "https://static-cdn.jtvnw.net/emoticons/v1/18/1.0",
    "TheTarFu": "https://static-cdn.jtvnw.net/emoticons/v1/70/1.0",
    "TheThing": "https://static-cdn.jtvnw.net/emoticons/v1/7427/1.0",
    "ThunBeast": "https://static-cdn.jtvnw.net/emoticons/v1/1898/1.0",
    "TinyFace": "https://static-cdn.jtvnw.net/emoticons/v1/67/1.0",
    "TooSpicy": "https://static-cdn.jtvnw.net/emoticons/v1/359/1.0",
    "TriHard": "https://static-cdn.jtvnw.net/emoticons/v1/171/1.0",
    "TTours": "https://static-cdn.jtvnw.net/emoticons/v1/38436/1.0",
    "twitchRaid": "https://static-cdn.jtvnw.net/emoticons/v1/62836/1.0",
    "TwitchRPG": "https://static-cdn.jtvnw.net/emoticons/v1/102157/1.0",
    "UleetBackup": "https://static-cdn.jtvnw.net/emoticons/v1/49/1.0",
    "UncleNox": "https://static-cdn.jtvnw.net/emoticons/v1/3666/1.0",
    "UnSane": "https://static-cdn.jtvnw.net/emoticons/v1/71/1.0",
    "VaultBoy": "https://static-cdn.jtvnw.net/emoticons/v1/54090/1.0",
    "VoHiYo": "https://static-cdn.jtvnw.net/emoticons/v1/81274/1.0",
    "Volcania": "https://static-cdn.jtvnw.net/emoticons/v1/166/1.0",
    "WholeWheat": "https://static-cdn.jtvnw.net/emoticons/v1/1896/1.0",
    "WinWaker": "https://static-cdn.jtvnw.net/emoticons/v1/167/1.0",
    "WTRuck": "https://static-cdn.jtvnw.net/emoticons/v1/1897/1.0",
    "WutFace": "https://static-cdn.jtvnw.net/emoticons/v1/28087/1.0",
    "YouWHY": "https://static-cdn.jtvnw.net/emoticons/v1/4337/1.0",
};

module.exports = {
    render: render,
};
