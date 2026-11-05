package httpapi

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/JosePabloOcampo1212/emprendedoras-incremento/backend/internal/auth"
)

const nombreCookie = "sesion"

type claveUsuario struct{}

// requiereSesion deja pasar solo si hay una sesión vigente.
func (s *Servidor) requiereSesion(siguiente http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(nombreCookie)
		if err != nil {
			responderError(w, http.StatusUnauthorized, "Debe iniciar sesión.")
			return
		}
		u, ok := auth.UsuarioDeSesion(r.Context(), s.db, c.Value)
		if !ok {
			responderError(w, http.StatusUnauthorized, "La sesión venció. Inicie sesión de nuevo.")
			return
		}
		ctx := context.WithValue(r.Context(), claveUsuario{}, u)
		siguiente(w, r.WithContext(ctx))
	}
}

// requiereRol deja pasar solo a usuarios con el rol indicado (principio de menor privilegio).
func (s *Servidor) requiereRol(rol string, siguiente http.HandlerFunc) http.HandlerFunc {
	return s.requiereSesion(func(w http.ResponseWriter, r *http.Request) {
		if usuarioDe(r).Rol != rol {
			responderError(w, http.StatusForbidden, "No tiene permiso para esta acción.")
			return
		}
		siguiente(w, r)
	})
}

// cabecerasSeguridad agrega cabeceras HTTP básicas de protección.
func cabecerasSeguridad(siguiente http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("X-Frame-Options", "DENY")
		h.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		h.Set("Content-Security-Policy", "default-src 'none'; img-src 'self'; frame-ancestors 'none'")
		siguiente.ServeHTTP(w, r)
	})
}

// registrarPeticiones escribe método, ruta, estado y duración. No registra cuerpos
// ni cookies, para no dejar datos personales en los logs.
func registrarPeticiones(siguiente http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		inicio := time.Now()
		rw := &escritorConEstado{ResponseWriter: w, estado: http.StatusOK}
		siguiente.ServeHTTP(rw, r)
		log.Printf("%s %s %d %s", r.Method, r.URL.Path, rw.estado, time.Since(inicio).Round(time.Millisecond))
	})
}

type escritorConEstado struct {
	http.ResponseWriter
	estado int
}

func (e *escritorConEstado) WriteHeader(estado int) {
	e.estado = estado
	e.ResponseWriter.WriteHeader(estado)
}
