
const React = require('react');

function render(message) {
    var tags = message.tags;
    if (tags === undefined) {
        return;
    }
    var badges = tags.badges;
    if (badges === undefined) {
        return;
    }
    badges = badges.split(",").map((x) => {
        return x.split('/')[0]
    });
    var nodes = [];
    for (var i = 0; i < badges.length; i++) {
        var imageURL = _images[badges[i]];
        if (imageURL !== undefined) {
            nodes.push(<img src={imageURL} />);
        }
    }
    return nodes;
}

const _images = {
  "admin": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABIAAAASCAYAAABWzo5XAAABfklEQVQ4Ea1Tv0vDQBh9SYyLIDhaEHQRf4CCq6uLglDwH3BWXMXFycXV/8BVoVgHnURQEBR0aUHdKjUtCApW/AFNcud9Z3LeNT+69IM073vvfa93l8RqHw1z9KDsHmTIiD49yC029bYr9ssF5em6Iv7xCN56UANZoGsQqx+C1Q+y5hWfH8QCsOcSmFcCmK+G0oBxRrEhvN8Fq+0DnAHBp6T902nAsmGPrcKZ2oqt6p66ImdyE3ZhWYVItwgkjrS0MoN49ErRP49vJPySE5qs2Bu5jCD+XlHDrHkisTU4AbqoYo6w7qXeOCMyOkOzxEujM7cHe2RF9vLQX84lph89lHrL+EQGRuEuXAnWEs42YPeT579iTmzLP5sHvp6UZmyNhOB2DTz8SYbQiAgmLbhbN0KkpCIjwBvHCC+L4N9epyS4xp/mlROauaI4rFVFcLEE9nqtBtjbjeAWxedSVZwOzDPSFcKWC2dmR7JhZVs8gey3Oz+oMzinT91ajj9T+gXD8I6HuPQ1DgAAAABJRU5ErkJggg==",
  "broadcaster": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABIAAAASCAYAAABWzo5XAAAAhklEQVQ4EWN8LiHxn4EKgIkKZoCNGMYGscDCSOzCBRgTL/3KwAAiz8zMIHb2LAOMT1YYMfHzgw1BtpFkg5gVFRlEDh5ENgPMJtkgwSVLMAwBCZBsEMPfv9Qx6H1MDHUM+vvoEcNrW1sMw0j3GtCI/58/M7wyNkYxjHE006KEBzYOWYGNzSAA6TUbUUpeebAAAAAASUVORK5CYII=",
  "global_mod": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABIAAAASCAYAAABWzo5XAAAAAXNSR0IArs4c6QAAAT1JREFUOBFjZMhX+M9ABcBEBTPARtDHoJXxUxj0pTTBNgpzCzBsz1jAcK/2EFZPsGAVhQraKpsybFOaz8DIyMhw4+VdBg1xZYaLz65j1YLXa1MOLwIbAtIJMuT1l3cMMYsLSDeobfdUFE2iPEIMkUZ+YLEc2zgGEIYBRnzRf6Z4E4MUvzhMLVZaqs4cLI7Ta/uyl4ENefbxJVYD0AUxDALFDsgQUJiAAtik148BFFaEAIpBIENWJ0yDG+I0NQqsHz2sYIYiWwCOfhNZPYZ06ygGT00HBiZoVIcuyIKpRwlUmCDIEGQLwC468/gSg7eWI9iQf///M4AMefv1A0wPAx8HL5wNYqAbAhJDibVnTScZtl8/wJC8vBwkRxJASdmgqDSV0yPJAJhilMAGCZ5+dAkmRxKNYRBJupEUU80gANncWuEJw+OJAAAAAElFTkSuQmCC",
  "moderator": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABIAAAASCAYAAABWzo5XAAABP0lEQVQ4EWM0Wcf1n4EKgIkKZoCNYCHVoECFJIY0zWq4Ns/tymA2SS5ykPRFMQRuGpBBtEGGIjYM5QYTkPWisIkySIVPm6HNdCGKRhBn9o02uBjBMJLkkmOYbL0JrgHG6LpYxLD/2UYYlwHDICZGZoZ///+CFfCzCTHMs98PVwxjVJ+OZzj35giMC6YxvJat3cigwKvOwMnMxbDC+TSKYhAn95g/hiEgcUb0BLnd8y5IHCtIPujE8OzbQ6xyGC4CKcYGIvea4TQEpB4ljBgZGBne/XyFYU7Qbj2G73++YogjC6C4KEu7gWG92xVkeTBbnFMGQwxdAG4QKGx85GLg8iBXwLzpLx/PAIpNfADFazCFWx8tZfjx5xvDsz8PGWB5CSaHi4a7CFnDlKt1DP+BkBSAEf2kaEZWC3cRsiA5bAA3hU+ysIFpygAAAABJRU5ErkJggg==",
  "staff": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABIAAAASCAYAAABWzo5XAAAA1ElEQVQ4EWNU4Df+z0AFwEQFM8BGEDQoLSeW4eazY3CMy2K8BoEMKa7KhOudNWUxnI3OwGkQNkN626aj64fzGXEFNsg7MAByCT5DQOpwughmCIjGZggs7GDqiDJIQlIMph5Mo3sbJEiUQYvWTmGAGYbNEJBBLCACHYAUIwN5BRmGg2c3IAuB2cixiOEiXDaim4IeASgGkWsIyBK413AZgm4zustgfHg6Qk43MEliDQGpR/EazAAQTYohIPVwg0AaYYBUQ0D64F6DGUIuDXcRuQbA9AEAtHhR9n5yEkgAAAAASUVORK5CYII=",
  // TODO: fetch subscriber badge for user dynamically
  //"subscriber": "",
  "turbo": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABIAAAASCAYAAABWzo5XAAAAyElEQVQ4EWNMcVz6n4EKgIkKZoCNGHwGsSB7rWd1IDIXL7skdD2KPFFe+/vnHwO6RhRTgByCBn368IOhPHIjuj4MPorX0GVBruAT5AALF/c6o0jDggHmUpwuAilgZmFi+PT+B9iA3uK9eL2H0yCQblDYIIOQdENkLgobp9dATn/17DNDV/4esIbcVnsGeTUhFM3IHJwGgRTBDKmb7cnAJwAJK5hmWNjA+HgNAinqXO4PDiuYBlw0XoNgMYNLM7I442juRw4OrGwA84AxLKWQUDUAAAAASUVORK5CYII=",
};

module.exports = {
    render: render,
    _images: _images,
};
