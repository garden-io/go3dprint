package function

import (
	"fmt"
)

// Handle echoes whatever you tell it
func Handle(req []byte) string {
	return fmt.Sprintf("Hello, Go. You said: %s", string(req))
}
