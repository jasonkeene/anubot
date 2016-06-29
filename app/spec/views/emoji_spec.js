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
            "Kappa": emoji._global.Kappa,
            "KappaHD": emoji._global.KappaHD,
            "MiniK": emoji._global.MiniK,
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
