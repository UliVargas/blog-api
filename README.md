# Blog Go

Un proyecto de blog desarrollado en Go utilizando arquitectura limpia y mejores prácticas de desarrollo.

## 🏗️ Estructura del Proyecto

Este proyecto sigue una arquitectura limpia (Clean Architecture) organizando el código en capas bien definidas:

```
blog-go/
├── cmd/                    # Puntos de entrada de la aplicación
│   └── api/               # Servidor API REST
│       └── main.go        # Punto de entrada principal
├── internal/              # Código interno de la aplicación
│   ├── config/           # Configuración de la aplicación
│   ├── handlers/         # Controladores HTTP (capa de presentación)
│   ├── middleware/       # Middleware HTTP (autenticación, CORS, etc.)
│   ├── models/          # Modelos de datos y estructuras
│   ├── repository/      # Capa de acceso a datos (implementaciones)
│   └── service/         # Lógica de negocio (casos de uso)
├── migrations/           # Scripts de migración de base de datos
├── pkg/                 # Paquetes reutilizables y utilidades
├── .env.example         # Ejemplo de variables de entorno
├── docker-compose.yml   # Configuración de Docker para desarrollo
├── go.mod              # Dependencias del módulo Go
└── go.sum              # Checksums de dependencias
```

## 📁 Descripción de Carpetas

### `/cmd/api/`

Contiene el punto de entrada principal de la aplicación. Aquí se inicializa el servidor HTTP, se configuran las rutas y se conectan todas las capas de la aplicación.

### `/internal/`

Código interno de la aplicación que no debe ser importado por otros proyectos:

- **`config/`**: Manejo de configuración de la aplicación (variables de entorno, configuración de base de datos, etc.)
- **`handlers/`**: Controladores HTTP que manejan las peticiones y respuestas. Actúan como la capa de presentación
- **`middleware/`**: Funciones middleware para autenticación, autorización, CORS, logging, etc.
- **`models/`**: Definición de estructuras de datos, modelos de dominio y DTOs
- **`repository/`**: Implementaciones de acceso a datos (interfaces con la base de datos)
- **`service/`**: Lógica de negocio y casos de uso de la aplicación

### `/migrations/`

Scripts SQL para migración y versionado de la base de datos. Permite mantener un historial de cambios en el esquema de la base de datos.

### `/pkg/`

Paquetes y utilidades que pueden ser reutilizados por otros proyectos. Código que no es específico de este dominio de negocio.

## 🛠️ Tecnologías Utilizadas

Basado en el análisis de `go.mod`, el proyecto utiliza:

- **Framework Web**: Gin (github.com/gin-gonic/gin)
- **Base de Datos**: PostgreSQL con GORM (gorm.io/driver/postgres, gorm.io/gorm)
- **Autenticación**: JWT (implícito por JWT_SECRET en .env.example)
- **CORS**: github.com/rs/cors
- **Configuración**: github.com/joho/godotenv
- **UUID**: github.com/google/uuid
- **Validación**: github.com/go-playground/validator/v10

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

## 🚀 Arquitectura

Este proyecto implementa los principios de **Clean Architecture**:

1. **Capa de Presentación** (`handlers/`, `middleware/`): Maneja HTTP y comunicación externa
2. **Capa de Aplicación** (`service/`): Contiene la lógica de negocio
3. **Capa de Infraestructura** (`repository/`, `config/`): Acceso a datos y configuración
4. **Capa de Dominio** (`models/`): Entidades y reglas de negocio

Esta estructura permite:

- **Separación de responsabilidades**
- **Facilidad de testing**
- **Mantenibilidad**
- **Escalabilidad**
- **Independencia de frameworks**

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

## 🧪 Testing

```bash
# Ejecutar todos los tests
go test ./...

# Ejecutar tests con cobertura
go test -cover ./...

# Ejecutar tests de un paquete específico
go test ./internal/service/
```

## 📦 Build

```bash
# Build para desarrollo
go build -o bin/api cmd/api/main.go

# Build para producción
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/api cmd/api/main.go
```

## 🔒 Seguridad

- **Variables de entorno**: El archivo `.env` está incluido en `.gitignore` para evitar exponer credenciales
- **JWT**: Usa claves secretas fuertes en producción
- **CORS**: Configurado para permitir orígenes específicos
- **Validación**: Todas las entradas son validadas usando `validator/v10`

## 📝 Estado del Proyecto

### ✅ Completado

- [x] Estructura de proyecto con Clean Architecture
- [x] Configuración de dependencias (go.mod)
- [x] Configuración de Docker para desarrollo
- [x] Variables de entorno y configuración
- [x] Gitignore configurado

### 🚧 En Desarrollo

- [ ] Implementación de handlers HTTP
- [ ] Modelos de datos y entidades
- [ ] Repositorios y acceso a datos
- [ ] Servicios de negocio
- [ ] Middleware de autenticación
- [ ] Migraciones de base de datos
- [ ] Tests unitarios
- [ ] Documentación de API

7. Implementar el punto de entrada en `/cmd/api/main.go`
