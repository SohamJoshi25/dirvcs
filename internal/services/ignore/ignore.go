package ignore

import (
	Path "dirvcs/internal/data/path"

	ignore "github.com/sabhiram/go-gitignore"
)

var Ignore *ignore.GitIgnore

func init() {
	ignore, err := ignore.CompileIgnoreFile(Path.IGNORE_PATH)
	if err != nil {
		Ignore = nil
	} else {
		Ignore = ignore
	}
}
