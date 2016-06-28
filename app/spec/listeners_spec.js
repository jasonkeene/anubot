"use strict";

const Listeners = require("../lib/listeners.js");

describe("Listeners", () => {
    var listeners,
        mockListener;

    beforeEach(() => {
        listeners = new Listeners();
        mockListener = jasmine.createSpy("mockListener");
    });

    it("can add listeners that respond to commands by name", () => {
        var uncalledMockListener = jasmine.createSpy("mockListener");

        listeners.add("test-command", mockListener);
        listeners.dispatch("test-command", "test-payload");
        expect(mockListener).toHaveBeenCalledWith("test-payload");
        expect(uncalledMockListener).not.toHaveBeenCalled();
    });

    it("does not dispatch commands that have no listeners", () => {
        listeners.add("test-command", mockListener);

        listeners.dispatch("bad-test-command", "test-payload");
        expect(mockListener).not.toHaveBeenCalled();
    });

    it("can remove listeners by command name", () => {
        // dispatch an event
        listeners.add("test-command", mockListener);
        listeners.dispatch("test-command", "test-payload");
        expect(mockListener).toHaveBeenCalledWith("test-payload");
        mockListener.calls.reset();

        // now remove this listener and it shouldn't get the second event
        listeners.remove("test-command", mockListener);
        listeners.dispatch("test-command", "test-payload");
        expect(mockListener).not.toHaveBeenCalled();
    });

    it("does not attempt to remove listeners from unknown commands", () => {
        listeners.add("test-command", mockListener);
        listeners.remove("bad-test-command", mockListener);
        listeners.dispatch("test-command", "test-payload");
        expect(mockListener).toHaveBeenCalledWith("test-payload");
    });
});
