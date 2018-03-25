all:
	${MAKE} clean
	${MAKE} cli gtk

cli:
	go build -ldflags "-s -w" insteadman-cli.go

clicross:
	./crossbuild-cli.sh

gtk:
	go build -ldflags "-s -w" insteadman-gtk.go

clean:
	rm -f insteadman-cli
	rm -f insteadman-gtk
	rm -rf build/*
