package helpers

import (
	"encoding/json"
	"net/http"
)

func JsonMarshal(data any) ([]byte, error) {
	return json.Marshal(data)
}

func WriteJSON(w http.ResponseWriter, body []byte) error {
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write(body)
	return err
}

func JsonEncode(w http.ResponseWriter, data any) error {
	body, err := JsonMarshal(data)
	if err != nil {
		return err
	}
	return WriteJSON(w, append(body, '\n'))
}
