PROJECTNAME="METRICS DEOPS"
USERGROUP=`id -gn`

help: Makefile
	@echo "Choose a command run in "$(PROJECTNAME)":"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'

## run: Run LocalServer application
run:
	@CGO_ENABLED=0 go run ./cmd/server/main.go .