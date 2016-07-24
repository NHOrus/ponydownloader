package main

import (
	"fmt"
	"strconv"
)

//Bool is a bool, because this is suggested method for go-flags to work with a boolean flag that could be flipped both ways
type Bool bool

//UnmarshalFlag implements flags.Unmarshaler interface for Bool
func (b *Bool) UnmarshalFlag(value string) error {
	t, err := strconv.ParseBool(value)

	if err != nil {
		return fmt.Errorf("only `true' and `false' are valid values, not `%s'", value)
	}

	*b = Bool(t)
	return nil
}

//MarshalFlag implements flags.Marshaler interface for Bool
func (b Bool) MarshalFlag() (string, error) {
	return strconv.FormatBool(bool(b)), nil
}
