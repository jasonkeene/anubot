"use strict";

const unpack = require("../lib/unpack.js");

describe("unpack", () => {
    it("unpacks JSON encoded data", () => {
        var raw = JSON.stringify({
            cmd: "test-command",
            payload: "test-payload",
            error: "test-error",
        });
        var [cmd, payload, error] = unpack(raw);
        expect(cmd).toEqual("test-command");
        expect(payload).toEqual("test-payload");
        expect(error).toEqual("test-error");
    });

    it("unpacks payload objects", () => {
        var raw = JSON.stringify({
            cmd: "test-command",
            payload: {
                field: "value",
            },
            error: "test-error",
        });
        var [cmd, payload, error] = unpack(raw);
        expect(cmd).toEqual("test-command");
        expect(payload).toEqual({
            field: "value",
        });
        expect(error).toEqual("test-error");
    });
});
