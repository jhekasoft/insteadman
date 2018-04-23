VERSION=3.0.2

all:
	${MAKE} deps
	${MAKE} clean
	${MAKE} cli gtk

deps:
	go get github.com/ghodss/yaml
	go get github.com/pyk/byten
	go get github.com/gotk3/gotk3/...
	go get github.com/stretchr/testify/assert

cli:
	go build -ldflags "-s -w" -o insteadman-cli ./cli

clicross:
	./clicrossbuild.sh ./cli insteadman-cli ${VERSION}

gtk:
	go build -ldflags "-s -w" -o insteadman-gtk ./gtk

gtkcrosswin32:
	CGO_ENABLED=1 \
	CC=i686-w64-mingw32-cc \
	GOOS=windows \
	GOARCH=386 \
	go build -ldflags "-H=windowsgui -s -w" -o insteadman-gtk.exe ./gtk

test:
	go test ./...

clean:
	rm -f insteadman-cli
	rm -f insteadman-gtk
	rm -f insteadman-gtk.exe
	rm -rf build/*

.PHONY: cli gtk
