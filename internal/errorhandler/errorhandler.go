package errorhandler

import (
	"fmt"
	"os"
)

//NewError returns a new error that contains the message of the given one and additionally the given text
func NewError(text string, err error) error {
	return fmt.Errorf("%s: %v", text, err)
}

func ExitIfError(err error, code int) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		os.Exit(code)
	}

}
