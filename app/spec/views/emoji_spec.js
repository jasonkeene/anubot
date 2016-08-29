"use strict";

const React = require("react"),
      emoji = require("../../lib/views/emoji.js");

describe("render", () => {
    it("returns the bare message if no emojis exist", () => {
        const msg = {
            body: "bare message",
        };
        expect(emoji.render(msg)).toEqual([msg.body]);
    });

    it("renders emoji as react elements", () => {
        var msg = {
            body: "test :) Kappa KappaHD Kappa :) test",
            tags: {
                emotes: "25:8-12,22-26/3286:14-20/499:5-6,28-29",
            },
        };
        expect(emoji.render(msg)).toEqual([
            "test ",
            React.createElement("img", {
                className: "emoji",
                key: "0-499",
                src: "https://static-cdn.jtvnw.net/emoticons/v1/499/1.0",
            }),
            " ",
            React.createElement("img", {
                className: "emoji",
                key: "1-25",
                src: "https://static-cdn.jtvnw.net/emoticons/v1/25/1.0",
            }),
            " ",
            React.createElement("img", {
                className: "emoji",
                key: "2-3286",
                src: "https://static-cdn.jtvnw.net/emoticons/v1/3286/1.0",
            }),
            " ",
            React.createElement("img", {
                className: "emoji",
                key: "3-25",
                src: "https://static-cdn.jtvnw.net/emoticons/v1/25/1.0",
            }),
            " ",
            React.createElement("img", {
                className: "emoji",
                key: "4-499",
                src: "https://static-cdn.jtvnw.net/emoticons/v1/499/1.0",
            }),
            " test",
        ]);
    });

    it("renders bttv emoji", () => {
        var msg = {
            body: "test FeelsGoodMan FeelsBadMan FeelsBadMan FeelsGoodMan test",
        };
        expect(emoji.render(msg)).toEqual([
            "test ",
            React.createElement("img", {
                className: "emoji",
                key: "0-FeelsGoodMan",
                src: emoji._bttv.FeelsGoodMan,
            }),
            " ",
            React.createElement("img", {
                className: "emoji",
                key: "0-FeelsBadMan",
                src: emoji._bttv.FeelsBadMan,
            }),
            " ",
            React.createElement("img", {
                className: "emoji",
                key: "2-FeelsBadMan",
                src: emoji._bttv.FeelsBadMan,
            }),
            " ",
            React.createElement("img", {
                className: "emoji",
                key: "6-FeelsGoodMan",
                src: emoji._bttv.FeelsGoodMan,
            }),
            " test",
        ]);
    });

    it("renders bttv emoji combined with regular emoji", () => {
        var msg = {
            body: "test KappaHD FeelsBadMan :) FeelsGoodMan test",
            tags: {
                emotes: "3286:5-11/499:25-26",
            },
        };
        expect(emoji.render(msg)).toEqual([
            "test ",
            React.createElement("img", {
                className: "emoji",
                key: "0-3286",
                src: "https://static-cdn.jtvnw.net/emoticons/v1/3286/1.0",
            }),
            " ",
            React.createElement("img", {
                className: "emoji",
                key: "2-FeelsBadMan",
                src: emoji._bttv.FeelsBadMan,
            }),
            " ",
            React.createElement("img", {
                className: "emoji",
                key: "1-499",
                src: "https://static-cdn.jtvnw.net/emoticons/v1/499/1.0",
            }),
            " ",
            React.createElement("img", {
                className: "emoji",
                key: "6-FeelsGoodMan",
                src: emoji._bttv.FeelsGoodMan,
            }),
            " test",
        ]);
    });
});
