---
title: "tsa info"
description: "The info command description and usage"
keywords: [ "display", "tsa", "information" ]
date: "2017-02-02"
menu:
  docs:
    parent: "tsa_cli"
github_edit: "https://github.com/kassisol/tsa/edit/master/docs/reference/commandline/info.md"
---

```markdown
Display information about TSA

Usage:
  tsa info [flags]
```

This command displays system wide information regarding the TSA installation.
Information displayed includes the server version, number of policies, groups, clusters and collections.

# Examples

## Display TSA information

```bash
$ tsa info
Certificate Authority:
 Type: root
 Expire: 2027-02-02
 Country: Canada
 State: Quebec
 Locality: Montreal
 Organization: Example
 Organizational Unit: IT department Certificate Authority
 Common Name: IT department Root CA
API:
 FQDN: tsa1.example.com
 Bind: 0.0.0.0
 Port: 443
Auth Type: none
Server Version: 0.0.0
Storage Driver: sqlite
Logging Driver: standard
TSA Root Dir: /var/lib/tsa
```
