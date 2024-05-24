package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/phonghaido/log-ingestor/data"
)

func main() {
	url := "http://localhost:3000/log"
	numReqs := 1000
	logLevels := []string{"info", "error", "debug", "warning"}

	var wg sync.WaitGroup
	wg.Add(numReqs)

	errCh := make(chan error, numReqs)

	for i := 0; i < numReqs; i++ {
		go func(index int) {
			defer wg.Done()
			randomUUID := uuid.New().String()

			logData := data.LogData{
				Level:      logLevels[rand.Intn(len(logLevels))],
				Message:    fmt.Sprintf("Log message of %s", randomUUID),
				ResourceID: fmt.Sprintf("resource-%s", randomUUID),
				Timestamp:  randomTimestamp(time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC), time.Now()).Format(time.RFC3339),
				TraceID:    fmt.Sprintf("trace-%s", randomUUID),
				SpanID:     fmt.Sprintf("span-%s", randomUUID),
				Commit:     fmt.Sprintf("commit-%s", randomUUID),
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
				errCh <- fmt.Errorf("unexpected error code: %d", resp.StatusCode)
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
