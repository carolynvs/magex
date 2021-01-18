package shx_test

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/carolynvs/magex/shx"
	"github.com/stretchr/testify/assert"
	"github.com/magefile/mage/mg"
)

func TestPreparedCommand_Run(t *testing.T) {
	stdout := shx.RecordStdout()
	defer stdout.Release()
	stderr := shx.RecordStderr()
	defer stderr.Release()

	err := shx.Command("go", "run", "echo.go", "hello world").Run()
	gotStdout := stdout.Output()
	gotStderr := stderr.Output()
	if err != nil {
		t.Fatal(err)
	}

	assert.Empty(t, gotStdout)
	assert.Empty(t, gotStderr)
}

func TestPreparedCommand_Run_Fail(t *testing.T) {
	stderr := shx.RecordStderr()
	defer stderr.Release()

	err := shx.Command("go", "run").Run()
	gotStderr := stderr.Output()
	if err == nil {
		t.Fatal("expected shx.Command to fail")
	}

	wantStderr := "go run: no go files listed\n"
	assert.Equal(t, wantStderr, gotStderr)
}

func TestPreparedCommand_Run_Verbose(t *testing.T) {
	os.Setenv(mg.VerboseEnv, "true")
	defer os.Unsetenv(mg.VerboseEnv)

	stdout := shx.RecordStdout()
	defer stdout.Release()
	stderr := shx.RecordStderr()
	defer stderr.Release()

	err := shx.Command("go", "run", "echo.go", "hello world").Run()
	gotStdout := stdout.Output()
	gotStderr := stderr.Output()
	if err != nil {
		t.Fatal(err)
	}

	wantStdout := "hello world\n"
	assert.Equal(t, wantStdout, gotStdout)

	wantStderr := "go run echo.go hello world"
	assert.Contains(t, gotStderr, wantStderr)
}

func TestPreparedCommand_RunE(t *testing.T) {
	os.Setenv(mg.VerboseEnv, "true")
	defer os.Unsetenv(mg.VerboseEnv)

	stdout := shx.RecordStdout()
	defer stdout.Release()

	err := shx.Command("go", "run", "echo.go", "hello world").RunE()
	gotStdout := stdout.Output()
	if err != nil {
		t.Fatal(err)
	}

	assert.Empty(t, gotStdout)
}

func TestPreparedCommand_RunE_Fail(t *testing.T) {
	stderr := shx.RecordStderr()
	defer stderr.Release()

	err := shx.Command("go", "run").RunE()
	gotStderr := stderr.Output()
	if err == nil {
		t.Fatal("expected the shx.Command to fail")
	}

	wantStderr := "go run: no go files listed\n"
	assert.Equal(t, wantStderr, gotStderr)
}

func TestPreparedCommand_RunE_Verbose(t *testing.T) {
	os.Setenv(mg.VerboseEnv, "true")
	defer os.Unsetenv(mg.VerboseEnv)

	stdout := shx.RecordStdout()
	defer stdout.Release()
	stderr := shx.RecordStderr()
	defer stderr.Release()

	err := shx.Command("go", "run", "echo.go", "hello world").RunE()
	gotStdout := stdout.Output()
	gotStderr := stderr.Output()
	if err != nil {
		t.Fatal(err)
	}

	assert.Empty(t, gotStdout)

	wantStderr := "go run echo.go hello world"
	assert.Contains(t, gotStderr, wantStderr)
}

func TestPreparedCommand_RunS(t *testing.T) {
	stdout := shx.RecordStdout()
	defer stdout.Release()
	stderr := shx.RecordStderr()
	defer stderr.Release()

	err := shx.Command("go", "run", "echo.go", "hello world").RunS()
	gotStdout := stdout.Output()
	gotStderr := stderr.Output()
	if err != nil {
		t.Fatal(err)
	}

	assert.Empty(t, gotStdout)
	assert.Empty(t, gotStderr)
}

