package handlers

import (
	"encoding/json"
	"io"
	"log"

	"github.com/phonghaido/log-ingestor/data"
)

func DecodeLogDataReqBody(body io.ReadCloser, reqPayload *data.LogData) error {
	if err := json.NewDecoder(body).Decode(reqPayload); err != nil {
		log.Fatal("Error decoding request body")
		return err
	}
	return nil
}
