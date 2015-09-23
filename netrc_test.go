package netrc_test

import (
	"io/ioutil"
	"testing"

	"github.com/dickeyxxx/netrc"
	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type NetrcSuite struct{}

var _ = Suite(&NetrcSuite{})

func (s *NetrcSuite) TestLogin(c *C) {
	f, err := netrc.Parse("./examples/login.netrc")
	c.Assert(err, IsNil)
	heroku := f.Machine("api.heroku.com")
	c.Check(heroku.Get("login"), Equals, "jeff@heroku.com")
	c.Check(heroku.Get("password"), Equals, "foo")
	body, _ := ioutil.ReadFile(f.Path)
	c.Check(f.Render(), Equals, string(body))
}

func (s *NetrcSuite) TestSetPassword(c *C) {
	f, err := netrc.Parse("./examples/login.netrc")
	c.Assert(err, IsNil)
	heroku := f.Machine("api.heroku.com")
	heroku.Set("password", "foobar")
	c.Check(f.Render(), Equals, "# this is my login netrc\nmachine api.heroku.com\n  login jeff@heroku.com # this is my username\n  password foobar\n")
}

func (s *NetrcSuite) TestSampleMulti(c *C) {
	f, err := netrc.Parse("./examples/sample_multi.netrc")
	c.Assert(err, IsNil)
	c.Check(f.Machine("m").Get("login"), Equals, "lm")
	c.Check(f.Machine("m").Get("password"), Equals, "pm")
	c.Check(f.Machine("n").Get("login"), Equals, "ln")
	c.Check(f.Machine("n").Get("password"), Equals, "pn")
	body, _ := ioutil.ReadFile(f.Path)
	c.Check(f.Render(), Equals, string(body))
}

func (s *NetrcSuite) TestNewlineless(c *C) {
	f, err := netrc.Parse("./examples/newlineless.netrc")
	c.Assert(err, IsNil)
	c.Check(f.Machine("m").Get("login"), Equals, "l")
	c.Check(f.Machine("m").Get("password"), Equals, "p")
	body, _ := ioutil.ReadFile(f.Path)
	c.Check(f.Render(), Equals, string(body))
}

func (s *NetrcSuite) TestBadDefaultOrder(c *C) {
	f, err := netrc.Parse("./examples/bad_default_order.netrc")
	c.Assert(err, IsNil)
	c.Check(f.Machine("mail.google.com").Get("login"), Equals, "joe@gmail.com")
	c.Check(f.Machine("mail.google.com").Get("password"), Equals, "somethingSecret")
	c.Check(f.Machine("ray").Get("login"), Equals, "demo")
	c.Check(f.Machine("ray").Get("password"), Equals, "mypassword")
	body, _ := ioutil.ReadFile(f.Path)
	c.Check(f.Render(), Equals, string(body))
}
