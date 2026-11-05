package httpapi

// Pruebas de integración de la HU-01 (registrar producto).

import (
	"bytes"
	"encoding/json"
	"image"
	"image/png"
	"mime/multipart"
	"net/http"
	"strings"
	"testing"
)

type formulario struct {
	campos map[string]string
	imagen []byte // nil = sin imagen
}

func camposValidos() map[string]string {
	return map[string]string{
		"nombre":      "Café de altura",
		"descripcion": "Café tostado artesanalmente, bolsa de 340 g.",
		"precio":      "4500",
		"categoriaId": "1",
		"imagenAlt":   "Bolsa de café sobre una mesa de madera",
	}
}

func pngDePrueba(t *testing.T) []byte {
	t.Helper()
	var buf bytes.Buffer
	if err := png.Encode(&buf, image.NewRGBA(image.Rect(0, 0, 4, 4))); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

func enviarProducto(t *testing.T, c *http.Client, url string, f formulario) (*http.Response, map[string]any) {
	t.Helper()
	var cuerpo bytes.Buffer
	mw := multipart.NewWriter(&cuerpo)
	for k, v := range f.campos {
		mw.WriteField(k, v)
	}
	if f.imagen != nil {
		parte, _ := mw.CreateFormFile("imagen", "foto.png")
		parte.Write(f.imagen)
	}
	mw.Close()

	req, _ := http.NewRequest(http.MethodPost, url+"/api/productos", &cuerpo)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	resp, err := c.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	var datos map[string]any
	json.NewDecoder(resp.Body).Decode(&datos)
	resp.Body.Close()
	return resp, datos
}

func clienteLogueado(t *testing.T, url, correo string) *http.Client {
	t.Helper()
	c := clienteConCookies(t)
	if resp := login(t, c, url, correo, claveDePrueba); resp.StatusCode != http.StatusOK {
		t.Fatalf("no se pudo iniciar sesión como %s: %d", correo, resp.StatusCode)
	}
	return c
}

func misProductosDe(t *testing.T, c *http.Client, url string) []map[string]any {
	t.Helper()
	resp, err := c.Get(url + "/api/mis-productos")
	if err != nil {
		t.Fatal(err)
	}
	var lista []map[string]any
	json.NewDecoder(resp.Body).Decode(&lista)
	return lista
}

// CP-01: registro válido con imagen queda "pendiente" y aparece en mis productos.
func TestCrearProducto_ValidoQuedaPendiente(t *testing.T) {
	s, ts := nuevoEntorno(t)
	c := clienteLogueado(t, ts.URL, "emprendedora1@asociacion.test")

	resp, datos := enviarProducto(t, c, ts.URL, formulario{campos: camposValidos(), imagen: pngDePrueba(t)})
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("esperaba 201, obtuvo %d: %v", resp.StatusCode, datos)
	}
	if datos["estado"] != "pendiente" {
		t.Fatalf("estado inesperado: %v", datos["estado"])
	}

	var hist int
	s.db.QueryRow(`SELECT COUNT(*) FROM historial_estados WHERE producto_id = $1`, datos["id"]).Scan(&hist)
	if hist != 1 {
		t.Fatalf("esperaba 1 registro en el historial, hay %d", hist)
	}

	var encontrado map[string]any
	for _, p := range misProductosDe(t, c, ts.URL) {
		if p["nombre"] == "Café de altura" {
			encontrado = p
		}
	}
	if encontrado == nil {
		t.Fatal("el producto no aparece en mis productos")
	}
	urlImagen, _ := encontrado["imagenUrl"].(string)
	if !strings.HasPrefix(urlImagen, "/uploads/") || !strings.HasSuffix(urlImagen, ".png") {
		t.Fatalf("ruta de imagen inesperada: %q", urlImagen)
	}

	// La imagen se puede descargar
	r, _ := http.Get(ts.URL + urlImagen)
	if r.StatusCode != http.StatusOK {
		t.Fatalf("la imagen no se sirve: %d", r.StatusCode)
	}
}

// CP-02: datos inválidos no se guardan y se informan por campo.
func TestCrearProducto_InvalidoNoGuarda(t *testing.T) {
	s, ts := nuevoEntorno(t)
	c := clienteLogueado(t, ts.URL, "emprendedora1@asociacion.test")

	var antes int
	s.db.QueryRow(`SELECT COUNT(*) FROM productos`).Scan(&antes)

	campos := camposValidos()
	campos["nombre"] = ""
	campos["precio"] = "-5"
	resp, datos := enviarProducto(t, c, ts.URL, formulario{campos: campos})

	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Fatalf("esperaba 422, obtuvo %d", resp.StatusCode)
	}
	errores, _ := datos["errores"].(map[string]any)
	if errores["nombre"] == nil || errores["precio"] == nil {
		t.Fatalf("faltan errores por campo: %v", errores)
	}

	var despues int
	s.db.QueryRow(`SELECT COUNT(*) FROM productos`).Scan(&despues)
	if antes != despues {
		t.Fatal("se guardó un producto inválido")
	}
}

