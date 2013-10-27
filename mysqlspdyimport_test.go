package main

import "testing"

type ConnStrTests struct {
	conf     *MysqlConf
	expected string
}

var connectionStrTests = []ConnStrTests{
	{
		&MysqlConf{host: "localhost", port: "3306", database: "db1"},
		"tcp(localhost:3306)/db1",
	},
	{
		&MysqlConf{username: "root", host: "localhost", port: "3306", database: "db1"},
		"root:@tcp(localhost:3306)/db1",
	},
	{
		&MysqlConf{
			username: "root",
			password: "secret",
			host:     "localhost",
			port:     "3306",
			database: "db1",
		},
		"root:secret@tcp(localhost:3306)/db1",
	},
}

func TestMysqlConfConnectionString(t *testing.T) {
	for _, test := range connectionStrTests {
		result := test.conf.ConnectionString()
		if result != test.expected {
			t.Errorf("Expected '%s', got '%s'", test.expected, result)
		}
	}
}

type CleanUpQueryTests struct {
	in            string
	expected      string
	expectedError string
}

var cleanUpQueryTests = []CleanUpQueryTests{
	{in: "insert something", expected: "insert something"},
	{in: "  insert something", expected: "insert something"},
	{in: "\tinsert something\n", expected: "insert something"},
	{in: "blah", expectedError: NotInsertStatement},
}

func testCleanUpQuery(t *testing.T) {
	for _, test := range cleanUpQueryTests {
		result, err := cleanUpQuery(test.in)
		if err != nil {
			if test.expectedError == "" {
				t.Errorf("Could not test '%s' - Errored: %s", test.in, err)
			}
		}

		if test.expectedError == err.Error() {
			t.Errorf("Expected error for '%s'.", test.in)
		}

		if result != test.expected {
			t.Errorf("Expected '%s', got '%s'", test.expected, result)
		}
	}
}
