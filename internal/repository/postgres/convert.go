// Package postgres implementa las interfaces de internal/repository usando
// internal/db (sqlc). Todo lo pgx-específico (pgtype.UUID, pgtype.Date, el
// fix de slices nil->[]) vive acá — ninguna otra capa importa pgtype.
package postgres

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func mustUUID(s string) (pgtype.UUID, error) {
	var u pgtype.UUID
	err := u.Scan(s)
	return u, err
}

func uuidToPtr(u pgtype.UUID) *string {
	if !u.Valid {
		return nil
	}
	s := u.String()
	return &s
}

func parseDate(s string) (pgtype.Date, error) {
	var d pgtype.Date
	if s == "" {
		return d, nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return d, err
	}
	d.Time = t
	d.Valid = true
	return d, nil
}

func dateToString(d pgtype.Date) string {
	if !d.Valid {
		return ""
	}
	return d.Time.Format("2006-01-02")
}

func dateToPtr(d pgtype.Date) *string {
	if !d.Valid {
		return nil
	}
	s := d.Time.Format("2006-01-02")
	return &s
}

func normalizeTags(tags []string) []string {
	if tags == nil {
		return []string{}
	}
	return tags
}
