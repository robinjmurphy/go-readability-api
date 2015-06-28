test:
	@go test ./readability

install: test
	@go get ./readability

.PHONY: install test
