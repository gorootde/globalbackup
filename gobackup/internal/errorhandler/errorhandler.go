package errorhandler

import (
	"fmt"
)

//NewError returns a new error that contains the message of the given one and additionally the given text
func NewError(text string, err error) error {
	return fmt.Errorf("%s: %v", text, err)
}
