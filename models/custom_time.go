package models

import (
	"encoding/json"
	"time"
)

// SafeTime maneja timestamps nulos o mal formateados de PostgreSQL
type SafeTime struct {
	time.Time
	Valid bool
}

func (st *SafeTime) UnmarshalJSON(b []byte) error {
	// Si viene null o cadena vacía
	if string(b) == "null" || string(b) == `""` {
		st.Valid = false
		return nil
	}

	// Intentamos parsear diferentes formatos
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		st.Valid = false
		return nil
	}

	if s == "" || s == "0001-01-01T00:00:00" {
		st.Valid = false
		return nil
	}

	// Intentar con formato ISO8601
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		// Intentar formato sin zona horaria
		t, err = time.Parse("2006-01-02T15:04:05", s)
		if err != nil {
			st.Valid = false
			return nil
		}
	}

	st.Time = t
	st.Valid = true
	return nil
}

func (st SafeTime) MarshalJSON() ([]byte, error) {
	if !st.Valid {
		return []byte(`null`), nil
	}
	return json.Marshal(st.Time.Format(time.RFC3339))
}
