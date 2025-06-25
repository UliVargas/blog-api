# Manejo de Errores de Base de Datos

Este documento explica cómo el sistema maneja los errores de base de datos y los convierte en mensajes amigables para el usuario.

## Función HandleDatabaseError

La función `HandleDatabaseError` en `internal/middleware/error_handler.go` detecta automáticamente diferentes tipos de errores de base de datos y los convierte en respuestas HTTP apropiadas con mensajes amigables.

### Tipos de Errores Manejados

#### 1. Violaciones de Restricciones Únicas

**Error de Base de Datos:**
```
ERROR: duplicate key value violates unique constraint "uni_users_email" (SQLSTATE 23505)
```

**Respuesta al Usuario:**
```json
{
  "error": "El email ya está registrado"
}
```
**Código HTTP:** 409 Conflict

#### 2. Violaciones de Clave Foránea

**Error de Base de Datos:**
```
ERROR: insert or update on table violates foreign key constraint
```

**Respuesta al Usuario:**
```json
{
  "error": "No se puede completar la operación debido a dependencias"
}
```
**Código HTTP:** 400 Bad Request

#### 3. Errores de Conexión

**Error de Base de Datos:**
```
connection refused
```

**Respuesta al Usuario:**
```json
{
  "error": "Error de conexión con la base de datos"
}
```
**Código HTTP:** 503 Service Unavailable

### Uso en Handlers

```go
// En lugar de usar HandleInternalError
if err := h.userService.Create(user.ToUser()); err != nil {
    middleware.HandleDatabaseError(c, err, "Error al crear el usuario")
    return
}
```

### Beneficios

1. **Mensajes Amigables**: Los usuarios reciben mensajes claros en lugar de errores técnicos
2. **Códigos HTTP Apropiados**: Cada tipo de error retorna el código HTTP correcto
3. **Centralizado**: Todo el manejo de errores de BD está en un solo lugar
4. **Extensible**: Fácil agregar nuevos tipos de errores
5. **Consistente**: Todos los handlers usan el mismo formato de respuesta

### Extensión

Para agregar nuevos tipos de errores, simplemente añade más condiciones en la función `HandleDatabaseError`:

```go
if strings.Contains(errorMsg, "nuevo_tipo_error") {
    userFriendlyMsg = "Mensaje amigable para el usuario"
    c.JSON(http.StatusAppropriate, ErrorResponse{
        Error: userFriendlyMsg,
    })
    return
}
```