func TestPreparedCommand_RunS_Verbose(t *testing.T) {
	os.Setenv(mg.VerboseEnv, "true")
	defer os.Unsetenv(mg.VerboseEnv)

	stdout := shx.RecordStdout()
	defer stdout.Release()
	stderr := shx.RecordStderr()
	defer stderr.Release()

	gotOutput, err := shx.Command("go", "run", "echo.go", "hello world").OutputS()
	gotStdout := stdout.Output()
	gotStderr := stderr.Output()
	if err != nil {
		t.Fatal(err)
	}

	wantOutput := "hello world"
	assert.Equal(t, wantOutput, gotOutput)
	assert.Empty(t, gotStdout)

	wantStderr := "go run echo.go hello world"
	assert.Contains(t, gotStderr, wantStderr)
}

func TestPreparedCommand_RunS_Fail(t *testing.T) {
	stderr := shx.RecordStderr()
	defer stderr.Release()

	err := shx.Command("go", "run").RunS()
	gotStderr := stderr.Output()
	if err == nil {
		t.Fatal("expected the shx.Command to fail")
	}

	assert.Empty(t, gotStderr)
}

func TestPreparedCommand_CollapseArgs(t *testing.T) {
	err := shx.Command("go", "", "run", "", "echo.go", "hello world", "").Run()
	if err == nil {
		t.Fatal("expected empty arguments to be preserved in the constructor")
	}

	err = shx.Command("go", "run").Args("", "echo.go", "", "hello world", "").Run()
	if err == nil {
		t.Fatal("expected empty arguments to be preserved when Args is called")
	}

	err = shx.Command("go", "", "run", "", "echo.go", "hello world", "").CollapseArgs().Run()
	if err != nil {
		t.Fatal("expected empty arguments to be removed")
	}
}

func TestPreparedCommand_Output(t *testing.T) {
	stdout := shx.RecordStdout()
	defer stdout.Release()
	stderr := shx.RecordStderr()
	defer stderr.Release()

	gotOutput, err := shx.Command("go", "run", "echo.go", "hello world").Output()
	gotStdout := stdout.Output()
	gotStderr := stderr.Output()
	if err != nil {
		t.Fatal(err)
	}

	wantOutput := "hello world"
	assert.Equal(t, wantOutput, gotOutput)
	assert.Empty(t, gotStdout)
	assert.Empty(t, gotStderr)
}

func TestPreparedCommand_Output_Fail(t *testing.T) {
	stderr := shx.RecordStderr()
	defer stderr.Release()

	gotOutput, err := shx.Command("go", "run").Output()
	gotStderr := stderr.Output()
	if err == nil {
		t.Fatal("expected shx.Command to fail")
	}

	wantStderr := "go run: no go files listed\n"
	assert.Equal(t, wantStderr, gotStderr)
	assert.Empty(t, gotOutput)
}

func TestPreparedCommand_Output_Verbose(t *testing.T) {
	os.Setenv(mg.VerboseEnv, "true")
	defer os.Unsetenv(mg.VerboseEnv)

	stdout := shx.RecordStdout()
	defer stdout.Release()
	stderr := shx.RecordStderr()
	defer stderr.Release()

	gotOutput, err := shx.Command("go", "run", "echo.go", "hello world").Output()
	gotStdout := stdout.Output()
	gotStderr := stderr.Output()
	if err != nil {
		t.Fatal(err)
	}

	wantOutput := "hello world"
	assert.Equal(t, wantOutput, gotOutput)

	wantStdout := "hello world\n"
	assert.Equal(t, wantStdout, gotStdout)

	wantStderr := "go run echo.go hello world"
	assert.Contains(t, gotStderr, wantStderr)
}

func TestPreparedCommand_OutputE(t *testing.T) {
	os.Setenv(mg.VerboseEnv, "true")
	defer os.Unsetenv(mg.VerboseEnv)

	stdout := shx.RecordStdout()
	defer stdout.Release()

	gotOutput, err := shx.Command("go", "run", "echo.go", "hello world").OutputE()
	gotStdout := stdout.Output()
	if err != nil {
		t.Fatal(err)
	}

	wantOutput := "hello world"
	assert.Equal(t, wantOutput, gotOutput)
	assert.Empty(t, gotStdout)
}

