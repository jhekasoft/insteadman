DESTDIR=
PREFIX=/usr
GO_VERSION := 1.15

all:
	${MAKE} insteadman
	${MAKE} insteadman-wails

insteadman:
	go build -ldflags "-s -w" -o ./build/insteadman ./cmd/insteadman

insteadman-wails:
	wails build

test:
	go test ./core/...
	go test ./cmd/...

clean:
	rm -f insteadman
	rm -f insteadman.exe
	rm -f insteadman-wails
	rm -f insteadman-wails.exe
	rm -rf build/*

insteadman-wails-win-setup:
	docker run --rm -i -v $(PWD)/build/windows:/work amake/innosetup setup.iss
