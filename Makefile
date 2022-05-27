SRC = lg/debug.go lg/dummy.go lg/logger.go \
	config.go dataset.go keys.go ui.go widget.go

cheatsheet: main.go $(SRC)
	go build .

debug: main.go $(SRC)
	go build -tags debug .

.PHONY: debug
