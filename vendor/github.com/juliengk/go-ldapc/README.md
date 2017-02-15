# Go LDAP Authentication Client

## Example

You can provide LDAP based authentication on your web page easily.

- Code:
```go
import (
    "github.com/juliengk/ldapc"
)

func main() {
	ldapclient := &ldapc.Client{
		Protocol:  ldapc.LDAP,
		Host:      "localhost",
		Port:      389,
		TLSConfig: nil,
		Bind: &ldapc.AuthBind{
			BindDN:       "uid=user1,ou=People,dc=test,dc=com",
			BindPassword: "admin",
			BaseDN:       "dc=test,dc=com",
			Filter:       "(&(objectClass=posixAccount)(uid=%s))",
		},
	}

	username := "user2"
	password := "user2"

	entry, err := ldapclient.Authenticate(username, password)
	if err != nil {
		fmt.Printf("LDAP Authenticate failed: %v\n", err)
	}

	// username and mail
	fmt.Printf("username: %v\n", entry.GetAttributeValue("uid"))
	fmt.Printf("mail: %v\n", entry.GetAttributeValue("mail"))
}
```

- Output:
```text
username: user2
mail: user2@test.com
```

In other cases Anonymous Bind, Direct Bind or Active Directory, example code [ldapc_test.go](./ldapc_test.go).
