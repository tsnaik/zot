package cli

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestConfigCmdBasics(t *testing.T) {
	Convey("Test config help", t, func() {
		args := []string{"--help"}
		configPath := makeConfigFile("showspinner = false")
		defer os.Remove(configPath)
		cmd := NewConfigCommand(configPath)
		buff := bytes.NewBufferString("")
		cmd.SetOut(buff)
		cmd.SetErr(ioutil.Discard)
		cmd.SetArgs(args)
		err := cmd.Execute()
		So(buff.String(), ShouldContainSubstring, "Usage")
		So(err, ShouldBeNil)
		Convey("with the shorthand", func() {
			args[0] = "-h"
			configPath := makeConfigFile("showspinner = false")
			defer os.Remove(configPath)
			cmd := NewConfigCommand(configPath)
			buff := bytes.NewBufferString("")
			cmd.SetOut(buff)
			cmd.SetErr(ioutil.Discard)
			cmd.SetArgs(args)
			err := cmd.Execute()
			So(buff.String(), ShouldContainSubstring, "Usage")
			So(err, ShouldBeNil)
		})
	})

	Convey("Test config no args", t, func() {
		args := []string{}
		configPath := makeConfigFile("showspinner = false")
		defer os.Remove(configPath)
		cmd := NewConfigCommand(configPath)
		buff := bytes.NewBufferString("")
		cmd.SetOut(buff)
		cmd.SetErr(ioutil.Discard)
		cmd.SetArgs(args)
		err := cmd.Execute()
		So(buff.String(), ShouldContainSubstring, "Usage")
		So(err, ShouldNotBeNil)
	})
}

func TestConfigCmdMain(t *testing.T) {
	Convey("Test fetch all config", t, func() {
		args := []string{"--list"}
		configPath := makeConfigFile("showspinner = false\nurl = https://test-url.com")
		defer os.Remove(configPath)
		cmd := NewConfigCommand(configPath)
		buff := bytes.NewBufferString("")
		cmd.SetOut(buff)
		cmd.SetErr(ioutil.Discard)
		cmd.SetArgs(args)
		err := cmd.Execute()
		So(buff.String(), ShouldContainSubstring, "url = https://test-url.com")
		So(buff.String(), ShouldContainSubstring, "showspinner = false")
		So(err, ShouldBeNil)

		Convey("with the shorthand", func() {
			args := []string{"-l"}
			configPath := makeConfigFile("showspinner = false\nurl = https://test-url.com")
			defer os.Remove(configPath)
			cmd := NewConfigCommand(configPath)
			buff := bytes.NewBufferString("")
			cmd.SetOut(buff)
			cmd.SetErr(ioutil.Discard)
			cmd.SetArgs(args)
			err := cmd.Execute()
			So(buff.String(), ShouldContainSubstring, "url = https://test-url.com")
			So(buff.String(), ShouldContainSubstring, "showspinner = false")
			So(err, ShouldBeNil)
		})
	})

	Convey("Test fetch a config", t, func() {
		args := []string{"url"}
		configPath := makeConfigFile("showspinner = false\nurl = https://test-url.com")
		defer os.Remove(configPath)
		cmd := NewConfigCommand(configPath)
		buff := bytes.NewBufferString("")
		cmd.SetOut(buff)
		cmd.SetErr(ioutil.Discard)
		cmd.SetArgs(args)
		err := cmd.Execute()
		So(buff.String(), ShouldEqual, "https://test-url.com\n")
		So(err, ShouldBeNil)
	})

	Convey("Test add a config", t, func() {
		args := []string{"url", "https://test-url.com"}
		configPath := makeConfigFile("showspinner = false")
		defer os.Remove(configPath)
		cmd := NewConfigCommand(configPath)
		buff := bytes.NewBufferString("")
		cmd.SetOut(buff)
		cmd.SetErr(ioutil.Discard)
		cmd.SetArgs(args)
		err := cmd.Execute()
		So(err, ShouldBeNil)

		actual, err := ioutil.ReadFile(configPath)
		if err != nil {
			panic(err)
		}
		actualStr := string(actual)
		So(actualStr, ShouldContainSubstring, "url = https://test-url.com")
		So(actualStr, ShouldContainSubstring, "showspinner = false")
		So(buff.String(), ShouldEqual, "")
	})

	Convey("Test overwrite a config", t, func() {
		args := []string{"url", "https://new-url.com"}
		configPath := makeConfigFile("showspinner = false\nurl = https://test-url.com")
		defer os.Remove(configPath)
		cmd := NewConfigCommand(configPath)
		buff := bytes.NewBufferString("")
		cmd.SetOut(buff)
		cmd.SetErr(ioutil.Discard)
		cmd.SetArgs(args)
		err := cmd.Execute()
		So(err, ShouldBeNil)

		actual, err := ioutil.ReadFile(configPath)
		if err != nil {
			panic(err)
		}
		actualStr := string(actual)
		So(actualStr, ShouldContainSubstring, "url = https://new-url.com")
		So(actualStr, ShouldContainSubstring, "showspinner = false")
		So(actualStr, ShouldNotContainSubstring, "url = https://test-url.com")
		So(buff.String(), ShouldEqual, "")
	})
}
