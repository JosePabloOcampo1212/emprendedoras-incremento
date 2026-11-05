package auth

import (
	"strings"
	"testing"
)

func TestHashClave_NoGuardaTextoPlano(t *testing.T) {
	hash, err := HashClave("clave-de-prueba-123")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(hash, "clave-de-prueba-123") {
		t.Fatal("el hash contiene la clave en texto plano")
	}
	if !strings.HasPrefix(hash, "pbkdf2-sha256$") {
		t.Fatalf("formato inesperado: %s", hash)
	}
}

func TestHashClave_SalDistintaCadaVez(t *testing.T) {
	a, _ := HashClave("misma-clave")
	b, _ := HashClave("misma-clave")
	if a == b {
		t.Fatal("dos hashes de la misma clave no deberían ser iguales")
	}
}

func TestVerificarClave(t *testing.T) {
	hash, _ := HashClave("correcta")

	casos := []struct {
		nombre string
		clave  string
		quiere bool
	}{
		{"clave correcta", "correcta", true},
		{"clave incorrecta", "incorrecta", false},
		{"clave vacía", "", false},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			ok, err := VerificarClave(c.clave, hash)
			if err != nil {
				t.Fatal(err)
			}
			if ok != c.quiere {
				t.Fatalf("quería %v, obtuvo %v", c.quiere, ok)
			}
		})
	}
}

func TestVerificarClave_HashMalFormado(t *testing.T) {
	if _, err := VerificarClave("x", "texto-cualquiera"); err == nil {
		t.Fatal("esperaba error con un hash mal formado")
	}
}
