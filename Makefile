run-server: up serve

up:
	docker compose up -d

serve:
	go run main.go start
