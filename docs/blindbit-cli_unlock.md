## blindbit-cli unlock

Unlocks the daemon

### Synopsis

This command has to be used before most commands can be used. The daemon is locked on startup with the encryption password. 
The encryption password was set during wallet creation.  


```
blindbit-cli unlock [flags]
```

### Options

```
  -h, --help   help for unlock
```

### Options inherited from parent commands

```
  -s, --socket string   Set the socket path. This is set to blindbitd default value (default "~/.blindbitd/run/blindbit.socket")
```

### SEE ALSO

* [blindbit-cli](blindbit-cli.md)	 - A cli application to interact with the blindbit daemon

###### Auto generated by spf13/cobra on 26-May-2024
