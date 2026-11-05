package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/JosePabloOcampo1212/emprendedoras-incremento/backend/internal/auth"
)

type credenciales struct {
	Correo string `json:"correo"`
	Clave  string `json:"clave"`
}

func (s *Servidor) iniciarSesion(w http.ResponseWriter, r *http.Request) {
	var c credenciales
	r.Body = http.MaxBytesReader(w, r.Body, 4096)
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		responderError(w, http.StatusBadRequest, "Datos inválidos.")
		return
	}
	c.Correo = strings.ToLower(strings.TrimSpace(c.Correo))
	if c.Correo == "" || c.Clave == "" {
		responderError(w, http.StatusBadRequest, "Escriba su correo y su clave.")
		return
	}

	token, u, err := auth.IniciarSesion(r.Context(), s.db, c.Correo, c.Clave)
	if errors.Is(err, auth.ErrCredenciales) {
		// Mismo mensaje si el correo no existe o la clave falla:
		// así no se revela qué correos están registrados.
		responderError(w, http.StatusUnauthorized, "Correo o clave incorrectos.")
		return
	}
	if err != nil {
		errorInterno(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     nombreCookie,
		Value:    token,
		Path:     "/",
		HttpOnly: true, // JavaScript no puede leerla
		Secure:   s.cookieSegura,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(auth.DuracionSesion.Seconds()),
	})
	responderJSON(w, http.StatusOK, u)
}

// sesionActual responde 200 con el usuario, o 204 si no hay sesión.
// No usa 401 porque la página lo consulta en cada visita, y un visitante
// sin sesión es un caso normal, no un error (DEF-05).
func (s *Servidor) sesionActual(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie(nombreCookie)
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	u, ok := auth.UsuarioDeSesion(r.Context(), s.db, c.Value)
	if !ok {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	responderJSON(w, http.StatusOK, u)
}

func (s *Servidor) cerrarSesion(w http.ResponseWriter, r *http.Request) {
	if c, err := r.Cookie(nombreCookie); err == nil {
		if err := auth.CerrarSesion(r.Context(), s.db, c.Value); err != nil {
			errorInterno(w, err)
			return
		}
	}
	http.SetCookie(w, &http.Cookie{
		Name: nombreCookie, Value: "", Path: "/", MaxAge: -1,
		HttpOnly: true, Secure: s.cookieSegura, SameSite: http.SameSiteStrictMode,
	})
	w.WriteHeader(http.StatusNoContent)
}
