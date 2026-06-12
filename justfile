dev:
        wails3 dev
test:
	CGO_CFLAGS="-Wno-deprecated-declarations" go test ./...
