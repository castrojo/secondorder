package db

import (
	"database/sql"
	"math"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	retryMaxDuration = 5 * time.Minute
	retryBaseDelay   = 100 * time.Millisecond
	retryMaxDelay    = 30 * time.Second
)

func isSQLiteBusy(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "database is locked") ||
		strings.Contains(msg, "SQLITE_BUSY")
}

func retryOnBusy(op string, fn func() error) error {
	var elapsed time.Duration
	for attempt := 0; ; attempt++ {
		err := fn()
		if err == nil || !isSQLiteBusy(err) {
			return err
		}

		delay := time.Duration(float64(retryBaseDelay) * math.Pow(2, float64(attempt)))
		if delay > retryMaxDelay {
			delay = retryMaxDelay
		}
		if elapsed+delay > retryMaxDuration {
			log.WithFields(log.Fields{
				"op":      op,
				"elapsed": elapsed.Round(time.Millisecond),
				"attempt": attempt + 1,
			}).Error("db: SQLITE_BUSY retry exhausted, giving up")
			return err
		}

		log.WithFields(log.Fields{
			"op":      op,
			"delay":   delay.Round(time.Millisecond),
			"attempt": attempt + 1,
		}).Warn("db: SQLITE_BUSY, retrying with backoff")

		time.Sleep(delay)
		elapsed += delay
	}
}

func (d *DB) Exec(query string, args ...any) (sql.Result, error) {
	var result sql.Result
	err := retryOnBusy("Exec", func() error {
		var e error
		result, e = d.DB.Exec(query, args...)
		return e
	})
	return result, err
}

func (d *DB) QueryRow(query string, args ...any) *Row {
	return &Row{db: d.DB, query: query, args: args}
}

// Row wraps sql.Row to add SQLITE_BUSY retry on Scan.
type Row struct {
	db    *sql.DB
	query string
	args  []any
}

func (r *Row) Scan(dest ...any) error {
	return retryOnBusy("QueryRow", func() error {
		return r.db.QueryRow(r.query, r.args...).Scan(dest...)
	})
}

func (d *DB) Query(query string, args ...any) (*sql.Rows, error) {
	var rows *sql.Rows
	err := retryOnBusy("Query", func() error {
		var e error
		rows, e = d.DB.Query(query, args...)
		return e
	})
	return rows, err
}

func (d *DB) Begin() (*sql.Tx, error) {
	var tx *sql.Tx
	err := retryOnBusy("Begin", func() error {
		var e error
		tx, e = d.DB.Begin()
		return e
	})
	return tx, err
}
