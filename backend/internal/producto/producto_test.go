package producto

import (
	"strings"
	"testing"
)

func entradaValida() Entrada {
	return Entrada{
		Nombre:      "Café de altura",
		Descripcion: "Café tostado artesanalmente en Coto Brus.",
		Precio:      4500,
		CategoriaID: 1,
		ImagenAlt:   "Bolsa de café de 340 gramos",
	}
}

func TestValidar_EntradaValida(t *testing.T) {
	if errs := Validar(entradaValida(), true); !errs.Vacio() {
		t.Fatalf("no esperaba errores: %v", errs)
	}
}

// Cada caso cambia un campo y revisa que el error aparezca en ese campo.
func TestValidar_Campos(t *testing.T) {
	casos := []struct {
		nombre      string
		modificar   func(*Entrada)
		tieneImagen bool
		campo       string
	}{
		{"nombre vacío", func(e *Entrada) { e.Nombre = "" }, false, "nombre"},
		{"nombre muy corto", func(e *Entrada) { e.Nombre = "Té" }, false, "nombre"},
		{"nombre muy largo", func(e *Entrada) { e.Nombre = strings.Repeat("a", 81) }, false, "nombre"},
		{"descripción vacía", func(e *Entrada) { e.Descripcion = "" }, false, "descripcion"},
		{"descripción muy corta", func(e *Entrada) { e.Descripcion = "Rico" }, false, "descripcion"},
		{"descripción muy larga", func(e *Entrada) { e.Descripcion = strings.Repeat("a", 501) }, false, "descripcion"},
		{"precio cero", func(e *Entrada) { e.Precio = 0 }, false, "precio"},
		{"precio negativo", func(e *Entrada) { e.Precio = -100 }, false, "precio"},
		{"precio excesivo", func(e *Entrada) { e.Precio = PrecioMax + 1 }, false, "precio"},
		{"sin categoría", func(e *Entrada) { e.CategoriaID = 0 }, false, "categoriaId"},
		{"imagen sin texto alternativo", func(e *Entrada) { e.ImagenAlt = "" }, true, "imagenAlt"},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			e := entradaValida()
			c.modificar(&e)
			errs := Validar(e, c.tieneImagen)
			if _, ok := errs[c.campo]; !ok {
				t.Fatalf("esperaba error en %q, obtuvo %v", c.campo, errs)
			}
			if len(errs) != 1 {
				t.Fatalf("esperaba solo un error, obtuvo %v", errs)
			}
		})
	}
}

func TestValidar_SinImagenNoPideAlt(t *testing.T) {
	e := entradaValida()
	e.ImagenAlt = ""
	if errs := Validar(e, false); !errs.Vacio() {
		t.Fatalf("sin imagen el texto alternativo no es obligatorio: %v", errs)
	}
}

func TestValidar_CuentaCaracteresNoBytes(t *testing.T) {
	e := entradaValida()
	e.Nombre = "Ñañ" // 3 caracteres, 5 bytes
	if errs := Validar(e, false); !errs.Vacio() {
		t.Fatalf("un nombre de 3 caracteres con tildes debería ser válido: %v", errs)
	}
}

func TestNormalizar(t *testing.T) {
	e := Entrada{Nombre: "  Jabón  ", Descripcion: " texto ", ImagenAlt: " alt "}
	e.Normalizar()
	if e.Nombre != "Jabón" || e.Descripcion != "texto" || e.ImagenAlt != "alt" {
		t.Fatalf("no quitó los espacios: %+v", e)
	}
}

func TestValidarImagen(t *testing.T) {
	casos := []struct {
		nombre string
		tipo   string
		tamano int64
		valida bool
	}{
		{"jpg pequeño", "image/jpeg", 500_000, true},
		{"png en el límite", "image/png", ImagenMaxBytes, true},
		{"png muy pesado", "image/png", ImagenMaxBytes + 1, false},
		{"gif no permitido", "image/gif", 1000, false},
		{"pdf no permitido", "application/pdf", 1000, false},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			msg := ValidarImagen(c.tipo, c.tamano)
			if (msg == "") != c.valida {
				t.Fatalf("quería válida=%v, mensaje=%q", c.valida, msg)
			}
		})
	}
}

func TestPuedeCambiar(t *testing.T) {
	casos := []struct {
		desde, hacia string
		quiere       bool
	}{
		{EstadoPendiente, EstadoPublicado, true},  // aprobar
		{EstadoPendiente, EstadoBorrador, true},   // rechazar
		{EstadoPublicado, EstadoArchivado, true},  // archivar
		{EstadoArchivado, EstadoPendiente, true},  // restaurar
		{EstadoBorrador, EstadoPendiente, true},   // reenviar
		{EstadoBorrador, EstadoPublicado, false},  // no se publica sin revisión
		{EstadoArchivado, EstadoPublicado, false}, // restaurar pasa por revisión
		{EstadoPublicado, EstadoPendiente, false},
	}
	for _, c := range casos {
		if got := PuedeCambiar(c.desde, c.hacia); got != c.quiere {
			t.Errorf("%s -> %s: quería %v, obtuvo %v", c.desde, c.hacia, c.quiere, got)
		}
	}
}
