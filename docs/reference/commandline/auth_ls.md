---
title: "tsa auth ls"
description: "The auth ls command description and usage"
keywords: [ "initialize", "config" ]
date: "2017-02-02"
menu:
  docs:
    parent: "tsa_cli_auth"
github_edit: "https://github.com/kassisol/tsa/edit/master/docs/reference/commandline/auth_ls.md"
---

```markdown
List authentication configurations

Usage:
  tsa auth ls [flags]

Aliases:
  ls, list
```

## Examples

```bash
$ tsa auth ls
KEY                    VALUE
auth_type              ldap
auth_host              ad1.example.com
auth_port              3269
auth_tls               true
auth_bind_username     username@example.com
auth_bind_password     secret
auth_search_base_user  ou=user,dc=example,dc=com
auth_search_filter     (&(objectCategory=user)(cn=%s))
auth_attr_members      memberOf
auth_group_allowed     cn=docker1,ou=group,dc=example,dc=com
```
