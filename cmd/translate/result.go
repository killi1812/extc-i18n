package translate

import (
	"slices"
)

// result is a map storing key location as map[group]value
type result map[string][]value

// value stores key value and its locations
type value struct {
	name string
	locs []location
}

// represents location of key in the project
type location struct {
	// path to folder
	path string
	// line line of code
	line int
	// position refers to starting position in line
	position int
}

// arguments for creating values
type valargs struct {
	name string
	// path to folder
	path string
	// line line of code
	line int
	// position refers to starting position in line
	position int
}

// Add will create a new group with new value
//
// if the group exist and only the value is new, value will be appended in the array
//
// else only adds new location to existing group-value pair
func (col result) Add(key string, val valargs) bool {
	index := slices.IndexFunc(col[key], func(e value) bool {
		return e.name == val.name
	})

	if index == -1 {
		col[key] = append(col[key], value{
			name: val.name,
			locs: []location{
				{
					line:     val.line,
					position: val.position,
					path:     val.path,
				},
			},
		})
	} else {
		col[key][index].locs = append(col[key][1].locs, location{
			line:     val.line,
			position: val.position,
			path:     val.path,
		},
		)
	}

	return true
}
