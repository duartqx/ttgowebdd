package colors

import (
	"fmt"
	"slices"
)

const (
	Blue    string = "\033[34m"
	Cyan    string = "\033[36m"
	Green   string = "\033[32m"
	Magenta string = "\033[35m"
	Red     string = "\033[31m"
	Reset   string = "\033[0m"
	Yellow  string = "\033[33m"
)

var clrs []string = []string{Blue, Cyan, Green, Magenta, Red, Yellow}

func ColorIt(color string, value interface{}) (string, error) {
	if !slices.Contains[[]string](clrs, color) {
		return Red, fmt.Errorf("Color not recognized!")
	}
	return fmt.Sprintf("%s%v%s", color, value, Reset), nil
}
