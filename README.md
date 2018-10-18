# kibanadmin

It's management tool for Kibana.

## How to use

```
$ make
```

## Refresh .kibana index-pattern

This command refresh targets index-patterns fields definition.  
You should reload page after execution.

Check targets index-patterns before execute

```
$ bin/kibana-refresh -i [regexp expression index-pattern] -v [kibana version] -c
```

Execute

```
$ bin/kibana-refresh -i [regexp expression index-pattern] -v [kibana version]
```

If there is basic authentication

```
$ bin/kibanadmin -i [regexp expression index-pattern] -v [kibana version] -u [username]
```

## License

MIT
