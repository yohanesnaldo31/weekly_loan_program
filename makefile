build-http:
	@echo "Formatting.."
	@go fmt ./...
	@echo "Building.."
	@go build -v -o weekly_loan cmd/http/main.go

run:
	@echo "Running.."
	make build-http
	@./weekly_loan