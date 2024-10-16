package predicate

import (
	"fmt"
	"regexp"
	"strconv"
)

var IsNotEmpty = Predicate[string]{
	Test: func(t string) bool {
		return t != ""
	},
	ErrorMessage: errorMessage("string cannot be empty"),
}
var AllowedCharacters = Predicate[string]{
	Test: func(t string) bool {
		valid := regexp.MustCompile(`^[0-9A-Za-z\-._~]+$`)
		return valid.MatchString(t)
	},
	ErrorMessage: errorMessage("\"the parameter can only contain: 0-9, A-Z, a-z, -, ., _, ~"),
}
var NonNegativePredicate = Predicate[string]{
	Test: func(t string) bool {
		toInt, err := strconv.Atoi(t)
		return err == nil && toInt >= 0
	},
	ErrorMessage: errorMessage("the parameter must be a positive number or zero"),
}
var EmailIsValid = Predicate[string]{
	Test: func(t string) bool {
		const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
		re := regexp.MustCompile(emailRegex)

		return re.MatchString(t)
	},
	ErrorMessage: errorMessage("email is invalid"),
}

type Predicate[T any] struct {
	Test         func(t T) bool
	ErrorMessage func(t T) string
}

func errorMessage(template string) func(s string) string {
	return func(s string) string {
		return fmt.Sprintf(template+": %s", s)
	}
}