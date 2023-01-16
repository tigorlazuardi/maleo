docs: docs-binary
	@DOCKER_HOST="" docker run --rm -it -v ${PWD}:/docs squidfunk/mkdocs-material:latest build

docs-serve: docs-binary
	@DOCKER_HOST="" docker run --rm -it -p 8000:8000 -v ${PWD}:/docs squidfunk/mkdocs-material:latest

docs-binary:
	@test -f './mkdocs.yml' || DOCKER_HOST="" docker run --rm -it -v ${PWD}:/docs squidfunk/mkdocs-material:latest new .
