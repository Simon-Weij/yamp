run:
    cd tui && INK_DISABLE_DEVTOOLS=1 bun build ./src/index.tsx --outfile ./dist/app.js --target bun --external react-devtools-core
    go run .
