package frinkconv_repl

import (
	"fmt"
	"regexp"
	"strconv"
)

var (
	compiledPattern *regexp.Regexp
)

func init() {
	compiledPattern = regexp.MustCompile(pattern)
}

func extractLastNumber(output string) (float64, error) {
	matches := compiledPattern.FindAllString(output, -1)
	if len(matches) == 0 {
		return 0.0, fmt.Errorf("could not parse %#+v with %#+v", output, pattern)
	}

	rawValue := matches[len(matches)-1]

	value, err := strconv.ParseFloat(rawValue, 64)
	if err != nil {
		return 0.0, err
	}

	return value, nil
}
