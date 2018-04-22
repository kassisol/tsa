---
title: "tsa init"
description: "The init command description and usage"
keywords: [ "initialize", "config" ]
date: "2017-02-02"
menu:
  docs:
    parent: "tsa_cli"
github_edit: "https://github.com/kassisol/tsa/edit/master/docs/reference/commandline/init.md"
---

```markdown
Initialize config

Usage:
  tsa init [flags]

Flags:
      --api-bind string   API Bind Interface (default "0.0.0.0")
      --api-fqdn string   API FQDN
      --api-port string   API Port (default "443")
      --city string       Locality
      --country string    Country
      --duration string   Duration (default "120")
      --email string      E-mail
      --org string        Organization
      --org-unit string   Organizational Unit
      --state string      State
```

## Examples

```bash
# tsa init
Country : Canada
State : Quebec
City : Montreal
Organization (O) : Example
Organizational Unit (OU) : IT department
Email : admin@example.com
API FQDN : tsa1.example.com
```
