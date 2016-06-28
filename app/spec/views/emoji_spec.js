"use strict";

const React = require("react"),
      emoji = require("../../lib/views/emoji.js");

describe("render", () => {
    it("returns the bare message if no emojis exist", () => {
        const msg = "bare message";
        expect(emoji.render(msg)).toEqual([msg]);
    });

    it("renders global emoji as react elements", () => {
        const fixtures = {
            "Kappa": "https://static-cdn.jtvnw.net/emoticons/v1/25/1.0",
            "KappaHD": "https://static-cdn.jtvnw.net/jtv_user_pictures/emoticon-2867-src-f02f9d40f66f0840-28x28.png",
            "MiniK": "https://static-cdn.jtvnw.net/jtv_user_pictures/emoticon-2868-src-5a7a81bb829e1a4c-28x28.png",
        };
        for (var name in fixtures) {
            var msg = "test " + name + " test";
            expect(emoji.render(msg)).toEqual([
                "test ",
                React.createElement(
                    "img",
                    {
                        className: "emoji",
                        src: fixtures[name],
                    }
                ),
                " test",
            ]);
        }
    });

    it("does not render emojis that are not surrounded by white space", () => {
        const msg = "test 'Kappa' test";
        expect(emoji.render(msg)).toEqual([msg]);
    });
});
