// nolint: lll
package cveinfo_test

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/anuvu/zot/pkg/api"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/resty.v1"
)

const (
	BaseURL1    = "http://127.0.0.1:8081"
	SecurePort1 = "8081"
	username    = "test"
	passphrase  = "test"
)

func TestCVESearch(t *testing.T) {
	Convey("Test CVE Search Feature", t, func() {
		config := api.NewConfig()
		config.HTTP.Port = SecurePort1
		htpasswdPath := makeHtpasswdFile()
		defer os.Remove(htpasswdPath)

		config.HTTP.Auth = &api.AuthConfig{
			HTPasswd: api.AuthHTPasswd{
				Path: htpasswdPath,
			},
		}
		c := api.NewController(config)
		/*dir, err := ioutil.TempDir("", "oci-repo-test")
		if err != nil {
			panic(err)
		}*/

		err := copyFiles("./testdata", cve.RootDir)
		if err != nil {
			panic(err)
		}

		defer os.RemoveAll(cve.RootDir)
		c.Config.Storage.RootDirectory = cve.RootDir
		c.Config.Extensions.Search.CVE.UpdateInterval = 1
		go func() {
			// this blocks
			if err := c.Run(); err != nil {
				return
			}
		}()

		// wait till ready
		for {
			_, err := resty.R().Get(BaseURL1)
			if err == nil {
				break
			}
			time.Sleep(100 * time.Millisecond)
		}

		defer func() {
			ctx := context.Background()
			_ = c.Server.Shutdown(ctx)
		}()

		// without creds, should get access error
		resp, err := resty.R().Get(BaseURL1 + "/query/")
		So(err, ShouldBeNil)
		So(resp, ShouldNotBeNil)
		So(resp.StatusCode(), ShouldEqual, 401)
		var e api.Error
		err = json.Unmarshal(resp.Body(), &e)
		So(err, ShouldBeNil)

		// with creds, should get expected status code
		resp, _ = resty.R().SetBasicAuth(username, passphrase).Get(BaseURL1)
		So(resp, ShouldNotBeNil)
		So(resp.StatusCode(), ShouldEqual, 404)

		resp, _ = resty.R().SetBasicAuth(username, passphrase).Get(BaseURL1 + "/query/")
		So(resp, ShouldNotBeNil)
		So(resp.StatusCode(), ShouldEqual, 200)

		resp, _ = resty.R().SetBasicAuth(username, passphrase).Get(BaseURL1 + "/query?query={CVEListForImage(repo:\"centos\"){Tag%20CVEList{Id%20Description%20Severity}}}")
		So(resp, ShouldNotBeNil)
		So(resp.StatusCode(), ShouldEqual, 200)

		resp, _ = resty.R().SetBasicAuth(username, passphrase).Get(BaseURL1 + "/query?query={ImageListForCVE(text:\"CVE-2018-20482\"){Name%20Tags}}")
		So(resp, ShouldNotBeNil)
		So(resp.StatusCode(), ShouldEqual, 200)

		resp, _ = resty.R().SetBasicAuth(username, passphrase).Get(BaseURL1 + "/query?query={CVEListForImage(repo:\"cntos\"){Tag%20CVEList{Id%20Description%20Severity}}}")
		So(resp, ShouldNotBeNil)
		So(resp.StatusCode(), ShouldEqual, 200)

		resp, _ = resty.R().SetBasicAuth(username, passphrase).Get(BaseURL1 + "/query?query={ImageListForCVE(text:\"CVE-201-20482\"){Name%20Tags}}")
		So(resp, ShouldNotBeNil)
		So(resp.StatusCode(), ShouldEqual, 200)

		// Testing Invalid Search URL
		resp, _ = resty.R().SetBasicAuth(username, passphrase).Get(BaseURL1 + "/query?query={CVEListForImage(repo:\"centos\"){Ta%20CVEList{Id%20Description%20Severity}}}")
		So(resp, ShouldNotBeNil)
		So(resp.StatusCode(), ShouldEqual, 422)

		resp, _ = resty.R().SetBasicAuth(username, passphrase).Get(BaseURL1 + "/query?query={ImageListForCVE(tet:\"CVE-2018-20482\"){Name%20Tags}}")
		So(resp, ShouldNotBeNil)
		So(resp.StatusCode(), ShouldEqual, 422)
	})
}

func makeHtpasswdFile() string {
	f, err := ioutil.TempFile("", "htpasswd-")
	if err != nil {
		panic(err)
	}

	// bcrypt(username="test", passwd="test")
	content := []byte("test:$2y$05$hlbSXDp6hzDLu6VwACS39ORvVRpr3OMR4RlJ31jtlaOEGnPjKZI1m\n")
	if err := ioutil.WriteFile(f.Name(), content, 0600); err != nil {
		panic(err)
	}

	return f.Name()
}
