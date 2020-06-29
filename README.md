# Yodel

[![Go Report Card](https://goreportcard.com/badge/github.com/amfl/yodel)](https://goreportcard.com/report/github.com/amfl/yodel)

Tool for doing LDAP lookups and comparing actual groups with desired groups,
according to some role defined in YAML.

**Status: [In Progress]**

- Configurable via env vars or yaml config
- Can read a user's current groups from LDAP
- Can read group/role database from YAML
- Can perform basic set operations on the above two sources