func TestPreparedCommand_OutputE_Fail(t *testing.T) {
	stderr := shx.RecordStderr()
	defer stderr.Release()

	gotOutput, err := shx.Command("go", "run").OutputE()
	gotStderr := stderr.Output()
	if err == nil {
		t.Fatal("expected the shx.Command to fail")
	}

	wantStderr := "go run: no go files listed\n"
	assert.Equal(t, wantStderr, gotStderr)
	assert.Empty(t, gotOutput)
}

func TestPreparedCommand_OutputE_Verbose(t *testing.T) {
	os.Setenv(mg.VerboseEnv, "true")
	defer os.Unsetenv(mg.VerboseEnv)

	stdout := shx.RecordStdout()
	defer stdout.Release()
	stderr := shx.RecordStderr()
	defer stderr.Release()

	gotOutput, err := shx.Command("go", "run", "echo.go", "hello world").OutputE()
	gotStdout := stdout.Output()
	gotStderr := stderr.Output()
	if err != nil {
		t.Fatal(err)
	}

	wantOutput := "hello world"
	assert.Equal(t, wantOutput, gotOutput)
	assert.Empty(t, gotStdout)

	wantStderr := "go run echo.go hello world"
	assert.Contains(t, gotStderr, wantStderr)
}

func TestPreparedCommand_OutputS(t *testing.T) {
	stdout := shx.RecordStdout()
	defer stdout.Release()
	stderr := shx.RecordStderr()
	defer stderr.Release()

	gotOutput, err := shx.Command("go", "run", "echo.go", "hello world").OutputS()
	gotStdout := stdout.Output()
	gotStderr := stderr.Output()
	if err != nil {
		t.Fatal(err)
	}

	wantOutput := "hello world"
	assert.Equal(t, wantOutput, gotOutput)
	assert.Empty(t, gotStdout)
	assert.Empty(t, gotStderr)
}

func TestPreparedCommand_OutputS_Verbose(t *testing.T) {
	os.Setenv(mg.VerboseEnv, "true")
	defer os.Unsetenv(mg.VerboseEnv)

	stdout := shx.RecordStdout()
	defer stdout.Release()
	stderr := shx.RecordStderr()
	defer stderr.Release()

	gotOutput, err := shx.Command("go", "run", "echo.go", "hello world").OutputS()
	gotStdout := stdout.Output()
	gotStderr := stderr.Output()
	if err != nil {
		t.Fatal(err)
	}

	wantOutput := "hello world"
	assert.Equal(t, wantOutput, gotOutput)
	assert.Empty(t, gotStdout)

	wantStderr := "go run echo.go hello world"
	assert.Contains(t, gotStderr, wantStderr)
}

func TestPreparedCommand_OutputS_Fail(t *testing.T) {
	stderr := shx.RecordStderr()
	defer stderr.Release()

	gotOutput, err := shx.Command("go", "run").OutputS()
	gotStderr := stderr.Output()
	if err == nil {
		t.Fatal("expected the shx.Command to fail")
	}

	assert.Empty(t, gotStderr)
	assert.Empty(t, gotOutput)
}

func ExamplePreparedCommand_RunV() {
	err := shx.Command("go", "run", "echo.go", "hello world").RunV()
	if err != nil {
		log.Fatal(err)
	}
	// Output: hello world
}

func ExamplePreparedCommand_In() {
	tmp, err := ioutil.TempDir("", "mage")
	if err != nil {
		log.Fatal(err)
	}

	contents := `package main

import "fmt"

func main() {
	fmt.Println("hello world")
}
`
	err = ioutil.WriteFile(filepath.Join(tmp, "test_main.go"), []byte(contents), 0644)
	if err != nil {
		log.Fatal(err)
	}

	// Run `go run test_main.go` in /tmp
	err = shx.Command("go", "run", "test_main.go").In(tmp).RunV()
	if err != nil {
		log.Fatal(err)
	}
	// Output: hello world
}

func ExamplePreparedCommand_RunS() {
	err := shx.Command("go", "run", "echo.go", "hello world").RunS()
	if err != nil {
		log.Fatal(err)
	}
	// Output:
}

func ExamplePreparedCommand_CollapseArgs() {
	err := shx.Command("go", "run", "echo.go", "hello", "", "world").CollapseArgs().RunV()
	if err != nil {
		log.Fatal(err)
	}

	// Output: hello world
}
