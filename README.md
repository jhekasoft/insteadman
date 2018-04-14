![InsteadMan](https://github.com/jhekasoft/insteadman3/raw/master/resources/images/logo32x32.png "InsteadMan") 
InsteadMan 3
============

Stable version: https://github.com/jhekasoft/insteadman

This is new version that is writing with go lang

GUI (GTK+ 3)
---

![InsteadMan GUI](https://github.com/jhekasoft/insteadman3/raw/master/resources/images/gtk-3_0_2-screenshot.png "InsteadMan GUI (GTK)")


CLI
---

![InsteadMan GUI](https://github.com/jhekasoft/insteadman3/raw/master/resources/images/cli-3_0_1-screenshot.png "InsteadMan CLI")

Download: https://github.com/jhekasoft/insteadman3/releases.

Build
-----

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
./insteadman-cli
```

Other build variants
--------------------

If you want to build only CLI or GTK-version, then install Go dependencies:

```bash
make dep
```

Then build only CLI-version:

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
