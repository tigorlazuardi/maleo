docs: docs-binary
	@docker run --rm -it -v ${PWD}:/docs squidfunk/mkdocs-material:latest build


docs-serve-local: docs-binary
	@DOCKER_HOST="" docker run --rm -it -p 8000:8000 -v ${PWD}:/docs squidfunk/mkdocs-material:latest

docs-serve: docs-binary
	@docker run --rm -it -p 8000:8000 -v ${PWD}:/docs squidfunk/mkdocs-material:latest

docs-binary:
	@test -f './mkdocs.yml' || DOCKER_HOST="" docker run --rm -it -v ${PWD}:/docs squidfunk/mkdocs-material:latest new .

test:
    @go test -v -cover ./...
    @go test -v -cover ./bucket/...