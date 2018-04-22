---
title: "Configuring and running TSA"
description: "Configuring and running the TSA server on various distributions"
tags: [ "tsa", "server", "configuration", "running", "process managers" ]
date: "2017-02-02"
url: "/docs/tsa/admin/configure/"
menu:
  docs:
    parent: tsa_admin
    weight: 0
github_edit: "https://github.com/kassisol/tsa/edit/master/docs/admin/configure.md"
toc: true
---

## Configuration

After successfully installing TSA, the `tsa` server runs with its default
configuration.

In a production environment, system administrators typically configure the
`tsa` server to start and stop according to an organization's requirements. In most
cases, the system administrator configures a process manager such as `systemd`
to manage the `tsa` daemon's start and stop.

