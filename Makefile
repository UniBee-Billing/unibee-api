ROOT_DIR    = $(shell pwd)
NAMESPACE   = "default"
DEPLOY_NAME = "template-single"
DOCKER_NAME = "template-single"

include ./hack/hack.mk

# Frontend build tasks
.PHONY: frontend.install
frontend.install:
	@echo "Installing frontend dependencies..."
	@cd embedded && npm install

.PHONY: frontend.build
frontend.build: frontend.install
	@echo "Building frontend..."
	@cd embedded && npm run build

.PHONY: frontend.clean
frontend.clean:
	@echo "Cleaning frontend build artifacts..."
	@rm -rf embedded/build
	@rm -rf resource/embedded/*

.PHONY: frontend.dev
frontend.dev:
	@echo "Starting frontend development server..."
	@cd embedded && npm start

# Update build target to include frontend
.PHONY: build.all
build.all: frontend.build build
