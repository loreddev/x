package forgejo

import (
	"forge.capytal.company/loreddev/x/blogo"
)

const pluginName = "blogo-forgejo"

type plugin struct {
	owner string
	repo  string
}
type Opts struct {
	Ref        string
}

func New(owner, repo, apiUrl string, opts ...Opts) blogo.Plugin {
	opt := Opts{}
	if len(opts) > 0 {
		opt = opts[0]
	}

	return &plugin{
	}
}

func (p *plugin) Name() string {
	return pluginName
}
