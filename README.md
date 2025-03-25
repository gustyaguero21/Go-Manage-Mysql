# 🚀 Proyecto en Go: Autenticación y Persistencia con MySQL

Este proyecto es una API desarrollada en Go que implementa autenticación con JWT y persistencia de datos utilizando MySQL y GORM.

## 📌 Tecnologías utilizadas
- **Go**: Lenguaje principal
- **GORM**: ORM para manejar MySQL
- **JWT**: Para autenticación segura
- **Variables de entorno**: Configuración segura de credenciales

## 🔧 Configuración

Antes de ejecutar el proyecto, configura tus variables de entorno en un archivo `.env`:

```env
DB_USER=tu_usuario
DB_PASSWORD=tu_contraseña
DB_HOST=localhost
DB_PORT=3306
DB_NAME=tu_base_de_datos

TOKEN=tu_token_secreto
TOKEN_VALID_TIME=1 //tiempo expresado en horas
```

## ▶️ Ejecución

1. Instala las dependencias:
   ```sh
   go mod tidy
   ```
2. Ejecuta la aplicación:
   ```sh
   go run cmd/api/main.go
   ```

## 📌 Funcionalidades
✅ Registro y autenticación de usuarios  
✅ Generación y validación de tokens JWT  
✅ CRUD de usuarios con GORM y MySQL  
✅ Manejo de configuración con variables de entorno  

---

💡 **Contribuciones y feedback son bienvenidos**. Si quieres probarlo o mejorarlo, ¡hablemos! 🚀

