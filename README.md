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

![InsteadMan GUI](https://github.com/jhekasoft/insteadman3/raw/master/resources/images/cli-3_0_2-screenshot.png "InsteadMan CLI")

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

Building CLI for all platforms (binaries will be placed to the `build` directory):

```bash
make cli-cross
```

Test
----

Run tests:

```bash
make test
```
