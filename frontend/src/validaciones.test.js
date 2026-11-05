import { validarProducto } from './validaciones.js'

const valido = {
  nombre: 'Café de altura',
  descripcion: 'Café tostado artesanalmente, bolsa de 340 g.',
  precio: '4500',
  categoriaId: '1',
  imagen: null,
  imagenAlt: '',
}

const archivo = (type, size) => ({ type, size })

describe('validarProducto', () => {
  it('acepta un producto válido sin imagen', () => {
    expect(validarProducto(valido)).toEqual({})
  })

  it.each([
    ['nombre vacío', { nombre: '' }, 'nombre'],
    ['nombre corto', { nombre: 'Té' }, 'nombre'],
    ['nombre largo', { nombre: 'a'.repeat(81) }, 'nombre'],
    ['descripción corta', { descripcion: 'Rico' }, 'descripcion'],
    ['descripción larga', { descripcion: 'a'.repeat(501) }, 'descripcion'],
    ['precio vacío', { precio: '' }, 'precio'],
    ['precio cero', { precio: '0' }, 'precio'],
    ['precio negativo', { precio: '-10' }, 'precio'],
    ['precio con texto', { precio: 'cuatro mil' }, 'precio'],
    ['precio con decimales', { precio: '4500.50' }, 'precio'],
    ['precio excesivo', { precio: '10000001' }, 'precio'],
    ['sin categoría', { categoriaId: '' }, 'categoriaId'],
  ])('%s', (_, cambio, campo) => {
    const errores = validarProducto({ ...valido, ...cambio })
    expect(Object.keys(errores)).toEqual([campo])
  })

  it('una imagen GIF no se permite', () => {
    const e = validarProducto({ ...valido, imagen: archivo('image/gif', 100), imagenAlt: 'algo' })
    expect(e.imagen).toMatch(/JPG o PNG/)
  })

  it('una imagen de más de 2 MB no se permite', () => {
    const e = validarProducto({ ...valido, imagen: archivo('image/png', 2 * 1024 * 1024 + 1), imagenAlt: 'algo' })
    expect(e.imagen).toMatch(/2 MB/)
  })

  it('con imagen, el texto alternativo es obligatorio', () => {
    const e = validarProducto({ ...valido, imagen: archivo('image/jpeg', 1000), imagenAlt: '  ' })
    expect(Object.keys(e)).toEqual(['imagenAlt'])
  })

  it('cuenta letras con tilde como un solo carácter', () => {
    expect(validarProducto({ ...valido, nombre: 'Ñañ' })).toEqual({})
  })
})
