all:
	${MAKE} deps
	${MAKE} clean
	${MAKE} cli gtk

deps:
	go get github.com/ghodss/yaml
	go get github.com/pyk/byten
	go get github.com/gotk3/gotk3/...
	go get github.com/stretchr/testify/assert

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
