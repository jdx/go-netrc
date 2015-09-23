package netrc

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"regexp"
	"unicode"
)

// ErrInvalidNetrc means there was an error parsing the netrc file
var ErrInvalidNetrc = errors.New("Invalid netrc")

// Netrc file
type Netrc struct {
	Path     string
	machines []*Machine
	tokens   []string
}

// Machine from the netrc file
type Machine struct {
	Name      string
	IsDefault bool
	tokens    []string
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
	for _, token := range n.tokens {
		b.WriteString(token)
	}
	for _, machine := range n.machines {
		for _, token := range machine.tokens {
			b.WriteString(token)
		}
	}
	return b.String()
}

func read(path string) (io.Reader, error) {
	// TODO: gpg decrypt
	return os.Open(path)
}

func lex(file io.Reader) []string {
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
	tokens := make([]string, 0, 100)
	for scanner.Scan() {
		tokens = append(tokens, scanner.Text())
	}
	return tokens
}

func parse(tokens []string) (*Netrc, error) {
	n := &Netrc{}
	n.machines = make([]*Machine, 0, 20)
	var machine *Machine
	for i, token := range tokens {
		// group tokens into machines
		if token == "machine" || token == "default" {
			// start new group
			machine = &Machine{}
			n.machines = append(n.machines, machine)
			if token == "default" {
				machine.IsDefault = true
			} else {
				machine.Name = tokens[i+2]
			}
		}
		if machine == nil {
			n.tokens = append(n.tokens, token)
		} else {
			machine.tokens = append(machine.tokens, token)
		}
	}
	return n, nil
}

// Get a property from a machine
func (m *Machine) Get(name string) string {
	i := 4
	if m.IsDefault {
		i = 2
	}
	for {
		if i+2 >= len(m.tokens) {
			return ""
		}
		if m.tokens[i] == name {
			return m.tokens[i+2]
		}
		i = i + 4
	}
}

// Set a property on the machine
func (m *Machine) Set(name, value string) {
	i := 4
	if m.IsDefault {
		i = 2
	}
	for i+2 < len(m.tokens) {
		if m.tokens[i] == name {
			m.tokens[i+2] = value
			return
		}
		i = i + 4
	}
	m.tokens = append(m.tokens, "  ", name, " ", value, "\n")
}
