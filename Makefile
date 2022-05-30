PREFIX ?= /usr/local
BINPREFIX ?= $(PREFIX)/bin

SRC = lg/debug.go lg/dummy.go lg/logger.go \
	config.go dataset.go keys.go ui.go widget.go

VALIDATOR = node_modules/ajv-cli/dist/index.js

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

$(VALIDATOR):
	npm install ajv-cli

validate: cheatsheet $(VALIDATOR)
 # Replace $schema key with $id in generated schema.
 # This is required for validating with `ajv`.
	./cheatsheet -schema | sed 's/$$schema/$$id/' > schema.json
	for data in data/*.json; do \
		node $(VALIDATOR) -s schema.json -d $$data; \
	done

.PHONY: debug install uninstall clean validate
