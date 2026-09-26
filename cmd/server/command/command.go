package command

import (
	"encoding/json"
	"fmt"
)

type ValidationError struct {
	Missing  []string
	BadValue []string
}

func (v *ValidationError) IsEmpty() bool {
	return len(v.Missing) == 0 && len(v.BadValue) == 0
}

func (v *ValidationError) AsJsonString() string {
	b, err := json.Marshal(v)
	if err != nil {
		fmt.Println("error parsing json validation message")
		return ""
	}

	return string(b)
}

type Command struct {
	Action string `json:"action"`
	Args   string `json:"args"`
}

func (c Command) Validate() error {
	v := ValidationError{
		Missing: make([]string, 0, 2),
	}

	if c.Action == "" {
		v.Missing = append(v.Missing, "action")
	}

	if v.IsEmpty() {
		return nil
	}

	return fmt.Errorf(v.AsJsonString())
}

func (c Command) ToMessage() (string, error) {
	cmd, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(string(cmd)), nil
}
