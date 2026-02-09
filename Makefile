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
	@echo "Starting both backend and frontend..."
	make -j 2 dev-backend dev-frontend

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
