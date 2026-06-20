package utils

import (
	"encoding/json"
	"strings"
)

type FlexBool bool

func (f *FlexBool) UnmarshalJSON(data []byte) error {
	var b bool
	if err := json.Unmarshal(data, &b); err == nil {
		*f = FlexBool(b)
		return nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	s = strings.ToLower(strings.TrimSpace(s))
	switch s {
	case "true", "1", "yes", "on":
		*f = true
	case "false", "0", "no", "off", "":
		*f = false
	default:
		*f = false
	}

	return nil
}

func (f FlexBool) MarshalJSON() ([]byte, error) {
	return json.Marshal(bool(f))
}

func (f FlexBool) Bool() bool {
	return bool(f)
}
