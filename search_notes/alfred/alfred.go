/* Modified from github.com/drgrib/alfred so that quicklookurl can be "null" without quotes */

package alfred

import (
	"encoding/json"
	. "fmt"
)

// Indent specifies the indent string used for the JSON output for String() and Run(). If set to "", no indentation will be used.
var Indent = "    "

// Rerun specifies the "rerun" value.
var Rerun float64

// Variables specifies the script filter level "variables" object.
var Variables = map[string]string{}

// Items specifies the "items" array. It can be accessed and iterated directly. It can also be appended to directly or appended to using the convenience function Add(item).
var Items = []Item{}

// Icon specifies the "icon" field of an Item.
type Icon struct {
	Type string `json:"type,omitempty"`
	Path string `json:"path,omitempty"`
}

// Mod specifies the values of an Item.Mods map for the "mods" object.
type Mod struct {
	Variables map[string]string `json:"variables,omitempty"`
	Valid     *bool             `json:"valid,omitempty"`
	Arg       string            `json:"arg,omitempty"`
	Subtitle  string            `json:"subtitle,omitempty"`
	Icon      *Icon             `json:"icon,omitempty"`
}

// Text specifies the "text" field of an Item.
type Text struct {
	Copy      string `json:"copy,omitempty"`
	Largetype string `json:"largetype,omitempty"`
}

// Item specifies the members of the "items" array.
type Item struct {
	Variables    map[string]string `json:"variables,omitempty"`
	UID          string            `json:"uid,omitempty"`
	Title        string            `json:"title"`
	Subtitle     string            `json:"subtitle,omitempty"`
	Arg          string            `json:"arg,omitempty"`
	Icon         *Icon             `json:"icon,omitempty"`
	Autocomplete string            `json:"autocomplete,omitempty"`
	Type         string            `json:"type,omitempty"`
	Valid        *bool             `json:"valid,omitempty"`
	Match        string            `json:"match,omitempty"`
	Mods         map[string]Mod    `json:"mods,omitempty"`
	Text         *Text             `json:"text,omitempty"`
	QuicklookURL string           `json:"quicklookurl"`
}

// Bool is a convenience function for filling optional bool values.
func Bool(b bool) *bool {
	return &b
}

// Add is a convenience function for adding new Item instances to Items.
func Add(item Item) {
	Items = append(Items, item)
}

type output struct {
	Rerun     float64           `json:"rerun,omitempty"`
	Variables map[string]string `json:"variables,omitempty"`
	Items     []Item            `json:"items"`
}

// String returns the JSON for currently populated values or the minimum required values.
func String() string {
	output := output{
		Rerun:     Rerun,
		Variables: Variables,
		Items:     Items,
	}
	var err error
	var b []byte
	if Indent == "" {
		b, err = json.Marshal(output)
	} else {
		b, err = json.MarshalIndent(output, "", Indent)
	}
	if err != nil {
		messageErr := Errorf("Error in parser. Please report this output to https://github.com/drgrib/alfred/issues: %v", err)
		panic(messageErr)
	}
	s := string(b)
	return s
}

// Run prints the result of String() to standard output for debugging or direct consumption by an Alfred script filter.
func Run() {
	Println(String())
}
