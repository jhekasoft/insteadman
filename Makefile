DESTDIR=
PREFIX=/usr
GO_VERSION := 1.15

all:
	${MAKE} insteadman
	${MAKE} insteadman-fyne

insteadman:
	go build -ldflags "-s -w" -o ./build/insteadman ./cmd/insteadman

insteadman-fyne:
	go build -ldflags "-s -w" -o ./build/insteadman-fyne ./cmd/insteadman-fyne

test:
	go test ./core/...
	go test ./cmd/...

clean:
	rm -f insteadman
	rm -f insteadman.exe
	rm -f insteadman-fyne
	rm -f insteadman-fyne.exe
	rm -rf build/*

insteadman-win:
	@echo "==> Building App in MinGW container..." && \
	docker run --rm -it \
		-v "$(PWD)":/tmp/build \
		x1unix/go-mingw:$(GO_VERSION) \
		/bin/sh -c "cd /tmp/build && go build -ldflags \"-s -w\" -o ./build/insteadman.exe ./cmd/insteadman"

insteadman-fyne-win:
	@echo "==> Building App in MinGW container..." && \
	docker run --rm -it \
		-v "$(PWD)":/tmp/build \
		x1unix/go-mingw:$(GO_VERSION) \
		/bin/sh -c "cd /tmp/build && go build \
		-ldflags \"-H=windowsgui -s -w\" -o ./build/insteadman-fyne.exe ./cmd/insteadman-fyne"

insteadman-fyne-win-setup:
	docker run --rm -i -v $(PWD)/build/windows:/work amake/innosetup setup.iss
