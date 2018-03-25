all:
	${MAKE} clean
	${MAKE} cli gtk

cli:
	go build -ldflags "-s -w" -o insteadman-cli ./cli

clicross:
	./crossbuild.sh ./cli insteadman-cli

gtk:
	go build -ldflags "-s -w" -o insteadman-gtk ./gtk

test:
	go test ./...

clean:
	rm -f insteadman-cli
	rm -f insteadman-gtk
	rm -rf build/*

.PHONY: cli gtk
