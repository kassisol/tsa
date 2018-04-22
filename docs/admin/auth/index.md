---
title: "Configuring Auth Drivers"
description: "Configuring Auth Drivers"
keywords: [ "tsa", "server", "configuration", "running", "process managers" ]
date: "2017-02-02"
url: "/docs/tsa/admin/auth/"
menu:
  docs:
    parent: tsa_admin_auth
    weight: 100
github_edit: "https://github.com/kassisol/tsa/edit/master/docs/admin/configure.md"
toc: true
---

# Configuring auth drivers

TSA includes multiple authentication mechanisms to help you get . These mechanisms are called auth drivers.

Each TSA server has a default auth driver, which uses unless you configure it to use a different auth driver.

## Configure the auth driver
When you start TSA server, you can configure it to use a different auth driver than the TSA server’s default. If the auth driver has configurable options, you can set them using the command `tsa auth add`.

To find the current auth driver for TSA serverr, run the following `tsa auth` command:

```bash
# tsa auth ls | grep type
auth_type           none
```

## Supported auth drivers
The following auth drivers are supported. See each driver's section below for its configurable options, if applicable.

| Driver | Description                                                                 |
|:--------|:---------------------------------------------------------------------------|
| none | No authentication will be available for user to request certificate. The default auth driver for TSA          |
| ldap | Use ldap to authenticate.          |

## Examples

### `none`
`none` is the default auth driver, and disables authentication. It has no options.

#### Examples
This example disables authentication with the `none` driver.

```bash
tsa enable none
```

### `ldap`
The `ldap` driver

#### Options
The `ldap` auth driver supports the following auth options:

| Option | Description | Min | Max                                                        |
|:-------|:------------|:----|:-----------------------------------------------------------|
| auth_host | | 1 | 1 |
| auth_port | | 1 | 1 |
| auth_tls | | 1 | 1 |
| auth_bind_username | | 1 | 1 |
| auth_bind_password | | 1 | 1 |
| auth_search_base_user | | 1 | 1 |
| auth_search_filter | | 1 | 1 |
| auth_attr_members | | 1 | 1 |
| auth_group_admin | | 1 | 100 |
| auth_group_user | | 1 | 100 |

#### Examples
This example configure authentication with the `ldap` driver.

```bash
tsa enable ldap
```

```bash
# tsa auth add auth_host ad1.example.com
# tsa auth add auth_port 3269
# tsa auth add auth_tls true
# tsa auth add auth_bind_username username@example.com
# tsa auth add auth_bind_password secret
# tsa auth add auth_search_base_user ou=user,dc=example,dc=com
# tsa auth add auth_search_filter "(&(objectCategory=user)(cn=%s))"
# tsa auth add auth_attr_members memberOf
# tsa auth add auth_group_admin cn=admindocker1,ou=group,dc=example,dc=com
# tsa auth add auth_group_user cn=docker1,ou=group,dc=example,dc=com
```
