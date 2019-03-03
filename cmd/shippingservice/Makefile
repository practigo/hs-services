service = shipping
target = $(service).dev

all: dev

$(target): *.go
	go build -o $(target)

dev: $(target)
	echo "Start $(service) service ..."
	./$(target)