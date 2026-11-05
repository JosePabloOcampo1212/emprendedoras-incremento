package producto

import (
	"strings"
	"testing"
)

func TestValidarMotivo(t *testing.T) {
	casos := []struct {
		nombre string
		motivo string
		valido bool
	}{
		{"vacío", "", false},
		{"solo espacios", "     ", false},
		{"muy corto", "Falta", false},
		{"adecuado", "Falta indicar el tamaño del producto.", true},
		{"muy largo", strings.Repeat("a", MotivoMax+1), false},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			if got := ValidarMotivo(c.motivo) == ""; got != c.valido {
				t.Fatalf("quería válido=%v para %q", c.valido, c.motivo)
			}
		})
	}
}
