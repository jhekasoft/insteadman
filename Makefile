DESTDIR=
PREFIX=/usr

all:
	${MAKE} insteadman
	${MAKE} insteadman-fyne

insteadman:
	go build -ldflags "-s -w" -o insteadman ./cmd/insteadman

insteadman-fyne:
	go build -ldflags "-s -w" -o insteadman-fyne ./cmd/insteadman-fyne

test:
	go test ./core/...
	go test ./cmd/...

clean:
	rm -f insteadman
	rm -f insteadman.exe
	rm -f insteadman-fyne
	rm -f insteadman-fyne.exe
	rm -rf build/*
