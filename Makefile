VERSION=3.1.2
DESTDIR=
PREFIX=/usr
GETTEXT_LANGS=ru uk
TARGETOS=$(shell uname -s)

CGO_LDFLAGS=""
CGO_CPPFLAGS=""
ifeq ($(TARGETOS),Darwin) # CGO flags for macOS
ifneq ($(JHBUILD_PREFIX),) # CGO flags for jhbuild
CGO_LDFLAGS="-lintl -L${HOME}/gtk/inst/lib"
CGO_CPPFLAGS="-I${HOME}/gtk/inst/include"
else # CGO flags for default shell
CGO_LDFLAGS="-lintl -L/usr/local/opt/gettext/lib"
CGO_CPPFLAGS="-I/usr/local/opt/gettext/include"
endif
endif

all:
	${MAKE} insteadman
	${MAKE} insteadman-gtk

insteadman-deps:
	GO111MODULE=auto go get github.com/ghodss/yaml
	GO111MODULE=auto go get github.com/pyk/byten
	GO111MODULE=auto go get github.com/fatih/color

insteadman-gtk-deps:
	GO111MODULE=auto go get github.com/ghodss/yaml
	GO111MODULE=auto go get github.com/pyk/byten
	GO111MODULE=auto go get github.com/gotk3/gotk3/...

	CGO_LDFLAGS=${CGO_LDFLAGS} \
	CGO_CPPFLAGS=${CGO_CPPFLAGS} \
	GO111MODULE=auto go get github.com/gosexy/gettext

insteadman:
	${MAKE} insteadman-deps
	GO111MODULE=auto go build -ldflags "-s -w -X main.version=${VERSION}" -o insteadman ./cli

insteadman-gtk:
	${MAKE} insteadman-gtk-deps

	CGO_LDFLAGS=${CGO_LDFLAGS} \
    CGO_CPPFLAGS=${CGO_CPPFLAGS} \
	GO111MODULE=auto go build -ldflags "-s -w -X main.version=${VERSION}" -o insteadman-gtk ./gtk

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

	for lang in $(GETTEXT_LANGS); do \
		install -d -m 0755 $(DESTDIR)$(PREFIX)/share/locale/$$lang/LC_MESSAGES; \
		install -m 0644 resources/locale/$$lang/LC_MESSAGES/insteadman.mo $(DESTDIR)$(PREFIX)/share/locale/$$lang/LC_MESSAGES/insteadman.mo; \
	done


uninstall:
	rm $(DESTDIR)$(PREFIX)/bin/insteadman
	rm $(DESTDIR)$(PREFIX)/bin/insteadman-gtk
	rm -rf $(DESTDIR)$(PREFIX)/share/insteadman/
	rm $(DESTDIR)$(PREFIX)/share/pixmaps/insteadman.png
	rm $(DESTDIR)$(PREFIX)/share/applications/insteadman.desktop

deps-dev:
	GO111MODULE=auto go get github.com/stretchr/testify/assert
	GO111MODULE=auto go get github.com/josephspurrier/goversioninfo/cmd/goversioninfo
	GO111MODULE=auto go get github.com/gosexy/gettext/go-xgettext

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

	CGO_LDFLAGS_ALLOW=".*" \
    CGO_LDFLAGS="-lintl" \
    PKG_CONFIG_PATH=/usr/i686-w64-mingw32/lib/pkgconfig \
    CGO_ENABLED=1 \
    CC=i686-w64-mingw32-cc \
    GOOS=windows \
    GOARCH=386 \
    go install github.com/gosexy/gettext

gtk-linux2win:
	./gtk-linux2win-build.sh ./gtk insteadman-gtk ${VERSION} "${GETTEXT_LANGS}"

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

gtk-prepare-i18n:
	#intltool-extract --type=gettext/glade resources/gtk/main.glade
	#intltool-extract --type=gettext/glade resources/gtk/settings.glade

	xgettext --sort-output --keyword=translatable -o resources/locale/insteadman-glade.pot \
		resources/gtk/main.glade resources/gtk/settings.glade

	go-xgettext -o resources/locale/insteadman-code.pot --package-name=insteadman -k=i18n.T gtk/*.go gtk/ui/*.go

	msgcat resources/locale/insteadman-glade.pot resources/locale/insteadman-code.pot > resources/locale/insteadman.pot

	# Init if there aren't insteadman.po files
	# msginit -l ru -o resources/locale/ru/LC_MESSAGES/insteadman.po -i resources/locale/insteadman.pot
	# msginit -l uk -o resources/locale/uk/LC_MESSAGES/insteadman.po -i resources/locale/insteadman.pot

	# Merge if there are insteadman.po files
	msgmerge -U resources/locale/ru/LC_MESSAGES/insteadman.po resources/locale/insteadman.pot
	msgmerge -U resources/locale/uk/LC_MESSAGES/insteadman.po resources/locale/insteadman.pot

gtk-compile-i18n:
	msgfmt resources/locale/ru/LC_MESSAGES/insteadman.po -o resources/locale/ru/LC_MESSAGES/insteadman.mo
	msgfmt resources/locale/uk/LC_MESSAGES/insteadman.po -o resources/locale/uk/LC_MESSAGES/insteadman.mo

test:
	GO111MODULE=auto go test ./core/...
	GO111MODULE=auto go test ./cli/...

gtk-test:
	CGO_LDFLAGS=${CGO_LDFLAGS} \
    CGO_CPPFLAGS=${CGO_CPPFLAGS} \
	GO111MODULE=auto go test ./gtk/...

report:
	GO111MODULE=auto go vet ./{core,cli}/...
	gocyclo -over 15 .
	golint ./...

clean:
	rm -f insteadman
	rm -f insteadman-gtk
	rm -f insteadman-gtk.exe
	rm -rf build/*
