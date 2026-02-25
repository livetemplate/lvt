package telemetry

import (
	"context"
	"database/sql"
	_ "embed"
	"encoding/json"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

//go:embed schema.sql
var schemaSQL string

// SQLiteStore implements Store using SQLite via modernc.org/sqlite.
type SQLiteStore struct {
	db *sql.DB
}

// OpenSQLite opens (or creates) a telemetry database at the given path.
func OpenSQLite(dbPath string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open telemetry db: %w", err)
	}
	s := &SQLiteStore{db: db}
	if err := s.ensureSchema(); err != nil {
		db.Close()
		return nil, err
	}
	return s, nil
}

func (s *SQLiteStore) ensureSchema() error {
	_, err := s.db.Exec(schemaSQL)
	if err != nil {
		return fmt.Errorf("create telemetry schema: %w", err)
	}
	return nil
}

func (s *SQLiteStore) Save(ctx context.Context, event *GenerationEvent) error {
	inputsJSON, err := json.Marshal(event.Inputs)
	if err != nil {
		return fmt.Errorf("marshal inputs: %w", err)
	}
	errorsJSON, err := json.Marshal(event.Errors)
	if err != nil {
		return fmt.Errorf("marshal errors: %w", err)
	}
	filesJSON, err := json.Marshal(event.FilesGenerated)
	if err != nil {
		return fmt.Errorf("marshal files: %w", err)
	}

	_, err = s.db.ExecContext(ctx,
		`INSERT INTO generation_events (id, timestamp, command, inputs, kit, lvt_version, success, validation, errors, duration_ms, files_generated)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		event.ID,
		event.Timestamp.UTC().Format(time.RFC3339),
		event.Command,
		string(inputsJSON),
		event.Kit,
		event.LvtVersion,
		event.Success,
		event.ValidationJSON,
		string(errorsJSON),
		event.DurationMs,
		string(filesJSON),
	)
	return err
}

func (s *SQLiteStore) Get(ctx context.Context, id string) (*GenerationEvent, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT id, timestamp, command, inputs, kit, lvt_version, success, validation, errors, duration_ms, files_generated
		 FROM generation_events WHERE id = ?`, id)
	return scanEvent(row)
}

func (s *SQLiteStore) List(ctx context.Context, opts ListOptions) ([]*GenerationEvent, error) {
	query := `SELECT id, timestamp, command, inputs, kit, lvt_version, success, validation, errors, duration_ms, files_generated
	          FROM generation_events WHERE 1=1`
	var args []any

	if !opts.Since.IsZero() {
		query += " AND timestamp >= ?"
		args = append(args, opts.Since.UTC().Format(time.RFC3339))
	}
	if opts.SuccessOnly != nil {
		query += " AND success = ?"
		args = append(args, *opts.SuccessOnly)
	}
	if opts.Command != "" {
		query += " AND command = ?"
		args = append(args, opts.Command)
	}
	query += " ORDER BY timestamp DESC"
	if opts.Limit > 0 {
		query += " LIMIT ?"
		args = append(args, opts.Limit)
	}

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*GenerationEvent
	for rows.Next() {
		e, err := scanEventRows(rows)
		if err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, rows.Err()
}

func (s *SQLiteStore) CountBySuccess(ctx context.Context, since time.Time) (total, successes int, err error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT COUNT(*), COALESCE(SUM(CASE WHEN success THEN 1 ELSE 0 END), 0)
		 FROM generation_events WHERE timestamp >= ?`,
		since.UTC().Format(time.RFC3339))
	err = row.Scan(&total, &successes)
	return
}

func (s *SQLiteStore) DeleteBefore(ctx context.Context, before time.Time) (int64, error) {
	res, err := s.db.ExecContext(ctx,
		`DELETE FROM generation_events WHERE timestamp < ?`,
		before.UTC().Format(time.RFC3339))
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (s *SQLiteStore) Close() error {
	return s.db.Close()
}

// scanner is satisfied by both *sql.Row and *sql.Rows.
type scanner interface {
	Scan(dest ...any) error
}

func scanInto(sc scanner) (*GenerationEvent, error) {
	var e GenerationEvent
	var ts, inputsJSON, errorsJSON, filesJSON string
	var kit, lvtVersion, validationJSON sql.NullString
	var durationMs sql.NullInt64

	err := sc.Scan(
		&e.ID, &ts, &e.Command, &inputsJSON,
		&kit, &lvtVersion, &e.Success, &validationJSON,
		&errorsJSON, &durationMs, &filesJSON,
	)
	if err != nil {
		return nil, err
	}

	e.Timestamp, err = time.Parse(time.RFC3339, ts)
	if err != nil {
		return nil, fmt.Errorf("parse timestamp %q: %w", ts, err)
	}
	e.Kit = kit.String
	e.LvtVersion = lvtVersion.String
	e.ValidationJSON = validationJSON.String
	e.DurationMs = durationMs.Int64

	if inputsJSON != "" && inputsJSON != "null" {
		if err := json.Unmarshal([]byte(inputsJSON), &e.Inputs); err != nil {
			return nil, fmt.Errorf("unmarshal inputs: %w", err)
		}
	}
	if errorsJSON != "" && errorsJSON != "null" {
		if err := json.Unmarshal([]byte(errorsJSON), &e.Errors); err != nil {
			return nil, fmt.Errorf("unmarshal errors: %w", err)
		}
	}
	if filesJSON != "" && filesJSON != "null" {
		if err := json.Unmarshal([]byte(filesJSON), &e.FilesGenerated); err != nil {
			return nil, fmt.Errorf("unmarshal files: %w", err)
		}
	}

	return &e, nil
}

func scanEvent(row *sql.Row) (*GenerationEvent, error) {
	return scanInto(row)
}

func scanEventRows(rows *sql.Rows) (*GenerationEvent, error) {
	return scanInto(rows)
}
