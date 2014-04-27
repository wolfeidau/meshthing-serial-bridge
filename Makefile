prefix := "/usr/local"

all: clean
	scripts/build.sh

install:
	install -m 0755 sniffer-bridge $(prefix)/bin

clean:
	rm bin/sniffer-bridge &> /dev/null || true

.PHONY: all	clean install
