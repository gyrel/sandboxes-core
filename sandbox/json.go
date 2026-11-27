package sandbox

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
)

// UnmarshalJSON decodes a sandbox record strictly and normalizes empty labels.
func (sb *Sandbox) UnmarshalJSON(data []byte) error {
	type record Sandbox
	var decoded record
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&decoded); err != nil {
		return err
	}
	if err := decoder.Decode(new(any)); err != io.EOF {
		return errors.New("trailing sandbox JSON")
	}
	if err := ValidateLabels(decoded.Labels); err != nil {
		return err
	}
	if len(decoded.Labels) == 0 {
		decoded.Labels = nil
	}
	*sb = Sandbox(decoded)
	return nil
}
