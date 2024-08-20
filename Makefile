AMOUNT?=5
PORT?=8000
all:
	@curl -s http://localhost:8000/test/$(AMOUNT) | jq
1:
	@curl -s http://localhost:8000/test/1 | jq
5:
	@curl -s http://localhost:8000/test/5 | jq
10:
	@curl -s http://localhost:8000/test/10 | jq
100:
	@curl -s http://localhost:8000/test/100 | jq

build:
	sudo docker build -t definitlynotsimon/test-endpoint .
run:
	sudo docker run -p $(PORT):8080 definitlynotsimon/test-endpoint

run_local:
	@go run api/main.go

push_image:
	sudo docker push definitlynotsimon/test-endpoint
