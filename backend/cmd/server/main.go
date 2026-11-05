// Servidor HTTP de la API.
// Uso: go run ./cmd/server
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/JosePabloOcampo1212/emprendedoras-incremento/backend/internal/config"
	"github.com/JosePabloOcampo1212/emprendedoras-incremento/backend/internal/db"
	"github.com/JosePabloOcampo1212/emprendedoras-incremento/backend/internal/httpapi"
)

func main() {
	cfg := config.Cargar()

	conexion, err := db.Abrir(context.Background(), cfg.URLBD)
	if err != nil {
		log.Fatal(err)
	}
	defer conexion.Close()

	if err := os.MkdirAll(cfg.DirUploads, 0o755); err != nil {
		log.Fatal(err)
	}

	srv := &http.Server{
		Addr:              ":" + cfg.Puerto,
		Handler:           httpapi.Nuevo(conexion, cfg.DirUploads, cfg.CookieSegura).Rutas(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
	}

	log.Printf("API escuchando en http://localhost:%s", cfg.Puerto)
	log.Fatal(srv.ListenAndServe())
}
