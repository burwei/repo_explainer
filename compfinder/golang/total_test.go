package golang

import (
	"strings"
	"testing"

	"github.com/burwei/repoexplainer/reportgen"
	"github.com/stretchr/testify/assert"
)

func TestComponentFinderFindComponent(t *testing.T) {
	testCases := []struct {
		name         string
		dirPath      string
		fileName     string
		fileContent  string
		expectedComp reportgen.ComponentMap
	}{
		{
			name:     "Simple struct without methods",
			dirPath:  "simple",
			fileName: "simple.go",
			fileContent: `
package simple

type SimpleStruct struct {
    ID int
}
`,
			expectedComp: reportgen.ComponentMap{
				"simple:SimpleStruct": reportgen.Component{
					File:    "simple.go",
					Package: "simple",
					Name:    "SimpleStruct",
					Type:    "struct",
					Fields:  []string{"ID int"},
				},
			},
		},
		{
			name:     "Struct with methods",
			dirPath:  "methods",
			fileName: "methods.go",
			fileContent: `
package methods

type StructWithMethods struct {
    Name string
}

func (s *StructWithMethods) GetName() string {
    return s.Name
}
`,
			expectedComp: reportgen.ComponentMap{
				"methods:StructWithMethods": reportgen.Component{
					File:    "methods.go",
					Package: "methods",
					Name:    "StructWithMethods",
					Type:    "struct",
					Fields:  []string{"Name string"},
					Methods: []string{"GetName() string"},
				},
			},
		},
		{
			name:     "Struct and Interface",
			dirPath:  "implementation",
			fileName: "implementation.go",
			fileContent: `
package implementation

type Interface interface {
    GetName() string
}

type Struct struct {
    Name string
}

func (s *Struct) GetName() string {
    return s.Name
}
`,
			expectedComp: reportgen.ComponentMap{
				"implementation:Interface": reportgen.Component{
					File:    "implementation.go",
					Package: "implementation",
					Name:    "Interface",
					Type:    "interface",
					Methods: []string{"GetName() string"},
				},
				"implementation:Struct": reportgen.Component{
					File:    "implementation.go",
					Package: "implementation",
					Name:    "Struct",
					Type:    "struct",
					Fields:  []string{"Name string"},
					Methods: []string{"GetName() string"},
				},
			},
		},
		{
			name:     "struct, interface and function",
			dirPath:  "allthree",
			fileName: "allthree.go",
			fileContent: `
package allthree

type Interface interface {
    GetName() string
}

type Struct struct {
    Name string
}

func (s *Struct) GetName() string {
    return s.Name
}

func Add(a, b int) int {
	return a + b
}
`,
			expectedComp: reportgen.ComponentMap{
				"allthree:Interface": reportgen.Component{
					File:    "allthree.go",
					Package: "allthree",
					Name:    "Interface",
					Type:    "interface",
					Methods: []string{"GetName() string"},
				},
				"allthree:Struct": reportgen.Component{
					File:    "allthree.go",
					Package: "allthree",
					Name:    "Struct",
					Type:    "struct",
					Fields:  []string{"Name string"},
					Methods: []string{"GetName() string"},
				},
				"allthree:Add": reportgen.Component{
					File:    "allthree.go",
					Package: "allthree",
					Name:    "Add(a, b int) int",
					Type:    "func",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cf := NewComponentFinder()
			cf.SetFile(tc.dirPath, tc.fileName)

			// Simulating line-by-line reading
			lines := strings.Split(tc.fileContent, "\n")
			for _, line := range lines {
				cf.FindComponent(line)
			}

			components := cf.GetComponents()

			assert.Equal(t, tc.expectedComp, components)
		})
	}
}
