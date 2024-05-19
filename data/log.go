package data

type Metadata struct {
	ParentResourceID string `json:"parentResourceId"`
}

type LogData struct {
	Level      string   `json:"level"`
	Message    string   `json:"message"`
	ResourceID string   `json:"resourceId"`
	Timestamp  string   `json:"timestamp"`
	TraceID    string   `json:"traceId"`
	SpanID     string   `json:"spanId"`
	Commit     string   `json:"commit"`
	Metadata   Metadata `json:"metadata"`
}
