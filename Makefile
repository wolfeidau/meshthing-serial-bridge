GOPATH := $(shell echo $$GOPATH)

# Get the git commit
GIT_COMMIT := $(shell git rev-parse HEAD)
GIT_DIRTY := $(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)
VERSION := $(shell grep "const Version " version.go | sed -E 's/.*"(.+)"$$/\1/' )

all: clean
	$(shell go build -ldflags "-X main.GitCommit ${GIT_COMMIT}${GIT_DIRTY}" -o sniffer-bridge)

clean:
	rm sniffer-bridge || true

dist-osx: 
	$(shell goxc -arch="amd64" -os="darwin" -d=pkg -pv="${VERSION}" -n=sniffer-bridge -build-ldflags="-X main.GitCommit ${GIT_COMMIT}${GIT_DIRTY}" go-install xc)

.PHONY: all	
