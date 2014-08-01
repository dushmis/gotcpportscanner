all:
	go fmt portscan.go
	go build portscan.go
run:
	./portscan