"use strict";

const React = require("react"),
      mentions = require("../../lib/views/mentions.js");

describe("render", () => {
    it("renders mentioned username as spans", () => {
        var nodes = mentions.render("test-user", [
            "hello test-user, how are you today, @test-user?",
        ]);
        expect(nodes).toEqual([
            "hello ",
            React.createElement("span", {
                className: "mention",
                key: "plain-mention-test-user-0",
            }, "test-user"),
            ", how are you today, ",
            React.createElement("span", {
                className: "mention",
                key: "at-mention-test-user-0",
            }, "@test-user"),
            "?",
        ]);
    });
});
