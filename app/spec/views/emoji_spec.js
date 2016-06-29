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
                src: "https://static-cdn.jtvnw.net/emoticons/v1/499/1.0",
            }),
            " ",
            React.createElement("img", {
                className: "emoji",
                src: "https://static-cdn.jtvnw.net/emoticons/v1/25/1.0",
            }),
            " ",
            React.createElement("img", {
                className: "emoji",
                src: "https://static-cdn.jtvnw.net/emoticons/v1/3286/1.0",
            }),
            " ",
            React.createElement("img", {
                className: "emoji",
                src: "https://static-cdn.jtvnw.net/emoticons/v1/25/1.0",
            }),
            " ",
            React.createElement("img", {
                className: "emoji",
                src: "https://static-cdn.jtvnw.net/emoticons/v1/499/1.0",
            }),
            " test",
        ]);
    });
});
