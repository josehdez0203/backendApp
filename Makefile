APP=fullstackApp
DB=${APP}
RED=red_${APP}
PASSWORD=postgres
USER=postgres
PUERTO=5432
SERVER=localhost
VERSION=postgres:14.5

DSN=${USER}@tcp(localhost:3306)/widgets?parseTime=true&tls=false

## start_front: starts the front end
start_front: build_front
	@echo "Starting the front end..."
	@env STRIPE_KEY=${STRIPE_KEY} STRIPE_SECRET=${STRIPE_SECRET} ./dist/gostripe -port=${GOSTRIPE_PORT} -dsn="${DSN}" &
	@echo "Front end running!"

postgres:
	@echo "Iniciando BD"
	docker run --name ${APP} --network ${RED} -p ${PUERTO}:${PUERTO} -e POSTGRES_USER=${USER} -e POSTGRES_PASSWORD=${PASSWORD} -d ${VERSION}
	
createdb:
	@echo "Iniciando BD ${DB}"
	docker exec -it ${APP} createdb --username=${USER} ${DB}

dropdb:
	@echo "Eliminando BD ${DB}"
	docker exec -it ${APP} dropdb --username=${USER} ${DB}

migrate-up:
	migrate -path internal/db/migrations -database "postgres://${USER}:${PASSWORD}@${SERVER}:${PUERTO}/${DB}?sslmode=disable" -verbose up

migrate-down:
	migrate -path internal/db/migrations -database "postgres://${USER}:${PASSWORD}@${SERVER}:${PUERTO}/${DB}?sslmode=disable" -verbose down

migrate-up1:
	migrate -path internal/db/migrations -database "postgres://${USER}:${PASSWORD}@${SERVER}:${PUERTO}/${DB}?sslmode=disable" -verbose up 1

migrate-down1:
	migrate -path internal/db/migrations -database "postgres://${USER}:${PASSWORD}@${SERVER}:${PUERTO}/${DB}?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

net-create:
	@echo "Creando red ${RED}"
	docker network create ${RED}

net-drop:
	@echo "Eliminado red ${RED}"
	docker network rm ${RED}

init_db: dropdb createdb migrate-up
	@echo "Reiniciada la BD ${DB}"


.PHONY: postgres createdb dropdb migrate-up migrate-down migrate-up1 migrate-down1 sqlc test server net-create net-drop
