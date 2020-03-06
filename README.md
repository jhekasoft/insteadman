![InsteadMan](https://github.com/jhekasoft/insteadman3/raw/master/resources/images/logo32x32.png "InsteadMan")
InsteadMan 3
============

[![goreportcard](https://goreportcard.com/badge/github.com/jhekasoft/insteadman3)](https://goreportcard.com/report/github.com/jhekasoft/insteadman3)

This is new version that is writing with go lang

Old version: <https://github.com/jhekasoft/insteadman2>

GUI (GTK+ 3)
---

![InsteadMan GUI](https://github.com/jhekasoft/insteadman3/raw/master/resources/images/gtk-3_0_2-screenshot.png "InsteadMan GUI (GTK)")

CLI
---

![InsteadMan GUI](https://github.com/jhekasoft/insteadman3/raw/master/resources/images/cli-3_0_2-screenshot.png "InsteadMan CLI")

Download: <https://github.com/jhekasoft/insteadman3/releases.>

Build
-----

Download `insteadman3`:

```bash
go get github.com/jhekasoft/insteadman3/...
```

Go to the `insteadman3` directory:

```bash
cd ~/go/src/github.com/jhekasoft/insteadman3/
```

Build CLI and GTK-versions:

```bash
make
```

Running
-------

Run GTK (GUI):

```bash
./insteadman-gtk
```

Run CLI:

```bash
./insteadman
```

Installing
----------

Install to the system path:

```bash
make install
```

Install to the destination dir:

```bash
make DESTDIR="package" install
```

Uninstalling
------------

Uninstall from the system path:

```bash
make uninstall
```

Uninstall from the destination dir:

```bash
make DESTDIR="package" uninstall
```

Other build variants
--------------------

Build only CLI-version:

```bash
make insteadman
```

Or only GTK-version:

```bash
make insteadman-gtk
```

Building CLI for all platforms (binaries will be placed to the `build` directory):

```bash
make insteadman-cross
```

Test
----

Run tests:

```bash
make test
```

Run GTK tests (with CGO):

```bash
make gtk-test
```
