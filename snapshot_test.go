package minimock_test

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/bradleyjkemp/cupaloy"
)

func mustRunMinimockWithArgs(
	t *testing.T,
	args ...string,
) error {
	t.Helper()

	var outBuffer bytes.Buffer

	timeoutContext, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	args = append([]string{"run", "cmd/minimock/minimock.go"}, args...)
	cmd := exec.CommandContext(timeoutContext, "go", args...)
	cmd.Stdout = &outBuffer
	cmd.Stderr = &outBuffer

	t.Log(cmd.String())

	err := cmd.Run()
	if err != nil {
		t.Log(outBuffer.String())
	}

	return err
}

func mustReadFile(
	t *testing.T,
	filename string,
) []byte {
	t.Helper()

	contents, err := os.ReadFile(filename)
	if err != nil {
		t.Errorf("failed to read the file: %v", err)
		t.Fail()
	}

	return contents
}

func TestSnapshot(t *testing.T) {
	snapshotter := cupaloy.New(
		cupaloy.SnapshotSubdirectory("snapshots"),
		cupaloy.SnapshotFileExtension(".go"),
	)

	type testCase struct {
		name               string
		args               []string
		expectedOutputFile string
		wantErr            bool
	}

	testCases := []testCase{
		{
			name:               "package reference",
			args:               []string{"-o", "./tests", "-i", "github.com/gojuno/minimock/v3.Tester"},
			expectedOutputFile: "./tests/tester_mock_test.go",
		},
		{
			name:               "relative reference",
			args:               []string{"-o", "./tests/formatter_mock.go", "-i", "./tests.Formatter"},
			expectedOutputFile: "./tests/formatter_mock.go",
		},
		{
			name:               "generics with any used as param and return type",
			args:               []string{"-o", "./tests/generic_inout.go", "-i", "./tests.genericInout"},
			expectedOutputFile: "./tests/generic_inout.go",
		},
		{
			name:               "generics with any used as return type",
			args:               []string{"-o", "./tests/generic_out.go", "-i", "./tests.genericOut"},
			expectedOutputFile: "./tests/generic_out.go",
		},
		{
			name:               "generics with any used as param type",
			args:               []string{"-o", "./tests/generic_in.go", "-i", "./tests.genericIn"},
			expectedOutputFile: "./tests/generic_in.go",
		},
		{
			name:               "generics with specific type used as a generic constraint",
			args:               []string{"-o", "./tests/generic_specific.go", "-i", "./tests.genericSpecific"},
			expectedOutputFile: "./tests/generic_specific.go",
		},
		{
			name:               "generics with simple union used as a generic constraint",
			args:               []string{"-o", "./tests/generic_simple_union.go", "-i", "./tests.genericSimpleUnion"},
			expectedOutputFile: "./tests/generic_simple_union.go",
		},
		{
			name:               "generics with complex union used as a generic constraint",
			args:               []string{"-o", "./tests/generic_complex_union.go", "-i", "./tests.genericComplexUnion"},
			expectedOutputFile: "./tests/generic_complex_union.go",
		},
		{
			name:               "generics with complex inline union used as a generic constraint",
			args:               []string{"-o", "./tests/generic_inline_union.go", "-i", "./tests.genericInlineUnion"},
			expectedOutputFile: "./tests/generic_inline_union.go",
		},
		{
			name:               "generics with complex inline union with many types",
			args:               []string{"-o", "./tests/generic_inline_with_many_options.go", "-i", "./tests.genericInlineUnionWithManyTypes"},
			expectedOutputFile: "./tests/generic_inline_with_many_options.go",
		},
		{
			name:               "package name specified",
			args:               []string{"-o", "./tests/package_name_specified_test.go", "-i", "github.com/gojuno/minimock/v3.Tester", "-p", "tests_test"},
			expectedOutputFile: "./tests/package_name_specified_test.go",
		},
		{
			name:    "error too many package names",
			args:    []string{"-o", "./tests/package_name_specified_test.go", "-i", "github.com/gojuno/minimock/v3.Tester", "-p", "tests_test, tests"},
			wantErr: true,
		},
		{
			name:    "error interface wildcard and package specified",
			args:    []string{"-o", "./tests/package_name_specified_test.go", "-i", "./tests.*", "-p", "tests_test"},
			wantErr: true,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			err := mustRunMinimockWithArgs(t, testCase.args...)
			if err != nil && testCase.wantErr {
				return
			}

			if testCase.wantErr != (err != nil) {
				t.FailNow()
			}

			snapshotter.SnapshotT(t, mustReadFile(t, testCase.expectedOutputFile))
		})
	}
}
