"""Prueba funcional de punta a punta del incremento (registro -> revisión -> catálogo).

Requiere: base de datos recién sembrada, backend encendido y la versión de producción
del frontend servida con `npm run build && npm run preview` (puerto 4173).
La clave de las cuentas se toma de SEED_CLAVE (variable de entorno o archivo .env).
Las capturas quedan en pruebas/e2e/capturas/.
"""
import os
import pathlib
from playwright.sync_api import sync_playwright, expect

AQUI = pathlib.Path(__file__).resolve().parent
RAIZ = AQUI.parent.parent


def leer_env(clave):
    if os.environ.get(clave):
        return os.environ[clave]
    archivo = RAIZ / ".env"
    if archivo.exists():
        for linea in archivo.read_text(encoding="utf-8").splitlines():
            if linea.startswith(clave + "="):
                return linea.split("=", 1)[1].strip().strip("\"'")
    raise SystemExit(f"Falta {clave} (en el entorno o en .env)")


BASE = os.environ.get("E2E_URL", "http://localhost:4173")
CLAVE = leer_env("SEED_CLAVE")
CAP = str(AQUI / "capturas") + "/"
IMAGEN = str(AQUI / "mermelada.png")
os.makedirs(CAP, exist_ok=True)
resultados = []


def paso(nombre, ok, detalle=""):
    resultados.append((nombre, ok, detalle))
    print(("PASA " if ok else "FALLA"), nombre, detalle)


def ingresar(page, correo):
    page.goto(BASE + "/ingresar")
    page.get_by_label("Correo").fill(correo)
    page.get_by_label("Clave").fill(CLAVE)
    page.get_by_role("button", name="Ingresar").click()


def salir(page):
    page.get_by_role("button", name="Salir").click()
    page.wait_for_url(BASE + "/")


