# Blog Go

Un proyecto de blog desarrollado en Go utilizando arquitectura limpia y mejores prácticas de desarrollo.

## 🏗️ Estructura del Proyecto

Este proyecto sigue una arquitectura limpia (Clean Architecture) organizando el código en capas bien definidas:

```
blog-go/
├── cmd/                           # Puntos de entrada de la aplicación
│   └── api/                      # Servidor API REST
│       └── main.go               # Punto de entrada principal
├── internal/                     # Código interno de la aplicación (Clean Architecture)
│   ├── application/              # Capa de Aplicación
│   │   └── service/             # Lógica de negocio y casos de uso
│   ├── domain/                   # Capa de Dominio
│   │   ├── dto/                 # Data Transfer Objects
│   │   └── model/               # Entidades y modelos de dominio
│   ├── infrastructure/          # Capa de Infraestructura
│   │   ├── config/              # Configuración de la aplicación
│   │   └── repository/          # Implementaciones de acceso a datos
│   └── presentation/            # Capa de Presentación
│       ├── handler/             # Controladores HTTP
│       └── middleware/          # Middleware HTTP (autenticación, CORS, etc.)
├── pkg/                         # Paquetes reutilizables y utilidades
│   ├── errors/                  # Manejo de errores personalizados
│   └── utils/                   # Utilidades generales
├── .env.example                 # Ejemplo de variables de entorno
├── docker-compose.yml           # Configuración de Docker para desarrollo
├── go.mod                       # Dependencias del módulo Go
└── go.sum                       # Checksums de dependencias
```

## 📁 Descripción de Carpetas

### `/cmd/api/`

Contiene el punto de entrada principal de la aplicación. Aquí se inicializa el servidor HTTP, se configuran las rutas y se conectan todas las capas de la aplicación.

### `/internal/`

Código interno de la aplicación organizado siguiendo **Clean Architecture** con separación clara de responsabilidades:

#### **Capa de Presentación** (`/presentation/`)

- **`handler/`**: Controladores HTTP que manejan peticiones y respuestas REST
- **`middleware/`**: Funciones middleware para autenticación, autorización, CORS, logging, etc.

#### **Capa de Aplicación** (`/application/`)

- **`service/`**: Lógica de negocio, casos de uso y orquestación entre capas

#### **Capa de Dominio** (`/domain/`)

- **`model/`**: Entidades de dominio y reglas de negocio centrales
- **`dto/`**: Data Transfer Objects para comunicación entre capas

#### **Capa de Infraestructura** (`/infrastructure/`)

- **`repository/`**: Implementaciones de acceso a datos (interfaces con la base de datos)
- **`config/`**: Configuración de la aplicación (variables de entorno, base de datos, etc.)

### `/migrations/`

Scripts SQL para migración y versionado de la base de datos. Permite mantener un historial de cambios en el esquema de la base de datos.

### `/pkg/`

Paquetes y utilidades que pueden ser reutilizados por otros proyectos. Código que no es específico de este dominio de negocio.

## 🛠️ Stack Tecnológico

### 🎭 Capa de Presentación

