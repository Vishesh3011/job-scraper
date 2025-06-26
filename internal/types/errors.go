package types

import "fmt"

var (
	ErrRecordNotFound error = fmt.Errorf("user not found")
)
