package spam

import (
	"errors"
	"strings"
)

func validateResponseSlashes(response string) error {
	if !strings.HasPrefix(response, "/") || strings.HasPrefix(response, "/me") || strings.HasPrefix(response, "/announce") {
		return nil
	} else if strings.HasPrefix(response, "/") {
		return errors.New("Slash commands except /me and /announce is disallowed. This response wont be ever sended.")
	} else {
		return nil
	}
}
