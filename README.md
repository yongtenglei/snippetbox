# snippetbox

## Tech used here
| item            | tech                                                                        |
|:-----------------:|:-----------------------------------------------------------------------------:|
| router          | [julienschmidt/httprouter](https://github.com/julienschmidt/httprouter)     |
| middleware     | [justinas/alice](https://github.com/justinas/alice)                         |
| form resolver   | [go-playground/form/v4](https://github.com/go-playground/form)
| cryptography    | [golang.org/x/crypto](https://pkg.go.dev/golang.org/x/crypto)
| session manager | [alexedwards/scs/v2](https://github.com/alexedwards/scs)                 |
| session storage | [alexedwards/scs/mysqlstore](https://github.com/alexedwards/scs/tree/master/mysqlstore) |
|Token-based CSRF mitagation | [justinas/nosurf](https://github.com/justinas/nosurf)|

## Setup MySQL
We need three tables (`snipptes`, `sessions` and `users`) in snippetbox database.

- `snippetbox` database
```mysql
CREATE DATABASE snippetbox CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE snippetbox;
```

- `snippets` table
```mysql
CREATE TABLE snippets (
 id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
 title VARCHAR(100) NOT NULL,
 content TEXT NOT NULL,
 created DATETIME NOT NULL,
 expires DATETIME NOT NULL
);

CREATE INDEX idx_snippets_created ON snippets(created);

-- insert things
INSERT INTO snippets (title, content, created, expires) VALUES (
 'An old silent pond',
 'An old silent pond...\nA frog jumps into the pond,\nsplash! Silence again.\n\n– Matsuo Bashō',
 UTC_TIMESTAMP(),
 DATE_ADD(UTC_TIMESTAMP(), INTERVAL 365 DAY)
);
INSERT INTO snippets (title, content, created, expires) VALUES (
 'Over the wintry forest',
 'Over the wintry\nforest, winds howl in rage\nwith no leaves to blow.\n\n– Natsume Soseki',
 UTC_TIMESTAMP(),
 DATE_ADD(UTC_TIMESTAMP(), INTERVAL 365 DAY)
);
INSERT INTO snippets (title, content, created, expires) VALUES (
 'First autumn morning',
 'First autumn morning\nthe mirror I stare into\nshows my father''s face.\n\n– Murakami Kijo',
 UTC_TIMESTAMP(),
 DATE_ADD(UTC_TIMESTAMP(), INTERVAL 7 DAY)
);
```

- `sessions` table
```mysql
CREATE TABLE sessions (
 token CHAR(43) PRIMARY KEY,
 data BLOB NOT NULL,
 expiry TIMESTAMP(6) NOT NULL
);
CREATE INDEX sessions_expiry_idx ON sessions (expiry);
```

- `users` table
```mysql
CREATE TABLE users (
 id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
 name VARCHAR(255) NOT NULL,
 email VARCHAR(255) NOT NULL,
 hashed_password CHAR(60) NOT NULL,
 created DATETIME NOT NULL
);
ALTER TABLE users ADD CONSTRAINT users_uc_email UNIQUE (email);
```

- You may want to create a dedicated user
```mysql
CREATE USER 'web'@'localhost';
GRANT SELECT, INSERT, UPDATE, DELETE ON snippetbox.* TO 'web'@'localhost';
-- Important: Make sure to swap 'pass' with a password of your own choosing.
ALTER USER 'web'@'localhost' IDENTIFIED BY 'pass';
``````

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

## Embedded file version
Typically, the program will serve and load static files from disc despite we construct a cache to speed this process. `Embedidng files` is the another way to go. You can find this version in `efs` branch.

To run that, you could say:
```bash
go build -o /tmp/web ./cmd/web/ # web is the binary name, ./cmd/web/ is where the code entry is.
cp -r ./tls /tmp # for certificate
cd /tmp/
./web
```

## Acknowledgment
Have a look [Let's GO 💖](https://lets-go.alexedwards.net/)
