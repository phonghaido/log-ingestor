package db

import (
	"database/sql"
	"encoding/json"

	"github.com/labstack/gommon/log"
	_ "github.com/lib/pq"
	"github.com/phonghaido/log-ingestor/data"
)

type PostgreSQL struct {
	Conn string
}

func NewPostgresql(conn string) *PostgreSQL {
	return &PostgreSQL{
		Conn: conn,
	}
}

func (p *PostgreSQL) Connect() (*sql.DB, error) {
	db, err := sql.Open("postgres", p.Conn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (p *PostgreSQL) Insert(db *sql.DB, data data.LogData) error {
	defer db.Close()
	metadataJSON, err := json.Marshal(data.Metadata)
	if err != nil {
		log.Infof("Error parsing metadata to json: %s\n", err.Error())
		return err
	}
	_, err = db.Exec("INSERT INTO logs (level, message, resource_id, timestamp, trace_id, span_id, commit, metadata) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		data.Level, data.Message, data.ResourceID, data.Timestamp, data.TraceID, data.SpanID, data.Commit, metadataJSON)

	if err != nil {
		return err
	}
	log.Printf("Insert log %s to database successfully.\n", data.TraceID)
	return nil
}
