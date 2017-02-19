"use strict";

const Listeners = require("../lib/listeners.js");

describe("Listeners", () => {
    var listeners,
        mockListener;

    beforeEach(() => {
        listeners = new Listeners();
        mockListener = jasmine.createSpy("mockListener");
    });

    describe("cmd listeners", () => {
        it("can add listeners that respond to commands", () => {
            listeners.cmd("test-command", mockListener);
            listeners.dispatch("test-command", "", "test-payload", "test-error");
            expect(mockListener).toHaveBeenCalledWith("test-payload", "test-error");
        });

        it("does not dispatch to listeners for other commands", () => {
            listeners.cmd("test-command", mockListener);
            listeners.dispatch("bad-test-command", "", "test-payload", "test-error");
            expect(mockListener).not.toHaveBeenCalled();
        });

        it("can remove cmd listeners", () => {
            // dispatch an event
            var id = listeners.cmd("test-command", mockListener);
            listeners.dispatch("test-command", "", "test-payload", "test-error");
            expect(mockListener).toHaveBeenCalledWith("test-payload", "test-error");
            mockListener.calls.reset();

            // now remove this listener and it shouldn't get the second event
            listeners.remove(id);
            listeners.dispatch("test-command", "", "test-payload", "test-error");
            expect(mockListener).not.toHaveBeenCalled();
        });
    });

    describe("request listeners", () => {
        it("can add listeners that respond to specific requests", () => {
            listeners.request("test-request-id", mockListener);
            listeners.dispatch("test-command", "test-request-id", "test-payload", "test-error");
            expect(mockListener).toHaveBeenCalledWith("test-payload", "test-error");
        });

        it("does not dispatch when request ids do not match", () => {
            listeners.request("test-request-id", mockListener);
            listeners.dispatch("test-command", "bad-request-id", "test-payload", "test-error");
            expect(mockListener).not.toHaveBeenCalledWith("test-payload", "test-error");
        });

        it("can remove request listeners", () => {
            // dispatch an event
            var id = listeners.request("test-request-id", mockListener);
            listeners.dispatch("test-command", "test-request-id", "test-payload", "test-error");
            expect(mockListener).toHaveBeenCalledWith("test-payload", "test-error");
            mockListener.calls.reset();

            // now remove this listener and it shouldn't get the second event
            listeners.remove(id);
            listeners.dispatch("test-command", "test-request-id", "test-payload", "test-error");
            expect(mockListener).not.toHaveBeenCalled();
        });
    });
});
