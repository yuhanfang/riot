package types

import (
	"encoding/json"
	"strconv"
	"strings"
)

// StringInt64 is a JSON string that represents an int64 value.
type StringInt64 int64

func (s *StringInt64) UnmarshalJSON(b []byte) error {
	var str string
	err := json.Unmarshal(b, &str)
	if err != nil {
		return err
	}
	intval, err := strconv.ParseInt(strings.TrimSpace(str), 10, 64)
	if err != nil {
		return err
	}
	*s = StringInt64(intval)
	return nil
}

func (s StringInt64) Int64() int64 {
	return int64(s)
}
