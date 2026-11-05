package httpapi

// Pruebas de integración de la HU-03 (catálogo público) y prueba no funcional de rendimiento.

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"testing"
	"time"
)

func pedirCatalogo(t *testing.T, url, query string) (*http.Response, []map[string]any) {
	t.Helper()
	resp, err := http.Get(url + "/api/catalogo" + query)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	var lista []map[string]any
	json.NewDecoder(resp.Body).Decode(&lista)
	return resp, lista
}

// CP-15: el catálogo es público y solo muestra publicados.
func TestCatalogo_SoloPublicados(t *testing.T) {
	_, ts := nuevoEntorno(t)
	crearPendiente(t, ts.URL) // un pendiente más que no debe aparecer

	resp, lista := pedirCatalogo(t, ts.URL, "")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("esperaba 200 sin sesión, obtuvo %d", resp.StatusCode)
	}
	if len(lista) != 3 { // el seed trae 3 publicados
		t.Fatalf("esperaba 3 productos publicados, obtuvo %d", len(lista))
	}
	for _, p := range lista {
		if _, tieneEstado := p["estado"]; tieneEstado {
			t.Fatal("el catálogo público no debe exponer el estado interno")
		}
		if _, tieneMotivo := p["motivoRechazo"]; tieneMotivo {
			t.Fatal("el catálogo público no debe exponer motivos de rechazo")
		}
	}
}

// CP-16: flujo completo registrar → aprobar → aparece en el catálogo.
func TestCatalogo_FlujoCompleto(t *testing.T) {
	_, ts := nuevoEntorno(t)
	id := crearPendiente(t, ts.URL)

	contiene := func() bool {
		_, lista := pedirCatalogo(t, ts.URL, "")
		for _, p := range lista {
			if int64(p["id"].(float64)) == id {
				return true
			}
		}
		return false
	}

	if contiene() {
		t.Fatal("un producto pendiente no debe verse en el catálogo")
	}
	revisar(t, clienteLogueado(t, ts.URL, "admin@asociacion.test"), ts.URL, id, "aprobar", "")
	if !contiene() {
		t.Fatal("después de aprobarlo, el producto debe aparecer en el catálogo")
	}
}

// CP-17: filtro por categoría.
func TestCatalogo_FiltroCategoria(t *testing.T) {
	_, ts := nuevoEntorno(t)

	var alimentos float64
	resp, _ := http.Get(ts.URL + "/api/categorias")
	var cats []map[string]any
	json.NewDecoder(resp.Body).Decode(&cats)
	for _, c := range cats {
		if c["nombre"] == "Alimentos" {
			alimentos = c["id"].(float64)
		}
	}

	_, lista := pedirCatalogo(t, ts.URL, fmt.Sprintf("?categoria=%d", int(alimentos)))
	if len(lista) != 2 { // el seed trae 2 alimentos publicados
		t.Fatalf("esperaba 2 productos de Alimentos, obtuvo %d", len(lista))
	}
	for _, p := range lista {
		if p["categoria"] != "Alimentos" {
			t.Fatalf("el filtro dejó pasar otra categoría: %v", p["categoria"])
		}
	}

	if resp, _ := pedirCatalogo(t, ts.URL, "?categoria=abc"); resp.StatusCode != http.StatusBadRequest {
		t.Errorf("categoría inválida: esperaba 400, obtuvo %d", resp.StatusCode)
	}
	if _, vacia := pedirCatalogo(t, ts.URL, "?categoria=9999"); len(vacia) != 0 {
		t.Errorf("categoría sin productos debe devolver lista vacía")
	}
}

// CP-18 (no funcional, rendimiento): con 500 productos publicados, el percentil 95
// de 30 consultas al catálogo debe quedar por debajo de 300 ms.
func TestCatalogo_Rendimiento(t *testing.T) {
	if testing.Short() {
		t.Skip("prueba de rendimiento omitida con -short")
	}
	s, ts := nuevoEntorno(t)

	_, err := s.db.Exec(`
		INSERT INTO productos (usuario_id, categoria_id, nombre, descripcion, precio, estado)
		SELECT (SELECT id FROM usuarios WHERE correo = 'emprendedora1@asociacion.test'),
		       (SELECT id FROM categorias ORDER BY id LIMIT 1),
		       'Producto sintético ' || g, 'Descripción sintética para prueba de carga.', 1000 + g, 'publicado'
		FROM generate_series(1, 500) AS g`)
	if err != nil {
		t.Fatal(err)
	}

	const n = 30
	tiempos := make([]time.Duration, 0, n)
	for range n {
		inicio := time.Now()
		resp, lista := pedirCatalogo(t, ts.URL, "")
		tiempos = append(tiempos, time.Since(inicio))
		if resp.StatusCode != http.StatusOK || len(lista) == 0 {
			t.Fatalf("respuesta inesperada: %d", resp.StatusCode)
		}
	}
	sort.Slice(tiempos, func(i, j int) bool { return tiempos[i] < tiempos[j] })
	p95 := tiempos[(n*95+99)/100-1] // percentil 95 (posición 29 de 30)
	t.Logf("catálogo con 503 publicados: mediana %v, p95 %v", tiempos[n/2], p95)
	if p95 > 300*time.Millisecond {
		t.Fatalf("p95 de %v supera 300 ms", p95)
	}
}

// DEF-04 (regresión): el catálogo usa ETag; si no cambió responde 304,
// y después de aprobar un producto el ETag cambia y se recibe la lista nueva.
func TestCatalogo_ETag(t *testing.T) {
	_, ts := nuevoEntorno(t)

	resp, _ := pedirCatalogo(t, ts.URL, "")
	etag := resp.Header.Get("ETag")
	if etag == "" || resp.Header.Get("Cache-Control") != "no-cache" {
		t.Fatalf("faltan cabeceras de caché: ETag=%q Cache-Control=%q", etag, resp.Header.Get("Cache-Control"))
	}

	pedirConEtag := func() *http.Response {
		req, _ := http.NewRequest(http.MethodGet, ts.URL+"/api/catalogo", nil)
		req.Header.Set("If-None-Match", etag)
		r, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		r.Body.Close()
		return r
	}

	if r := pedirConEtag(); r.StatusCode != http.StatusNotModified {
		t.Fatalf("sin cambios: esperaba 304, obtuvo %d", r.StatusCode)
	}

	id := crearPendiente(t, ts.URL)
	revisar(t, clienteLogueado(t, ts.URL, "admin@asociacion.test"), ts.URL, id, "aprobar", "")

	if r := pedirConEtag(); r.StatusCode != http.StatusOK {
		t.Fatalf("después de aprobar: esperaba 200 con datos nuevos, obtuvo %d", r.StatusCode)
	}
}
