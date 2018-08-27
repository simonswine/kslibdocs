package docgen

import (
	"github.com/GeertJohan/go.rice/embedded"
	"time"
)

func init() {

	// define files
	file2 := &embedded.EmbeddedFile{
		Filename:    "argument_type.md",
		FileModTime: time.Unix(1535372921, 0),
		Content:     string("`<%= propertyName %>` takes an array of type `<%= typeName %>`. You can create\nan instance of `<%= typeName %>` with `<%= typeDef %>.new()`.\n\nsee [<%= typeDef %>]({{< relref \"/<%= typeURL %>\" >}})"),
	}

	// define dirs
	dir1 := &embedded.EmbeddedDir{
		Filename:   "",
		DirModTime: time.Unix(1535370387, 0),
		ChildFiles: []*embedded.EmbeddedFile{
			file2, // "argument_type.md"

		},
	}

	// link ChildDirs
	dir1.ChildDirs = []*embedded.EmbeddedDir{}

	// register embeddedBox
	embedded.RegisterEmbeddedBox(`templates`, &embedded.EmbeddedBox{
		Name: `templates`,
		Time: time.Unix(1535370387, 0),
		Dirs: map[string]*embedded.EmbeddedDir{
			"": dir1,
		},
		Files: map[string]*embedded.EmbeddedFile{
			"argument_type.md": file2,
		},
	})
}
