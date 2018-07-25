package survey

import (
	"github.com/Netflix/go-expect"
	"github.com/hinshun/vt10x"
	"github.com/stretchr/testify/require"
	"gopkg.in/AlecAivazis/survey.v1"
	"gopkg.in/AlecAivazis/survey.v1/terminal"
	"testing"

	"os"
	"os/exec"
)

// This test uses Select's Prompt() directly, which it seems we shouldn't need to do.
// But it *does* work, so maybe we could use Prompt() directly, worst case.
// But it would be nicer for our code to use the documented Ask() or AskOne()
func TestCLIPrompt(t *testing.T) {
	// Multiplex stdin/stdout to a virtual terminal to respond to ANSI escape
	// sequences (i.e. cursor position report).
	c, state, err := vt10x.NewVT10XConsole()
	require.Nil(t, err)
	defer c.Close()

	donec := make(chan struct{})
	go func() {
		defer close(donec)
		c.ExpectString("Choose a color")
		// You can either send chars that make the selection and <cr>
		c.Send("b")
		//c.Send("l")
		//c.SendLine("")

		// Or just send the line with the selection
		//c.SendLine("blue")

		// Or even use up-and-down arrow stuff.
		c.Send(string(terminal.KeyBackspace)) // back off the "b" we already sent
		c.Send(string(terminal.KeyArrowDown)) // Arrow down to "blue"
		//c.Send(string(terminal.KeyArrowUp))

		c.SendLine("") // Send what we have
		c.ExpectEOF()
		println("Past expect eof")
	}()

	sprompt := &survey.Select{
		Message: "Choose a color:",
		Options: []string{"red", "blue", "green"},
	}
	sprompt.WithStdio(Stdio(c))

	sanswer, err := sprompt.Prompt()
	require.Nil(t, err)
	require.Equal(t, "blue", sanswer)

	// Close the slave end of the pty, and read the remaining bytes from the master end.
	c.Tty().Close()
	<-donec

	println("color = ", sanswer)

	// Dump the terminal's screen.
	t.Log(expect.StripTrailingEmptyLines(state.String()))
}

// This one *attempts* to use the documented survey.AskOne(), but I can't seem to get the
// console working right.
func TestCLIAsk(t *testing.T) {
	// Multiplex stdin/stdout to a virtual terminal to respond to ANSI escape
	// sequences (i.e. cursor position report).
	c, state, err := vt10x.NewVT10XConsole()
	require.Nil(t, err)
	defer c.Close()

	donec := make(chan struct{})
	go func() {
		defer close(donec)
		c.ExpectString("Choose a color")
		c.SendLine("blue")
		c.ExpectEOF()
		println("Past expect eof")
	}()

	color := ""
	sprompt := &survey.Select{
		Message: "Choose a color:",
		Options: []string{"red", "blue", "green"},
	}

	setStdioFunc := func(options *survey.AskOptions) error {
		options.Stdio = Stdio(c)

		return nil
	}

	err = survey.AskOne(sprompt, &color, nil, setStdioFunc)
	require.Nil(t, err)
	require.Equal(t, "blue", color)

	// Close the slave end of the pty, and read the remaining bytes from the master end.
	c.Tty().Close()
	<-donec

	println("color = ", color)

	// Dump the terminal's screen.
	t.Log(expect.StripTrailingEmptyLines(state.String()))
}

func TestSimpleBinary(t *testing.T) {
	// Multiplex stdin/stdout to a virtual terminal to respond to ANSI escape
	// sequences (i.e. cursor position report).
	c, state, err := vt10x.NewVT10XConsole()
	require.Nil(t, err)
	defer c.Close()

	donec := make(chan struct{})
	go func() {
		defer close(donec)
		c.ExpectString("What is your name?")
		c.SendLine("Johnny Appleseed II")
		c.ExpectString("Choose a color")
		c.SendLine("blue")
		c.ExpectEOF()
	}()

	cmd := exec.Command("go", "run", "./examples/simple.go")
	cmd.Stdin = c.Tty()
	cmd.Stdout = c.Tty()
	cmd.Stderr = c.Tty()

	err = cmd.Run()
	require.Nil(t, err)

	// Close the slave end of the pty, and read the remaining bytes from the master end.
	c.Tty().Close()
	<-donec

	// Dump the terminal's screen.
	t.Log(expect.StripTrailingEmptyLines(state.String()))
}

