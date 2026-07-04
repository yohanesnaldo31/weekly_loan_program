build-http:
	@echo "Formatting.."
	@go fmt ./...
	@echo "Building.."
	@go build -v -o weekly_loan cmd/http/main.go

run:
	@echo "Running.."
	make build-http
	@./weekly_loan

build-cron: 
	@echo "Formatting.."
	@go fmt ./...
	@echo "Building.."
	@go build -v -o weekly_loan_cron cmd/cron/main.go

run-cron:
	@echo "Running.."
	make build-cron
	@./weekly_loan_cron	