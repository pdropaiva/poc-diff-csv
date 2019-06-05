include env
export $(shell sed 's/=.*//' env)

GOENV=development

GOPATH := $(if $(GOPATH),$(GOPATH),$(HOME)/go)

poc:
	@echo
	@echo "Starting the PoC..."
	@echo
	@	GOENV=$(GOENV) \
		APP_KEY=$(APP_KEY) \
		OLD_EXPORT_ID=$(OLD_EXPORT_ID) \
		NEW_EXPORT_ID=$(NEW_EXPORT_ID) \
		EXPORT_URL=$(EXPORT_URL) \
		go run main.go

test:
	go test -cover ./...

lint:
	@$(GOPATH)/bin/golint ./... || echo 'Error: Please install `golint` before performing the setup!'
