---
title: "Installation from binary"
linkTitle: "From binary"
description: "Instructions for installing TSA as a binary. Mostly meant for hackers who want to try out TSA on a variety of environments."
keywords: [ "binary", "installation", "tsa", "documentation", "linux" ]
date: "2017-01-27"
url: "/docs/tsa/install/binary/"
menu:
  docs:
    parent: "tsa_install"
    weight: 110
github_edit: "https://github.com/kassisol/tsa/edit/master/docs/install/binary.md"
toc: true
---

TSA is composed of 2 binaries, the daemon and the client.

* `tsad`
* `tsa`

## Get the TSA binaries

To get the list of stable release version numbers from GitHub, view the `kassisol/tsa`
[releases page](https://github.com/kassisol/tsa/releases).

To download a specific release version, use the following
URL patterns:

```
https://github.com/kassisol/tsa/releases/download/x.x.x/tsad
https://github.com/kassisol/tsa/releases/download/x.x.x/tsa
```

#### Install the Linux binary

After downloading, TSA requires this binary to be installed in your host's `$PATH`.
For example, to install the binaries in `/usr/local/sbin`:

```bash
# mv tsad /usr/local/sbin/
# mv tsa /usr/local/bin/
```

> If you already have TSA installed on your host, make sure you
> stop TSA before installing (`killall tsa`), and install the binary
> in the same location. You can find the location of the current installation
> with `dirname $(which tsa)`.


## Upgrade TSA

To upgrade your manual installation of TSA on Linux, first kill the tsa
server:

```
# killall tsa
```

Then follow the [regular installation steps](#get-the-linux-binaries).

## Next steps

Continue with the [Admin Guide](../admin/index.md).
