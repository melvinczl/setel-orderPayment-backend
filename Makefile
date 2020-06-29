.PHONY: build clean deploy

build:
	dep ensure -v
	env GOOS=linux go build -ldflags="-s -w" -o bin/createOrder order/createOrder.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/getOrder order/getOrder.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/updateOrder order/updateOrder.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/processPayment payment/processPayment.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/getPayment payment/getPayment.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose
