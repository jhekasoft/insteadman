![InsteadMan](https://github.com/jhekasoft/insteadman3/raw/master/resources/images/logo32x32.png "InsteadMan") 
InsteadMan 3
============

Stable version: https://github.com/jhekasoft/insteadman

This is new version that is writing with go lang

GUI (GTK+ 3)
---

![InsteadMan GUI](https://github.com/jhekasoft/insteadman3/raw/master/resources/images/gtk-3_0_2-screenshot "InsteadMan GUI (GTK)")


CLI
---

![InsteadMan GUI](https://github.com/jhekasoft/insteadman3/raw/master/resources/images/cli-3_0_1-screenshot.png "InsteadMan CLI")

Download: https://github.com/jhekasoft/insteadman3/releases.

Build
-----

Install `go` dependencies.

1. Required:

```bash
go get github.com/ghodss/yaml
go get github.com/pyk/byten
```

2. Only for GTK-version (non-released):

```bash
go get github.com/gotk3/gotk3/gtk
```

3. Only for tests:

```bash
go get github.com/stretchr/testify/assert
```

Then you can build CLI and GTK-versions:

```bash
make
```

Or only CLI-version:

```bash
make cli
```

Or only GTK-version:

```bash
make gtk
```

Building CLI for all platforms (binaries will be placed to the `build` directory`):

```bash
make clicross
```

Test
----

Run tests:

```bash
make test
```