- **Framework Web**: [Gin](https://gin-gonic.com/) - Framework HTTP web de alto rendimiento
- **CORS**: Configuración para permitir solicitudes cross-origin
- **Middleware**: Autenticación JWT, logging, validación
- **Validación**: [validator](https://github.com/go-playground/validator) para validación de datos

### ⚙️ Capa de Aplicación

- **Autenticación**: JWT (JSON Web Tokens) para manejo de sesiones
- **Lógica de Negocio**: Servicios para casos de uso específicos
- **Orquestación**: Coordinación entre capas

### 🏛️ Capa de Dominio

- **Entidades**: Modelos de datos centrales
- **DTOs**: Objetos de transferencia de datos
- **UUID**: [google/uuid](https://github.com/google/uuid) para identificadores únicos

### 🔧 Capa de Infraestructura

- **Base de Datos**: PostgreSQL con [GORM](https://gorm.io/) - ORM para Go
- **Configuración**: [godotenv](https://github.com/joho/godotenv) para variables de entorno
- **Repositorios**: Implementación de acceso a datos

### 🐳 Herramientas de Desarrollo

- **Contenedores**: Docker y Docker Compose para desarrollo local
- **Migraciones**: Sistema de migraciones de base de datos
- **Testing**: Framework de pruebas integrado de Go

## 🐳 Desarrollo con Docker

El proyecto incluye un `docker-compose.yml` que proporciona:

- **PostgreSQL 17 Alpine**: Base de datos principal
- **Persistencia de datos**: Volumen para mantener los datos entre reinicios
- **Puerto 5432**: Expuesto para conexión local

## ⚙️ Configuración

"

### Variables de Entorno

1. Copia el archivo de ejemplo:

```bash
cp .env.example .env
```

2. Configura las variables en `.env`:

```env
# Cadena de conexión a PostgreSQL
DBDSN="host=localhost user=postgres password=postgres dbname=blog_db port=5432 sslmode=disable"

# Clave secreta para JWT (usa una clave segura en producción)
JWTSECRET="your-super-secret-key-here"

# Puerto del servidor
PORT=":8080"
```

### Base de Datos

Puedes usar Docker para levantar PostgreSQL:

```bash
docker-compose up -d
```

O instalar PostgreSQL localmente y crear la base de datos:

```sql
CREATE DATABASE blog_db;
```

## 🏛️ Arquitectura Clean Architecture

Este proyecto implementa los principios de **Clean Architecture** de Robert C. Martin, organizando el código en capas concéntricas con dependencias que apuntan hacia adentro:

### 📊 Flujo de Dependencias

```
┌─────────────────────────────────────────────────────────┐
│                    PRESENTATION                         │
│  ┌─────────────┐              ┌─────────────────────┐   │
│  │   Handler   │              │    Middleware       │   │
│  └─────────────┘              └─────────────────────┘   │
└─────────────────────┬───────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────┐
│                   APPLICATION                           │
│  ┌─────────────────────────────────────────────────┐     │
│  │                Service                          │     │
│  └─────────────────────────────────────────────────┘     │
└─────────────────────┬───────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────┐
│                    DOMAIN                               │
│  ┌─────────────┐              ┌─────────────────────┐   │
│  │    Model    │              │        DTO          │   │
│  └─────────────┘              └─────────────────────┘   │
└─────────────────────┬───────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────┐
│                INFRASTRUCTURE                           │
│  ┌─────────────┐              ┌─────────────────────┐   │
│  │ Repository  │              │       Config        │   │
│  └─────────────┘              └─────────────────────┘   │
└─────────────────────────────────────────────────────────┘
```

### 🎯 Principios Implementados

1. **🎭 Capa de Presentación** (`presentation/`):

   - Maneja HTTP requests/responses
   - Middleware de autenticación y validación
   - **Depende de**: Application

2. **⚙️ Capa de Aplicación** (`application/`):

   - Orquesta casos de uso y lógica de negocio
   - Coordina entre Domain e Infrastructure
   - **Depende de**: Domain

3. **🏛️ Capa de Dominio** (`domain/`):

   - Entidades centrales y reglas de negocio
   - DTOs para transferencia de datos
   - **No depende de nada** (núcleo independiente)

4. **🔧 Capa de Infraestructura** (`infrastructure/`):
   - Acceso a datos y servicios externos
   - Configuración y detalles técnicos
   - **Depende de**: Domain (implementa interfaces)

### ✨ Beneficios de esta Arquitectura

- **🔒 Independencia de Frameworks**: El dominio no conoce Gin, GORM, etc.
- **🧪 Testabilidad**: Cada capa puede probarse independientemente
- **🔄 Mantenibilidad**: Cambios en una capa no afectan otras
- **📈 Escalabilidad**: Fácil agregar nuevas funcionalidades
- **🔀 Flexibilidad**: Cambiar base de datos o framework web sin afectar lógica
- **👥 Separación de Responsabilidades**: Cada capa tiene un propósito claro

## 🚀 Instalación y Ejecución

### Prerrequisitos

- Go 1.24.4 o superior
- PostgreSQL 12+ (o Docker)
- Git

### Pasos de Instalación

1. **Clona el repositorio:**

```bash
git clone https://github.com/UliVargas/blog-go.git
cd blog-go
```

2. **Instala las dependencias:**

```bash
go mod download
```

3. **Configura las variables de entorno:**

```bash
cp .env.example .env
# Edita .env con tus configuraciones
```

4. **Levanta la base de datos (con Docker):**

```bash
docker-compose up -d
```

5. **Ejecuta las migraciones:**

```bash
# Aquí irían los comandos de migración cuando estén implementados
```

6. **Ejecuta la aplicación:**

```bash
go run cmd/api/main.go
```

La API estará disponible en `http://localhost:8080`

## 📋 Comandos Útiles

### 🏗️ Desarrollo

```bash
# Ejecutar la aplicación en modo desarrollo
go run cmd/api/main.go

# Ejecutar con hot reload (usando air)
air

# Verificar dependencias
go mod tidy
go mod verify
```

### 🧪 Testing

```bash
# Ejecutar todos los tests
go test ./...

# Ejecutar tests con coverage detallado
go test -cover -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Ejecutar tests de una capa específica
go test ./internal/application/...
go test ./internal/domain/...
go test ./internal/infrastructure/...
go test ./internal/presentation/...

# Ejecutar tests con verbose output
go test -v ./...
```

### 🏗️ Build y Deploy

```bash
# Build para desarrollo
go build -o bin/api cmd/api/main.go

# Build para producción (optimizado)
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags='-w -s' -o bin/api cmd/api/main.go

# Build multiplataforma
GOOS=windows GOARCH=amd64 go build -o bin/api.exe cmd/api/main.go
GOOS=darwin GOARCH=amd64 go build -o bin/api-mac cmd/api/main.go
```

### 🗄️ Base de Datos

```bash
# Crear migración
migrate create -ext sql -dir migrations nombre_migracion

# Aplicar migraciones
migrate -path migrations -database "postgres://user:password@localhost:5432/blogdb?sslmode=disable" up

# Rollback última migración
migrate -path migrations -database "postgres://user:password@localhost:5432/blogdb?sslmode=disable" down 1

# Ver estado de migraciones
migrate -path migrations -database "postgres://user:password@localhost:5432/blogdb?sslmode=disable" version
```

### 🔍 Análisis de Código

```bash
# Linting
golangci-lint run

# Formateo de código
go fmt ./...
goimports -w .

# Análisis de vulnerabilidades
go list -json -m all | nancy sleuth

# Análisis de dependencias
go mod graph
go mod why github.com/gin-gonic/gin
```

## 🔒 Seguridad

### 🛡️ Medidas de Seguridad Implementadas

#### 🔐 Autenticación y Autorización

- **JWT Tokens**: Implementación segura con expiración configurable
- **Secretos Fuertes**: Variables de entorno para claves sensibles
- **Middleware de Auth**: Validación de tokens en rutas protegidas
- **Rotación de Tokens**: Soporte para refresh tokens

#### 🌐 Seguridad Web

- **CORS**: Configuración restrictiva para orígenes autorizados
- **Rate Limiting**: Protección contra ataques de fuerza bruta
- **Headers de Seguridad**: Implementación de headers HTTP seguros
- **HTTPS**: Configuración para producción con TLS

#### 📝 Validación de Datos

- **Input Validation**: Validación estricta en capa de presentación
- **Sanitización**: Limpieza de datos antes del procesamiento
- **Type Safety**: Uso de DTOs tipados para transferencia de datos
- **SQL Injection**: Protección mediante GORM y prepared statements

#### 🔒 Gestión de Secretos

- **Variables de Entorno**: Nunca commitear archivos `.env` con datos sensibles
- **Configuración Segura**: Separación de configuración por ambiente
- **Encriptación**: Hash seguro de contraseñas con bcrypt
- **Logs Seguros**: No logging de información sensible

### ⚠️ Mejores Prácticas

```bash
# Generar JWT secret seguro
openssl rand -base64 32

# Verificar vulnerabilidades
go list -json -m all | nancy sleuth

# Análisis de seguridad
gosec ./...
```

### 🚨 Checklist de Seguridad

- [ ] Variables de entorno configuradas correctamente
- [ ] JWT secret fuerte y único por ambiente
- [ ] CORS configurado para dominios específicos
- [ ] Rate limiting habilitado
- [ ] Logs no contienen información sensible
- [ ] Base de datos con credenciales seguras
- [ ] HTTPS habilitado en producción
- [ ] Validación de inputs en todas las rutas

## 📝 Estado del Proyecto

### ✅ Funcionalidades Completadas

#### 🏗️ Arquitectura y Estructura

- ✅ **Clean Architecture**: Implementación completa con 4 capas
- ✅ **Separación de Responsabilidades**: Cada capa con propósito específico
- ✅ **Inversión de Dependencias**: Dependencias apuntando hacia el dominio
- ✅ **Estructura de Directorios**: Organización clara y escalable

#### 🔧 Infraestructura

- ✅ **Base de Datos**: PostgreSQL con GORM
- ✅ **Migraciones**: Sistema automatizado de migraciones
- ✅ **Docker**: Configuración para desarrollo local
- ✅ **Variables de Entorno**: Configuración flexible por ambiente

#### 🎭 Capa de Presentación

- ✅ **API REST**: Endpoints con Gin framework
- ✅ **Middleware**: Autenticación, CORS, validación
- ✅ **Validación**: Input validation con validator/v10
- ✅ **Manejo de Errores**: Respuestas HTTP consistentes

#### ⚙️ Capa de Aplicación

- ✅ **Servicios de Negocio**: Lógica de casos de uso
- ✅ **Autenticación JWT**: Login y registro de usuarios
- ✅ **Gestión de Posts**: CRUD completo de artículos
- ✅ **Orquestación**: Coordinación entre capas

#### 🏛️ Capa de Dominio

- ✅ **Entidades**: User y Post models
- ✅ **DTOs**: Objetos de transferencia tipados
- ✅ **Reglas de Negocio**: Validaciones centralizadas

### 🚧 Próximas Mejoras

#### 🧪 Testing y Calidad

- 🚧 **Tests Unitarios**: Cobertura por capas
- 🚧 **Tests de Integración**: E2E testing
- 🚧 **Mocks e Interfaces**: Testabilidad mejorada
- 🚧 **CI/CD Pipeline**: Automatización de pruebas

#### 📚 Documentación y Monitoreo

- 🚧 **Swagger/OpenAPI**: Documentación automática
- 🚧 **Logging Estructurado**: Logs con contexto
- 🚧 **Métricas**: Prometheus y Grafana
- 🚧 **Health Checks**: Endpoints de salud

#### 🚀 Funcionalidades del Roadmap

##### 👥 Sistema de Usuarios Avanzado

- 🚧 **Perfiles Completos**: Bio, avatar, firstName, lastName
- 🚧 **Roles Granulares**: Admin, Author, Reader con permisos específicos
- 🚧 **Gestión de Usuarios**: CRUD completo para administradores
- 🚧 **Estados de Usuario**: Activación/desactivación de cuentas

##### 📝 Sistema de Posts Completo

- 🚧 **Estados de Post**: Draft, Published, Archived
- 🚧 **Slug Automático**: Generación desde título
- 🚧 **Excerpt**: Resúmenes automáticos (300 chars)
- 🚧 **Featured Images**: Imágenes destacadas
- 🚧 **Contador de Vistas**: Tracking de popularidad
- 🚧 **Fechas de Publicación**: Control de timing
- 🚧 **Búsqueda Avanzada**: Por título, contenido, autor

##### 🗂️ Organización de Contenido

- 🚧 **Categorías**: Sistema completo con descripción
- 🚧 **Tags Flexibles**: Etiquetado libre
- 🚧 **Relaciones M:N**: Posts con múltiples categorías/tags
- 🚧 **Navegación**: Por categoría y tag
- 🚧 **Slugs**: URLs amigables para SEO

##### 💬 Sistema de Comentarios

- 🚧 **Comentarios Anidados**: Hasta 3 niveles de respuestas
- 🚧 **Moderación**: Aprobación por administradores
- 🚧 **Edición Limitada**: 15 minutos para editar
- 🚧 **Notificaciones**: Email al autor del post
- 🚧 **Gestión**: CRUD para propietarios y admins

##### 🔐 Autenticación Avanzada

- 🚧 **Refresh Tokens**: Renovación automática de sesiones
- 🚧 **Recuperación de Contraseña**: Reset por email
- 🚧 **Validación Robusta**: Contraseñas de 8+ caracteres
- 🚧 **Logout Seguro**: Invalidación de tokens

##### 🛡️ Seguridad y Performance

- 🚧 **Rate Limiting**: Por usuario y endpoint
- 🚧 **Cache Inteligente**: Redis para respuestas frecuentes
- 🚧 **Upload Seguro**: Gestión de imágenes
- 🚧 **Audit Logs**: Trazabilidad completa
- 🚧 **Validación Estricta**: En todas las capas

##### 📊 Funcionalidades Adicionales

- 🚧 **Paginación Avanzada**: Con filtros y ordenamiento
- 🚧 **Dashboard Admin**: Panel de administración
- 🚧 **Estadísticas**: Métricas de posts y usuarios
- 🚧 **Exportación**: Backup de contenido
- 🚧 **API Versioning**: Versionado de endpoints

### 🛣️ Roadmap de Endpoints API

#### 🔐 Autenticación

```
POST   /api/auth/register        # Registro de usuarios
POST   /api/auth/login           # Inicio de sesión
POST   /api/auth/logout          # Cerrar sesión
POST   /api/auth/refresh         # Refrescar token
POST   /api/auth/forgot-password # Solicitar reset de contraseña
POST   /api/auth/reset-password  # Resetear contraseña
```

#### 👥 Usuarios

```
GET    /api/users               # Listar usuarios (admin)
GET    /api/users/:id           # Obtener usuario específico
GET    /api/users/profile       # Mi perfil actual
PUT    /api/users/profile       # Actualizar mi perfil
DELETE /api/users/:id           # Eliminar usuario (admin)
```

#### 📝 Posts

```
GET    /api/posts               # Listar posts públicos
GET    /api/posts/:slug         # Obtener post por slug
POST   /api/posts               # Crear post (author/admin)
PUT    /api/posts/:id           # Actualizar post (owner/admin)
DELETE /api/posts/:id           # Eliminar post (owner/admin)
POST   /api/posts/:id/publish   # Publicar post
POST   /api/posts/:id/unpublish # Despublicar post
```

#### 🗂️ Categorías

```
GET    /api/categories          # Listar categorías
GET    /api/categories/:slug    # Obtener categoría
GET    /api/categories/:slug/posts # Posts de una categoría
POST   /api/categories          # Crear categoría (admin)
PUT    /api/categories/:id      # Actualizar categoría (admin)
DELETE /api/categories/:id      # Eliminar categoría (admin)
```

#### 💬 Comentarios

```
GET    /api/posts/:postId/comments    # Comentarios de un post
POST   /api/posts/:postId/comments    # Crear comentario
PUT    /api/comments/:id              # Actualizar comentario (owner)
DELETE /api/comments/:id              # Eliminar comentario (owner/admin)
POST   /api/comments/:id/approve      # Aprobar comentario (admin)
```

#### 🏷️ Tags

```
GET    /api/tags                # Listar tags
GET    /api/tags/:slug/posts    # Posts con un tag
POST   /api/tags                # Crear tag
DELETE /api/tags/:id            # Eliminar tag (admin)
```

### 📋 Modelos de Datos Planificados

#### 👤 Usuario Extendido

```json
{
  "id": "uuid",
  "username": "string (único)",
  "email": "string (único)",
  "firstName": "string",
  "lastName": "string",
  "bio": "text (opcional)",
  "avatarUrl": "string (opcional)",
  "role": "enum ['admin', 'author', 'reader']",
  "isActive": "boolean",
  "createdAt": "timestamp",
  "updatedAt": "timestamp"
}
```

#### 📄 Post Completo

```json
{
  "id": "uuid",
  "title": "string",
  "slug": "string (único)",
  "content": "text",
  "excerpt": "string (max 300 chars)",
  "featuredImage": "string (opcional)",
  "status": "enum ['draft', 'published', 'archived']",
  "authorId": "uuid (FK -> User)",
  "publishedAt": "timestamp (opcional)",
  "viewCount": "integer (default 0)",
  "categories": "array de Category",
  "tags": "array de Tag",
  "createdAt": "timestamp",
  "updatedAt": "timestamp"
}
```

#### 💬 Comentario

```json
{
  "id": "uuid",
  "content": "text",
  "postId": "uuid (FK -> Post)",
  "userId": "uuid (FK -> User)",
  "parentId": "uuid (FK -> Comment, opcional)",
  "isApproved": "boolean (default false)",
  "createdAt": "timestamp",
  "updatedAt": "timestamp"
}
```

#### 🗂️ Categoría

```json
{
  "id": "uuid",
  "name": "string (único)",
  "slug": "string (único)",
  "description": "text (opcional)",
  "createdAt": "timestamp"
}
```

#### 🏷️ Tag

```json
{
  "id": "uuid",
  "name": "string (único)",
  "slug": "string (único)"
}
```

---

## 🤝 Contribución

### 📋 Guías de Contribución

1. **Fork** el repositorio
2. **Crea** una rama para tu feature (`git checkout -b feature/nueva-funcionalidad`)
3. **Sigue** los principios de Clean Architecture
4. **Escribe** tests para tu código
5. **Asegúrate** de que todos los tests pasen
6. **Commit** tus cambios (`git commit -am 'Agrega nueva funcionalidad'`)
7. **Push** a la rama (`git push origin feature/nueva-funcionalidad`)
8. **Crea** un Pull Request

### 🎯 Estándares de Código

- **Formato**: Usar `go fmt` y `goimports`
- **Linting**: Pasar `golangci-lint run`
- **Tests**: Mantener cobertura > 80%
- **Documentación**: Comentar funciones públicas
- **Commits**: Usar [Conventional Commits](https://www.conventionalcommits.org/)

### 🏗️ Arquitectura

- **Respetar** la separación de capas
- **No** crear dependencias circulares
- **Usar** interfaces para desacoplar
- **Seguir** principios SOLID
- **Mantener** el dominio independiente

## 📄 Licencia

Este proyecto está bajo la Licencia MIT. Ver el archivo [LICENSE](LICENSE) para más detalles.

## 👨‍💻 Autor

**Ulises Vargas**

- GitHub: [@UliVargas](https://github.com/UliVargas)
- Email: [uli.vargas02@gmail.com](mailto:uli.vargas02@gmail.com)

---

<div align="center">

**⭐ Si este proyecto te fue útil, considera darle una estrella ⭐**

_Desarrollado con ❤️ usando Clean Architecture y Go_

</div>
