// Package producto define las reglas del contenido "producto":
// validación de datos (HU-01) y ciclo de vida de publicación (HU-02, HU-04).
package producto

import (
	"strings"
	"unicode/utf8"
)

// Estados del ciclo de vida del contenido.
const (
	EstadoBorrador  = "borrador"
	EstadoPendiente = "pendiente"
	EstadoPublicado = "publicado"
	EstadoArchivado = "archivado"
)

// Límites de validación. Si cambian, actualizar también los criterios de aceptación.
const (
	NombreMin      = 3
	NombreMax      = 80
	DescripcionMin = 10
	DescripcionMax = 500
	PrecioMax      = 10_000_000 // colones
	AltMax         = 150
	ImagenMaxBytes = 2 * 1024 * 1024 // 2 MB
)

// Entrada son los datos que envía el formulario de registro.
type Entrada struct {
	Nombre      string `json:"nombre"`
	Descripcion string `json:"descripcion"`
	Precio      int64  `json:"precio"`
	CategoriaID int64  `json:"categoriaId"`
	ImagenAlt   string `json:"imagenAlt"`
}

// Errores agrupa los mensajes por campo, para mostrarlos junto a cada campo.
type Errores map[string]string

func (e Errores) Vacio() bool { return len(e) == 0 }

// Normalizar quita espacios sobrantes antes de validar y guardar.
func (e *Entrada) Normalizar() {
	e.Nombre = strings.TrimSpace(e.Nombre)
	e.Descripcion = strings.TrimSpace(e.Descripcion)
	e.ImagenAlt = strings.TrimSpace(e.ImagenAlt)
}

// Validar revisa los campos del producto. tieneImagen indica si se adjuntó imagen,
// porque en ese caso el texto alternativo pasa a ser obligatorio (accesibilidad).
func Validar(e Entrada, tieneImagen bool) Errores {
	errs := Errores{}

	switch n := utf8.RuneCountInString(e.Nombre); {
	case n == 0:
		errs["nombre"] = "El nombre es obligatorio."
	case n < NombreMin:
		errs["nombre"] = "El nombre debe tener al menos 3 caracteres."
	case n > NombreMax:
		errs["nombre"] = "El nombre no puede pasar de 80 caracteres."
	}

	switch n := utf8.RuneCountInString(e.Descripcion); {
	case n == 0:
		errs["descripcion"] = "La descripción es obligatoria."
	case n < DescripcionMin:
		errs["descripcion"] = "La descripción debe tener al menos 10 caracteres."
	case n > DescripcionMax:
		errs["descripcion"] = "La descripción no puede pasar de 500 caracteres."
	}

	switch {
	case e.Precio <= 0:
		errs["precio"] = "El precio debe ser un número mayor que cero."
	case e.Precio > PrecioMax:
		errs["precio"] = "El precio no puede ser mayor a ₡10 000 000."
	}

	if e.CategoriaID <= 0 {
		errs["categoriaId"] = "Seleccione una categoría."
	}

	if tieneImagen {
		switch n := utf8.RuneCountInString(e.ImagenAlt); {
		case n == 0:
			errs["imagenAlt"] = "Describa la imagen para personas que usan lector de pantalla."
		case n > AltMax:
			errs["imagenAlt"] = "La descripción de la imagen no puede pasar de 150 caracteres."
		}
	}

	return errs
}

// ValidarImagen revisa el tipo real del archivo (detectado por contenido, no por extensión)
// y su tamaño. Devuelve "" si es válida.
func ValidarImagen(tipoDetectado string, tamano int64) string {
	if tipoDetectado != "image/jpeg" && tipoDetectado != "image/png" {
		return "La imagen debe ser JPG o PNG."
	}
	if tamano > ImagenMaxBytes {
		return "La imagen no puede pesar más de 2 MB."
	}
	return ""
}

// transiciones permitidas en el ciclo de publicación.
var transiciones = map[string][]string{
	EstadoBorrador:  {EstadoPendiente},                 // la emprendedora corrige y reenvía
	EstadoPendiente: {EstadoPublicado, EstadoBorrador}, // aprobar o rechazar (HU-02)
	EstadoPublicado: {EstadoArchivado},                 // archivar (HU-04)
	EstadoArchivado: {EstadoPendiente},                 // restaurar vuelve a revisión (HU-04)
}

// PuedeCambiar indica si un producto puede pasar de un estado a otro.
func PuedeCambiar(desde, hacia string) bool {
	for _, permitido := range transiciones[desde] {
		if permitido == hacia {
			return true
		}
	}
	return false
}
