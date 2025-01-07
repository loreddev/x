package forgejo

import (
	"io/fs"
)

func (p *plugin) Source() (fs.FS, error) {
	return nil, nil
}
