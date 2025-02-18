BIN_DIR := bin
BIN := $(BIN_DIR)/ghcp

all: create_bin ghcp

.PHONY: create_bin
create_bin:
	mkdir -p $(BIN_DIR)

.PHONY: ghcp
ghcp:
	go build -o $(BIN) main.go

.PHONY: clean
clean:
	rm -rf $(BIN_DIR)
