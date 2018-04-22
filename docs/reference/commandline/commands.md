---
title: "TSA commands"
description: "TSA's CLI command description and usage"
tags: [ "TSA", "TSA documentation", "CLI", "command line" ]
date: "2017-02-14"
menu:
  docs:
    parent: "tsa_cli"
    weight: -200
github_edit: "https://github.com/kassisol/tsa/edit/master/docs/reference/commandline/commands.md"
toc: true
---

This section contains reference information on using TSA's command line
client. Each command has a reference page along with samples.

### TSA management commands

| Command | Description                                                                |
|:--------|:---------------------------------------------------------------------------|
| [info](info.md) | Display information about TSA                                      |
| [init](init.md) | Initialize config                                                  |
| [passwd](passwd.md) | Change admin password                                          |
| [version](version.md) | Show the TSA version information                             |

### Access commands

| Command | Description                                                                |
|:--------|:---------------------------------------------------------------------------|
| [access ls](access_ls.md) | List sessions                                            |
| [access login](access_login.md) | Get TSA access token                               |
| [access rm](access_rm.md) | Remove session                                           |
| [access status](access_status.md) | Session status                                   |
| [access unuse](access_unuse.md) | Unuse session                                      |
| [access use](access_use.md) | Use session                                            |

### Authentication commands

| Command | Description                                                                |
|:--------|:---------------------------------------------------------------------------|
| [auth add](auth_add.md) | Add auth configuration                                     |
| [auth disable](auth_disable.md) | Disable authentication                             |
| [auth enable](auth_enable.md) | Enable authentication                                |
| [auth ls](cert_ls.md) | List authentication configurations                           |
| [auth rm](cert_rm.md) | Remove authentication configuration                          |

### Certificate commands

| Command | Description                                                                |
|:--------|:---------------------------------------------------------------------------|
| [cert ls](cert_ls.md) | List certificates issued                                     |
| [cert revoke](cert_revoke.md) | Revoke certificate                                   |

### Server commands

| Command | Description                                                                |
|:--------|:---------------------------------------------------------------------------|
| [server add](server_add.md) | Add TSA server                                         |
| [server ls](server_ls.md) | List TSA servers                                         |
| [server rm](server_rm.md) | Remove TSA server                                        |
