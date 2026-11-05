package httpapi

// HU-03 Catálogo público (sin sesión).

import (
	"net/http"
	"strconv"

	"github.com/JosePabloOcampo1212/emprendedoras-incremento/backend/internal/producto"
)

// publico es lo que se muestra a visitantes: sin estado, motivo ni ids internos de usuario.
type publico struct {
	ID          int64  `json:"id"`
	Nombre      string `json:"nombre"`
	Descripcion string `json:"descripcion"`
	Precio      int64  `json:"precio"`
	CategoriaID int64  `json:"categoriaId"`
	Categoria   string `json:"categoria"`
	ImagenURL   string `json:"imagenUrl,omitempty"`
	ImagenAlt   string `json:"imagenAlt,omitempty"`
	Negocio     string `json:"negocio"`
}

func (s *Servidor) catalogo(w http.ResponseWriter, r *http.Request) {
	var categoriaID int64
	if v := r.URL.Query().Get("categoria"); v != "" {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil || id <= 0 {
			responderError(w, http.StatusBadRequest, "Categoría inválida.")
			return
		}
		categoriaID = id
	}

	lista, err := producto.ListarPublicados(r.Context(), s.db, categoriaID)
	if err != nil {
		errorInterno(w, err)
		return
	}

	salida := make([]publico, 0, len(lista))
	for _, p := range lista {
		salida = append(salida, publico{
			ID: p.ID, Nombre: p.Nombre, Descripcion: p.Descripcion, Precio: p.Precio,
			CategoriaID: p.CategoriaID, Categoria: p.Categoria,
			ImagenURL: p.ImagenURL, ImagenAlt: p.ImagenAlt, Negocio: p.Negocio,
		})
	}

	// El catálogo es público y cambia poco: se permite caché corta en el navegador
	// para ahorrar consultas y datos móviles.
	w.Header().Set("Cache-Control", "public, max-age=30")
	responderJSON(w, http.StatusOK, salida)
}
