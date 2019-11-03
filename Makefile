.PHONY: lint tests mockgen

lint:
	golint -set_exit_status ./...

tests:
	go test ./... -count=1 -parallel=4 -race

generate:
	go generate

