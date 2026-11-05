package httpapi

// HU-02 Aprobar o rechazar productos (rol administradora).

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/JosePabloOcampo1212/emprendedoras-incremento/backend/internal/producto"
)

func (s *Servidor) listarPendientes(w http.ResponseWriter, r *http.Request) {
	lista, err := producto.ListarPendientes(r.Context(), s.db)
	if err != nil {
		errorInterno(w, err)
		return
	}
	responderJSON(w, http.StatusOK, lista)
}

func (s *Servidor) aprobarProducto(w http.ResponseWriter, r *http.Request) {
	s.cambiarEstado(w, r, producto.EstadoPublicado, "Aprobado")
}

func (s *Servidor) rechazarProducto(w http.ResponseWriter, r *http.Request) {
	var cuerpo struct {
		Motivo string `json:"motivo"`
	}
	r.Body = http.MaxBytesReader(w, r.Body, 4096)
	if err := json.NewDecoder(r.Body).Decode(&cuerpo); err != nil {
		responderError(w, http.StatusBadRequest, "Datos inválidos.")
		return
	}
	if msg := producto.ValidarMotivo(cuerpo.Motivo); msg != "" {
		responderErrores(w, producto.Errores{"motivo": msg})
		return
	}
	s.cambiarEstado(w, r, producto.EstadoBorrador, cuerpo.Motivo)
}

func (s *Servidor) cambiarEstado(w http.ResponseWriter, r *http.Request, hacia, nota string) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		responderError(w, http.StatusNotFound, "Producto no encontrado.")
		return
	}

	err = producto.CambiarEstado(r.Context(), s.db, id, hacia, usuarioDe(r).ID, nota)
	switch {
	case errors.Is(err, producto.ErrNoEncontrado):
		responderError(w, http.StatusNotFound, "Producto no encontrado.")
	case errors.Is(err, producto.ErrTransicionInvalida):
		responderError(w, http.StatusConflict, "Este producto ya fue revisado. Actualice la lista.")
	case err != nil:
		errorInterno(w, err)
	default:
		responderJSON(w, http.StatusOK, map[string]any{"id": id, "estado": hacia})
	}
}
