package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
)

type state struct {
	Modules []moduleState `json:"modules"`
}

// read populates the state object from a statefile.
func (s *state) read(stateFile io.Reader) error {

	// read statefile contents
	b, err := ioutil.ReadAll(stateFile)
	if err != nil {
		return err
	}

	// parse into struct
	err = json.Unmarshal(b, s)
	if err != nil {
		return err
	}

	return nil
}

// resources returns a slice of the Resources found in the statefile.
func (s *state) resources() []*Resource {
	inst := make([]*Resource, 0)

	for _, m := range s.Modules {
		for _, k := range m.resourceKeys() {

			// Terraform stores resources in a name->map map, but we need the name to
			// decide which groups to include the resource in. So wrap it in a higher-
			// level object with both properties.
			r, err := NewResource(k, m.ResourceStates[k])
			if err != nil {
				continue
			}

			if r.IsSupported() {
				inst = append(inst, r)
			}
		}
	}

	var hostIdentifiers = make(map[string]bool)
	for _, r := range inst {
		if hostIdentifiers[r.Host()] {
			fmt.Fprintf(os.Stderr, "Error: duplicate host \"%s\"\n", r.Host())
			os.Exit(1)
		}
		hostIdentifiers[r.Host()] = true
	}

	return inst
}

type moduleState struct {
	ResourceStates map[string]resourceState `json:"resources"`
}

// resourceKeys returns a sorted slice of the key names of the resources in this
// module. Do this instead of range over ResourceStates, to ensure that the
// output is consistent.
func (ms *moduleState) resourceKeys() []string {
	lk := len(ms.ResourceStates)
	keys := make([]string, lk, lk)
	i := 0

	for k := range ms.ResourceStates {
		keys[i] = k
		i += 1
	}

	sort.Strings(keys)
	return keys
}

type resourceState struct {

	// Populated from statefile
	Type    string        `json:"type"`
	Primary instanceState `json:"primary"`
}

type instanceState struct {
	ID         string            `json:"id"`
	Attributes map[string]string `json:"attributes,omitempty"`
}
