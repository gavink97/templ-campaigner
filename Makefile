.PHONY: tailwind-watch
tw:
	npx tailwindcss -o ./public/css/styles.css --watch

.PHONY: tailwind-build
tb:
	npx tailwindcss -o ./public/css/style.min.css --minify

.PHONY: temple
tg:
	templ generate

.PHONY: sharp
sharp:
	node build/sharp.js

.PHONY: test
gotest:
	go test -race -v -timeout 30s ./...

.PHONY: dev
dev:
	go build -o .tmp/templ-campaigner ./cmd/templ-campaigner/main.go \
	&& air
