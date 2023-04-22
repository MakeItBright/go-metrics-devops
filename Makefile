PROJECTNAME="METRICS DEOPS"
USERGROUP=`id -gn`

help: Makefile
	@echo "Choose a command run in "$(PROJECTNAME)":"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'

## run: Run LocalServer application
run:
	@CGO_ENABLED=0 go run ./cmd/server/main.go .

## build: Build  application
buildserver:
	## @mkdir -p ./build
	@CGO_ENABLED=0 go build -o ./cmd/server/server  ./cmd/server/main.go

## build: Build  application
buildagent:
	## @mkdir -p ./build
	@CGO_ENABLED=0 go build -o ./cmd/agent/agent  ./cmd/agent/main.go