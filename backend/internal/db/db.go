// Package db abre la conexión a PostgreSQL y aplica el esquema.
package db

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // registra el driver "pgx" para database/sql
)

//go:embed schema.sql
var esquema string

// Abrir conecta, verifica que la base responda y aplica el esquema.
func Abrir(ctx context.Context, url string) (*sql.DB, error) {
	conexion, err := sql.Open("pgx", url)
	if err != nil {
		return nil, err
	}
	conexion.SetMaxOpenConns(10)
	conexion.SetMaxIdleConns(5)
	conexion.SetConnMaxLifetime(30 * time.Minute)

	ctxPing, cancelar := context.WithTimeout(ctx, 5*time.Second)
	defer cancelar()
	if err := conexion.PingContext(ctxPing); err != nil {
		conexion.Close()
		return nil, fmt.Errorf("no se pudo conectar a PostgreSQL (¿está corriendo `docker compose up -d`?): %w", err)
	}

	if _, err := conexion.ExecContext(ctx, esquema); err != nil {
		conexion.Close()
		return nil, fmt.Errorf("aplicar esquema: %w", err)
	}
	return conexion, nil
}