// CP-03: el precio con texto se rechaza.
func TestCrearProducto_PrecioNoNumerico(t *testing.T) {
	_, ts := nuevoEntorno(t)
	c := clienteLogueado(t, ts.URL, "emprendedora1@asociacion.test")
	campos := camposValidos()
	campos["precio"] = "cuatro mil"
	resp, datos := enviarProducto(t, c, ts.URL, formulario{campos: campos})
	errores, _ := datos["errores"].(map[string]any)
	if resp.StatusCode != http.StatusUnprocessableEntity || errores["precio"] == nil {
		t.Fatalf("esperaba error de precio, obtuvo %d %v", resp.StatusCode, datos)
	}
}

// CP-04: un archivo que no es JPG/PNG se rechaza aunque se llame .png.
func TestCrearProducto_ImagenTipoNoPermitido(t *testing.T) {
	_, ts := nuevoEntorno(t)
	c := clienteLogueado(t, ts.URL, "emprendedora1@asociacion.test")
	gif := []byte("GIF89a\x01\x00\x01\x00\x80\x00\x00\x00\x00\x00\xff\xff\xff!\xf9\x04\x01\x00\x00\x00\x00,\x00\x00\x00\x00\x01\x00\x01\x00\x00\x02\x02D\x01\x00;")
	resp, datos := enviarProducto(t, c, ts.URL, formulario{campos: camposValidos(), imagen: gif})
	errores, _ := datos["errores"].(map[string]any)
	if resp.StatusCode != http.StatusUnprocessableEntity || errores["imagen"] == nil {
		t.Fatalf("esperaba error de imagen, obtuvo %d %v", resp.StatusCode, datos)
	}
}

// CP-05: imagen mayor a 2 MB se rechaza.
func TestCrearProducto_ImagenMuyPesada(t *testing.T) {
	_, ts := nuevoEntorno(t)
	c := clienteLogueado(t, ts.URL, "emprendedora1@asociacion.test")
	grande := append(pngDePrueba(t), make([]byte, 3*1024*1024)...)
	resp, datos := enviarProducto(t, c, ts.URL, formulario{campos: camposValidos(), imagen: grande})
	errores, _ := datos["errores"].(map[string]any)
	if resp.StatusCode != http.StatusUnprocessableEntity || errores["imagen"] == nil {
		t.Fatalf("esperaba error de tamaño, obtuvo %d %v", resp.StatusCode, datos)
	}
}

// CP-06: imagen sin texto alternativo se rechaza.
func TestCrearProducto_ImagenSinAlt(t *testing.T) {
	_, ts := nuevoEntorno(t)
	c := clienteLogueado(t, ts.URL, "emprendedora1@asociacion.test")
	campos := camposValidos()
	campos["imagenAlt"] = "   "
	resp, datos := enviarProducto(t, c, ts.URL, formulario{campos: campos, imagen: pngDePrueba(t)})
	errores, _ := datos["errores"].(map[string]any)
	if resp.StatusCode != http.StatusUnprocessableEntity || errores["imagenAlt"] == nil {
		t.Fatalf("esperaba error de texto alternativo, obtuvo %d %v", resp.StatusCode, datos)
	}
}

// CP-07: categoría inexistente.
func TestCrearProducto_CategoriaInexistente(t *testing.T) {
	_, ts := nuevoEntorno(t)
	c := clienteLogueado(t, ts.URL, "emprendedora1@asociacion.test")
	campos := camposValidos()
	campos["categoriaId"] = "999"
	resp, datos := enviarProducto(t, c, ts.URL, formulario{campos: campos})
	errores, _ := datos["errores"].(map[string]any)
	if resp.StatusCode != http.StatusUnprocessableEntity || errores["categoriaId"] == nil {
		t.Fatalf("esperaba error de categoría, obtuvo %d %v", resp.StatusCode, datos)
	}
}

// CP-08: permisos. Sin sesión 401; administradora 403.
func TestCrearProducto_Permisos(t *testing.T) {
	_, ts := nuevoEntorno(t)

	resp, _ := enviarProducto(t, clienteConCookies(t), ts.URL, formulario{campos: camposValidos()})
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("sin sesión: esperaba 401, obtuvo %d", resp.StatusCode)
	}

	admin := clienteLogueado(t, ts.URL, "admin@asociacion.test")
	resp, _ = enviarProducto(t, admin, ts.URL, formulario{campos: camposValidos()})
	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("administradora: esperaba 403, obtuvo %d", resp.StatusCode)
	}
}

// CP-09: cada emprendedora ve solo sus productos.
func TestMisProductos_SoloPropios(t *testing.T) {
	_, ts := nuevoEntorno(t)
	c1 := clienteLogueado(t, ts.URL, "emprendedora1@asociacion.test")
	c2 := clienteLogueado(t, ts.URL, "emprendedora2@asociacion.test")

	for _, p := range misProductosDe(t, c1, ts.URL) {
		if p["negocio"] != "Dulces de la Abuela" {
			t.Fatalf("emprendedora1 ve un producto ajeno: %v", p)
		}
	}
	for _, p := range misProductosDe(t, c2, ts.URL) {
		if p["negocio"] != "Tejidos Montaña Verde" {
			t.Fatalf("emprendedora2 ve un producto ajeno: %v", p)
		}
	}
}

// El directorio de imágenes no debe listar archivos.
func TestUploads_SinListado(t *testing.T) {
	_, ts := nuevoEntorno(t)
	r, _ := http.Get(ts.URL + "/uploads/")
	if r.StatusCode != http.StatusNotFound {
		t.Fatalf("esperaba 404, obtuvo %d", r.StatusCode)
	}
}
