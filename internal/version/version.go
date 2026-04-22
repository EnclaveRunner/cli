package version

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// RemoteVersionURL is the raw GitHub URL holding the latest version string.
const RemoteVersionURL = "https://raw.githubusercontent.com/EnclaveRunner/cli/main/Version"

// fetchRemote retrieves the remote Version file contents.
func fetchRemote() (string, error) {
	resp, err := http.Get(RemoteVersionURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("unexpected status: %s", resp.Status)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(b)), nil
}

// Compare returns 1 if a>b, 0 if equal, -1 if a<b.
func Compare(a, b string) int {
	av := normalize(a)
	bv := normalize(b)
	for i := 0; i < 3; i++ {
		if av[i] > bv[i] {
			return 1
		}
		if av[i] < bv[i] {
			return -1
		}
	}
	return 0
}

func normalize(s string) [3]int {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "v")
	parts := strings.Split(s, ".")
	var out [3]int
	for i := 0; i < 3 && i < len(parts); i++ {
		n, _ := strconv.Atoi(parts[i])
		out[i] = n
	}
	return out
}

// CheckRemote compares local to remote and returns remote + whether it's newer.
func CheckRemote(local string) (remote string, newer bool, err error) {
	r, err := fetchRemote()
	if err != nil {
		return "", false, err
	}
	c := Compare(local, r)
	return r, c == -1, nil
}
