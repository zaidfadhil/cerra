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

release-tags:
	@if [ -z "$(VERSION)" ]; then \
		echo "make release-tags VERSION=v0.0.0"; \
		exit 1; \
	fi
	@echo "Creating git tags for all modules with version $(VERSION)"
	@find . -name "go.mod" -exec dirname {} \; | while read dir; do \
		if [ "$$dir" = "." ]; then \
			echo "Creating main module tag: $(VERSION)"; \
			git tag -a $(VERSION) -m "$(VERSION)"; \
		else \
			module_name=$$(basename $$dir); \
			echo "Creating $$module_name module tag: $$module_name/$(VERSION)"; \
			git tag -a $$module_name/$(VERSION) -m "$$module_name/$(VERSION)"; \
		fi \
	done
	@echo "to push all tags run: git push origin --tags"
