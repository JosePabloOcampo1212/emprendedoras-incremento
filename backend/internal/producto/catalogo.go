package producto

// HU-03 Catálogo público.

import (
	"context"
	"database/sql"
)

// LimiteCatalogo evita respuestas enormes; con más productos habría que paginar.
const LimiteCatalogo = 200

// ListarPublicados devuelve solo productos publicados. Si categoriaID es 0,
// devuelve todas las categorías.
func ListarPublicados(ctx context.Context, db *sql.DB, categoriaID int64) ([]Vista, error) {
	filas, err := db.QueryContext(ctx, consultaBase+`
		WHERE p.estado = $1
		  AND ($2 = 0 OR p.categoria_id = $2)
		ORDER BY p.actualizado DESC, p.id DESC
		LIMIT $3`, EstadoPublicado, categoriaID, LimiteCatalogo)
	if err != nil {
		return nil, err
	}
	return escanear(filas)
}
