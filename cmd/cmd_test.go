package cmd

import (
	"flag"
	"os"
	"testing"

	"gopkg.in/urfave/cli.v2"

	"github.com/sirupsen/logrus"
)

func TestMain(m *testing.M) {
	c := m.Run()

	os.Exit(c)
}

func TestGetLogLevel(t *testing.T) {
	tests := []struct {
		name     string
		expected logrus.Level
	}{
		{
			name:     "debug",
			expected: logrus.DebugLevel,
		},
		{
			name:     "info",
			expected: logrus.InfoLevel,
		},
		{
			name:     "warn",
			expected: logrus.WarnLevel,
		},
		{
			name:     "warning",
			expected: logrus.WarnLevel,
		},
		{
			name:     "error",
			expected: logrus.ErrorLevel,
		},
		{
			name:     "panic",
			expected: logrus.PanicLevel,
		},
		{
			name:     "fatal",
			expected: logrus.FatalLevel,
		},
		{
			name:     "other",
			expected: logrus.DebugLevel + 1,
		},
	}

	for _, test := range tests {
		lvl := getLogLevel(test.name)
		if lvl != test.expected {
			t.Errorf("expected: %v, got %v", test.expected, lvl)
		}
	}
}

func TestSetup(t *testing.T) {
	f, err := os.OpenFile("test.env", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		if err = f.Close(); err != nil {
			t.Error(err)
			return
		}
		if err = os.Remove("test.env"); err != nil {
			t.Error(err)
			return
		}
	}()

	if _, err = f.WriteString("TEST=yes"); err != nil {
		t.Error(err)
		return
	}

	flags := []struct {
		f, val string
	}{
		{
			f:   "dotenv",
			val: "true",
		},
		{
			f:   "dotenv-loc",
			val: "test.env",
		},
	}

	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	for _, flag := range globalFlags {
		flag.Apply(fs)
	}
	for _, flag := range flags {
		err = fs.Set(flag.f, flag.val)
		if err != nil {
			t.Error(err)
			return
		}
	}

	ctx := cli.NewContext(&cli.App{}, fs, nil)

	if err = setup(ctx); err != nil {
		t.Errorf("Could not run setup: %v", err)
		return
	}

	if os.Getenv("TEST") != "yes" {
		t.Errorf("Expected: %s, got: %s", "yes", os.Getenv("TEST"))
	}
}

func TestSetupDev(t *testing.T) {
	f, err := os.OpenFile("test.env", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		if err = f.Close(); err != nil {
			t.Error(err)
			return
		}
		if err = os.Remove("test.env"); err != nil {
			t.Error(err)
			return
		}
	}()

	if _, err = f.WriteString("TEST=yes"); err != nil {
		t.Error(err)
		return
	}

	flags := []struct {
		f, val string
	}{
		{
			f:   "dev",
			val: "true",
		},
		{
			f:   "dotenv-loc",
			val: "test.env",
		},
	}

	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	for _, flag := range globalFlags {
		flag.Apply(fs)
	}
	for _, flag := range flags {
		err = fs.Set(flag.f, flag.val)
		if err != nil {
			t.Error(err)
			return
		}
	}

	ctx := cli.NewContext(&cli.App{}, fs, nil)

	if err = setup(ctx); err != nil {
		t.Errorf("Could not run setup: %v", err)
		return
	}

	if os.Getenv("TEST") != "yes" {
		t.Errorf(".env: expected: %s, got: %s", "yes", os.Getenv("TEST"))
		return
	}

	if os.Getenv("LOG_LEVEL") != "debug" {
		t.Errorf("LOG_LEVEL: expected: %s, got: %s", "debug", os.Getenv("LOG_LEVEL"))
		return
	}
}
