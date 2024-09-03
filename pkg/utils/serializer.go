package utils

import (
	"bytes"
	"encoding/json"
)

type Message map[string]interface{}

func SerializeToBytes(msg Message) ([]byte, error) {
	var b bytes.Buffer
	encoder := json.NewEncoder(&b)
	err := encoder.Encode(msg)
	return b.Bytes(), err
}
