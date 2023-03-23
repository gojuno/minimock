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

func mustRunMinimockWithParams(
	t *testing.T,
	interfacePattern string,
	outputFile string,
) {
	t.Helper()

	var outBuffer bytes.Buffer

	timeoutContext, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	cmd := exec.CommandContext(timeoutContext, "go", "run", "cmd/minimock/minimock.go", "-o", outputFile, "-i", interfacePattern)
	cmd.Stdout = &outBuffer
	cmd.Stderr = &outBuffer

	t.Log(cmd.String())

	if err := cmd.Run(); err != nil {
		t.Log(outBuffer.String())
		t.Fail()
	}
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
		outputFile         string
		inputInterface     string
		expectedOutputFile string
	}

	testCases := []testCase{
		{
			name:               "package reference",
			outputFile:         "./tests",
			inputInterface:     "github.com/gojuno/minimock/v3.Tester",
			expectedOutputFile: "tests/tester_mock_test.go",
		},
		{
			name:               "relative reference",
			inputInterface:     "./tests.Formatter",
			outputFile:         "./tests/formatter_mock.go",
			expectedOutputFile: "./tests/formatter_mock.go",
		},
		{
			name:               "generics with any used as param and return type",
			inputInterface:     "./tests.genericInout",
			outputFile:         "./tests/generic_inout.go",
			expectedOutputFile: "./tests/generic_inout.go",
		},
		{
			name:               "generics with any used as return type",
			inputInterface:     "./tests.genericOut",
			outputFile:         "./tests/generic_out.go",
			expectedOutputFile: "./tests/generic_out.go",
		},
		{
			name:               "generics with any used as param type",
			inputInterface:     "./tests.genericIn",
			outputFile:         "./tests/generic_in.go",
			expectedOutputFile: "./tests/generic_in.go",
		},
		{
			name:               "generics with specific type used as a generic constraint",
			inputInterface:     "./tests.genericSpecific",
			outputFile:         "./tests/generic_specific.go",
			expectedOutputFile: "./tests/generic_specific.go",
		},
		{
			name:               "generics with simple union used as a generic constraint",
			inputInterface:     "./tests.genericSimpleUnion",
			outputFile:         "./tests/generic_simple_union.go",
			expectedOutputFile: "./tests/generic_simple_union.go",
		},
		{
			name:               "generics with complex union used as a generic constraint",
			inputInterface:     "./tests.genericComplexUnion",
			outputFile:         "./tests/generic_complex_union.go",
			expectedOutputFile: "./tests/generic_complex_union.go",
		},
		{
			name:               "generics with complex inline union used as a generic constraint",
			inputInterface:     "./tests.genericInlineUnion",
			outputFile:         "./tests/generic_inline_union.go",
			expectedOutputFile: "./tests/generic_inline_union.go",
		},
		{
			name:               "generics with complex inline union with many types",
			inputInterface:     "./tests.genericInlineUnionWithManyTypes",
			outputFile:         "./tests/generic_inline_with_many_options.go",
			expectedOutputFile: "./tests/generic_inline_with_many_options.go",
		},
		{
			name:               "2 generic arguments with different types",
			inputInterface:     "./tests.genericMultipleTypes",
			outputFile:         "./tests/generic_multiple_args_with_different_types.go",
			expectedOutputFile: "./tests/generic_multiple_args_with_different_types.go",
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			mustRunMinimockWithParams(t, testCase.inputInterface, testCase.outputFile)
			snapshotter.SnapshotT(t, mustReadFile(t, testCase.expectedOutputFile))
		})
	}
}
