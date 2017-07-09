package finger

import (
	"net"
	"regexp"
	"strings"
)

// Query is a valid query a finger server can receive or a client can send.
type Query struct {
	Username string   // Username, can be blank
	Hostname []string // Hostname (zero or more)

	RemoteAddr net.Addr // RemoteAddr set by Server, no effect on clients
}

// ParseQuery parses the line to determine the finger query. According to
// the spec, the line should end with a newline. The input here can omit this.
// If the newline exists, it will still parse.
func ParseQuery(line string) (*Query, error) {
	values := lineRegexp.FindStringSubmatch(line)
	if values == nil {
		return nil, nil
	}

	// The username is always available at index 1, even if blank
	var result Query
	result.Username = values[1]

	// If we have non-empty host text, then parse those
	if values[2] != "" {
		parts := strings.Split(values[2], "@")
		result.Hostname = parts[1:]
	}

	return &result, nil
}

/*
   From RFC 1288 the BNF is as follows. This is matchable via a regexp.

        {Q1}    ::= [{W}|{W}{S}{U}]{C}
        {Q2}    ::= [{W}{S}][{U}]{H}{C}
        {U}     ::= username
        {H}     ::= @hostname | @hostname{H}
        {W}     ::= /W
        {S}     ::= <SP> | <SP>{S}
        {C}     ::= <CRLF>
*/
var lineRegexp = regexp.MustCompile(`` + // This just forces alignment below
	`\s*` + // [{S}]
	`(?P<U>[\w-]+)?` + // [{U}]
	`(?P<H>(@[\w-]+)+)*`, // {H}
)
