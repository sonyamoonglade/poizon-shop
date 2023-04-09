.SILENT:
.PHONY:
.DEFAULT_GOAL := run

run:
	./scripts/run.sh

unit-test:
	go test -count=2 -short ./...

ci: unit-test
	go vet ./...
	./scripts/e2e.sh
