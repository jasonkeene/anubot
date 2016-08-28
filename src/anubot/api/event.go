package api

// event is the structure sent over websocket connections by both ends.
type event struct {
	Cmd       string      `json:"cmd"`        // used to dispatch event to handlers
	Payload   interface{} `json:"payload"`    // actual data being transmitted
	RequestID string      `json:"request_id"` // used for req/resp to group events together
	Error     *apiError   `json:"error"`      // used to indicate an error has occured
}
