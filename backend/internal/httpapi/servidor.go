// Package httpapi expone la API REST que consume el frontend.
package httpapi

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/JosePabloOcampo1212/emprendedoras-incremento/backend/internal/auth"
)

type Servidor struct {
	db           *sql.DB
	dirUploads   string
	cookieSegura bool
}

func Nuevo(db *sql.DB, dirUploads string, cookieSegura bool) *Servidor {
	return &Servidor{db: db, dirUploads: dirUploads, cookieSegura: cookieSegura}
}

// Rutas arma el enrutador. Usa los patrones "MÉTODO /ruta" de Go 1.22+.
func (s *Servidor) Rutas() http.Handler {
	mux := http.NewServeMux()

	// Públicas
	mux.HandleFunc("GET /api/salud", s.salud)
	mux.HandleFunc("GET /api/categorias", s.listarCategorias)

	// Sesión
	mux.HandleFunc("POST /api/sesion", s.iniciarSesion)
	mux.HandleFunc("GET /api/sesion", s.requiereSesion(s.sesionActual))
	mux.HandleFunc("DELETE /api/sesion", s.requiereSesion(s.cerrarSesion))

	// HU-01 Registrar producto (rol emprendedora)
	mux.HandleFunc("POST /api/productos", s.requiereRol(auth.RolEmprendedora, s.crearProducto))
	mux.HandleFunc("GET /api/mis-productos", s.requiereRol(auth.RolEmprendedora, s.misProductos))

	// HU-02 Aprobar o rechazar (rol administradora)
	mux.HandleFunc("GET /api/admin/pendientes", s.requiereRol(auth.RolAdministradora, s.listarPendientes))
	mux.HandleFunc("POST /api/admin/productos/{id}/aprobar", s.requiereRol(auth.RolAdministradora, s.aprobarProducto))
	mux.HandleFunc("POST /api/admin/productos/{id}/rechazar", s.requiereRol(auth.RolAdministradora, s.rechazarProducto))

	// HU-03 Catálogo público
	mux.HandleFunc("GET /api/catalogo", s.catalogo)

	// Imágenes subidas (solo lectura)
	mux.Handle("GET /uploads/", sinListado(http.StripPrefix("/uploads/", http.FileServer(http.Dir(s.dirUploads)))))

	return registrarPeticiones(cabecerasSeguridad(mux))
}

func (s *Servidor) salud(w http.ResponseWriter, r *http.Request) {
	responderJSON(w, http.StatusOK, map[string]string{"estado": "ok"})
}

type categoria struct {
	ID     int64  `json:"id"`
	Nombre string `json:"nombre"`
}

func (s *Servidor) listarCategorias(w http.ResponseWriter, r *http.Request) {
	filas, err := s.db.QueryContext(r.Context(), `SELECT id, nombre FROM categorias ORDER BY nombre`)
	if err != nil {
		errorInterno(w, err)
		return
	}
	defer filas.Close()

	lista := []categoria{}
	for filas.Next() {
		var c categoria
		if err := filas.Scan(&c.ID, &c.Nombre); err != nil {
			errorInterno(w, err)
			return
		}
		lista = append(lista, c)
	}
	responderJSON(w, http.StatusOK, lista)
}

// sinListado evita que /uploads/ muestre la lista de todos los archivos.
func sinListado(siguiente http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}
		siguiente.ServeHTTP(w, r)
	})
}

// ---------- utilidades de respuesta ----------

func responderJSON(w http.ResponseWriter, estado int, datos any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(estado)
	json.NewEncoder(w).Encode(datos)
}

func responderError(w http.ResponseWriter, estado int, mensaje string) {
	responderJSON(w, estado, map[string]string{"error": mensaje})
}

// errorInterno registra el detalle en el log y devuelve un mensaje genérico,
// para no exponer información interna al cliente.
func errorInterno(w http.ResponseWriter, err error) {
	log.Printf("error interno: %v", err)
	responderError(w, http.StatusInternalServerError, "Ocurrió un error. Intente de nuevo.")
}

// usuarioDe obtiene el usuario que dejó el middleware en el contexto.
func usuarioDe(r *http.Request) auth.Usuario {
	u, _ := r.Context().Value(claveUsuario{}).(auth.Usuario)
	return u
}
