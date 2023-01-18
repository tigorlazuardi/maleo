set positional-arguments

sync-deps:
	@GOSUMDB=off ./scripts/sync-deps.sh

docs-deploy version: docs-binary
	@mike deploy --push --update-aliases "$1" latest

docs-build version: docs-binary
	@mike deploy --update-aliases "$1" latest

docs: docs-binary
	@mkdocs serve

docs-binary:
	@if ! command -v mkdocs >/dev/null 2>&1; then \
		echo "==> [just]: Installing mkdocs"; \
		pip install mkdocs-material; \
	fi
	@if ! command -v mike >/dev/null 2>&1; then \
		echo "==> [just]: Installing mike"; \
		pip install mike; \
	fi
	@if ! test -f mkdocs.yml; then \
		echo "==> [just]: deploying first docs"; \
		mkdocs new .; \
		mike deploy --push --update-aliases 0.1 latest; \
		mike set-default --push latest; \
	fi

test:
    @go test -v -cover ./...
    @go test -v -cover ./bucket/...
    @go test -v -cover ./bucket/maleos3-v2/...
    @go test -v -cover ./bucket/maleominio-v7/...
    @go test -v -cover ./loader/...
    @go test -v -cover ./locker/...
    @go test -v -cover ./locker/maleogoredis-v8/...
    @go test -v -cover ./locker/maleogoredis-v9/...
    @go test -v -cover ./locker/maleogomemcache/...
    @go test -v -cover ./queue/...
    @go test -v -cover ./maleohttp/...
