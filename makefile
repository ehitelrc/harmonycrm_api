pull:
	sudo docker compose pull

up:
	sudo docker compose up -d

all: pull up

down:
	sudo docker compose down

logs:
	sudo docker compose logs -f

restart: down up

build: docker buildx build --platform linux/amd64 -t ehitelrc/harmony_service:latest . --push 
