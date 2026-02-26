package atcommands

import "fmt"

// === ATError ===
type ATError struct {
	Code    int
	Message string
	Command string
}

func (e ATError) Error() string {
	if e.Command != "" {
		return fmt.Sprintf("AT command '%s' error %d: %s", e.Command, e.Code, e.Message)
	}
	return fmt.Sprintf("AT error %d: %s", e.Code, e.Message)
}

func NewATError(code int, command string) ATError {
	message := errorCodes[code]
	if message == "" {
		message = "Unknown error"
	}
	return ATError{
		Code:    code,
		Message: message,
		Command: command,
	}
}

var errorCodes = map[int]string{
	-1: "invalid command format",
	-2: "invalid command",
	-3: "Not yet defined",
	-4: "invalid parameter",
	-5: "Not yet defined",
}
