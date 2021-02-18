package function

import (
	"fmt"
)

// Handle a serverless request
func Handle(req []byte) string {
	return fmt.Sprintf(`<html><body>Hello, checkout <a href="https://skyline.github.com/alexellis/2020">my GitHub Skyline from 2020</a></body></html>`)
}
