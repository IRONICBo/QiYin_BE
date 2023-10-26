package common

// msg is a mapping of message.
var msg = map[int]string{
	SUCCESS:        "success",
	ERROR:          "error",
	INVALID_PARAMS: "request params error",
	UNAUTHORIZED:   "unauthorized",
	FORBIDDEN:      "forbidden",
}

// GetMsg get the message by code.
func GetMsg(code int) string {
	m, ok := msg[code]
	if ok {
		return m
	}

	return msg[ERROR]
}
