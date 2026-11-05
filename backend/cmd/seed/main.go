// Carga datos sintéticos de prueba en la base de datos.
// Uso: go run ./cmd/seed
package main

import (
	"context"
	"log"

	"github.com/JosePabloOcampo1212/emprendedoras-incremento/backend/internal/config"
	"github.com/JosePabloOcampo1212/emprendedoras-incremento/backend/internal/db"
	"github.com/JosePabloOcampo1212/emprendedoras-incremento/backend/internal/seed"
)

func main() {
	cfg := config.Cargar()
	if cfg.SeedClave == "" {
		log.Fatal("Falta SEED_CLAVE. Copie .env.example a .env y defina una clave.")
	}

	conexion, err := db.Abrir(context.Background(), cfg.URLBD)
	if err != nil {
		log.Fatal(err)
	}
	defer conexion.Close()

	if err := seed.Cargar(context.Background(), conexion, cfg.SeedClave); err != nil {
		log.Fatal(err)
	}

	log.Print("Datos de prueba cargados en PostgreSQL")
	for _, u := range seed.Usuarios {
		log.Printf("  %-32s rol: %s", u.Correo, u.Rol)
	}
	log.Print("La clave de todas las cuentas es el valor de SEED_CLAVE en su .env")
}
