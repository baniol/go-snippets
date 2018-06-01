generate-keys:
	./generate-keys.sh

generate:
	go run main.go generate

validate:
	go build -o jwt
	export TOKEN=$$(./jwt generate)
	echo $(TOKEN)