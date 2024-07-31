run:
	go run cmd/app/main.go

tidy:
	go mod tidy

#Create docker container
docker-up:
	docker-compose up -d

#Delete docker container
docker-down:
	docker-compose down

#Stop docker container without delete it
docker-stop:
	docker-compose stop

#Run postgres in terminal
run-db:
	docker exec -it postgres-db psql -U username -d mlvt
