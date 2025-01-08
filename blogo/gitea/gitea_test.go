package gitea_test

import (
	"io"
	"testing"

	"forge.capytal.company/loreddev/x/blogo"
	"forge.capytal.company/loreddev/x/blogo/gitea"
)

func TestSource(t *testing.T) {
	plugin := gitea.New("loreddev", "x", "https://forge.capytal.company")

	s := plugin.(blogo.SourcerPlugin)

	fs, err := s.Source()
	if err != nil {
		t.Fatalf("Failed to source file system: %s %v", err.Error(), err)
	}

	file, err := fs.Open("blogo/LICENSE")
	if err != nil {
		t.Fatalf("Failed to open file: %s %v", err.Error(), err)
	}

	contents, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("Failed to read contents of file: %s %v", err.Error(), err)
	}

	t.Logf("Successfully read contents of file: %s", string(contents))
}
