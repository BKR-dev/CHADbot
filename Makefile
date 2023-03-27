CC = go
LD_FLAGS = -ldflags="-s -w"
CMD = build
SRC = main.go
BIN = terminator-shitpost

all: clean debug

debug:
	$(CC) $(CMD)

release:
	$(CC) $(CMD) -o $(BIN) $(LD_FLAGS) $(SRC) 

kill:
	pkill terminator

clean:
	$(CC) clean
