package producto

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

// Vista es un producto tal como lo devuelve la API.
type Vista struct {
	ID            int64     `json:"id"`
	Nombre        string    `json:"nombre"`
	Descripcion   string    `json:"descripcion"`
	Precio        int64     `json:"precio"`
	CategoriaID   int64     `json:"categoriaId"`
	Categoria     string    `json:"categoria"`
	ImagenURL     string    `json:"imagenUrl,omitempty"`
	ImagenAlt     string    `json:"imagenAlt,omitempty"`
	Estado        string    `json:"estado"`
	MotivoRechazo string    `json:"motivoRechazo,omitempty"`
	Negocio       string    `json:"negocio,omitempty"`
	Creado        time.Time `json:"creado"`
}

var ErrCategoriaInexistente = errors.New("la categoría no existe")

// consultaBase devuelve las columnas que lee escanear, en el mismo orden.
const consultaBase = `
	SELECT p.id, p.nombre, p.descripcion, p.precio, p.categoria_id, c.nombre,
	       COALESCE(p.imagen_ruta, ''), COALESCE(p.imagen_alt, ''), p.estado,
	       COALESCE(p.motivo_rechazo, ''), COALESCE(u.negocio, ''), p.creado
	FROM productos p
	JOIN categorias c ON c.id = p.categoria_id
	JOIN usuarios u   ON u.id = p.usuario_id`

func escanear(filas *sql.Rows) ([]Vista, error) {
	defer filas.Close()
	lista := []Vista{}
	for filas.Next() {
		var v Vista
		if err := filas.Scan(&v.ID, &v.Nombre, &v.Descripcion, &v.Precio, &v.CategoriaID, &v.Categoria,
			&v.ImagenURL, &v.ImagenAlt, &v.Estado, &v.MotivoRechazo, &v.Negocio, &v.Creado); err != nil {
			return nil, err
		}
		lista = append(lista, v)
	}
	return lista, filas.Err()
}

// Crear guarda un producto nuevo en estado "pendiente" y registra el cambio
// en el historial, todo en una sola transacción.
func Crear(ctx context.Context, db *sql.DB, usuarioID int64, e Entrada, imagenRuta string) (int64, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var existe bool
	if err := tx.QueryRowContext(ctx,
		`SELECT EXISTS (SELECT 1 FROM categorias WHERE id = $1)`, e.CategoriaID,
	).Scan(&existe); err != nil {
		return 0, err
	}
	if !existe {
		return 0, ErrCategoriaInexistente
	}

	var id int64
	err = tx.QueryRowContext(ctx, `
		INSERT INTO productos (usuario_id, categoria_id, nombre, descripcion, precio, imagen_ruta, imagen_alt, estado)
		VALUES ($1, $2, $3, $4, $5, NULLIF($6, ''), NULLIF($7, ''), $8)
		RETURNING id`,
		usuarioID, e.CategoriaID, e.Nombre, e.Descripcion, e.Precio, imagenRuta, e.ImagenAlt, EstadoPendiente,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	if _, err := tx.ExecContext(ctx, `
		INSERT INTO historial_estados (producto_id, estado_antes, estado_nuevo, usuario_id, nota)
		VALUES ($1, NULL, $2, $3, 'Registro inicial')`,
		id, EstadoPendiente, usuarioID,
	); err != nil {
		return 0, err
	}

	return id, tx.Commit()
}

// ListarDeUsuario devuelve los productos de una emprendedora, del más reciente al más antiguo.
func ListarDeUsuario(ctx context.Context, db *sql.DB, usuarioID int64) ([]Vista, error) {
	filas, err := db.QueryContext(ctx, consultaBase+`
		WHERE p.usuario_id = $1
		ORDER BY p.creado DESC, p.id DESC`, usuarioID)
	if err != nil {
		return nil, err
	}
	return escanear(filas)
}
