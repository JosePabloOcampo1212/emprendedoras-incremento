// Package config carga la configuración desde variables de entorno.
// Si existe un archivo .env (en la carpeta actual o en alguna carpeta superior),
// sus valores se cargan primero. Las variables ya definidas en el sistema
// tienen prioridad sobre el archivo.
package config

import (
	"bufio"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Puerto       string // puerto HTTP del backend
	URLBD        string // conexión a PostgreSQL
	DirUploads   string // carpeta donde se guardan las imágenes subidas
	SeedClave    string // clave de las cuentas de prueba (solo la usa el seed)
	CookieSegura bool   // true en producción (HTTPS)
}

func Cargar() Config {
	cargarEnvCercano()

	return Config{
		Puerto:       valor("PUERTO", "8080"),
		URLBD:        urlBaseDatos(),
		DirUploads:   valor("DIR_UPLOADS", "data/uploads"),
		SeedClave:    os.Getenv("SEED_CLAVE"),
		CookieSegura: valor("COOKIE_SEGURA", "false") == "true",
	}
}

// urlBaseDatos usa DATABASE_URL si existe; si no, la arma con las mismas
// variables POSTGRES_* que usa docker-compose, para no repetir la clave.
func urlBaseDatos() string {
	if u := os.Getenv("DATABASE_URL"); u != "" {
		return u
	}
	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(valor("POSTGRES_USER", "emprendedoras"), os.Getenv("POSTGRES_PASSWORD")),
		Host:     valor("DB_HOST", "localhost") + ":" + valor("DB_PUERTO", "5432"),
		Path:     valor("POSTGRES_DB", "emprendedoras"),
		RawQuery: "sslmode=" + valor("DB_SSLMODE", "disable"),
	}
	return u.String()
}

func valor(clave, porDefecto string) string {
	if v := strings.TrimSpace(os.Getenv(clave)); v != "" {
		return v
	}
	return porDefecto
}

// cargarEnvCercano busca un .env en la carpeta actual y hasta 4 niveles arriba.
// Así funciona igual desde backend/ que desde las carpetas de las pruebas.
func cargarEnvCercano() {
	dir, err := os.Getwd()
	if err != nil {
		return
	}
	for range 5 {
		ruta := filepath.Join(dir, ".env")
		if _, err := os.Stat(ruta); err == nil {
			cargarArchivoEnv(ruta)
			return
		}
		dir = filepath.Dir(dir)
	}
}

// cargarArchivoEnv lee líneas CLAVE=valor. Ignora comentarios y líneas vacías.
func cargarArchivoEnv(ruta string) {
	f, err := os.Open(ruta)
	if err != nil {
		return
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		linea := strings.TrimSpace(sc.Text())
		if linea == "" || strings.HasPrefix(linea, "#") {
			continue
		}
		clave, val, ok := strings.Cut(linea, "=")
		if !ok {
			continue
		}
		clave = strings.TrimSpace(clave)
		val = strings.Trim(strings.TrimSpace(val), `"'`)
		if _, existe := os.LookupEnv(clave); !existe {
			os.Setenv(clave, val)
		}
	}
}
