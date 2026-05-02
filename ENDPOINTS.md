# API Endpoints — Macabi Backend

Base URL: `http://<host>:<port>`

> **Roles**: `user` → `admin` → `super_admin` (jerárquico)

---

## Auth — sin token

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| `POST` | `/auth/login` | Login |
| `POST` | `/auth/forgot-password` | Solicitar reset de contraseña |
| `POST` | `/auth/reset-password` | Confirmar nuevo password con token |
| `POST` | `/auth/accept-invitation` | Activar cuenta desde invitación |

---

## Usuarios

> Requieren JWT en header `Authorization: Bearer <token>`

| Método | Endpoint | Roles | Descripción |
|--------|----------|-------|-------------|
| `GET` | `/api/me` | todos | Usuario autenticado actual |
| `PATCH` | `/api/me/password` | todos | Cambiar contraseña |
| `GET` | `/api/users` | admin+ | Listar usuarios (paginado) |
| `PUT` | `/api/users/:id` | admin+ | Editar usuario |
| `PATCH` | `/api/users/:id/role` | super_admin | Cambiar rol |
| `PATCH` | `/api/users/:id/status` | admin+ | Activar / desactivar usuario |
| `POST` | `/api/users/invitations` | admin+ | Invitar usuario |
| `GET` | `/api/users/invitations` | admin+ | Listar invitaciones pendientes |
| `POST` | `/api/users/invitations/:id/resend` | admin+ | Reenviar invitación |
| `DELETE` | `/api/users/invitations/:id` | admin+ | Revocar invitación |

---

## Proyectos 🆕

| Método | Endpoint | Roles | Descripción |
|--------|----------|-------|-------------|
| `GET` | `/api/projects` | todos | Listar proyectos |
| `GET` | `/api/projects/:id` | todos | Detalle de un proyecto |
| `POST` | `/api/projects` | super_admin | Crear proyecto |
| `PUT` | `/api/projects/:id` | super_admin | Editar proyecto |
| `DELETE` | `/api/projects/:id` | super_admin | Eliminar proyecto |

**Query params (GET /api/projects):** `?page=1&page_size=10`

**Body (POST / PUT):**
```json
{
  "name": "Matok",
  "description": "Proyecto para niños",
  "admin_user_id": "<uuid>"
}
```

---

## Templates de vianda

| Método | Endpoint | Roles | Descripción |
|--------|----------|-------|-------------|
| `GET` | `/api/meal-templates` | todos | Listar templates |
| `POST` | `/api/meal-templates` | admin+ | Crear template |
| `PUT` | `/api/meal-templates/:id` | admin+ | Editar template |
| `DELETE` | `/api/meal-templates/:id` | admin+ | Eliminar template |
| `POST` | `/api/meal-templates/:id/garnish-options` | admin+ | 🆕 Agregar guarnición |
| `DELETE` | `/api/meal-templates/:id/garnish-options/:garnishId` | admin+ | 🆕 Quitar guarnición |

**Body (POST garnish-options):**
```json
{ "name": "Papas" }
```

**Respuesta de template** (incluye guarniciones):
```json
{
  "id": "...",
  "title": "Pollo al horno",
  "type": "almuerzo",
  "garnish_options": [
    { "id": "...", "name": "Papas" },
    { "id": "...", "name": "Puré de calabaza" }
  ]
}
```

---

## Viandas (instancias por fecha)

| Método | Endpoint | Roles | Descripción |
|--------|----------|-------|-------------|
| `GET` | `/api/meals` | todos | Viandas disponibles del día |
| `POST` | `/api/meals` | admin+ | Crear vianda |
| `DELETE` | `/api/meals/:id` | admin+ | Eliminar vianda |

**Query params (GET /api/meals):**
- `date=2026-05-02` (requerido, formato `YYYY-MM-DD`)
- `project_id=<uuid>` 🆕 (opcional, filtra por proyecto)

**Body (POST /api/meals):** 🆕 requiere `project_id`
```json
{
  "template_id": "<uuid>",
  "project_id": "<uuid>",
  "available_count": 30,
  "date": "2026-05-02"
}
```

---

## Reservas

| Método | Endpoint | Roles | Descripción |
|--------|----------|-------|-------------|
| `POST` | `/api/bookings` | todos | Reservar vianda |
| `GET` | `/api/bookings/mine` | todos | Mis reservas |
| `DELETE` | `/api/bookings/:id` | todos | Cancelar reserva |
| `GET` | `/api/admin/bookings/daily-summary` | admin+ | Resumen del día |

**Query params (GET /api/admin/bookings/daily-summary):**
- `date=2026-05-02` (requerido)

**Body (POST /api/bookings):** 🆕 acepta `garnish_option_id`
```json
{
  "meal_id": "<uuid>",
  "garnish_option_id": "<uuid>"
}
```
> `garnish_option_id` es **obligatorio** si la vianda tiene guarniciones; debe ser uno de los IDs del template.

---

## Notas para el frontend

- Una persona **no puede reservar almuerzo y cena el mismo día** — si ya tenía reserva del mismo tipo y reserva otra del mismo tipo, se reemplaza automáticamente. Si intenta reservar de distinto tipo el mismo día, recibe `409 Conflict`.
- Las viandas solo están disponibles **sábados y domingos**.
- Los templates ahora incluyen el array `garnish_options`. Si está vacío, no se pide guarnición al reservar.
- `project_id` en `/api/meals` es el UUID del proyecto retornado por `/api/projects`.
