# snippetbox

## TLS enabled
This project enabled TLS by default, generate your own self-signed certificate for dev purpose.

```bash
cd $PROJECT_PATH/snippetbox
mkdir tls
cd tls

go run /usr/local/go/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost
```
Then you will get two file, `cert.pem` and `key.pem`. 

`cert.pem` is a self-signed TLS certificate for the host `localhost` containing the public key - which it stores in a cert.pem file.

`key.pem` is the file stores the private key.

### Find your command
You need to know the place on your computer where the source code for the GO standard library is installed.

The `generate_cert.go` tool, can be found under `/usr/local/go/src/crypto/tls` if you use Linux. It should be located under `/usr/local/Cellar/go/<version>/libexec/src/crypto/tls` or a similar path if you use MacOS instead and install GO via Homebrew.

