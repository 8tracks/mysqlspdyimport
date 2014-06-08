# mysqlspdyimport

Parallel-ized, concurrent-ized, STDIN mysql import tool. Just pipe `INSERT`
statements in and watch your database fill up with your favorite mysql query
tool.

## Primary Motivators

1. At 8tracks we use different databases on EC2 (mysql and postgres). We dump
   database/tables with SQL statements to transfer data between databases.
2. Continue processing a file even if there are bad queries like unescaped
   quotes or backslashes.
3. Speed up imports by using concurrency and parallelism of go.


## Assumptions

* The insert queries do not conflict with each other. If you're importing into
  a table that requires unique records but your queries are not unique, this
  *could* cause problems.
* Insert order does not matter.


## Install For Development

You'll need to setup go on your machine. Checkout
http://golang.org/doc/install.

	go get github.com/8tracks/mysqlspdyimport


## Basic Usage

Use this tool like you would if you were to pipe a bunch of SQL commands into
the `mysql` commandline tool.

 By default, the following commands will use 3 goroutines and import data into
 the localhost mysql instance for the given database.

	cat dump.sql | mysqlspdyimport -d DATABASE

Increase the number of goroutines to 100.

	cat dump.sql | mysqlspdyimport -d DATABASE -c 100

The data piped in can be any query, but `mysqlspdyimport` will only run queries
that start with `INSERT`.


