# Yodel

[![Docker Automated Build](https://img.shields.io/docker/cloud/automated/amfl/yodel)](https://hub.docker.com/r/amfl/yodel)
[![Go Report Card](https://goreportcard.com/badge/github.com/amfl/yodel)](https://goreportcard.com/report/github.com/amfl/yodel)

Tool for doing LDAP lookups and comparing actual groups with desired groups,
according to some role defined in YAML.

**Status: [In Progress]**

- Configurable via env vars or yaml config
- Can read a user's current groups from LDAP
- Can read group/role database from YAML
- Can perform basic set operations on the above two sources

## Usage

```bash
yodel LDAP_USERNAME [ROLE_NAME]
```

- If `ROLE_NAME` is omitted, yodel will output the user's current LDAP groups.
- If `ROLE_NAME` is provided, yodel will output the groups the user is missing
  from LDAP in order to have the role.
