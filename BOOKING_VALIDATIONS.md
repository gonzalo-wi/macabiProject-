# Validaciones de reserva — `POST /bookings`

## Reglas de negocio

1. **Un usuario no puede reservar en dos proyectos distintos el mismo día.**
2. **Un usuario no puede tener almuerzo y cena el mismo día** (independientemente del proyecto).
3. **No se puede reservar si la vianda está agotada.**
4. **El plazo de reserva cierra el viernes a las 23:59** para el fin de semana siguiente.

---

## Respuestas de error `409 Conflict`

| Mensaje en `error` | Causa |
|---|---|
| `"vianda agotada"` | Sin stock disponible |
| `"ya tenés una reserva para este día; no podés reservar almuerzo y cena el mismo día"` | Misma fecha, mismo proyecto, distinto tipo de comida (almuerzo vs cena) |
| `"ya tenés una reserva en otro proyecto para ese día; no podés asistir a dos proyectos el mismo día"` | Misma fecha, proyecto diferente al de la reserva existente |

## Otras respuestas de error

| HTTP | Mensaje en `error` | Causa |
|---|---|---|
| `422` | `"el plazo para reservar ya cerró (viernes 23:59)"` | Intento de reserva fuera de plazo |
| `400` | `"esta vianda requiere que selecciones una guarnición"` | Falta `garnish_option_id` obligatorio |
| `400` | `"la guarnición seleccionada no pertenece a esta vianda"` | `garnish_option_id` de otra vianda |
| `404` | `"vianda no encontrada"` | `meal_id` inexistente |

---

## Comportamiento de swap

Si el usuario ya tiene una reserva para ese día en el **mismo proyecto y mismo tipo de comida**, la reserva anterior se **reemplaza automáticamente** por la nueva (se restaura el stock de la vianda anterior y se descuenta el de la nueva). No se devuelve error.

---

## Recomendación para el frontend

- Una vez que el usuario tiene una reserva confirmada para un día, **deshabilitar o filtrar los menús de otros proyectos** para esa misma fecha.
- Al listar menús disponibles (`GET /meals?date=...`), se puede pasar `project_id` como filtro para mostrar solo los del proyecto en cuestión.
