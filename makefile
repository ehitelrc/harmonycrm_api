pull:
        sudo docker compose pull

up-hs:
        sudo docker compose pull && sudo docker compose up harmony_service  -d

up-messages:
        sudo docker compose pull && sudo docker compose up harmony_service_messages_in -d

up-website:
        sudo docker compose pull && sudo docker compose up harmony_web -d

all: pull up

down:
        sudo docker compose down

logs:
        sudo docker compose logs -f

restart: down up