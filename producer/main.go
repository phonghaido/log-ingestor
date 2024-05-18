package main

import (
	"log"
	"net/http"

	"github.com/phonghaido/log-ingestor/data"
	"github.com/phonghaido/log-ingestor/producer/handlers"
)

var (
	kafkaConfig   = data.NewKafkaConfig("localhost:9092", "logs", 1, 1)
	kafkaProducer = handlers.NewKafkaProducer(*kafkaConfig)
)

func main() {
	mux := http.NewServeMux()

	srv := &http.Server{
		Addr:    ":3000",
		Handler: mux,
	}
	if err := kafkaProducer.CreateKafkaTopic(); err != nil {
		log.Fatalf("Failed to create topic %s", err)
	}

	log.Println("Server listing on port 3000...")

	mux.HandleFunc("/log", HandlePostLog)

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe: %v\n", err)
	}
}

func HandlePostLog(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Fatal("Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	var reqPayload data.LogData
	log.Printf("Body: %v\n", r.Body)
	if err := handlers.DecodeLogDataReqBody(r.Body, &reqPayload); err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	if err := kafkaProducer.ProduceLogKafka(reqPayload); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Success"))
}
