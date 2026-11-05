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

func (s *Servidor) sesionActual(w http.ResponseWriter, r *http.Request) {
	responderJSON(w, http.StatusOK, usuarioDe(r))
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
