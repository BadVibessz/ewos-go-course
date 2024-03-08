package json

import (
	"bytes"
	"encoding/json"
)

func PrettifyJSON(in string) string {
	var out bytes.Buffer

	if err := json.Indent(&out, []byte(in), "", "\t"); err != nil {
		return in
	}

	return out.String()
}
