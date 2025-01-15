package exporter

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"strings"
)

type AuthEndpoint struct {
	username string
	password string
}
type AuthEndpoints map[string]AuthEndpoint

func ParseAuthFile(path string) (*AuthEndpoints, error) {
	endpoints := make(AuthEndpoints)

	auth_f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	line_n := 0
	lines := bufio.NewScanner(auth_f)
	lines.Split(bufio.ScanLines)

	for lines.Scan() {
		line_n++

		f := strings.SplitN(lines.Text(), " ", 3)

		_, err := url.ParseRequestURI(f[0])
		if err != nil {
			return nil, fmt.Errorf("error parsing URL \"%s\" at %s:%d: %w", f[0], path, line_n, err)

		}

		endpoints[f[0]] = AuthEndpoint{
			username: f[1],
			password: f[2],
		}
	}
	return &endpoints, nil
}
