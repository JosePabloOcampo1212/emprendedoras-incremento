// Package auth maneja las claves y las sesiones.
package auth

import (
	"crypto/pbkdf2"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Iteraciones de PBKDF2-SHA256 (recomendación OWASP 2023: 600 000).
const iteraciones = 600_000

const largoSal = 16
const largoHash = 32

// HashClave devuelve "pbkdf2-sha256$iteraciones$sal$hash" en base64.
// Nunca se guarda la clave en texto plano.
func HashClave(clave string) (string, error) {
	sal := make([]byte, largoSal)
	if _, err := rand.Read(sal); err != nil {
		return "", err
	}
	h, err := pbkdf2.Key(sha256.New, clave, sal, iteraciones, largoHash)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("pbkdf2-sha256$%d$%s$%s",
		iteraciones,
		base64.RawStdEncoding.EncodeToString(sal),
		base64.RawStdEncoding.EncodeToString(h),
	), nil
}

var ErrFormatoHash = errors.New("formato de hash inválido")

// VerificarClave compara en tiempo constante.
func VerificarClave(clave, guardado string) (bool, error) {
	partes := strings.Split(guardado, "$")
	if len(partes) != 4 || partes[0] != "pbkdf2-sha256" {
		return false, ErrFormatoHash
	}
	iter, err := strconv.Atoi(partes[1])
	if err != nil || iter <= 0 {
		return false, ErrFormatoHash
	}
	sal, err := base64.RawStdEncoding.DecodeString(partes[2])
	if err != nil {
		return false, ErrFormatoHash
	}
	esperado, err := base64.RawStdEncoding.DecodeString(partes[3])
	if err != nil {
		return false, ErrFormatoHash
	}
	calculado, err := pbkdf2.Key(sha256.New, clave, sal, iter, len(esperado))
	if err != nil {
		return false, err
	}
	return subtle.ConstantTimeCompare(calculado, esperado) == 1, nil
}
