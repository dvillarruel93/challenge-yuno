# Etapa 1: Construcción del binario
FROM golang:1.23-alpine AS builder

# Establece el directorio de trabajo
WORKDIR /app

# Copia los archivos de dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copia el código fuente
COPY . .

# Compila el proyecto apuntando al punto de entrada
RUN go build -o api ./cmd/api/main.go

# Etapa 2: Imagen final
FROM alpine:3.18

# Instalar bash (necesario para wait-for-it.sh)
RUN apk add --no-cache bash

# Instalar dependencias necesarias, incluyendo curl
RUN apk add --no-cache curl

# Copia el binario desde la etapa de compilación
COPY --from=builder /app/api /app/api

# Copia el script wait-for-it.sh al contenedor
COPY wait-for-it.sh /usr/local/bin/wait-for-it

# Establece permisos de ejecución para el script
RUN chmod +x /usr/local/bin/wait-for-it

# Establece el directorio de trabajo
WORKDIR /app

CMD ["./api"]