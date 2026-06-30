generate:
	wails3 generate build-assets --binaryname yamp --dir build --typescript
	wails3 generate bindings --ts

dev:
	just generate
	just frontend-install
	wails3 dev

build:
	just generate
	cd frontend && pnpm install && pnpm build
	wails3 build

test:
	CGO_CFLAGS="-Wno-deprecated-declarations" go test ./...

format-check:
	git ls-files '*.go' ':!:build/**' | xargs gofmt -w

format:
	gofmt -w .

tidy-check:
	go mod tidy
	git diff --exit-code go.mod go.sum

lint:
	golangci-lint run

frontend-install:
	cd frontend && pnpm install --frozen-lockfile

frontend-check:
	cd frontend && pnpm format:check && pnpm check && pnpm lint

frontend-dev:
	cd frontend && pnpm dev

frontend-build:
	cd frontend && pnpm build

pre-ci:
	rm -rf build
	rm -rf frontend/bindings

ci: pre-ci generate frontend-install frontend-check frontend-build test format-check tidy-check lint 

check: format-check tidy-check lint frontend-check

all: ci build

vet:
	go vet ./...

