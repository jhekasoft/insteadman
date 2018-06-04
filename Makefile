VERSION=3.1.1
DESTDIR=
PREFIX=/usr

all:
	${MAKE} insteadman
	${MAKE} insteadman-gtk

insteadman-deps:
	go get github.com/ghodss/yaml
	go get github.com/pyk/byten

insteadman-gtk-deps:
	go get github.com/ghodss/yaml
	go get github.com/pyk/byten
	go get github.com/gotk3/gotk3/...

insteadman:
	${MAKE} insteadman-deps
	go build -ldflags "-s -w -X main.version=${VERSION}" -o insteadman ./cli

insteadman-gtk:
	${MAKE} insteadman-gtk-deps
	go build -ldflags "-s -w -X main.version=${VERSION}" -o insteadman-gtk ./gtk

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

deps-dev:
	go get github.com/stretchr/testify/assert
	go get github.com/josephspurrier/goversioninfo/cmd/goversioninfo

insteadman-cross: insteadman-deps
	./cli-cross-build.sh ./cli insteadman ${VERSION}

gtk-linux64: insteadman-gtk-deps
	./gtk-package-build.sh ./gtk insteadman-gtk ${VERSION} linux amd64

gtk-linux32: insteadman-gtk-deps
	./gtk-package-build.sh ./gtk insteadman-gtk ${VERSION} linux 386

gtk-linux2win-deps:
	CGO_LDFLAGS_ALLOW=".*" \
    PKG_CONFIG_PATH=/usr/i686-w64-mingw32/lib/pkgconfig \
    CGO_ENABLED=1 \
    CC=i686-w64-mingw32-cc \
    GOOS=windows \
    GOARCH=386 \
    go install github.com/gotk3/gotk3/gtk

gtk-linux2win:
	./gtk-linux2win-build.sh ./gtk insteadman-gtk ${VERSION}

gtk-darwin64:
	./gtk-package-build.sh ./gtk insteadman-gtk ${VERSION} darwin amd64

gtk-darwin-bundle: # build it from 'jhbuild shell'
	${MAKE} PREFIX="${JHBUILD_PREFIX}" install
	./gtk-darwin-bundle-prepare.sh ./gtk insteadman-gtk ${VERSION}
	gtk-mac-bundler ./resources/darwin/bundle-gtk/insteadman.bundle

	# Create DMG
	#test -f "./build/InsteadMan-${VERSION}.dmg" && rm "./build/InsteadMan-${VERSION}.dmg"
	#${HOME}/app/create-dmg/create-dmg \
	create-dmg \
    --volname "InsteadMan ${VERSION}" \
    --volicon "./resources/images/logo.icns" \
    --background "./resources/darwin/bundle-gtk/dmg_background.png" \
    --window-pos 200 120 \
	--window-size 508 337 \
	--icon-size 64 \
	--icon "InsteadMan.app" 114 200 \
	--hide-extension "InsteadMan.app" \
	--app-drop-link 390 200 \
	"./build/InsteadMan-${VERSION}.dmg" \
	"./build/InsteadMan.app"

test:
	go test ./...

clean:
	rm -f insteadman
	rm -f insteadman-gtk
	rm -f insteadman-gtk.exe
	rm -rf build/*
