package models

import "testing"

func TestGetCommand(t *testing.T) {

	// create a test table to iterate over
	tables := []struct {
		command string
		test    string
	}{
		{
			"foo bar",
			"foo bar",
		},
	}

	// iterate around the test tables and perform the tests
	for _, table := range tables {

		// create a SpawnBuild object to work wtih
		sb := SpawnBuild{
			Command: table.command,
		}

		fullCmd := sb.GetCommand()

		if fullCmd != table.test {
			t.Error("Command has not been derived correctly")
		}
	}
}

func TestGetCommandParts(t *testing.T) {

	// create a test table to iterate over
	tables := []struct {
		command  string
		testCmd  string
		testArgs string
	}{
		{
			"foo bar",
			"foo",
			"bar",
		},
	}

	for _, table := range tables {

		// create a SpawnBuild object to work wtih
		sb := SpawnBuild{
			Command: table.command,
		}

		cmd, args := sb.GetCommandParts()

		if cmd != table.testCmd {
			t.Error("Command is not as expected")
		}

		if args != table.testArgs {
			t.Error("Argument list is not as expected")
		}
	}
}
