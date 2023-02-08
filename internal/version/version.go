package version

import (
	"encoding/json"
	"fmt"
	"os"
)

var (
	major = "1"
	minor = "0"
	build = "0"
)

func PrintVersion() {
	if err := json.NewEncoder(os.Stdout).Encode(struct {
		Major string
		Minor string
		Build string
	}{
		Major: major,
		Minor: minor,
		Build: build,
	}); err != nil {
		fmt.Printf("error while decode version info: %v\n", err)
	}
}
