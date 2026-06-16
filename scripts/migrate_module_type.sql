-- Alinea module_type con el backend/front actual:
--   attendance, meal, transport, custom
--
-- Antes (según DB): meal, transport, inventory, custom
-- - Agrega attendance
-- - Elimina inventory (filas pasan a custom)
-- - No incluye materials
--
-- Ejecutar en la base Macabi (psql, DBeaver, etc.) con el usuario dueño del schema.

BEGIN;

-- 1) Normalizar filas que usan el valor viejo "inventory"
UPDATE event_modules
SET type = 'custom'
WHERE type::text = 'inventory';

-- 2) Nuevo enum
CREATE TYPE module_type_new AS ENUM (
  'attendance',
  'meal',
  'transport',
  'custom'
);

-- 3) Migrar columna (fallará si queda algún valor fuera de la lista nueva)
ALTER TABLE event_modules
  ALTER COLUMN type TYPE module_type_new
  USING type::text::module_type_new;

-- 4) Reemplazar tipo viejo
DROP TYPE module_type;
ALTER TYPE module_type_new RENAME TO module_type;

COMMIT;

-- Verificación:
-- SELECT enumlabel FROM pg_enum e
-- JOIN pg_type t ON e.enumtypid = t.oid
-- WHERE t.typname = 'module_type'
-- ORDER BY e.enumsortorder;
