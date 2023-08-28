# netrc [![GoDoc](https://godoc.org/github.com/jdx/netrc?status.svg)](http://godoc.org/github.com/jdx/netrc) [![CircleCI](https://dl.circleci.com/status-badge/img/gh/jdx/netrc/tree/master.svg?style=svg)](https://dl.circleci.com/status-badge/redirect/gh/jdx/netrc/tree/master)

A netrc parser for Go.

# Usage

Getting credentials for a host.

```go
usr, err := user.Current()
n, err := netrc.Parse(filepath.Join(usr.HomeDir, ".netrc"))
fmt.Println(n.Machine("api.heroku.com").Get("password"))
```

Setting credentials on a host.

```go
usr, err := user.Current()
n, err := netrc.Parse(filepath.Join(usr.HomeDir, ".netrc"))
n.Machine("api.heroku.com").Set("password", "newapikey")
n.Save()
```
