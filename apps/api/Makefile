.SILENT:
.PHONY:
.DEFAULT_GOAL := run-dev

run-dev:
	./scripts/run-dev.sh

unit-test:
	go test -count=2 -short ./...

ci: unit-test
	go vet ./...
	./scripts/e2e.sh

HANDLER_SRC=internal/handler/handler.go
HANDLER_DST=internal/handler/mocks/mocks.go
mocks:
	mockgen -source ${HANDLER_SRC} -destination ${HANDLER_DST}
