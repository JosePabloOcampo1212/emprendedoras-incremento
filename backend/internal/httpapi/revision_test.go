package httpapi

// Pruebas de integración de la HU-02 (aprobar o rechazar).

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func pendientes(t *testing.T, c *http.Client, url string) (int, []map[string]any) {
	t.Helper()
	resp, err := c.Get(url + "/api/admin/pendientes")
	if err != nil {
		t.Fatal(err)
	}
	var lista []map[string]any
	json.NewDecoder(resp.Body).Decode(&lista)
	return resp.StatusCode, lista
}

func revisar(t *testing.T, c *http.Client, url string, id int64, accion, motivo string) (int, map[string]any) {
	t.Helper()
	cuerpo := "{}"
	if accion == "rechazar" {
		b, _ := json.Marshal(map[string]string{"motivo": motivo})
		cuerpo = string(b)
	}
	resp, err := c.Post(url+"/api/admin/productos/"+strconv.FormatInt(id, 10)+"/"+accion,
		"application/json", strings.NewReader(cuerpo))
	if err != nil {
		t.Fatal(err)
	}
	var datos map[string]any
	json.NewDecoder(resp.Body).Decode(&datos)
	return resp.StatusCode, datos
}

// crea un producto nuevo como emprendedora1 y devuelve su id
func crearPendiente(t *testing.T, url string) int64 {
	t.Helper()
	c := clienteLogueado(t, url, "emprendedora1@asociacion.test")
	resp, datos := enviarProducto(t, c, url, formulario{campos: camposValidos()})
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("no se pudo crear el producto: %d %v", resp.StatusCode, datos)
	}
	return int64(datos["id"].(float64))
}

func estadoDe(t *testing.T, s *Servidor, id int64) (estado, motivo string) {
	t.Helper()
	var m *string
	if err := s.db.QueryRow(`SELECT estado, motivo_rechazo FROM productos WHERE id = $1`, id).Scan(&estado, &m); err != nil {
		t.Fatal(err)
	}
	if m != nil {
		motivo = *m
	}
	return
}

// CP-10: solo la administradora ve la lista de pendientes.
func TestPendientes_SoloAdministradora(t *testing.T) {
	_, ts := nuevoEntorno(t)

	if cod, _ := pendientes(t, clienteLogueado(t, ts.URL, "emprendedora1@asociacion.test"), ts.URL); cod != http.StatusForbidden {
		t.Errorf("emprendedora: esperaba 403, obtuvo %d", cod)
	}
	if cod, _ := pendientes(t, clienteConCookies(t), ts.URL); cod != http.StatusUnauthorized {
		t.Errorf("sin sesión: esperaba 401, obtuvo %d", cod)
	}

	cod, lista := pendientes(t, clienteLogueado(t, ts.URL, "admin@asociacion.test"), ts.URL)
	if cod != http.StatusOK {
		t.Fatalf("administradora: esperaba 200, obtuvo %d", cod)
	}
	for _, p := range lista {
		if p["estado"] != "pendiente" {
			t.Fatalf("la lista incluye un producto no pendiente: %v", p)
		}
	}
	if len(lista) != 1 { // el seed trae un producto pendiente
		t.Fatalf("esperaba 1 pendiente del seed, hay %d", len(lista))
	}
}

// CP-11: aprobar publica el producto y deja historial.
func TestAprobar_PublicaProducto(t *testing.T) {
	s, ts := nuevoEntorno(t)
	id := crearPendiente(t, ts.URL)
	admin := clienteLogueado(t, ts.URL, "admin@asociacion.test")

	cod, datos := revisar(t, admin, ts.URL, id, "aprobar", "")
	if cod != http.StatusOK || datos["estado"] != "publicado" {
		t.Fatalf("esperaba 200 publicado, obtuvo %d %v", cod, datos)
	}
	if estado, _ := estadoDe(t, s, id); estado != "publicado" {
		t.Fatalf("estado en BD: %s", estado)
	}

	var hist int
	s.db.QueryRow(`SELECT COUNT(*) FROM historial_estados WHERE producto_id = $1 AND estado_nuevo = 'publicado'`, id).Scan(&hist)
	if hist != 1 {
		t.Fatalf("falta el registro de aprobación en el historial")
	}
}

