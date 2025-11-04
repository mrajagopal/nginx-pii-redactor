# nginx-pii-redactor

This is a tool to redact NGINX configuration using the nginx-go-crossplane library.

## Building
```
make build
```

## Usage
The `build` produces a binary, `pii-redactor`, which takes an NGINX configuraiton file as input: `./pii-redactor <nginx-conf-file>
```
./pii-redactor nginx.conf
```

Or use make run that takes a sample configuration file:
```
make run
```

Use `make clean` to tidy up:
```
make clean
```