func TestSimple2SelectBinary(t *testing.T) {
	// Multiplex stdin/stdout to a virtual terminal to respond to ANSI escape
	// sequences (i.e. cursor position report).
	c, state, err := vt10x.NewVT10XConsole()
	require.Nil(t, err)
	defer c.Close()

	donec := make(chan struct{})
	go func() {
		defer close(donec)
		c.ExpectString("What is your name?")
		c.SendLine("Johnny Appleseed II")
		c.ExpectString("Choose a color")
		c.SendLine("blue")
		c.ExpectString("Choose your extras")
		c.SendLine("bacon")
		c.ExpectEOF()
	}()

	cmd := exec.Command("go", "run", "./examples/simple_2select.go")
	cmd.Stdin = c.Tty()
	cmd.Stdout = c.Tty()
	cmd.Stderr = c.Tty()

	err = cmd.Run()
	require.Nil(t, err)

	// Close the slave end of the pty, and read the remaining bytes from the master end.
	c.Tty().Close()
	<-donec

	// Dump the terminal's screen.
	t.Log(expect.StripTrailingEmptyLines(state.String()))
}

func TestSimpleBinarySelectInMiddle(t *testing.T) {
	// Multiplex stdin/stdout to a virtual terminal to respond to ANSI escape
	// sequences (i.e. cursor position report).
	c, state, err := vt10x.NewVT10XConsole()
	require.Nil(t, err)
	defer c.Close()

	donec := make(chan struct{})
	go func() {
		defer close(donec)
		c.ExpectString("What is your name?")
		c.SendLine("Johnny Appleseed II")
		c.ExpectString("Choose a color")
		c.SendLine("blue")
		c.ExpectString("Choose your extras")
		c.SendLine("bacon")
		c.ExpectEOF()
	}()

	cmd := exec.Command("go", "run", "./examples/simple_1select_input_follow.go")
	cmd.Stdin = c.Tty()
	cmd.Stdout = c.Tty()
	cmd.Stderr = c.Tty()

	err = cmd.Run()
	require.Nil(t, err)

	// Close the slave end of the pty, and read the remaining bytes from the master end.
	c.Tty().Close()
	<-donec

	// Dump the terminal's screen.
	t.Log(expect.StripTrailingEmptyLines(state.String()))
}

func TestSimpleBinarySelectAtEnd(t *testing.T) {
	// Multiplex stdin/stdout to a virtual terminal to respond to ANSI escape
	// sequences (i.e. cursor position report).
	c, state, err := vt10x.NewVT10XConsole()
	require.Nil(t, err)
	defer c.Close()

	donec := make(chan struct{})
	go func() {
		defer close(donec)
		c.ExpectString("What is your name?")
		c.SendLine("Johnny Appleseed II")
		c.ExpectString("Choose your extras")
		c.SendLine("bacon")
		c.ExpectString("Choose a color")
		c.SendLine("blue")
		c.ExpectEOF()
	}()

	cmd := exec.Command("go", "run", "./examples/simple_1select_input_precedes.go")
	cmd.Stdin = c.Tty()
	cmd.Stdout = c.Tty()
	cmd.Stderr = c.Tty()

	err = cmd.Run()
	require.Nil(t, err)

	// Close the slave end of the pty, and read the remaining bytes from the master end.
	c.Tty().Close()
	<-donec

	// Dump the terminal's screen.
	t.Log(expect.StripTrailingEmptyLines(state.String()))
}

func TestDdevBinary(t *testing.T) {

	accessKeyID := os.Getenv("DDEV_DRUD_S3_AWS_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("DDEV_DRUD_S3_AWS_SECRET_ACCESS_KEY")
	if accessKeyID == "" || secretAccessKey == "" {
		t.Skip("No DDEV_DRUD_S3_AWS_ACCESS_KEY_ID and  DDEV_DRUD_S3_AWS_SECRET_ACCESS_KEY env vars have been set. Skipping DrudS3 specific test.")
	}

	// Multiplex stdin/stdout to a virtual terminal to respond to ANSI escape
	// sequences (i.e. cursor position report).
	c, state, err := vt10x.NewVT10XConsole()
	require.Nil(t, err)
	defer c.Close()

	donec := make(chan struct{})
	go func() {
		defer close(donec)
		c.ExpectString("Project name")
		c.SendLine("d7-kickstart")
		c.ExpectString("Docroot Location")
		c.SendLine("")
		c.ExpectString("Project Type")
		c.SendLine("")
		c.ExpectString("AWS access key id")
		c.SendLine(accessKeyID)
		c.ExpectString("AWS secret access key")
		c.SendLine(secretAccessKey)
		c.ExpectString("AWS S3 Bucket Name")
		c.SendLine("ddev-local-tests")
		c.ExpectString("Choose an environment")
		c.SendLine("production")
		c.ExpectEOF()
	}()

	cmd := exec.Command("ddev", "config", "drud-s3")
	cmd.Stdin = c.Tty()
	cmd.Stdout = c.Tty()
	cmd.Stderr = c.Tty()

	err = cmd.Run()
	require.Nil(t, err)

	// Close the slave end of the pty, and read the remaining bytes from the master end.
	c.Tty().Close()
	<-donec

	// Dump the terminal's screen.
	t.Log(expect.StripTrailingEmptyLines(state.String()))
}
