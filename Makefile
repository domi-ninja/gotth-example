# make (no args) -> dev
default: dev;

.PHONY: dev
dev: 
	air -v || go install github.com/air-verse/air@latest
	air

.PHONY: build
build:
	make tailwind-build
	go build -tags production -o ./bin/app ./cmd/server

.PHONY: vet
vet:
	go vet ./...

.PHONY: staticcheck
staticcheck:
	staticcheck ./...

.PHONY: clean
clean:
	rm -r ./bin

.PHONY: test
test:

	go test -race -v -timeout 30s ./...

.PHONY: db
db:	
# install our two db tools if not present 
	sqlc version || go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	goose --version || go install github.com/pressly/goose/v3/cmd/goose@latest

# run the goose migrations for a local sqlite db in project root without changing directories
	goose -dir ./sql/migrations sqlite3 ./data.db up

# generate the db bindings
	rm ./internal/db_generated/* || echo "db_generated clean"
	sqlc generate


# we are calling the binary version of tailwindcss which you need to have on your path

# linux: tailwindcss (the compiled version)'s watch command depends on a binary called "watchman", so install it if you get the warning  
# Then run this in the background, separate from the server, because it does not seem to work with air and '&' backgrounding
.PHONY: tailwind-watch
tailwind-watch:
	npx tailwindcss@3.x.x -c ./frontend/tailwind.config.js -i ./frontend/tailwind-input-v3.css -o ./frontend/assets/style.css --watch

.PHONY: tailwind-build
tailwind-build:
	npx tailwindcss@3.x.x -c ./frontend/tailwind.config.js -i ./frontend/tailwind-input-v3.css -o ./frontend/assets/style.css --minify

