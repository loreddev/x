package forgejo

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"forge.capytal.company/loreddev/x/blogo"
)

const pluginName = "blogo-gitea"

type plugin struct {
	client *client

	owner string
	repo  string
}

type Opts struct {
	HTTPClient *http.Client
	Ref        string
}

func New(owner, repo, apiUrl string, opts ...Opts) blogo.Plugin {
	opt := Opts{}
	if len(opts) > 0 {
		opt = opts[0]
	}

	if opt.HTTPClient == nil {
		opt.HTTPClient = http.DefaultClient
	}

	u, err := url.Parse(apiUrl)
	if err != nil {
		panic(
			fmt.Sprintf(
				"%s: %q is not a valid URL. Err: %q",
				pluginName,
				apiUrl,
				err.Error(),
			),
		)
	}

	if u.Path == "" || u.Path == "/" {
		u.Path = "/api/v1"
	} else {
		u.Path = strings.TrimSuffix(u.Path, "/api/v1")
	}

	client := newClient(u.String(), opt.HTTPClient)

	return &plugin{
		client: client,
	}
}

func (p *plugin) Name() string {
	return pluginName
}
