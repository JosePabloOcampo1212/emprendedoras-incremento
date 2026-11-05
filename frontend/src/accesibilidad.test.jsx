// Prueba no funcional de accesibilidad: revisa cada pantalla con axe-core
// (el motor de reglas WCAG que usa Lighthouse). jsdom no calcula colores,
// por eso la regla de contraste se revisa aparte con Lighthouse en el navegador.

import { render, screen } from '@testing-library/react'
import axe from 'axe-core'
import { MemoryRouter } from 'react-router-dom'
import { vi } from 'vitest'
import { api } from './api.js'
import Catalogo from './pages/Catalogo.jsx'
import Login from './pages/Login.jsx'
import MisProductos from './pages/MisProductos.jsx'
import Pendientes from './pages/Pendientes.jsx'

vi.mock('./sesion/useSesion.js', async (original) => ({
  ...(await original()),
  useSesion: () => ({ usuario: { rol: 'emprendedora', negocio: 'Dulces de la Abuela' }, iniciar: vi.fn() }),
}))

const producto = {
  id: 1,
  nombre: 'Bufanda tejida',
  negocio: 'Tejidos Montaña Verde',
  categoria: 'Textiles',
  categoriaId: 4,
  precio: 12000,
  descripcion: 'Lana tejida a mano.',
  imagenUrl: '/uploads/a.png',
  imagenAlt: 'Bufanda color tierra',
  estado: 'pendiente',
}

beforeEach(() => {
  vi.restoreAllMocks()
  vi.spyOn(api, 'categorias').mockResolvedValue([{ id: 4, nombre: 'Textiles' }])
  vi.spyOn(api, 'catalogo').mockResolvedValue([producto])
  vi.spyOn(api, 'misProductos').mockResolvedValue([producto])
  vi.spyOn(api, 'pendientes').mockResolvedValue([producto])
})

async function revisar(contenedor) {
  const resultado = await axe.run(contenedor, {
    runOnly: { type: 'tag', values: ['wcag2a', 'wcag2aa', 'wcag21aa', 'wcag22aa'] },
    rules: { 'color-contrast': { enabled: false } },
  })
  return resultado.violations.map((v) => `${v.id}: ${v.help}`)
}

const pantallas = [
  ['Ingreso', <Login key="l" />, 'Ingresar'],
  ['Catálogo', <Catalogo key="c" />, 'Bufanda tejida'],
  ['Mis productos', <MisProductos key="m" />, 'Bufanda tejida'],
  ['Por revisar', <Pendientes key="p" />, 'Bufanda tejida'],
]

describe('Accesibilidad (axe-core, WCAG 2.2 A/AA)', () => {
  it.each(pantallas)('%s no tiene infracciones', async (_, pagina, textoEsperado) => {
    const { container } = render(<MemoryRouter>{pagina}</MemoryRouter>)
    await screen.findAllByText(textoEsperado)
    expect(await revisar(container)).toEqual([])
  })
})
