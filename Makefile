GOPATH := $(shell echo $$GOPATH)

# Get the git commit
GIT_COMMIT := $(shell git rev-parse HEAD)
GIT_DIRTY := $(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)
VERSION := $(shell grep "const Version " version.go | sed -E 's/.*"(.+)"$$/\1/' )

all: clean
	scripts/build.sh

clean:
	rm bin/sniffer-bridge &> /dev/null || true

.PHONY: all	
