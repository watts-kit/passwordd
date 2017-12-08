## a simpe command line client implementation
This command line client uses the library at '../passwordclib'

usage should be self-explaining
```
$ ./passwordc --help
usage: passwordc [<flags>] <command> [<args> ...]

a cli client to set/get passwords

Flags:
  --help     Show context-sensitive help (also try --help-long and --help-man).
  --version  Show application version.

Commands:
  help [<command>...]
    Show help.

  set <key>
    set a secret

  get <key>
    get a secret
```
