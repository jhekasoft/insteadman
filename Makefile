VERSION=3.0.5
DESTDIR=
PREFIX=/usr

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
	go build -ldflags "-s -w -X main.version=${VERSION}" -o insteadman ./cli

cli-cross:
	./cli-cross-build.sh ./cli insteadman ${VERSION}

gtk:
	go build -ldflags "-s -w -X main.version=${VERSION}" -o insteadman-gtk ./gtk

gtk-linux64:
	./gtk-package-build.sh ./gtk insteadman-gtk ${VERSION} linux amd64

gtk-linux32:
	./gtk-package-build.sh ./gtk insteadman-gtk ${VERSION} linux 386

gtk-linux2win:
	./gtk-linux2win-build.sh ./gtk insteadman-gtk ${VERSION}

gtk-darwin64:
	./gtk-package-build.sh ./gtk insteadman-gtk ${VERSION} darwin amd64

test:
	go test ./...

clean:
	rm -f insteadman
	rm -f insteadman-gtk
	rm -f insteadman-gtk.exe
	rm -rf build/*

install: all
	install -d -m 0755 $(DESTDIR)$(PREFIX)/bin/
	install -m 0755 insteadman $(DESTDIR)$(PREFIX)/bin/insteadman
	install -m 0755 insteadman-gtk $(DESTDIR)$(PREFIX)/bin/insteadman-gtk

	install -d -m 0755 $(DESTDIR)$(PREFIX)/share/insteadman/skeleton/
	install -m 0644 skeleton/* $(DESTDIR)$(PREFIX)/share/insteadman/skeleton/

	install -d -m 0755 $(DESTDIR)$(PREFIX)/share/insteadman/resources/gtk/
	install -d -m 0755 $(DESTDIR)$(PREFIX)/share/insteadman/resources/images/
	install -m 0644 resources/gtk/*.glade $(DESTDIR)$(PREFIX)/share/insteadman/resources/gtk/
	install -m 0644 resources/images/logo.png $(DESTDIR)$(PREFIX)/share/insteadman/resources/images/

	install -d -m 0755 $(DESTDIR)$(PREFIX)/share/pixmaps/
	install -d -m 0755 $(DESTDIR)$(PREFIX)/share/applications/
	install -m 0644 resources/images/logo128x128.png $(DESTDIR)$(PREFIX)/share/pixmaps/insteadman.png
	install -m 0644 resources/unix/insteadman.desktop $(DESTDIR)$(PREFIX)/share/applications/insteadman.desktop

uninstall:
	rm $(DESTDIR)$(PREFIX)/bin/insteadman
	rm $(DESTDIR)$(PREFIX)/bin/insteadman-gtk
	rm -rf $(DESTDIR)$(PREFIX)/share/insteadman/
	rm $(DESTDIR)$(PREFIX)/share/pixmaps/insteadman.png
	rm $(DESTDIR)$(PREFIX)/share/applications/insteadman.desktop


.PHONY: cli gtk
