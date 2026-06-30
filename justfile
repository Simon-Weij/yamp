generate-build-assets:
	wails3 generate build-assets --binaryname yamp --dir build

dev:
	just generate-build-assets
	just frontend-install
	wails3 dev

build:
	just generate-build-assets
	wails3 build

test:
	CGO_CFLAGS="-Wno-deprecated-declarations" go test ./...

format-check:
	test -z "$(gofmt -l .)"

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

ci: generate-build-assets test format-check tidy-check lint frontend-install frontend-check

check: format-check tidy-check lint frontend-check

all: ci build

vet:
	go vet ./...
