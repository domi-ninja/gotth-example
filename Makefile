# make (no args) -> dev
default: dev;

.PHONY: dev
dev: 
	air -v || go install github.com/air-verse/air@latest
	air

.PHONY: db
db:	
# install our two db tools if not present 
	sqlc version || go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	goose --version || go install github.com/pressly/goose/v3/cmd/goose@latest

# run the goose migrations for a local sqlite db in project root without changing directories
	./goose.sh up

# generate the db bindings
	rm ./internal/db_generated/* || echo "db_generated clean"
	sqlc generate

.PHONY: build
build:
	make db
#   Here is some ideas for automatically testing your prod build before shipping:
#	make vet
#	make test
	make tailwind-build
	go build -tags production -o ./bin/app ./cmd/server

# Using the v3 version of tailwindcss because v4 seems cursed.
# If you want to avoid npm, there is a blob version somewhere.
# Please run `make tailwind-watch`  in the background, separate from the server, because it does not seem to work with air and '&' backgrounding
.PHONY: tailwind-watch
tailwind-watch:
	npx tailwindcss@3.x.x -c ./frontend/tailwind.config-v3.js -i ./frontend/tailwind-input-v3.css -o ./frontend/assets/style.css --watch

.PHONY: tailwind-build
tailwind-build:
	npx tailwindcss@3.x.x -c ./frontend/tailwind.config-v3.js -i ./frontend/tailwind-input-v3.css -o ./frontend/assets/style.css --minify

# code analysis tools
.PHONY: vet
vet:
	go vet ./...
	go install honnef.co/go/tools/cmd/staticcheck@latest
	staticcheck ./...

# run tests, also test for race conditions
.PHONY: test
test:
	go test -race -v -timeout 30s ./...
