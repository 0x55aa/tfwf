package tfwf

import (
	"fmt"
)


var settings map[string]string


func init() {
    settings = map[string]string{"template_dir": "templates", "static_dir": "static", }
}


func AddSettings(m map[string]string) {
	for k, v := range m {
		settings[k] = v
	}
}

//test
func PrintSettings() {
	for k, v := range settings {
		fmt.Printf("%s: %s\n", k, v)
	}
}
