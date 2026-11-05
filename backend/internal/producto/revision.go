package producto

// HU-02 Revisión de productos por la administradora.

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"unicode/utf8"
)

var (
	ErrNoEncontrado       = errors.New("producto no encontrado")
	ErrTransicionInvalida = errors.New("el producto no se puede pasar a ese estado")
)

const (
	MotivoMin = 10
	MotivoMax = 300
)

// ValidarMotivo revisa el motivo de rechazo. Devuelve "" si es válido.
func ValidarMotivo(motivo string) string {
	n := utf8.RuneCountInString(strings.TrimSpace(motivo))
	switch {
	case n == 0:
		return "Escriba el motivo para que la emprendedora sepa qué corregir."
	case n < MotivoMin:
		return "El motivo debe tener al menos 10 caracteres."
	case n > MotivoMax:
		return "El motivo no puede pasar de 300 caracteres."
	}
	return ""
}

// ListarPendientes devuelve los productos en revisión, los más antiguos primero
// (así se atienden en orden de llegada).
func ListarPendientes(ctx context.Context, db *sql.DB) ([]Vista, error) {
	filas, err := db.QueryContext(ctx, consultaBase+`
		WHERE p.estado = $1
		ORDER BY p.creado ASC, p.id ASC`, EstadoPendiente)
	if err != nil {
		return nil, err
	}
	return escanear(filas)
}

// CambiarEstado mueve un producto a otro estado si la transición está permitida,
// y deja registro en el historial. Bloquea la fila (FOR UPDATE) para que dos
// revisiones al mismo tiempo no se pisen.
func CambiarEstado(ctx context.Context, db *sql.DB, productoID int64, hacia string, usuarioID int64, nota string) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var actual string
	err = tx.QueryRowContext(ctx,
		`SELECT estado FROM productos WHERE id = $1 FOR UPDATE`, productoID,
	).Scan(&actual)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNoEncontrado
	}
	if err != nil {
		return err
	}

	if !PuedeCambiar(actual, hacia) {
		return ErrTransicionInvalida
	}

	// El motivo solo se guarda cuando el producto se devuelve a borrador.
	motivo := sql.NullString{}
	if hacia == EstadoBorrador {
		motivo = sql.NullString{String: strings.TrimSpace(nota), Valid: true}
	}

	if _, err := tx.ExecContext(ctx, `
		UPDATE productos SET estado = $1, motivo_rechazo = $2, actualizado = now()
		WHERE id = $3`, hacia, motivo, productoID,
	); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, `
		INSERT INTO historial_estados (producto_id, estado_antes, estado_nuevo, usuario_id, nota)
		VALUES ($1, $2, $3, $4, NULLIF($5, ''))`,
		productoID, actual, hacia, usuarioID, strings.TrimSpace(nota),
	); err != nil {
		return err
	}

	return tx.Commit()
}
