.PHONY: up down run build keys clean seed

up: 
	docker compose up -d

down: 
	docker compose down -v

run: 
	go run cmd/main.go

build: 
	go build -o bin/server cmd/server/main.go

clean:
	rm -rf bin/ keys/	

keys:
	"C:/Program Files/OpenSSL-Win64/bin/openssl.exe" genrsa -out keys/private.pem 2048
	"C:/Program Files/OpenSSL-Win64/bin/openssl.exe" rsa -in keys/private.pem -pubout -out keys/public.pem

seed: 
	go run cmd/seed/main.go