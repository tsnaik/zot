package cli

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	zotErrors "github.com/anuvu/zot/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func makeConfigFile(content string) string {
	f, err := ioutil.TempFile("", "config-*.properties")
	if err != nil {
		panic(err)
	}

	defer f.Close()

	text := []byte(content)
	if err := ioutil.WriteFile(f.Name(), text, 0644); err != nil {
		panic(err)
	}

	return f.Name()
}

func TestSearchCveCmd(t *testing.T) {
	Convey("Test cve help", t, func() {
		args := []string{"--help"}
		configPath := makeConfigFile("showSpinner = false")
		defer os.Remove(configPath)
		cmd := NewCveCommand(new(mockService), configPath)
		buff := bytes.NewBufferString("")
		cmd.SetOut(buff)
		cmd.SetErr(ioutil.Discard)
		cmd.SetArgs(args)
		err := cmd.Execute()
		So(buff.String(), ShouldContainSubstring, "Usage")
		So(err, ShouldBeNil)
		Convey("with the shorthand", func() {
			args[0] = "-h"
			configPath := makeConfigFile("showSpinner = false")
			defer os.Remove(configPath)
			cmd := NewCveCommand(new(mockService), configPath)
			buff := bytes.NewBufferString("")
			cmd.SetOut(buff)
			cmd.SetErr(ioutil.Discard)
			cmd.SetArgs(args)
			err := cmd.Execute()
			So(buff.String(), ShouldContainSubstring, "Usage")
			So(err, ShouldBeNil)
		})
	})
	Convey("Test cve no url no config", t, func() {
		args := []string{"--cve-id", "dummyIdRandom"}
		configPath := makeConfigFile("showSpinner = false")
		defer os.Remove(configPath)
		cmd := NewCveCommand(new(mockService), configPath)
		buff := bytes.NewBufferString("")
		cmd.SetOut(buff)
		cmd.SetErr(ioutil.Discard)
		cmd.SetArgs(args)
		err := cmd.Execute()
		So(err, ShouldNotBeNil)
	})

	Convey("Test cve no params", t, func() {
		args := []string{"--url", "someUrl"}
		configPath := makeConfigFile("showSpinner = false")
		defer os.Remove(configPath)
		cmd := NewCveCommand(new(mockService), configPath)
		buff := bytes.NewBufferString("")
		cmd.SetOut(buff)
		cmd.SetErr(ioutil.Discard)
		cmd.SetArgs(args)
		err := cmd.Execute()
		So(err, ShouldEqual, zotErrors.ErrInvalidFlagsCombination)
	})
	Convey("Test invalid arg combination", t, func() {
		args := []string{"--cve-id", "dummyIdRandom", "--image-name", "dummyImageName", "--url", "someUrl"}
		configPath := makeConfigFile("showSpinner = false")
		defer os.Remove(configPath)
		cmd := NewCveCommand(new(mockService), configPath)
		buff := bytes.NewBufferString("")
		cmd.SetOut(buff)
		cmd.SetErr(ioutil.Discard)
		cmd.SetArgs(args)
		err := cmd.Execute()
		So(err, ShouldEqual, zotErrors.ErrInvalidFlagsCombination)
	})
	Convey("Test cve invalid url", t, func() {
		args := []string{"--image-name", "dummyImageName", "--url", "invalidUrl"}
		configPath := makeConfigFile("showSpinner = false")
		defer os.Remove(configPath)
		cmd := NewCveCommand(new(searchService), configPath)
		buff := bytes.NewBufferString("")
		cmd.SetOut(buff)
		cmd.SetErr(ioutil.Discard)
		cmd.SetArgs(args)
		err := cmd.Execute()
		So(err, ShouldEqual, zotErrors.ErrInvalidURL)
	})
	Convey("Test cve invalid url port", t, func() {
		args := []string{"--image-name", "dummyImageName", "--url", "http://localhost:99999"}
		configPath := makeConfigFile("showSpinner = false")
		defer os.Remove(configPath)
		cmd := NewCveCommand(new(searchService), configPath)
		buff := bytes.NewBufferString("")
		cmd.SetOut(buff)
		cmd.SetErr(ioutil.Discard)
		cmd.SetArgs(args)
		err := cmd.Execute()
		So(err, ShouldNotBeNil)
	})
	Convey("Test cve invalid url port with user", t, func() {
		args := []string{"--image-name", "dummyImageName", "--url", "http://localhost:99999", "--user", "test:test"}
		configPath := makeConfigFile("showSpinner = false")
		defer os.Remove(configPath)
		cmd := NewCveCommand(new(searchService), configPath)
		buff := bytes.NewBufferString("")
		cmd.SetOut(buff)
		cmd.SetErr(ioutil.Discard)
		cmd.SetArgs(args)
		err := cmd.Execute()
		So(err, ShouldNotBeNil)
	})
	Convey("Test cve by image name", t, func() {
		args := []string{"--image-name", "dummyImageName", "--url", "someUrl"}
		configPath := makeConfigFile("showSpinner = false")
		defer os.Remove(configPath)
		cmd := NewCveCommand(new(mockService), configPath)
		buff := bytes.NewBufferString("")
		cmd.SetOut(buff)
		cmd.SetErr(ioutil.Discard)
		cmd.SetArgs(args)
		err := cmd.Execute()
		So(strings.TrimSpace(buff.String()), ShouldEqual, "")
		So(err, ShouldBeNil)
		Convey("using shorthand", func() {
			args := []string{"-I", "dummyImageNameShort", "--url", "someUrl"}
			buff := bytes.NewBufferString("")
			configPath := makeConfigFile("showSpinner = false")
			defer os.Remove(configPath)
			cmd := NewCveCommand(new(mockService), configPath)
			cmd.SetOut(buff)
			cmd.SetErr(ioutil.Discard)
			cmd.SetArgs(args)
			err := cmd.Execute()
			So(strings.TrimSpace(buff.String()), ShouldEqual, "")
			So(err, ShouldBeNil)
		})
	})

	Convey("Test cve url from config", t, func() {
		args := []string{"--image-name", "dummyImageName"}

		configPath := makeConfigFile("showSpinner = false\nurl=url-test.com")
		defer os.Remove(configPath)

		cmd := NewCveCommand(new(mockService), configPath)
		buff := bytes.NewBufferString("")
		cmd.SetOut(buff)
		cmd.SetErr(ioutil.Discard)
		cmd.SetArgs(args)
		err := cmd.Execute()
		So(strings.TrimSpace(buff.String()), ShouldEqual, "")
		So(err, ShouldBeNil)
	})

	Convey("Test cve affected images by cve id", t, func() {
		args := []string{"--cve-id", "dummyCveID", "--url", "someUrlImage"}
		configPath := makeConfigFile("showSpinner = false")
		defer os.Remove(configPath)
		imageCmd := NewCveCommand(new(mockService), configPath)
		buff := bytes.NewBufferString("")
		imageCmd.SetOut(buff)
		imageCmd.SetErr(ioutil.Discard)
		imageCmd.SetArgs(args)
		err := imageCmd.Execute()
		So(strings.TrimSpace(buff.String()), ShouldEqual, "")
		So(err, ShouldBeNil)
		Convey("using shorthand", func() {
			args := []string{"-i", "dummyCveIDShort", "--url", "someUrlImage"}
			buff := bytes.NewBufferString("")
			configPath := makeConfigFile("showSpinner = false")
			defer os.Remove(configPath)
			imageCmd := NewCveCommand(new(mockService), configPath)
			imageCmd.SetOut(buff)
			imageCmd.SetErr(ioutil.Discard)
			imageCmd.SetArgs(args)
			err := imageCmd.Execute()

			So(strings.TrimSpace(buff.String()), ShouldEqual, "")
			So(err, ShouldBeNil)
		})
	})
}

type mockService struct{}

func (service mockService) findCveByImageName(imageName, serverURL,
	username, password string) (CVEListForImageStruct, error) {
	return CVEListForImageStruct{}, nil
}

func (service mockService) findImagesByCveID(cveID, serverURL,
	username, password string) (ImageListForCVEStruct, error) {
	return ImageListForCVEStruct{}, nil
}
