"use strict";

const React = require("react"),
      badges = require("../../lib/views/badges.js");

describe("render", () => {
    it("renders badges", () => {
        var msg = {
            tags: {
                badges: "admin/1,turbo/1",
            },
        };
        expect(badges.render(msg)).toEqual([
            React.createElement(
                "img",
                {
                    key: "admin",
                    src: badges._images.admin,
                }
            ),
            React.createElement(
                "img",
                {
                    key: "turbo",
                    src: badges._images.turbo,
                }
            ),
        ]);
    });

    it("ignores invalid badges", () => {
        var msg = {
            tags: {
                badges: "invalid-badge/1",
            },
        };
        expect(badges.render(msg)).toEqual([]);
    });

    it("returns undefined if there are no tags", () => {
        var msg = {};
        expect(badges.render(msg)).toBeUndefined();
    });

    it("returns undefined if there are no badges", () => {
        var msg = {
            tags: {},
        };
        expect(badges.render(msg)).toBeUndefined();
    });
});
