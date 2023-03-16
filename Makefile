POSTGRES_HOST=database
POSTGRES_PORT=5432
POSTGRES_USER=developer
POSTGRES_PASSWORD=2002
POSTGRES_DATABASE=request_db

-include .env
  
DB_URL="postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DATABASE)?sslmode=disable"

print:
  echo "$(DB_URL)"
  
start:
  go run cmd/main.go

migrateup:
  migrate -path migrations -database "$(DB_URL)" -verbose up

migrateup1:
  migrate -path migrations -database "$(DB_URL)" -verbose up 1

migratedown:
  migrate -path migrations -database "$(DB_URL)" -verbose down

migratedown1:
  migrate -path migrations -database "$(DB_URL)" -verbose down 1

.PHONY: start migrateup migratedown