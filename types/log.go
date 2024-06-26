package types

type Metadata struct {
	ParentResourceID string `json:"parentResourceId"`
}

type LogData struct {
	ID         string   `json:"id,omitempty"`
	Level      string   `json:"level"`
	Message    string   `json:"message"`
	ResourceID string   `json:"resourceId"`
	Timestamp  string   `json:"timestamp"`
	TraceID    string   `json:"traceId"`
	SpanID     string   `json:"spanId"`
	Commit     string   `json:"commit"`
	Metadata   Metadata `json:"metadata"`
}

type SearchData struct {
	Level      string
	ResourceID string
	TraceID    string
	SpanID     string
	StartDate  string
	EndDate    string
}
