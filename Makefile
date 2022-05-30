PREFIX ?= /usr/local
BINPREFIX ?= $(PREFIX)/bin

SRC = lg/debug.go lg/dummy.go lg/logger.go \
	config.go dataset.go keys.go ui.go widget.go

cheatsheet: main.go $(SRC)
	go build .

debug: main.go $(SRC)
	go build -tags debug .

install:
	mkdir -p $(DESTDIR)$(BINPREFIX)
	cp -p cheatsheet $(DESTDIR)$(BINPREFIX)

uninstall:
	rm -f $(DESTDIR)$(BINPREFIX)/cheatsheet

clean:
	rm -f cheatsheet

.PHONY: debug install uninstall clean
