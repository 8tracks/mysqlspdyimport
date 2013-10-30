# mysqlspdyimport

Parallel-ized, concurrent-ized, STDIN mysql import tool. Just pipe `INSERT` statements in and watch your data import with your favorite mysql tool.

## Primary Motivators

1. At 8tracks we use different databases on EC2 (mysql and postgres). We dump database/tables with SQL statements to transfer data between databases.
2. Continue processing a file even if there are bad queries like unescaped quotes or backslashes.
3. Speed up imports by using concurrency and parallelism of go.

## Install

You'll need to setup go on your machine. Checkout http://golang.org/doc/install.

Then you can `go get` it.

    go get github.com/8tracks/mysqlspdyimport

Or you can checkout the repo and build it yourself.

    $ cd $GOPROJ/src
    $ git clone https://github.com/8tracks/mysqlspdyimport.git github.com/8tracks/mysqlspdyimport
    $ cd github.com/8tracks/mysqlspdyimport
    $ go install

## Usage


