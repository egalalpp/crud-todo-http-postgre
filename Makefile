include .env
export


postgre-run:
	sudo systemctl start postgresql

postgre-stop:
	sudo systemctl stop postgresql

migration-up:
	export $(grep -v '^#' .env | xargs) && migrate -path ./migrations -database "postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=$DB_SSLMODE" up 1

migration-down:
	export $(grep -v '^#' .env | xargs) && migrate -path ./migrations -database "postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=$DB_SSLMODE" down 1

start-service:
	go run cmd/main.go