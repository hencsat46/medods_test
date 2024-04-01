docker:
	docker-compose up
run:
	go build -o ./bin/main ./cmd/main.go; ./bin/main localhost $(secret)
docker-rm:
	sudo docker container prune; sudo docker image rm $$(sudo docker image ls -q)
image:

