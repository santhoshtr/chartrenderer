main_package_path = cmd/main.go
binary_name = chartadapter

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

.PHONY: no-dirty
no-dirty:
	@test -z "$(shell git status --porcelain)"


# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## audit: run quality control checks
.PHONY: audit
audit: test
	go mod tidy -diff
	go mod verify
	test -z "$(shell gofmt -l .)"
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

## test: run all tests
.PHONY: test
test:
	go test -v -race -buildvcs ./...

## test/cover: run all tests and display coverage
.PHONY: test/cover
test/cover:
	go test -v -race -buildvcs -coverprofile=/tmp/coverage.out ./...
	go tool cover -html=/tmp/coverage.out


# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## tidy: tidy modfiles and format .go files
.PHONY: tidy
tidy:
	go mod tidy -v
	go fmt ./...

## build: build the application
.PHONY: build
build: wasm
	@mkdir -p build
	go build -o=build/${binary_name} ${main_package_path}


## run: run the  application
.PHONY: run
run: build
	go run github.com/eliben/static-server@latest -port 8080 -silent .

## dev: run the application with reloading on file changes
.PHONY: dev
dev:
	go run github.com/cosmtrek/air@v1.43.0 \
		--build.cmd "make build" \
		--build.bin "build/${binary_name}" \
		--build.delay "100" \
		--build.exclude_dir "external,bin,tmp,node_modules,docs" \
		--build.exclude_file "assets/dist/page.css" \
		--build.include_ext "go, tpl, tmpl, html, css, scss, js, ts, sql, jpeg, jpg, gif, png, bmp, svg, webp, ico" \
		--misc.clean_on_exit "true"


# ==================================================================================== #
# OPERATIONS
# ==================================================================================== #

## push: push changes to the remote Git repository
.PHONY: push
push: confirm audit no-dirty
	git push

## wasm: build the wasm file
.PHONY: wasm
wasm: build/${binary_name}.wasm assets/wasm-exec.js

build/chartadapter.wasm:
	@mkdir -p build
	GOOS=js GOARCH=wasm tinygo build -o $@ -no-debug wasm/main.go

assets/wasm-exec.js:
	@mkdir -p build
	@cp /usr/local/lib/tinygo/targets/wasm_exec.js $@

## clean: clean up generated files
.PHONY: clean
clean:
	@rm -rf build/* assets/wasm-exec.js build/chartadapter.wasm
	@go clean -cache -testcache -modcache