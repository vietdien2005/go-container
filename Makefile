# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean

build: build-process build-user build-reexec build-mount

build-process:
	GOOS=linux $(GOBUILD) -o ./dist/go_process -v ./src/process/main.go

build-user:
	GOOS=linux $(GOBUILD) -o ./dist/go_user -v ./src/user/main.go

build-reexec:
	GOOS=linux $(GOBUILD) -o ./dist/go_reexec -v ./src/reexec/main.go

build-mount:
	GOOS=linux $(GOBUILD) -o ./dist/go_mount -v ./src/mount/main.go

clean:
	$(GOCLEAN)
	rm -f ./dist/*