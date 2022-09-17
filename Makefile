PREFIX ?= /usr/local
BINPREFIX ?= $(PREFIX)/bin

SRC = lg/debug.go lg/dummy.go lg/logger.go \
	config.go dataset.go keys.go ui.go

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
	rm -rf node_modules

validate: cheatsheet
	@for data in data/*.yml; do \
		./cheatsheet -validate $$data; \
	done

.PHONY: debug install uninstall clean validate
