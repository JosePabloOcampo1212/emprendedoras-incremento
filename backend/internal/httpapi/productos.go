package httpapi

// HU-01 Registrar producto (rol emprendedora).

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/JosePabloOcampo1212/emprendedoras-incremento/backend/internal/producto"
)

// Tamaño máximo de toda la petición: imagen de 2 MB + campos de texto.
const maxCuerpoProducto = producto.ImagenMaxBytes + 512*1024

func (s *Servidor) crearProducto(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxCuerpoProducto)
	if err := r.ParseMultipartForm(maxCuerpoProducto); err != nil {
		var demasiado *http.MaxBytesError
		if errors.As(err, &demasiado) {
			responderErrores(w, producto.Errores{"imagen": "La imagen no puede pesar más de 2 MB."})
			return
		}
		responderError(w, http.StatusBadRequest, "Datos inválidos.")
		return
	}
	defer r.MultipartForm.RemoveAll()

	errs := producto.Errores{}

	precio, err := strconv.ParseInt(strings.TrimSpace(r.FormValue("precio")), 10, 64)
	if err != nil {
		precio = 0 // la validación lo reporta como precio inválido
	}
	categoriaID, _ := strconv.ParseInt(r.FormValue("categoriaId"), 10, 64)

	entrada := producto.Entrada{
		Nombre:      r.FormValue("nombre"),
		Descripcion: r.FormValue("descripcion"),
		Precio:      precio,
		CategoriaID: categoriaID,
		ImagenAlt:   r.FormValue("imagenAlt"),
	}
	entrada.Normalizar()

	// Imagen opcional: se revisa el tipo real leyendo sus primeros bytes.
	var contenido []byte
	var extension string
	archivo, cabecera, err := r.FormFile("imagen")
	tieneImagen := err == nil
	if tieneImagen {
		defer archivo.Close()
		contenido, err = io.ReadAll(archivo)
		if err != nil {
			errorInterno(w, err)
			return
		}
		tipo := http.DetectContentType(contenido)
		if msg := producto.ValidarImagen(tipo, cabecera.Size); msg != "" {
			errs["imagen"] = msg
		}
		extension = map[string]string{"image/jpeg": ".jpg", "image/png": ".png"}[tipo]
	}

	for campo, msg := range producto.Validar(entrada, tieneImagen) {
		errs[campo] = msg
	}
	if !errs.Vacio() {
		responderErrores(w, errs)
		return
	}

	// Nombre de archivo aleatorio: no se usa el nombre original (evita rutas maliciosas
	// y no expone datos de la persona que subió la imagen).
	var rutaPublica, rutaDisco string
	if tieneImagen {
		nombre := nombreAleatorio() + extension
		rutaDisco = filepath.Join(s.dirUploads, nombre)
		if err := os.MkdirAll(s.dirUploads, 0o755); err != nil {
			errorInterno(w, err)
			return
		}
		if err := os.WriteFile(rutaDisco, contenido, 0o644); err != nil {
			errorInterno(w, err)
			return
		}
		rutaPublica = "/uploads/" + nombre
	}

	id, err := producto.Crear(r.Context(), s.db, usuarioDe(r).ID, entrada, rutaPublica)
	if err != nil {
		if rutaDisco != "" {
			os.Remove(rutaDisco) // no dejar imágenes huérfanas
		}
		if errors.Is(err, producto.ErrCategoriaInexistente) {
			responderErrores(w, producto.Errores{"categoriaId": "La categoría seleccionada no existe."})
			return
		}
		errorInterno(w, err)
		return
	}

	responderJSON(w, http.StatusCreated, map[string]any{"id": id, "estado": producto.EstadoPendiente})
}

func (s *Servidor) misProductos(w http.ResponseWriter, r *http.Request) {
	lista, err := producto.ListarDeUsuario(r.Context(), s.db, usuarioDe(r).ID)
	if err != nil {
		errorInterno(w, err)
		return
	}
	responderJSON(w, http.StatusOK, lista)
}

// responderErrores devuelve 422 con un mensaje general y los errores por campo.
func responderErrores(w http.ResponseWriter, errs producto.Errores) {
	responderJSON(w, http.StatusUnprocessableEntity, map[string]any{
		"error":   "Revise los campos marcados.",
		"errores": errs,
	})
}

func nombreAleatorio() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
