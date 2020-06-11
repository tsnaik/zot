package accesscontrol //nolint:testpackage

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/anuvu/zot/pkg/log"

	. "github.com/smartystreets/goconvey/convey"
)

func makeAccessControlConfigFile(content string) string {
	f, err := ioutil.TempFile("", "access-config-")
	if err != nil {
		panic(err)
	}

	if err := ioutil.WriteFile(f.Name(), []byte(content), 0600); err != nil {
		panic(err)
	}

	return f.Name()
}

func TestAccessControl(t *testing.T) {
	logger := log.NewLogger("debug", "/dev/null")

	Convey("one cred, one static path", t, func() {
		content := `{"repositories":[{"name":"oci-repo-test",
		"users":[{"username":"alice","allowAPIs":[{"path":"/v2/","method":"GET"}]}]}]}`
		accessConfigPath := makeAccessControlConfigFile(content)
		defer os.Remove(accessConfigPath)

		Convey("can access /v2/", func() {
			actual := IsAuthorized("alice", "GET", "/v2/", accessConfigPath, logger)
			So(actual, ShouldBeTrue)
		})

		Convey("cannot access /v2/_catalog", func() {
			actual := IsAuthorized("alice", "GET", "/v2/_catalog", accessConfigPath, logger)
			So(actual, ShouldBeFalse)
		})
	})

	Convey("empty config path should skip the check altogether", t, func() {
		actual := IsAuthorized("alice", "GET", "/v2/", "", logger)
		So(actual, ShouldBeTrue)

		actual = IsAuthorized("alice", "GET", "/v2/_catalog", "", logger)
		So(actual, ShouldBeTrue)
	})

	Convey("one cred, one path with repo name", t, func() {
		content := `{"repositories":[{"name":"oci-repo-test",
		"users":[{"username":"alice","allowAPIs":[{"path":"/v2/{name}/tags/list","method":"GET"}]}]}]}`
		accessConfigPath := makeAccessControlConfigFile(content)
		defer os.Remove(accessConfigPath)

		Convey("oci-repo-test is allowed", func() {
			actual := IsAuthorized("alice", "GET", "/v2/oci-repo-test/tags/list", accessConfigPath, logger)
			So(actual, ShouldBeTrue)
		})

		Convey("any other repo should be denied", func() {
			actual := IsAuthorized("alice", "GET", "/v2/other-repo/tags/list", accessConfigPath, logger)
			So(actual, ShouldBeFalse)
		})
	})
}

func TestAccessControlTwoUsers(t *testing.T) {
	logger := log.NewLogger("debug", "/dev/null")

	Convey("two creds, two paths with repo name", t, func() {
		content := `{"repositories":[{"name":"repo-one","users":[{"username":"alice",
		"allowAPIs":[{"path":"/v2/{name}/tags/list","method":"GET"}]}]},
		{"name":"repo-two","users":[{"username":"bob",
		"allowAPIs":[{"path":"/v2/{name}/tags/list","method":"GET"}]}]}]}`
		accessConfigPath := makeAccessControlConfigFile(content)
		defer os.Remove(accessConfigPath)

		Convey("alice can access repo-one but cannot access repo-two", func() {
			actual := IsAuthorized("alice", "GET", "/v2/repo-one/tags/list", accessConfigPath, logger)
			So(actual, ShouldBeTrue)

			actual = IsAuthorized("alice", "GET", "/v2/repo-two/tags/list", accessConfigPath, logger)
			So(actual, ShouldBeFalse)
		})

		Convey("bob can access repo-two but cannot access repo-one", func() {
			actual := IsAuthorized("bob", "GET", "/v2/repo-two/tags/list", accessConfigPath, logger)
			So(actual, ShouldBeTrue)

			actual = IsAuthorized("bob", "GET", "/v2/repo-one/tags/list", accessConfigPath, logger)
			So(actual, ShouldBeFalse)
		})

		Convey("unknown user, chuck can't access anything", func() {
			actual := IsAuthorized("chuck", "GET", "/v2/repo-one/tags/list", accessConfigPath, logger)
			So(actual, ShouldBeFalse)

			actual = IsAuthorized("chuck", "GET", "/v2/repo-two/tags/list", accessConfigPath, logger)
			So(actual, ShouldBeFalse)

			actual = IsAuthorized("chuck", "GET", "/anything", accessConfigPath, logger)
			So(actual, ShouldBeFalse)
		})
	})
}
