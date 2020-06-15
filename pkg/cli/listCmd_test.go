package cli

import (
	"bytes"
	"testing"

	"github.com/kardianos/osext"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSearchCmd(t *testing.T) {
	Convey("Test search help", t, func() {
		args := []string{"--help"}
		dir, _ := osext.ExecutableFolder()
		configPath := dir + "/.cli.properties"

		cmd := NewListCmd(configPath)
		buff := bytes.NewBufferString("")
		cmd.SetOut(buff)
		cmd.SetArgs(args)
		err := cmd.Execute()
		So(buff.String(), ShouldContainSubstring, "Usage")
		So(err, ShouldBeNil)
		Convey("with the shorthand", func() {
			args[0] = "-h"
			dir, _ := osext.ExecutableFolder()
			configPath := dir + "/.cli.properties"

			cmd := NewListCmd(configPath)
			buff := bytes.NewBufferString("")
			cmd.SetOut(buff)
			cmd.SetArgs(args)
			err := cmd.Execute()
			So(buff.String(), ShouldContainSubstring, "Usage")
			So(err, ShouldBeNil)
		})
	})

	Convey("Test search invalid subcommand", t, func() {
		args := []string{"randomSubCommand"}
		dir, _ := osext.ExecutableFolder()
		configPath := dir + "/.cli.properties"

		cmd := NewListCmd(configPath)
		buff := bytes.NewBufferString("")
		cmd.SetOut(buff)
		cmd.SetArgs(args)
		err := cmd.Execute()
		So(buff.String(), ShouldContainSubstring, "usage")
		So(err, ShouldNotBeNil)
	})

	Convey("Test search invalid flag", t, func() {
		args := []string{"--random"}
		dir, _ := osext.ExecutableFolder()
		configPath := dir + "/.cli.properties"

		cmd := NewListCmd(configPath)
		buff := bytes.NewBufferString("")
		cmd.SetOut(buff)
		cmd.SetArgs(args)
		err := cmd.Execute()
		So(buff.String(), ShouldContainSubstring, "unknown flag")
		So(err, ShouldNotBeNil)
		Convey("and a shorthand", func() {
			args[0] = "-r"
			dir, _ := osext.ExecutableFolder()
			configPath := dir + "/.cli.properties"

			cmd := NewListCmd(configPath)
			buff := bytes.NewBufferString("")
			cmd.SetOut(buff)
			cmd.SetArgs(args)
			err := cmd.Execute()
			So(buff.String(), ShouldContainSubstring, "unknown shorthand flag")
			So(err, ShouldNotBeNil)
		})
	})
}
