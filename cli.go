package main

import (
	"encoding/json"
	"fmt"
	"io"
)

func cmdList(stdout io.Writer, stderr io.Writer, s *state) int {
	groups := make(map[string][]string, 0)
	for _, res := range s.resources() {
		for _, grp := range res.Groups() {

			_, ok := groups[grp]
			if !ok {
				groups[grp] = []string{}
			}

			groups[grp] = append(groups[grp], res.Host())
		}
	}

	return output(stdout, stderr, groups)
}

func cmdHost(stdout io.Writer, stderr io.Writer, s *state, hostname string, remap map[string]string) int {
	for _, res := range s.resources() {
		if hostname == res.Host() {
			return output(stdout, stderr, res.Attributes(remap))
		}
	}

	fmt.Fprintf(stderr, "No such host: %s\n", hostname)
	return 1
}

// output marshals an arbitrary JSON object and writes it to stdout, or writes
// an error to stderr, then returns the appropriate exit code.
func output(stdout io.Writer, stderr io.Writer, whatever interface{}) int {
	b, err := json.Marshal(whatever)
	if err != nil {
		fmt.Fprintf(stderr, "Error encoding JSON: %s\n", err)
		return 1
	}

	_, err = stdout.Write(b)
	if err != nil {
		fmt.Fprintf(stderr, "Error writing JSON: %s\n", err)
		return 1
	}

	return 0
}
