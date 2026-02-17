.PHONY: init
init:
	@echo "Initializing projects..."
	cd proto && buf dep update
	cd backend && buf dep update && go mod download
	cd frontend && pnpm install

.PHONY: api
api:
	@echo "Generating API code..."
	cd backend && make api

.PHONY: dev-backend
dev-backend:
	cd backend && kratos run

.PHONY: dev-frontend
dev-frontend:
	cd frontend && pnpm run dev

.PHONY: dev
dev:
	@echo "Starting full Docker dev stack with hot reload..."
	docker-compose -f docker-compose.dev.yml up --build

.PHONY: dev-down
dev-down:
	docker-compose -f docker-compose.dev.yml down

.PHONY: build
build:
	@echo "Building projects..."
	cd backend && make build
	cd frontend && pnpm run build

.PHONY: docker-up
docker-up:
	docker-compose up -d

.PHONY: docker-down
docker-down:
	docker-compose down
