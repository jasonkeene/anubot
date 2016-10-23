"use strict";

const unpack = require("../lib/unpack.js");

describe("unpack", () => {
    it("unpacks JSON encoded data", () => {
        var raw = JSON.stringify({
            cmd: "test-command",
            request_id: "test-request-id",
            payload: "test-payload",
            error: "test-error",
        });
        var [cmd, request_id, payload, error] = unpack(raw);
        expect(cmd).toEqual("test-command");
        expect(request_id).toEqual("test-request-id");
        expect(payload).toEqual("test-payload");
        expect(error).toEqual("test-error");
    });

    it("unpacks payload objects", () => {
        var raw = JSON.stringify({
            cmd: "test-command",
            request_id: "test-request-id",
            payload: {
                field: "value",
            },
            error: "test-error",
        });
        var [cmd, request_id, payload, error] = unpack(raw);
        expect(cmd).toEqual("test-command");
        expect(request_id).toEqual("test-request-id");
        expect(payload).toEqual({
            field: "value",
        });
        expect(error).toEqual("test-error");
    });
});
