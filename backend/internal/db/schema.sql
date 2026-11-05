-- Esquema de la base de datos (PostgreSQL 16).
-- Se ejecuta al iniciar el servidor; usa IF NOT EXISTS para no borrar datos.

CREATE TABLE IF NOT EXISTS usuarios (
    id         BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    nombre     TEXT NOT NULL,
    correo     TEXT NOT NULL UNIQUE,
    clave_hash TEXT NOT NULL,
    rol        TEXT NOT NULL CHECK (rol IN ('emprendedora', 'administradora')),
    negocio    TEXT
);

CREATE TABLE IF NOT EXISTS sesiones (
    token      TEXT PRIMARY KEY,
    usuario_id BIGINT NOT NULL REFERENCES usuarios(id) ON DELETE CASCADE,
    expira     TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS categorias (
    id     BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    nombre TEXT NOT NULL UNIQUE
);

-- Ciclo de vida del contenido: borrador -> pendiente -> publicado -> archivado
CREATE TABLE IF NOT EXISTS productos (
    id             BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    usuario_id     BIGINT NOT NULL REFERENCES usuarios(id),
    categoria_id   BIGINT NOT NULL REFERENCES categorias(id),
    nombre         VARCHAR(80) NOT NULL,
    descripcion    VARCHAR(500) NOT NULL,
    precio         INTEGER NOT NULL CHECK (precio > 0 AND precio <= 10000000), -- colones
    imagen_ruta    TEXT,
    imagen_alt     VARCHAR(150),
    estado         TEXT NOT NULL DEFAULT 'pendiente'
                   CHECK (estado IN ('borrador', 'pendiente', 'publicado', 'archivado')),
    motivo_rechazo TEXT,
    creado         TIMESTAMPTZ NOT NULL DEFAULT now(),
    actualizado    TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Historial de cambios de estado (auditoría del ciclo de publicación)
CREATE TABLE IF NOT EXISTS historial_estados (
    id           BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    producto_id  BIGINT NOT NULL REFERENCES productos(id) ON DELETE CASCADE,
    estado_antes TEXT,
    estado_nuevo TEXT NOT NULL,
    usuario_id   BIGINT NOT NULL REFERENCES usuarios(id),
    nota         TEXT,
    fecha        TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_productos_estado ON productos(estado);
CREATE INDEX IF NOT EXISTS idx_productos_usuario ON productos(usuario_id);
CREATE INDEX IF NOT EXISTS idx_sesiones_expira ON sesiones(expira);
