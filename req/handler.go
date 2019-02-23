package req

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// Handle sends a get request to `echo`
func Handle(req []byte) string {
	response, err := http.Get("http://echo/")
	if err != nil {
		return fmt.Sprintf("%s", err)
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return fmt.Sprintf("%s", err)
		}
		return fmt.Sprintf("%s\n", string(contents))
	}
}
