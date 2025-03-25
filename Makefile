install:
	@find . -name "go.mod" -exec dirname {} \; | while read dir; do \
		echo "Running go get -u ./... in $$dir"; \
		( cd $$dir && go get -u ./... ); \
	done

tidy:
	@find . -name "go.mod" -exec dirname {} \; | while read dir; do \
		echo "Running go mod tidy in $$dir"; \
		( cd $$dir && go mod tidy ); \
	done

lint:
	golangci-lint run ./...

test:
	@echo "mode: atomic" > coverage.out
	@find . -name "go.mod" -exec dirname {} \; | while read dir; do \
		echo "Running tests in $$dir"; \
		( cd $$dir && go test -count=1 -v -cover -race -covermode=atomic -coverprofile=coverage.tmp ./... ); \
		tail -n +2 $$dir/coverage.tmp >> coverage.out; \
		rm $$dir/coverage.tmp; \
	done
