package httpapi

// Pruebas de integración: levantan la API completa contra PostgreSQL,
// en un esquema temporal con los datos sintéticos del seed.

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/JosePabloOcampo1212/emprendedoras-incremento/backend/internal/auth"
	"github.com/JosePabloOcampo1212/emprendedoras-incremento/backend/internal/dbprueba"
	"github.com/JosePabloOcampo1212/emprendedoras-incremento/backend/internal/seed"
)

const claveDePrueba = "clave-solo-para-pruebas"

func nuevoEntorno(t *testing.T) (*Servidor, *httptest.Server) {
	t.Helper()
	conexion := dbprueba.Nueva(t)

	if err := seed.Cargar(context.Background(), conexion, claveDePrueba); err != nil {
		t.Fatal(err)
	}
	s := Nuevo(conexion, filepath.Join(t.TempDir(), "uploads"), false)
	ts := httptest.NewServer(s.Rutas())
	t.Cleanup(ts.Close)
	return s, ts
}

func clienteConCookies(t *testing.T) *http.Client {
	jar, _ := cookiejar.New(nil)
	return &http.Client{Jar: jar}
}

func login(t *testing.T, c *http.Client, url, correo, clave string) *http.Response {
	t.Helper()
	cuerpo := `{"correo":"` + correo + `","clave":"` + clave + `"}`
	resp, err := c.Post(url+"/api/sesion", "application/json", strings.NewReader(cuerpo))
	if err != nil {
		t.Fatal(err)
	}
	return resp
}

func TestSalud(t *testing.T) {
	_, ts := nuevoEntorno(t)
	resp, err := http.Get(ts.URL + "/api/salud")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("estado %d", resp.StatusCode)
	}
}

func TestCabecerasSeguridad(t *testing.T) {
	_, ts := nuevoEntorno(t)
	resp, _ := http.Get(ts.URL + "/api/salud")
	for _, h := range []string{"X-Content-Type-Options", "X-Frame-Options", "Content-Security-Policy"} {
		if resp.Header.Get(h) == "" {
			t.Errorf("falta la cabecera %s", h)
		}
	}
}

func TestCategorias(t *testing.T) {
	_, ts := nuevoEntorno(t)
	resp, _ := http.Get(ts.URL + "/api/categorias")
	var lista []categoria
	json.NewDecoder(resp.Body).Decode(&lista)
	if len(lista) != len(seed.Categorias) {
		t.Fatalf("esperaba %d categorías, obtuvo %d", len(seed.Categorias), len(lista))
	}
}

func TestLogin_Correcto(t *testing.T) {
	_, ts := nuevoEntorno(t)
	c := clienteConCookies(t)

	resp := login(t, c, ts.URL, "emprendedora1@asociacion.test", claveDePrueba)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("estado %d", resp.StatusCode)
	}

	var cookieSesion *http.Cookie
	for _, ck := range resp.Cookies() {
		if ck.Name == nombreCookie {
			cookieSesion = ck
		}
	}
	if cookieSesion == nil || !cookieSesion.HttpOnly {
		t.Fatal("la cookie de sesión debe existir y ser HttpOnly")
	}

	// Con la cookie, /api/sesion devuelve al usuario
	resp2, _ := c.Get(ts.URL + "/api/sesion")
	var u auth.Usuario
	json.NewDecoder(resp2.Body).Decode(&u)
	if u.Rol != auth.RolEmprendedora {
		t.Fatalf("rol inesperado: %q", u.Rol)
	}
}

func TestLogin_CorreoSinImportarMayusculas(t *testing.T) {
	_, ts := nuevoEntorno(t)
	resp := login(t, clienteConCookies(t), ts.URL, "  Admin@Asociacion.TEST ", claveDePrueba)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("estado %d", resp.StatusCode)
	}
}

func TestLogin_MismoMensajeParaCorreoOClaveIncorrectos(t *testing.T) {
	_, ts := nuevoEntorno(t)

	leerError := func(resp *http.Response) string {
		var m map[string]string
		json.NewDecoder(resp.Body).Decode(&m)
		return m["error"]
	}

	r1 := login(t, clienteConCookies(t), ts.URL, "admin@asociacion.test", "clave-mala")
	r2 := login(t, clienteConCookies(t), ts.URL, "noexiste@asociacion.test", claveDePrueba)

	if r1.StatusCode != http.StatusUnauthorized || r2.StatusCode != http.StatusUnauthorized {
		t.Fatalf("esperaba 401 en ambos, obtuvo %d y %d", r1.StatusCode, r2.StatusCode)
	}
	if m1, m2 := leerError(r1), leerError(r2); m1 != m2 {
		t.Fatalf("los mensajes deben ser iguales: %q vs %q", m1, m2)
	}
}

func TestSesion_SinCookieDa401(t *testing.T) {
	_, ts := nuevoEntorno(t)
	resp, _ := http.Get(ts.URL + "/api/sesion")
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("estado %d", resp.StatusCode)
	}
}

func TestCerrarSesion_InvalidaElToken(t *testing.T) {
	_, ts := nuevoEntorno(t)
	c := clienteConCookies(t)
	login(t, c, ts.URL, "emprendedora1@asociacion.test", claveDePrueba)

	req, _ := http.NewRequest(http.MethodDelete, ts.URL+"/api/sesion", nil)
	if resp, _ := c.Do(req); resp.StatusCode != http.StatusNoContent {
		t.Fatalf("estado %d al cerrar sesión", resp.StatusCode)
	}
	if resp, _ := c.Get(ts.URL + "/api/sesion"); resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("después de cerrar sesión esperaba 401, obtuvo %d", resp.StatusCode)
	}
}

// Prueba del control de acceso por rol con una ruta de prueba.
func TestRequiereRol(t *testing.T) {
	s, ts := nuevoEntorno(t)

	soloAdmin := httptest.NewServer(s.requiereRol(auth.RolAdministradora,
		func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) }))
	defer soloAdmin.Close()

	pedir := func(correo string) int {
		c := clienteConCookies(t)
		resp := login(t, c, ts.URL, correo, claveDePrueba)
		req, _ := http.NewRequest(http.MethodGet, soloAdmin.URL, nil)
		for _, ck := range resp.Cookies() {
			req.AddCookie(ck)
		}
		r, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		return r.StatusCode
	}

	if got := pedir("emprendedora1@asociacion.test"); got != http.StatusForbidden {
		t.Errorf("emprendedora: esperaba 403, obtuvo %d", got)
	}
	if got := pedir("admin@asociacion.test"); got != http.StatusOK {
		t.Errorf("administradora: esperaba 200, obtuvo %d", got)
	}
}
