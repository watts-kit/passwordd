## simple client code
The go code is a ready to use library:
```
// get the password for the given key
// the string will contain the password if error is nil
GetPassword(key string) (string, error)

// set the password for the given key
// password has been set if error is nil
SetPassword(key string, value string) (string, error)
```

The python code contains a simple example
