package predicate

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

// CSVContainsNumeric /*
/*
Determines whether the file `multipart.File` contains at least one numeric value,
'Test' returns 'true' if at least one numeric is found.
*/
var CSVContainsNumeric = Predicate[io.Reader]{
	Test: func(t io.Reader) bool {

		var buf bytes.Buffer

		// Copy the contents of the reader into the buffer
		if _, err := io.Copy(&buf, t); err != nil {
			return false
		}

		// Create a new reader from the buffer
		reader := csv.NewReader(&buf)

		// Read all records from the new reader
		records, err := reader.ReadAll()
		if err != nil {
			return false
		}

		// Check if any records exist
		if len(records) == 0 {
			return false
		}

		for _, row := range records {
			for _, value := range row {
				trimmedValue := strings.TrimSpace(value)
				if _, err := strconv.ParseFloat(trimmedValue, 64); err == nil {
					return true
				}
			}
		}
		return false
	},

	ErrorMessage: func(t io.Reader) string {
		return "does not contain numerics" // Customize the error message as needed
	},
}

// IsNotEmpty /*
/*
Checks whether a `string` is empty.
'Test' returns 'true' if the string is non-empty e.g. not "".
*/
var IsNotEmpty = Predicate[string]{
	Test: func(t string) bool {
		return t != ""
	},
	ErrorMessage: ErrorMessage("string cannot be empty"),
}

//AllowedCharacters /*
/*
Checks whether a `string` contains the characters that are specified in the Microsoft Azure REST API design guidelines:
https://learn.microsoft.com/en-us/azure/architecture/best-practices/api-design
https://github.com/microsoft/api-guidelines/blob/vNext/azure/Guidelines.md
"✅ DO restrict the characters in service-defined path segments to 0-9  A-Z  a-z  -  .  _  ~, with : allowed only as described below to designate an action operation."
'Test' returns 'true' if a string contains characters that are specified in the regex pattern.
*/
var AllowedCharacters = Predicate[string]{
	Test: func(t string) bool {
		valid := regexp.MustCompile(`^[0-9A-Za-z\-._~]+$`)
		return valid.MatchString(t)
	},
	ErrorMessage: ErrorMessage("\"the parameter can only contain: 0-9, A-Z, a-z, -, ., _, ~"),
}

//NonNegative /*
/*
Checks whether a string when parsed to an integer is a positive integer
'Test' returns 'true' if a string is greater than or equal to 0.
*/
var NonNegative = Predicate[string]{
	Test: func(t string) bool {
		toInt, err := strconv.Atoi(t)
		return err == nil && toInt >= 0
	},
	ErrorMessage: ErrorMessage("the parameter must be a positive number or zero"),
}

//EmailIsValid /*
/*
Checks whether the email provided contains alphanumeric characters only, and is in the format [x]@[y].[z]
'Test' returns 'true' if this is the case.
*/
var EmailIsValid = Predicate[string]{
	Test: func(t string) bool {
		const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
		re := regexp.MustCompile(emailRegex)

		return re.MatchString(t)
	},
	ErrorMessage: ErrorMessage("email is invalid"),
}

type Predicate[T any] struct {
	Test         func(t T) bool
	ErrorMessage func(t T) string
}

func ErrorMessage(template string) func(s string) string {
	return func(s string) string {
		return fmt.Sprintf(template+": %s", s)
	}
}
