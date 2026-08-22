check:
	go vet ./...
	go build .

run:
	go run ./main.go

devrun:
	go run ./main.go --dev

dev:
	docker compose -f docker-compose.override.yml up --build -d
	docker compose -f docker-compose.override.yml logs -f city-pulse

prod:
	docker compose -f docker-compose.yml up --build -d

stop:
	docker compose down
