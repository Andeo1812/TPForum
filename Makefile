.PHONY: all clear build run-linter run-tests

all: run-linter build

PKG = ./...
LINTER_CFG = ./configs/linter.yml

run-linter:
	$(GOPATH)/bin/golangci-lint run $(PKG) --config=$(LINTER_CFG)
	go fmt $(PKG)

build:
	go build -o main.out cmd/main/main.go

# easyjson
generate:
	go generate ${PKG}

build-docker:
	docker build -t forum-tp .

run-tests:
	./technopark-dbms-forum func -u http://localhost:5000/ -r report.html

run:
	docker run  --memory 2G --log-opt max-size=5M --log-opt max-file=3 -p 80:80 -p 5432:5432 --name forum-tp -t forum-tp

run-build: build-docker run

restart-app:
	make build
	./main.out

restart-global:
	mkdir logs
	docker stop forum-tp
	make rm-docker
	make run-build

# Utils
rm-docker:
	docker rm -vf $$(docker ps -a -q) || true

clear:
	sudo rm -rf main.out *.log logs/