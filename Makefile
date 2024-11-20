.PHONY: build run clean

build:
	podman build -t diceware .

run:
	podman run -it diceware

clean:
	podman rmi diceware
