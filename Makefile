dist_dir=$(shell pwd)/dist

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: clean

clean:
	rm -rf $(dist_dir)

generate:
	@go generate ./...
	@echo "[OK] Files added to embed box!"

run:
	@./dist/hello-world_darwin_amd64/hello-world

snapshot: lint
	goreleaser --rm-dist --snapshot

# release
release: lint
	goreleaser --rm-dist