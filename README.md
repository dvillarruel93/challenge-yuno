Yuno - Go API challenge

EL proyecto ha sido construído en [GO](https://github.com/golang/go) en la versión 1.23, usando el framework [Echo Context](https://github.com/labstack/echo).

Fue hecho siguiendo la arquitectura de [Domain-Driven Design](https://en.wikipedia.org/wiki/Domain-driven_design) debido a su simplicidad y fácil mantenimiento.  
El proyecto está dividido en 4 capas:
* Handler: Primera capa donde los requests llegan y se validan según corresponda, todo está ubicado en el directorio `cmd/api/v1`.
* Usecase: Segunda capa donde está la lógica de negocio. Ubicada en el directorio `internal/business/usecases`.
* Repository: Tercera capa encargada de la conexión y comunicación con la base de datos. Ubicada en el directorio `internal/platform/repositories`.
* Model: Cuarta capa donde se encuentran definidos los modelos de negocio. Ubicada en `internal/business/domain`.
* Service: Quinta capa donde se encuentran los servicios con los que vamos a comunicarnos. Ubicada en `internal/services`.

### Validator

El proyecto fue diseñado para validar cualquier estructura con [Go Package Validator](https://pkg.go.dev/github.com/go-playground/validator/v10)
que básicamente implementa validaciones de valores para estructuras y campos individuales basándose en etiquetas.  
Las etiquetas están definidas en la capa del Model  

### Tests

Al proyecto se le agregaron tests unitarios en la capa del handler y en la capa del repository (debido a problemas de tiempo, sólo se agregó en la parte del kvs)
Los mocks fueron generados con [Mockery](https://vektra.github.io/mockery/latest/), una herramienta que facilita la creación de las funciones mockeadas según las definiciones que existan en las interfaces.

### Error handler

Los errores son manejados con el mismo framework [Echo Context web framework](https://github.com/labstack/echo) siguiendo su propia estructura [error structure](https://echo.labstack.com/docs/error-handling).  

## How to

El proyecto está listo para ser levantado en un docker y correr la DB en un postgreSQL, para hacerlo simplemente hay que ejecutar el comando:
```
make compose-up
```
y luego de dejar todo listo en el contenedor, se verá el siguiente mensaje
```
http server started on [::]:8080
```
que indicará que la api está lista para ser usada.  

Todos los endpoints están definidos en la siguiente [Colección de Postman](https://www.postman.com/dvillarruel/workspace/accelone/collection/4793868-dd52f49f-4b07-447b-b37a-c611ba3799e7?action=share&creator=4793868)
