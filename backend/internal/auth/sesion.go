package auth

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"
)

const DuracionSesion = 8 * time.Hour

// Usuario es lo que el resto de la aplicación necesita saber de quien inició sesión.
type Usuario struct {
	ID      int64  `json:"id"`
	Nombre  string `json:"nombre"`
	Correo  string `json:"correo"`
	Rol     string `json:"rol"`
	Negocio string `json:"negocio,omitempty"`
}

const (
	RolEmprendedora   = "emprendedora"
	RolAdministradora = "administradora"
)

var ErrCredenciales = errors.New("correo o clave incorrectos")

// IniciarSesion valida las credenciales y crea un token de sesión aleatorio.
func IniciarSesion(ctx context.Context, db *sql.DB, correo, clave string) (string, Usuario, error) {
	var u Usuario
	var hash string
	var negocio sql.NullString
	err := db.QueryRowContext(ctx,
		`SELECT id, nombre, correo, rol, negocio, clave_hash FROM usuarios WHERE correo = $1`,
		correo,
	).Scan(&u.ID, &u.Nombre, &u.Correo, &u.Rol, &negocio, &hash)
	if errors.Is(err, sql.ErrNoRows) {
		return "", Usuario{}, ErrCredenciales
	}
	if err != nil {
		return "", Usuario{}, err
	}
	u.Negocio = negocio.String

	ok, err := VerificarClave(clave, hash)
	if err != nil {
		return "", Usuario{}, err
	}
	if !ok {
		return "", Usuario{}, ErrCredenciales
	}

	token, err := nuevoToken()
	if err != nil {
		return "", Usuario{}, err
	}
	_, err = db.ExecContext(ctx,
		`INSERT INTO sesiones (token, usuario_id, expira) VALUES ($1, $2, $3)`,
		token, u.ID, time.Now().Add(DuracionSesion).UTC(),
	)
	if err != nil {
		return "", Usuario{}, err
	}
	return token, u, nil
}

// UsuarioDeSesion devuelve el usuario dueño de un token vigente.
func UsuarioDeSesion(ctx context.Context, db *sql.DB, token string) (Usuario, bool) {
	if token == "" {
		return Usuario{}, false
	}
	var u Usuario
	var negocio sql.NullString
	err := db.QueryRowContext(ctx, `
		SELECT u.id, u.nombre, u.correo, u.rol, u.negocio
		FROM sesiones s JOIN usuarios u ON u.id = s.usuario_id
		WHERE s.token = $1 AND s.expira > now()`,
		token,
	).Scan(&u.ID, &u.Nombre, &u.Correo, &u.Rol, &negocio)
	if err != nil {
		return Usuario{}, false
	}
	u.Negocio = negocio.String
	return u, true
}

func CerrarSesion(ctx context.Context, db *sql.DB, token string) error {
	_, err := db.ExecContext(ctx, `DELETE FROM sesiones WHERE token = $1`, token)
	return err
}

func nuevoToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
