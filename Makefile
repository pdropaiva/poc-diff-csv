include env
export $(shell sed 's/=.*//' env)

GOENV=development

GOPATH := $(if $(GOPATH),$(GOPATH),$(HOME)/go)

v1:
	@echo
	@echo "Starting the PoC v1..."
	@echo
	@	GOENV=$(GOENV) \
		APP_KEY=$(APP_KEY) \
		OLD_EXPORT_ID=$(OLD_EXPORT_ID) \
		NEW_EXPORT_ID=$(NEW_EXPORT_ID) \
		EXPORT_URL=$(EXPORT_URL) \
		go run ./cmd/v1

v2:
	@echo
	@echo "Starting the PoC v2..."
	@echo
	@	GOENV=$(GOENV) \
		APP_KEY=$(APP_KEY) \
		OLD_EXPORT_ID=$(OLD_EXPORT_ID) \
		NEW_EXPORT_ID=$(NEW_EXPORT_ID) \
		EXPORT_URL=$(EXPORT_URL) \
		go run ./cmd/v2

test:
	go test -cover ./...

lint:
	@$(GOPATH)/bin/golint ./... || echo 'Error: Please install `golint` before performing the setup!'
