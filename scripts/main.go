package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/phonghaido/log-ingestor/types"
)

func main() {
	url := "http://localhost:3000/log"
	numReqs := 10
	logLevels := []string{"info", "error", "debug", "warning"}

	var wg sync.WaitGroup
	wg.Add(numReqs)

	errCh := make(chan error, numReqs)

	for i := 0; i < numReqs; i++ {
		go func(index int) {
			defer wg.Done()
			randomUUID := uuid.New().String()

			metadata := types.Metadata{
				ParentResourceID: fmt.Sprintf("parent-resource-%s", randomUUID),
			}
			logData := types.LogData{
				Level:      logLevels[rand.Intn(len(logLevels))],
				Message:    fmt.Sprintf("Log message of %s", randomUUID),
				ResourceID: fmt.Sprintf("resource-%s", randomUUID),
				Timestamp:  randomTimestamp(time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC), time.Now()).Format(time.RFC3339),
				TraceID:    fmt.Sprintf("trace-%s", randomUUID),
				SpanID:     fmt.Sprintf("span-%s", randomUUID),
				Commit:     fmt.Sprintf("commit-%s", randomUUID),
				Metadata:   metadata,
			}

			logBytes, err := json.Marshal(logData)
			if err != nil {
				errCh <- err
				return
			}

			resp, err := http.Post(url, "application/json", bytes.NewBuffer(logBytes))
			if err != nil {
				errCh <- err
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					errCh <- fmt.Errorf("error reading response body: %s", err.Error())
				}
				errCh <- fmt.Errorf("unexpected error code: %d, response body: %s", resp.StatusCode, string(body))
				return
			}
		}(i)
	}
	wg.Wait()

	close(errCh)

	for err := range errCh {
		fmt.Println("Error:", err)
	}
}

func randomTimestamp(start, end time.Time) time.Time {
	delta := end.Sub(start)
	nanos := rand.Int63n(delta.Nanoseconds())
	return start.Add(time.Duration(nanos))
}
