VERSION=3.0.3
DESTDIR=

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
	go build -ldflags "-s -w" -o insteadman ./cli

cli-cross:
	./cli-cross-build.sh ./cli insteadman ${VERSION}

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
	rm -f insteadman
	rm -f insteadman-gtk
	rm -f insteadman-gtk.exe
	rm -rf build/*

install:
	install -d -m 0755 $(DESTDIR)/usr/bin/
	install -m 0755 insteadman $(DESTDIR)/usr/bin/insteadman
	install -m 0755 insteadman-gtk $(DESTDIR)/usr/bin/insteadman-gtk

	install -d -m 0755 $(DESTDIR)/usr/share/insteadman/skeleton/
	install -m 0644 skeleton/* $(DESTDIR)/usr/share/insteadman/skeleton/

	install -d -m 0755 $(DESTDIR)/usr/share/insteadman/resources/gtk/
	install -d -m 0755 $(DESTDIR)/usr/share/insteadman/resources/images/
	install -m 0644 resources/gtk/*.glade $(DESTDIR)/usr/share/insteadman/resources/gtk/
	install -m 0644 resources/images/logo.png $(DESTDIR)/usr/share/insteadman/resources/images/

	install -d -m 0755 $(DESTDIR)/usr/share/pixmaps/
	install -d -m 0755 $(DESTDIR)/usr/share/applications/
	install -m 0644 resources/images/logo128x128.png $(DESTDIR)/usr/share/pixmaps/insteadman.png
	install -m 0644 resources/unix/insteadman.desktop $(DESTDIR)/usr/share/applications/insteadman.desktop

uninstall:
	rm $(DESTDIR)/usr/bin/insteadman
	rm $(DESTDIR)/usr/bin/insteadman-gtk
	rm -rf $(DESTDIR)/usr/share/insteadman/
	rm $(DESTDIR)/usr/share/pixmaps/insteadman.png
	rm $(DESTDIR)/usr/share/applications/insteadman.desktop


.PHONY: cli gtk
