// Package seed carga datos SINTÉTICOS de prueba. Ningún nombre, correo
// ni negocio corresponde a personas reales. Los correos usan el dominio
// reservado .test (RFC 2606), que no existe en internet.
package seed

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/JosePabloOcampo1212/emprendedoras-incremento/backend/internal/auth"
)

var Categorias = []string{"Alimentos", "Artesanía", "Cosmética natural", "Textiles"}

type usuarioPrueba struct {
	Nombre, Correo, Rol, Negocio string
}

var Usuarios = []usuarioPrueba{
	{"Administradora de prueba", "admin@asociacion.test", auth.RolAdministradora, ""},
	{"Emprendedora Uno", "emprendedora1@asociacion.test", auth.RolEmprendedora, "Dulces de la Abuela"},
	{"Emprendedora Dos", "emprendedora2@asociacion.test", auth.RolEmprendedora, "Tejidos Montaña Verde"},
}

type productoPrueba struct {
	correo, categoria, nombre, descripcion, estado string
	precio                                         int64
}

var productos = []productoPrueba{
	{"emprendedora1@asociacion.test", "Alimentos", "Cajeta de leche", "Cajeta artesanal hecha con leche de la zona, caja de 12 unidades.", "publicado", 3500},
	{"emprendedora1@asociacion.test", "Alimentos", "Mermelada de mora", "Mermelada de mora silvestre sin preservantes, frasco de 250 g.", "publicado", 2800},
	{"emprendedora2@asociacion.test", "Textiles", "Bufanda tejida a mano", "Bufanda de lana tejida a mano, colores tierra, 1,5 m de largo.", "publicado", 12000},
	{"emprendedora2@asociacion.test", "Artesanía", "Portavasos de fibra natural", "Juego de 6 portavasos hechos con fibra de plátano.", "pendiente", 6000},
}

// Cargar inserta categorías, usuarios y productos de ejemplo.
// Es idempotente: si ya existen, no los duplica.
func Cargar(ctx context.Context, db *sql.DB, clave string) error {
	if len(clave) < 8 {
		return fmt.Errorf("SEED_CLAVE debe tener al menos 8 caracteres")
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, c := range Categorias {
		if _, err := tx.ExecContext(ctx, `INSERT INTO categorias (nombre) VALUES ($1) ON CONFLICT (nombre) DO NOTHING`, c); err != nil {
			return err
		}
	}

	hash, err := auth.HashClave(clave)
	if err != nil {
		return err
	}
	for _, u := range Usuarios {
		_, err := tx.ExecContext(ctx,
			`INSERT INTO usuarios (nombre, correo, clave_hash, rol, negocio)
			 VALUES ($1, $2, $3, $4, NULLIF($5, ''))
			 ON CONFLICT (correo) DO NOTHING`,
			u.Nombre, u.Correo, hash, u.Rol, u.Negocio)
		if err != nil {
			return err
		}
	}

	var hay int
	if err := tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM productos`).Scan(&hay); err != nil {
		return err
	}
	if hay == 0 {
		for _, p := range productos {
			_, err := tx.ExecContext(ctx, `
				INSERT INTO productos (usuario_id, categoria_id, nombre, descripcion, precio, estado)
				VALUES ((SELECT id FROM usuarios WHERE correo = $1),
				        (SELECT id FROM categorias WHERE nombre = $2), $3, $4, $5, $6)`,
				p.correo, p.categoria, p.nombre, p.descripcion, p.precio, p.estado)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}
