package utilities

import "regexp"

func HasNumber(str string) bool {
	re := regexp.MustCompile("^[^0-9]+$")
    return !re.MatchString(str)
}

func HasSymbol(str string) bool {
	re := regexp.MustCompile("^[a-zA-Z]+$")
    return !re.MatchString(str)
}