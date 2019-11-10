.PHONY: all clean

OUTPUT=ap

all:clean
	go build -o bin/${OUTPUT}

clean:
	rm -rf bin/${OUTPUT}