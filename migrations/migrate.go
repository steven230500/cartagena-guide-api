// Package migrations embebe los .sql en el binario y los aplica al arrancar
// la API — saca el paso manual de `migrate up` de cada deploy.
package migrations

import (
	"embed"
	"errors"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed *.sql
var FS embed.FS

// Run aplica todas las migraciones pendientes contra databaseURL. Idempotente:
// si ya está todo al día, migrate.ErrNoChange no se trata como error.
func Run(databaseURL string) error {
	source, err := iofs.New(FS, ".")
	if err != nil {
		return err
	}

	migrateURL := strings.Replace(databaseURL, "postgres://", "pgx5://", 1)
	m, err := migrate.NewWithSourceInstance("iofs", source, migrateURL)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}
