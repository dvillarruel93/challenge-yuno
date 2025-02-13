# Variables
APP_NAME := yuno-api
DOCKER_IMAGE := $(APP_NAME)
DOCKER_COMPOSE := docker-compose.yml

# Ayuda
help:
	@echo "Comandos disponibles:"
	@echo "  build           - Compila la aplicación Go"
	@echo "  run             - Ejecuta la aplicación localmente"
	@echo "  docker-build    - Construye la imagen Docker"
	@echo "  docker-run      - Ejecuta la aplicación dentro de un contenedor Docker"
	@echo "  compose-up      - Levanta la app con Docker Compose"
	@echo "  compose-down    - Detiene y elimina los contenedores de Docker Compose"
	@echo "  clean           - Limpia archivos binarios y contenedores"

# Compilar la aplicación
build:
	@echo "Compilando la aplicación..."
	go build -o bin/$(APP_NAME) ./cmd/app/main.go

# Ejecutar la aplicación localmente
run: build
	@echo "Ejecutando la aplicación localmente..."
	./bin/$(APP_NAME)

# Construir la imagen Docker
docker-build:
	@echo "Construyendo la imagen Docker..."
	docker build -t $(DOCKER_IMAGE) .

# Ejecutar la aplicación en Docker
docker-run: docker-build
	@echo "Ejecutando la aplicación en un contenedor Docker..."
	docker run --rm -e $(DOCKER_IMAGE)

# Levantar la app con Docker Compose
compose-up:
	@echo "Levantando la aplicación con Docker Compose..."
	docker-compose -f $(DOCKER_COMPOSE) up --build

# Detener y eliminar los contenedores de Docker Compose
compose-down:
	@echo "Deteniendo y eliminando los contenedores de Docker Compose..."
	docker-compose -f $(DOCKER_COMPOSE) down

# Limpiar binarios y contenedores
clean:
	@echo "Limpiando binarios y contenedores..."
	rm -rf bin/
	docker system prune -f
