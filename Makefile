all: nat66-docker

nat66-docker: docker.go ip6tables.go main.go
	go build -o nat66-docker