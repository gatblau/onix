# Warden HTTP Proxy

Warden is an HTTP proxy that can be used to relay traffic between secured network tiers.

In doing so, the proxy should be able to inspect traffic and execute ad-hoc logic depending on the specific use case.

The ability to perform custom inspections will be added in the future, for now the proxy support basic mode, i.e. it
only handles explicit CONNECT requests.

## Building the proxy

Using artisan as follows:

```bash
$ art run build-linux
$ cp bin/linux/amd64/warden /usr/local/bin
```

## Running the proxy

As of now only basic mode is available:

```bash
$ warden launch -va http://127.0.0.1:8080
```

## Testing the proxy

Perform an HTTP request to www.google.com using wget. 

Instruct wget to go via the proxy listening on 8080:

```bash
$ wget \
    -e use_proxy=yes \
    -e http_proxy="http://127.0.0.1:8080" \
    www.google.com
```

### Proxy output

The proxy should output the following:

```bash
... INFO: Got request / www.google.com GET http://www.google.com/
... INFO: Sending request GET http://www.google.com/
... INFO: Received response 200 OK
... INFO: Copying response to client 200 OK [200]
... INFO: Copied 14728 bytes to client error=<nil>

```
