[![Build Status](https://travis-ci.org/go-goracle/goracle.svg?branch=v1)](https://travis-ci.org/go-goracle/goracle)
[![GoDoc](https://godoc.org/gopkg.in/goracle.v1?status.svg)](http://godoc.org/gopkg.in/goracle.v1)

# Warning #
*goracle won't work with Go 1.6*, I'm afraid - the change would be quite hard work,
[rana/ora](gopkg.in/rana/ora.v3) has much nicer internal structure than goracle.
(This is understandable: goracle is the port of [cx_Oracle](http://cx-oracle.sourceforge.net/html/index.html),
a Python module, and ora has been written from scratch in Go, for Go).

_So for every new project, please use gopkg.in/rana/ora.v3 !_

# goracle #
[goracle](driver.go) is a package which is a
[database/sql/driver.Driver](http://golang.org/pkg/database/sql/driver/#Driver)
compliant wrapper for goracle/oracle - passes github.com/bradfitz/go-sql-test
(as github.com/tgulacsi/go-sql-test).

*goracle/oracle* is a package is a translated version of
[cx_Oracle](http://cx-oracle.sourceforge.net/html/index.html)
(by Anthony Tuininga) converted from C (Python module) to Go.

## Versions ##
### v1 ###
The initial version, coming from `github.com/tgulacsi/goracle`.

## There ##
CHAR, VARCHAR2, NUMBER, DATETIME, INTERVAL simple AND array bind/define.
CURSOR, CLOB, BLOB

## Not working ##
Cannot input PLS_INTEGER, only INTEGER (this is OK, as PLS_INTEGER is a
PL/SQL type, not an SQL one).

## Not working ATM ##
Nothing I know of.

## Not tested (yet) ##
LONG, LONG RAW, BFILE

# Usage and intentions #
I haven't had the pressure to force me understanding database/sql - yet.
I've ported cx_Oracle because I'm using Python with Oracle most of,
and no featureful OCI binding has existed for Go that time.
Thus I'm fluent with cx_Oracle and that means goracle/oracle.

BUT I'd start and stick with database/sql as long as it is possible
- my impression is that Go's standard library is very high quality.

Of course if you need to use Oracle's non-standard features
(out bind variables, returning cursors, sending and receiving
PL/SQL associative tables...) then goracle/oracle is the straight choice.

For simple (connection, Ping, Select) usage, and testing connection
(DSN can be tricky), see [conntest](conntest/main.go).

# Changes #
With b0219c8f we can reuse statements with different number of bind variables!

# Debug #
You can build the test executable (for debugging with gdb, for example) with
`go test -c`

You can build a tracing version with the "trace" build tag
(go build -tags=trace) that will print out everything before calling OCI
C functions.

See [c](./c) for example.


# Install #
It is `go get`'able  with `go get gopkg.in/goracle.v1`
iff you have
[Oracle DB](http://www.oracle.com/technetwork/database/enterprise-edition/index.html) installed
OR the Oracle's
[InstantClient](http://www.oracle.com/technetwork/database/features/instant-client/index-097480.html)
*both* the Basic Client and the SDK (for the header files), too!
- installed

For environment variables, you can try [env](./env)

## Linux ##
AND you have set proper environment variables:

    export CGO_CFLAGS=-I$(dirname $(find $ORACLE_HOME -type f -name oci.h))
    export CGO_LDFLAGS="-L$(dirname $(find $ORACLE_HOME -type f -name libclntsh.so\*)) -lclntsh"
    go get gopkg.in/goracle.v1

For example, with my [XE](http://www.oracle.com/technetwork/products/express-edition/downloads/index.html):

    ORACLE_HOME=/u01/app/oracle/product/11.2.0/xe
    CGO_CFLAGS=-I/u01/app/oracle/product/11.2.0/xe/rdbms/public
    CGO_LDFLAGS="-L/u01/app/oracle/product/11.2.0/xe/lib -lclntsh"

With InstantClient:

    CGO_CFLAGS=-I/usr/include/oracle/11.2/client64
    CGO_LDFLAGS="-L/usr/include/oracle/11.2/client64 -lclntsh"

### RHEL 5 ###
If your git is too old, gopkg.in may present too much hops. You can do

	mkdir -p $GOPATH/src/gopkg.in/ && cd $GOPATH/src/gopkg.in && \
		git checkout https://github.com/go-errgo/errgo.git errgo.v1 && \
		cd errgo.v1 && git checkout -f v1
	mkdir -p $GOPATH/src/gopkg.in/ && cd $GOPATH/src/gopkg.in && \
		git checkout https://github.com/inconshreveable/log15.git log15.v2 && \
		cd log15.v2 && git checkout -f v2


## Mac OS X ##
For Mac OS X I did the following:

You have to get both the Instant Client Package Basic and the Instant Client Package SDK (for the header files).

Then set the env vars as this (note the SDK here was unpacked into the base directory of the Basic package)

    export CGO_CFLAGS=-I/Users/dfils/src/oracle/instantclient_11_2/sdk/include
    export CGO_LDFLAGS="-L/Users/dfils/src/oracle/instantclient_11_2 -lclntsh"
    export DYLD_LIBRARY_PATH=/Users/dfils/src/oracle/instantclient_11_2:$DYLD_LIBRARY_PATH

Perhaps this export would work too, but I did not try it.  I understand this is another way to do this

    export DYLD_FALLBACK_LIBRARY_PATH=/Users/dfils/src/oracle/instantclient_11_2

The DYLD vars are needed to run the binary, not to compile it.

## Windows 7 64-bit ##
Thanks to Johann Kropf!
### Requirements ###
  * mingw-w64
  * msys
  * go
  * Oracle InstantClient Basic 64bit
  * Oracle InstantClient SDK 64bit
  * `gopkg.in/goracle.v1` under `%GOPATH%`

Set `CGO_CFLAGS=-IC:\Oracle64Instant\sdk\include`
Set `CGO_LDFLAGS=-LC:\Oracle64Instant\sdk\lib -loci`
