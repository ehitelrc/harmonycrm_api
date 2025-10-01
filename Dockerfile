# Usa una imagen base ligera de Go
FROM golang:alpine

# Establece el directorio de trabajo en el contenedor
WORKDIR /app

# Copia los archivos del proyecto al contenedor
COPY . .


# Descarga las dependencias
RUN go mod download


# Compila el proyecto desde el directorio `cmd`
RUN go build -o main .

# Expone el puerto en el que corre la API (ajusta según tu configuración)
EXPOSE 8098

# Comando para ejecutar la aplicación
CMD ["./main"]