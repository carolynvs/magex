package shx_test

import (
	"os"
	"strings"
	"testing"

	"github.com/carolynvs/magex/shx"
	"github.com/magefile/mage/mg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPreparedCommand_Run(t *testing.T) {
	stdout := shx.RecordStdout()
	defer stdout.Release()
	stderr := shx.RecordStderr()
	defer stderr.Release()

	err := shx.Run("go", "run", "echo.go", "hello world")
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

	err := shx.Run("go", "run")
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

	err := shx.Run("go", "run", "echo.go", "hello world")
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

	err := shx.RunE("go", "run", "echo.go", "hello world")
	gotStdout := stdout.Output()
	if err != nil {
		t.Fatal(err)
	}

	assert.Empty(t, gotStdout)
}

func TestPreparedCommand_RunE_Fail(t *testing.T) {
	stderr := shx.RecordStderr()
	defer stderr.Release()

	err := shx.RunE("go", "run")
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

	err := shx.RunE("go", "run", "echo.go", "hello world")
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

	err := shx.RunS("go", "run", "echo.go", "hello world")
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

	gotOutput, err := shx.OutputS("go", "run", "echo.go", "hello world")
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

	err := shx.RunS("go", "run")
	gotStderr := stderr.Output()
	if err == nil {
		t.Fatal("expected the shx.Command to fail")
	}

	assert.Empty(t, gotStderr)
}

func TestPreparedCommand_CollapseArgs(t *testing.T) {
	err := shx.Run("go", "", "run", "", "echo.go", "hello world", "")
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

	gotOutput, err := shx.Output("go", "run", "echo.go", "hello world")
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

	gotOutput, err := shx.Output("go", "run")
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

	gotOutput, err := shx.Output("go", "run", "echo.go", "hello world")
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

	gotOutput, err := shx.OutputE("go", "run", "echo.go", "hello world")
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

	gotOutput, err := shx.OutputE("go", "run")
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

	gotOutput, err := shx.OutputE("go", "run", "echo.go", "hello world")
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

	gotOutput, err := shx.OutputS("go", "run", "echo.go", "hello world")
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

	gotOutput, err := shx.OutputS("go", "run", "echo.go", "hello world")
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

	gotOutput, err := shx.OutputS("go", "run")
	gotStderr := stderr.Output()
	if err == nil {
		t.Fatal("expected the shx.Command to fail")
	}

	assert.Empty(t, gotStderr)
	assert.Empty(t, gotOutput)
}

func TestPreparedCommand_Stdin(t *testing.T) {
	stdin := strings.NewReader("hello world")
	gotOutput, err := shx.Command("go", "run", "echo.go", "-").Stdin(stdin).OutputE()
	require.NoError(t, err, "command failed")

	assert.Equal(t, "hello world", gotOutput)
}
