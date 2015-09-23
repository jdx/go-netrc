package netrc

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"regexp"
	"strings"
	"unicode"
)

// ErrInvalidNetrc means there was an error parsing the netrc file
var ErrInvalidNetrc = errors.New("Invalid netrc")

// Netrc file
type Netrc struct {
	Path     string
	Default  *Machine
	machines []*Machine
	preChars string
}

// Machine from the netrc file
type Machine struct {
	Name       string
	Login      string
	Password   string
	nameWS     string
	loginWS    string
	passwordWS string
	postChars  string
	isDefault  bool
}

// Parse the netrc file at the given path
// It returns a Netrc instance
func Parse(path string) (*Netrc, error) {
	file, err := read(path)
	if err != nil {
		return nil, err
	}
	netrc, err := parse(lex(file))
	if err != nil {
		return nil, err
	}
	netrc.Path = path
	return netrc, nil
}

// Machine gets a machine by name
func (n *Netrc) Machine(name string) *Machine {
	for _, m := range n.machines {
		if m.Name == name {
			return m
		}
	}
	return nil
}

// Render out the netrc file to a string
func (n *Netrc) Render() string {
	var b bytes.Buffer
	b.WriteString(n.preChars)
	for _, machine := range n.machines {
		b.WriteString("machine " + machine.Name + machine.nameWS)
		b.WriteString("login " + machine.Login + machine.loginWS)
		b.WriteString("password " + machine.Password + machine.passwordWS)
		b.WriteString(machine.postChars)
	}
	return b.String()
}

func read(path string) (io.Reader, error) {
	// TODO: gpg decrypt
	return os.Open(path)
}

func lex(file io.Reader) *bufio.Scanner {
	commentRe := regexp.MustCompile("\\s*#")
	scanner := bufio.NewScanner(file)
	scanner.Split(func(data []byte, eof bool) (int, []byte, error) {
		if eof && len(data) == 0 {
			return 0, nil, nil
		}
		inWhitespace := unicode.IsSpace(rune(data[0]))
		for i, c := range data {
			if c == '#' {
				// line has a comment
				i = commentRe.FindIndex(data)[0]
				if i == 0 {
					// currently in a comment
					i = bytes.IndexByte(data, '\n')
					if i == -1 {
						// no newline at end
						if !eof {
							return 0, nil, nil
						}
						i = len(data)
					}
					for i < len(data) {
						i++
						if !unicode.IsSpace(rune(data[i])) {
							return i, data[0:i], nil
						}
					}
				}
				return i, data[0:i], nil
			}
			if unicode.IsSpace(rune(c)) != inWhitespace {
				return i, data[0:i], nil
			}
		}
		if eof {
			return len(data), data, nil
		}
		return 0, nil, nil
	})
	return scanner
}

func parse(scanner *bufio.Scanner) (*Netrc, error) {
	n := &Netrc{}
	n.machines = make([]*Machine, 0, 20)
	var machine *Machine
	for scanner.Scan() {
		token := scanner.Text()
		if token == "default" {
			machine = &Machine{}
			n.Default = machine
			n.machines = append(n.machines, machine)
		}
		if token == "machine" {
			machine = &Machine{}
			n.machines = append(n.machines, machine)
			scanner.Scan()
			if strings.Contains(scanner.Text(), "\n") {
				return nil, ErrInvalidNetrc
			}
			scanner.Scan()
			machine.Name = scanner.Text()
			scanner.Scan()
			machine.nameWS = scanner.Text()
		} else if token == "login" {
			scanner.Scan()
			if strings.Contains(scanner.Text(), "\n") {
				return nil, ErrInvalidNetrc
			}
			scanner.Scan()
			machine.Login = scanner.Text()
			scanner.Scan()
			machine.loginWS = scanner.Text()
		} else if token == "password" {
			scanner.Scan()
			if strings.Contains(scanner.Text(), "\n") {
				return nil, ErrInvalidNetrc
			}
			scanner.Scan()
			machine.Password = scanner.Text()
			scanner.Scan()
			machine.passwordWS = scanner.Text()
		} else {
			if machine == nil {
				n.preChars += scanner.Text()
			} else {
				machine.postChars += scanner.Text()
			}
		}
	}
	return n, nil
}