// CP-12: rechazar exige motivo; con motivo vuelve a borrador y la emprendedora lo ve.
func TestRechazar_ConMotivo(t *testing.T) {
	s, ts := nuevoEntorno(t)
	id := crearPendiente(t, ts.URL)
	admin := clienteLogueado(t, ts.URL, "admin@asociacion.test")

	cod, datos := revisar(t, admin, ts.URL, id, "rechazar", "  ")
	if cod != http.StatusUnprocessableEntity {
		t.Fatalf("sin motivo: esperaba 422, obtuvo %d", cod)
	}
	if errs, _ := datos["errores"].(map[string]any); errs["motivo"] == nil {
		t.Fatalf("falta el error del campo motivo: %v", datos)
	}
	if estado, _ := estadoDe(t, s, id); estado != "pendiente" {
		t.Fatalf("sin motivo el estado no debe cambiar, quedó %s", estado)
	}

	motivo := "Falta indicar el peso o tamaño del producto."
	cod, _ = revisar(t, admin, ts.URL, id, "rechazar", motivo)
	if cod != http.StatusOK {
		t.Fatalf("con motivo: esperaba 200, obtuvo %d", cod)
	}
	estado, guardado := estadoDe(t, s, id)
	if estado != "borrador" || guardado != motivo {
		t.Fatalf("esperaba borrador con motivo, obtuvo %s / %q", estado, guardado)
	}

	// La emprendedora ve el motivo en su lista
	c := clienteLogueado(t, ts.URL, "emprendedora1@asociacion.test")
	visto := false
	for _, p := range misProductosDe(t, c, ts.URL) {
		if int64(p["id"].(float64)) == id && p["motivoRechazo"] == motivo {
			visto = true
		}
	}
	if !visto {
		t.Fatal("la emprendedora no ve el motivo de rechazo")
	}
}

// CP-13: un producto ya revisado no se puede revisar otra vez.
func TestRevisar_DosVecesDaConflicto(t *testing.T) {
	_, ts := nuevoEntorno(t)
	id := crearPendiente(t, ts.URL)
	admin := clienteLogueado(t, ts.URL, "admin@asociacion.test")

	revisar(t, admin, ts.URL, id, "aprobar", "")
	if cod, _ := revisar(t, admin, ts.URL, id, "aprobar", ""); cod != http.StatusConflict {
		t.Errorf("aprobar de nuevo: esperaba 409, obtuvo %d", cod)
	}
	if cod, _ := revisar(t, admin, ts.URL, id, "rechazar", "Motivo suficientemente largo"); cod != http.StatusConflict {
		t.Errorf("rechazar un publicado: esperaba 409, obtuvo %d", cod)
	}
}

// CP-14: producto inexistente o id inválido da 404; emprendedora no puede aprobar.
func TestRevisar_NoEncontradoYPermisos(t *testing.T) {
	_, ts := nuevoEntorno(t)
	admin := clienteLogueado(t, ts.URL, "admin@asociacion.test")

	if cod, _ := revisar(t, admin, ts.URL, 999999, "aprobar", ""); cod != http.StatusNotFound {
		t.Errorf("inexistente: esperaba 404, obtuvo %d", cod)
	}
	resp, _ := admin.Post(ts.URL+"/api/admin/productos/abc/aprobar", "application/json", nil)
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("id inválido: esperaba 404, obtuvo %d", resp.StatusCode)
	}

	id := crearPendiente(t, ts.URL)
	emp := clienteLogueado(t, ts.URL, "emprendedora1@asociacion.test")
	if cod, _ := revisar(t, emp, ts.URL, id, "aprobar", ""); cod != http.StatusForbidden {
		t.Errorf("emprendedora aprobando: esperaba 403, obtuvo %d", cod)
	}
}
