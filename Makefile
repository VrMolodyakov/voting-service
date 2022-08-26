test:
	go test ./... -cover
start:
	docker-compose up --build
stop:
	docker-compose down