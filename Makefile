LOCAL_BIN=$(CURDIR)/bin

.PHONY: .bindeps
.bindeps:
	mkdir -p bin
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.50.1

.PHONY: lint
lint: .bindeps
	$(LOCAL_BIN)/golangci-lint run --fix

.PHONY: clean
clean:
	rm -rf $(LOCAL_BIN)
