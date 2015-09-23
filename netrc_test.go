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
	c.Check(heroku.Login, Equals, "jeff@heroku.com")
	c.Check(heroku.Password, Equals, "foo")
	body, _ := ioutil.ReadFile(f.Path)
	c.Check(f.Render(), Equals, string(body))
}

func (s *NetrcSuite) TestSampleMulti(c *C) {
	f, err := netrc.Parse("./examples/sample_multi.netrc")
	c.Assert(err, IsNil)
	c.Check(f.Machine("m").Login, Equals, "lm")
	c.Check(f.Machine("m").Password, Equals, "pm")
	c.Check(f.Machine("n").Login, Equals, "ln")
	c.Check(f.Machine("n").Password, Equals, "pn")
	body, _ := ioutil.ReadFile(f.Path)
	c.Check(f.Render(), Equals, string(body))
}

func (s *NetrcSuite) TestNewlineless(c *C) {
	f, err := netrc.Parse("./examples/newlineless.netrc")
	c.Assert(err, IsNil)
	c.Check(f.Machine("m").Login, Equals, "l")
	c.Check(f.Machine("m").Password, Equals, "p")
	body, _ := ioutil.ReadFile(f.Path)
	c.Check(f.Render(), Equals, string(body))
}

func (s *NetrcSuite) TestBadDefaultOrder(c *C) {
	f, err := netrc.Parse("./examples/bad_default_order.netrc")
	c.Assert(err, IsNil)
	c.Check(f.Machine("mail.google.com").Login, Equals, "joe@gmail.com")
	c.Check(f.Machine("mail.google.com").Password, Equals, "somethingSecret")
	body, _ := ioutil.ReadFile(f.Path)
	c.Check(f.Render(), Equals, string(body))
}
