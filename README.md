# Yodel

Tool for doing LDAP lookups and comparing actual groups with desired groups,
according to some role defined in YAML.

**Status: [In Progress]**

- Configurable via env vars or yaml config
- Can read a user's current groups from LDAP
- Can read group/role database from YAML
- Can perform basic set operations on the above two sources

## Usage

Yodel will emit the result of a single expression in yodelscript.

```bash
# Print the groups for user Joe
ldap:Joe
# Print groups owned by "Joe" which "Bob" does not have
ldap:Joe - ldap:Bob
# Print groups "Joe" is missing in order to have the "Developer" role
db_role:Developer - ldap:Joe
# Print groups that "Joe" has in excess of the "Developer" role
ldap:Joe - db_role:Developer
# Make sure that Joe, Bob, and Frank together have enough access to be a FooAdmin
db_role:FooAdmin - (ldap:Joe + ldap:Bob + ldap:Frank)
```
