package main

import (
	"fmt"
	"strconv"
)

//Bool is a convenience type, build around bool, because this is suggested method for go-flags
//to work with a boolean flag that could be flipped both ways. If this isn't done, passing "false" does nothing
type Bool bool

//UnmarshalFlag implements flags.Unmarshaler interface for Bool
func (b *Bool) UnmarshalFlag(value string) error {
	t, err := strconv.ParseBool(value)

	if err != nil {
		return fmt.Errorf("`%s' is not a boolean value, try \"true\"", value)
	}

	*b = Bool(t)
	return nil
}

//MarshalFlag implements flags.Marshaler interface for Bool
func (b Bool) MarshalFlag() (string, error) {
	return strconv.FormatBool(bool(b)), nil
}
