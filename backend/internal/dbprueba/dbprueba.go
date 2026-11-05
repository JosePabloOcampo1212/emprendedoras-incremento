// Package dbprueba prepara una base de datos aislada para las pruebas de integración.
//
// Cada prueba trabaja en un esquema propio de PostgreSQL (prueba_xxxx) que se
// borra al terminar, así las pruebas no tocan los datos de desarrollo ni se
// pisan entre sí. Usa TEST_DATABASE_URL si existe; si no, la misma conexión del .env.
package dbprueba

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/JosePabloOcampo1212/emprendedoras-incremento/backend/internal/config"
	"github.com/JosePabloOcampo1212/emprendedoras-incremento/backend/internal/db"
)

// Nueva devuelve una conexión con el esquema de la aplicación ya creado.
// Si PostgreSQL no está disponible, la prueba se marca como omitida (SKIP).
func Nueva(t *testing.T) *sql.DB {
	t.Helper()
	ctx := context.Background()

	base := os.Getenv("TEST_DATABASE_URL")
	if base == "" {
		base = config.Cargar().URLBD
	}

	admin, err := sql.Open("pgx", base)
	if err != nil {
		t.Fatal(err)
	}
	ctxPing, cancelar := context.WithTimeout(ctx, 3*time.Second)
	defer cancelar()
	if err := admin.PingContext(ctxPing); err != nil {
		admin.Close()
		t.Skipf("PostgreSQL no disponible, se omite la prueba de integración (levántelo con `docker compose up -d`): %v", err)
	}

	esquema := "prueba_" + sufijo()
	if _, err := admin.ExecContext(ctx, "CREATE SCHEMA "+esquema); err != nil {
		admin.Close()
		t.Fatal(err)
	}

	u, err := url.Parse(base)
	if err != nil {
		t.Fatal(err)
	}
	q := u.Query()
	q.Set("search_path", esquema)
	u.RawQuery = q.Encode()

	conexion, err := db.Abrir(ctx, u.String())
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		conexion.Close()
		admin.ExecContext(context.Background(), "DROP SCHEMA "+esquema+" CASCADE")
		admin.Close()
	})
	return conexion
}

func sufijo() string {
	b := make([]byte, 6)
	rand.Read(b)
	return hex.EncodeToString(b)
}
