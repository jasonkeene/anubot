"use strict";

const unpack = require("../lib/unpack.js");

describe("unpack", () => {
    it("unpacks JSON encoded data", () => {
        var raw = JSON.stringify({
            cmd: "test-command",
            payload: "test-payload",
        });
        var [cmd, payload] = unpack(raw);
        expect(cmd).toEqual("test-command");
        expect(payload).toEqual("test-payload");
    });

    it("unpacks payload objects", () => {
        var raw = JSON.stringify({
            cmd: "test-command",
            payload: {
                field: "value",
            },
        });
        var [cmd, payload] = unpack(raw);
        expect(cmd).toEqual("test-command");
        expect(payload).toEqual({
            field: "value",
        });
    });
});
