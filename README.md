# gocks5

Very simple socks 5 proxy

# Build and start

```
cp gocks5.yml.example gocks5.yml
```

Edit gocks5.yml (change username, password, port etc.)

```
go install
go build
./gocks -c gocks5.yml
```

You can copy gocks5 to /usr/bin and config to /etc and start app:

```
gocks -c /etc/gocks5.yml -l /var/log/gocks5.log
```
