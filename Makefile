TAG	:= 1.0.0

build:
	cd sample && env KO_DOCKER_REPO="ko.local/m-moris/dapr-cron-issue" ko build -t $(TAG) -t latest -B ./
up:
	docker-compose -f docker-compose.yaml up

down:
	docker-compose -f docker-compose.yaml down