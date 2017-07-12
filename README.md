# DBWeb

DBWeb is a web based database admin tool like phpmyadmin. It' written via 
[xorm](http://github.com/go-xorm/xorm), [tango](http://github.com/lunny/tango), [nodb](http://github.com/lunny/nodb).

# Screenshot

![dbweb](screenshot.png)

# UI Languages

Now support English and 简体中文.

# Database Supports

* MySQL
* PostgreSQL
* sqlite3 : build tag -sqlite3

# Installation

```Go
go get github.com/go-xorm/dbweb
go install github.com/go-xorm/dbweb
```
# Build via make

If you want to embbed the `langs`, `public` and `templats` to the binary, use the below command.
You have to install `make` before this.

```Shell
TAGS="bindata" make generate build
```

Notice: If you want to serve via HTTPS, you still put your *.pem files on the `home` directory.

# Run

```Shell
./dbweb -home=$GOPATH/src/github.com/go-xorm/dbweb/
```

```Shell
./dbweb -help

dbweb version 0.2

  -debug=false: enable debug mode
  -help=false: show help
  -https=false: enable https
  -home=./: set the home dir which contain templates,static,langs,certs
  -port=8989: listen port
```

Then visit http://localhost:8989/

The default user is `admin` and password is also `admin`. You can change it after you logged in.

## Changelog

#### Author by arstercz
20170712

1. remove change password features.
2. add google totp to verify password.
3. limit sql to execute.
    `delete/update` sql must have `where` condition;
    `select` sql must have `where` or `limit` condition;
     disable to execute `drop/truncate <table>`, `use <database>`, `create <database/schema>`, `drop <database/schema>`, `grant/revoke ..`;
    format sql statement;
4. return error if table size is greater than 200 MB when you execute `alter table`;
5. limit the user only access the databases that in `usercfg.conf`, the `all` means you have all database privileges;

