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

deps-dev:
	go get github.com/josephspurrier/goversioninfo/cmd/goversioninfo

cli:
	go build -ldflags "-s -w" -o insteadman-cli ./cli

cli-cross:
	./cli-cross-build.sh ./cli insteadman-cli ${VERSION}

gtk:
	go build -ldflags "-s -w" -o insteadman-gtk ./gtk

gtk-linux64:
	./gtk-linux-build.sh ./gtk insteadman-gtk ${VERSION} amd64

gtk-linux32:
	./gtk-linux-build.sh ./gtk insteadman-gtk ${VERSION} 386

gtk-linux2win:
	./gtk-linux2win-build.sh ./gtk insteadman-gtk ${VERSION}

test:
	go test ./...

clean:
	rm -f insteadman-cli
	rm -f insteadman-gtk
	rm -f insteadman-gtk.exe
	rm -rf build/*

.PHONY: cli gtk
