
GIT_COMMIT = $(shell git rev-parse --short HEAD)
GIT_TAG = $(shell git describe --tags --abbrev=0 2>/dev/null || echo "notag")

SOURCE_FILES = $(shell find cmd -name '*.go')

OS = $(shell uname -s | tr '[:upper:]' '[:lower:]')
ARCH = $(shell uname -m | sed 's/x86_64/amd64/')


LDFLAGS = -s -w -X main.VCSCommit=$(GIT_COMMIT) -X main.VCSTag=$(GIT_TAG)


all: $(SOURCE_FILES)
	GOOS=$(OS) GOARCH=$(ARCH) go build \
	  -o bin/priv8-$(OS)-$(ARCH) \
	  -ldflags "$(LDFLAGS) -X main.Executable=priv8-$(OS)-$(ARCH)" \
	  $^
	ln -sf $(CURDIR)/bin/priv8-$(OS)-$(ARCH) bin/priv8

clean:
	rm -rf ./build ./bin

