package markdown

import (
	"io"
	"io/fs"
	"strings"

	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"

	"forge.capytal.company/loreddev/x/blogo"
)

const pluginName = "blogo-markdown"

type plugin struct {
	parser   parser.Parser
	renderer renderer.Renderer
}

func New() blogo.Plugin {
	m := goldmark.New(
		goldmark.WithExtensions(
			extension.NewLinkify(),
			meta.Meta,
		),
	)

	return &plugin{
		parser:   m.Parser(),
		renderer: m.Renderer(),
	}
}

func (p *plugin) Name() string {
	return pluginName
}

func (p *plugin) Render(f fs.File, w io.Writer) error {
	stat, err := f.Stat()
	if err != nil || !strings.HasSuffix(stat.Name(), ".md") {
		return blogo.ErrRendererNotSupportedFile
	}

	src, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	txt := text.NewReader(src)

	ast := p.parser.Parse(txt)

	return p.renderer.Render(w, src, ast)
}