with sync_playwright() as p:
    nav = p.chromium.launch(channel=os.environ.get("E2E_NAVEGADOR", "chrome"))
    page = nav.new_page(viewport={"width": 1280, "height": 860}, locale="es-CR")

    # 1. Catálogo inicial
    page.goto(BASE + "/")
    expect(page.get_by_role("status")).to_have_text("3 productos")
    page.screenshot(path=CAP + "01-catalogo.png", full_page=True)
    paso("Catálogo público muestra 3 publicados del seed", True)

    # 2. Login con errores
    page.goto(BASE + "/ingresar")
    page.get_by_role("button", name="Ingresar").click()
    expect(page.get_by_label("Correo")).to_have_attribute("aria-invalid", "true")
    page.get_by_label("Correo").fill("emprendedora1@asociacion.test")
    page.get_by_label("Clave").fill("clave-equivocada")
    page.get_by_role("button", name="Ingresar").click()
    expect(page.get_by_role("alert")).to_have_text("Correo o clave incorrectos.")
    page.screenshot(path=CAP + "02-login-error.png")
    paso("Login rechaza clave incorrecta con mensaje genérico", True)

    # 3. Emprendedora: errores de validación
    ingresar(page, "emprendedora1@asociacion.test")
    page.wait_for_url(BASE + "/mis-productos")
    page.get_by_role("button", name="Enviar a revisión").click()
    expect(page.get_by_label("Nombre del producto")).to_have_attribute("aria-invalid", "true")
    expect(page.get_by_label("Nombre del producto")).to_be_focused()
    page.screenshot(path=CAP + "03-registro-errores.png", full_page=True)
    paso("HU-01 CA1: errores junto a cada campo y foco en el primero", True)

    # 4. Registro solo con teclado (CA4). El archivo se adjunta con set_input_files
    #    porque el diálogo del sistema operativo no es parte de la página.
    page.reload()
    page.wait_for_selector("text=Registrar un producto")
    page.keyboard.press("Tab")  # enlace "Saltar al contenido"
    activo = ""
    for _ in range(15):
        page.keyboard.press("Tab")
        activo = page.evaluate("document.activeElement.name || ''")
        if activo == "nombre":
            break
    paso("HU-01 CA4: se llega al formulario con Tab", activo == "nombre", f"campo activo: {activo}")
    page.keyboard.type("Mermelada de mora con chile")
    page.keyboard.press("Tab")
    page.keyboard.type("Mermelada de mora con un toque de chile dulce, frasco de 250 g.")
    page.keyboard.press("Tab")
    page.keyboard.type("2800")
    page.keyboard.press("Tab")
    # En macOS la flecha abre el menú nativo del select en vez de cambiar la opción;
    # se elige sobre el campo que ya tiene el foco, sin usar el mouse.
    page.locator(":focus").select_option(label="Alimentos")
    page.get_by_label("Imagen (opcional)").set_input_files(IMAGEN)
    page.get_by_label("Descripción de la imagen").focus()
    page.keyboard.type("Frasco de mermelada morada con tapa de tela")
    # Enter en el último campo envía el formulario (en macOS Tab no siempre llega a los botones)
    enfocado = page.evaluate("document.activeElement.name")
    page.keyboard.press("Enter")
    expect(page.get_by_role("status")).to_contain_text("Producto enviado a revisión")
    expect(page.get_by_text("Mermelada de mora con chile", exact=True)).to_be_visible()
    page.screenshot(path=CAP + "04-registro-exito.png", full_page=True)
    paso("HU-01 CA2/CA4: registro completado con teclado, queda 'En revisión'", True, f"Enter en el campo: {enfocado}")
    salir(page)

    # 5. Administradora: pendientes
    ingresar(page, "admin@asociacion.test")
    page.wait_for_url(BASE + "/admin/pendientes")
    expect(page.get_by_role("heading", name="Mermelada de mora con chile")).to_be_visible()
    expect(page.get_by_role("img", name="Frasco de mermelada morada con tapa de tela")).to_be_visible()
    page.screenshot(path=CAP + "05-pendientes.png", full_page=True)
    paso("HU-02 CA1: la administradora ve los pendientes con imagen y texto alternativo", True)

    # 6. Rechazar sin motivo, luego con motivo
    tarjeta = page.locator("li", has=page.get_by_role("heading", name="Portavasos de fibra natural"))
    tarjeta.get_by_role("button", name="Rechazar").click()
    tarjeta.get_by_role("button", name="Confirmar rechazo").click()
    expect(tarjeta.get_by_label("Motivo del rechazo")).to_have_attribute("aria-invalid", "true")
    page.screenshot(path=CAP + "06-rechazo-sin-motivo.png", full_page=True)
    tarjeta.get_by_label("Motivo del rechazo").fill("Indique cuántas unidades trae el juego y sus medidas.")
    tarjeta.get_by_role("button", name="Confirmar rechazo").click()
    expect(page.get_by_role("status")).to_contain_text("fue devuelto")
    paso("HU-02 CA3: rechazo exige motivo y devuelve el producto", True)

    # 7. Aprobar
    t2 = page.locator("li", has=page.get_by_role("heading", name="Mermelada de mora con chile"))
    t2.get_by_role("button", name="Aprobar").click()
    expect(page.get_by_role("status")).to_contain_text("fue aprobado")
    expect(page.get_by_text("No hay productos pendientes de revisión.")).to_be_visible()
    page.screenshot(path=CAP + "07-aprobado.png", full_page=True)
    paso("HU-02 CA2: aprobar publica el producto", True)
    salir(page)

    # 8. Catálogo con filtro (sin recargar)
    page.goto(BASE + "/")
    expect(page.get_by_role("status")).to_have_text("4 productos")
    url_antes = page.evaluate("performance.getEntriesByType('navigation').length")
    page.get_by_label("Filtrar por categoría").select_option(label="Alimentos")
    expect(page.get_by_role("status")).to_have_text("3 productos en Alimentos")
    navegaciones = page.evaluate("performance.getEntriesByType('navigation').length")
    expect(page.get_by_role("heading", name="Mermelada de mora con chile")).to_be_visible()
    page.screenshot(path=CAP + "08-catalogo-filtro.png", full_page=True)
    paso("HU-03 CA1/CA2: aparece el producto aprobado y el filtro no recarga la página",
         navegaciones == url_antes, f"url: {page.url}")

    # 9. Emprendedora 2 ve el motivo
    ingresar(page, "emprendedora2@asociacion.test")
    page.wait_for_url(BASE + "/mis-productos")
    expect(page.get_by_text("Indique cuántas unidades trae el juego")).to_be_visible()
    page.screenshot(path=CAP + "09-motivo-visible.png", full_page=True)
    paso("HU-02 CA3: la emprendedora ve el motivo de devolución", True)

    # 10. Acceso por URL a página de otro rol
    page.goto(BASE + "/admin/pendientes")
    expect(page.get_by_role("heading", name="Sin acceso")).to_be_visible()
    paso("Una emprendedora no puede abrir la página de revisión", True)

    # 11. Vista móvil del catálogo
    movil = nav.new_page(viewport={"width": 390, "height": 844}, locale="es-CR")
    movil.goto(BASE + "/")
    expect(movil.get_by_role("status")).to_have_text("4 productos")
    ancho = movil.evaluate("document.documentElement.scrollWidth")
    movil.screenshot(path=CAP + "10-catalogo-movil.png", full_page=True)
    paso("Catálogo en móvil sin desplazamiento horizontal", ancho <= 390, f"ancho: {ancho}px")

    nav.close()

fallas = [r for r in resultados if not r[1]]
print(f"\n{len(resultados) - len(fallas)}/{len(resultados)} pasos correctos")
raise SystemExit(1 if fallas else 0)
