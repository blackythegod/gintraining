r:
	docker run -d -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgrespw -e POSTGRES_DB=gindb -p 5436:5432  --rm --name postgres-dev postgres

d:
	docker stop gintraining
	docker stop postgres-dev
build:
	docker build -t gin-app .
