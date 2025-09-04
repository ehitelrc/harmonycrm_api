// models/date.go
package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

const layoutISO = "2006-01-02"

// Date almacena sólo fecha (YYYY-MM-DD)
type Date struct {
	time.Time
}

// ---------- JSON ----------
func (d *Date) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "" || s == "null" {
		d.Time = time.Time{}
		return nil
	}
	t, err := time.Parse(layoutISO, s)
	if err != nil {
		return err
	}
	d.Time = t
	return nil
}

func (d Date) MarshalJSON() ([]byte, error) {
	if d.Time.IsZero() {
		return []byte(`null`), nil
	}
	return json.Marshal(d.Time.Format(layoutISO))
}

// ---------- DB (gorm / database/sql) ----------
func (d Date) Value() (driver.Value, error) {
	if d.Time.IsZero() {
		return nil, nil // guarda NULL en la columna
	}
	// Devolver time.Time permite al driver de Postgres/pgx mapearlo a DATE
	return d.Time, nil
}

func (d *Date) Scan(value any) error {
	if value == nil {
		d.Time = time.Time{}
		return nil
	}
	switch v := value.(type) {
	case time.Time:
		// Postgres suele entregar time.Time para DATE
		d.Time = v
		return nil
	case []byte:
		t, err := time.Parse(layoutISO, string(v))
		if err != nil {
			return fmt.Errorf("Date.Scan parse []byte: %w", err)
		}
		d.Time = t
		return nil
	case string:
		t, err := time.Parse(layoutISO, v)
		if err != nil {
			return fmt.Errorf("Date.Scan parse string: %w", err)
		}
		d.Time = t
		return nil
	default:
		return fmt.Errorf("Date.Scan: tipo no soportado %T", value)
	}
}
