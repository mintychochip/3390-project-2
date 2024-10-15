package predicate

import (
	"regexp"
	"strconv"
)

var AllowedCharacters = Predicate[string]{
	Test: func(t string) bool {
		valid := regexp.MustCompile(`^[0-9A-Za-z\-\._~]+$`)
		return valid.MatchString(t)
	},
	Error: "the parameter can only contain: 0-9, A-Z, a-z, -, ., _, ~",
}
var NonNegativePredicate = Predicate[string]{
	Test: func(t string) bool {
		toInt, err := strconv.Atoi(t)
		return err == nil && toInt >= 0
	},
	Error: "invalid parameter: must be non-negative number",
}
var EmailIsValid = Predicate[string]{
	Test: func(t string) bool {
		const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
		re := regexp.MustCompile(emailRegex)

		return re.MatchString(t)
	},
	Error: "the email is invalid",
}

type Predicate[T any] struct {
	Test  func(t T) bool
	Error string
}
