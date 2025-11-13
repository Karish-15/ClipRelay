COMPOSE_FILE=docker-compose.local.yml
FE_PORT=5173
FE_DIR=client

.PHONY: cliprelay-start cliprelay-end

cliprelay-start:
	@echo "Starting ClipRelay backend..."
	docker compose -f $(COMPOSE_FILE) up --build -d
	@echo "Starting ClipRelay frontend (from $(FE_DIR))..."
	cd $(FE_DIR) && python3 -m http.server $(FE_PORT) & echo $$! > ../fe.pid
	@echo "Frontend running at http://localhost:$(FE_PORT)"

cliprelay-end:
	@echo "Stopping ClipRelay backend..."
	docker compose -f $(COMPOSE_FILE) down
	@echo "Stopping frontend server..."
	@bash -c 'if [ -f fe.pid ]; then kill $$(cat fe.pid) && rm fe.pid; fi'
	@echo "ClipRelay stopped."
