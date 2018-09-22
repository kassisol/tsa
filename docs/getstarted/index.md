---
title: "Get started with TSA"
description: "Getting started with TSA"
keywords: [ "getting started", "TSA" ]
date: "2017-02-02"
url: "/docs/getstarted/tsa/"
menu:
  docs:
    parent: getstarted
    weight: -85
---

## Install

Refer to the [page](../installation/index.md) specific to your Linux distribution.

## Start server
### Standalone

```bash
# tsad --tls --tlsgen --fqdn tsa1.example.com
⇛ https server started on [::]:443
```

### Container based

```bash
$ docker container run -d --mount "type=bind,source=/path/to/host,target=/var/lib/tsa" -p 443:443 kassisol/tsa:x.x.x --tls --tlsgen --fqdn tsa1.example.com
```

## Configure

### Initialize TSA config
#### Add the server

```bash
$ tsa server add tsa1 https://tsa1.example.com
```

#### Change the default password

Default admin password is "admin"

```bash
$ tsa passwd tsa1
```

#### Login to the instance

```bash
$ tsa access login tsa1
```

#### Initialize the instance

Answer all question, the below is an example:

```bash
$ tsa init
Country : CA
State : Quebec
City : Montreal
Organization (O) : Example
Organizational Unit (OU) : IT department
```

```bash
$ tsa info
Certificate Authority:
 Type: root
 Expire: 2027-02-02
 Country: CA
 State: Quebec
 Locality: Montreal
 Organization: Example
 Organizational Unit: IT department Certificate Authority
 Common Name: IT department Root CA
Certificates: 2
 Valid: 2
 Expired: 0
 Revoked: 0
API:
 FQDN: tsa1.example.com
 Bind Address: 0.0.0.0
 Bind Port: 443
Auth Type: none
Server Version: 0.0.0
Storage Driver: sqlite
Logging Driver: standard
TSA Root Dir: /var/lib/tsa
```

### Configure authentication
```bash
$ tsa auth ls
KEY                 VALUE
auth_type           none
```

```bash
$ tsa auth enable ldap
```

```bash
$ tsa auth ls
KEY                 VALUE
auth_type           ldap
```

```bash
$ tsa auth add auth_host ad1.example.com
$ tsa auth add auth_port 3269
$ tsa auth add auth_tls true
$ tsa auth add auth_bind_username username@example.com
$ tsa auth add auth_bind_password secret
$ tsa auth add auth_search_base_user ou=user,dc=example,dc=com
$ tsa auth add auth_search_filter "(&(objectCategory=user)(cn=%s))"
$ tsa auth add auth_attr_members memberOf
$ tsa auth add auth_group_admin cn=admindocker1,ou=group,dc=example,dc=com
$ tsa auth add auth_group_user cn=docker1,ou=group,dc=example,dc=com
```

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
auth_group_admin       cn=admindocker1,ou=group,dc=example,dc=com
auth_group_user        cn=docker1,ou=group,dc=example,dc=com
```
