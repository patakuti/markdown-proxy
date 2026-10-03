package markdown

import (
	"fmt"
	"strings"
	"testing"

	"github.com/yuin/goldmark/ast"
)

func TestGFMIDs_Duplicates(t *testing.T) {
	ids := newGFMIDs()
	gen := func(s string) string { return string(ids.Generate([]byte(s), ast.KindHeading)) }

	got := []string{gen("User"), gen("User"), gen("User"), gen("Claude"), gen("Claude"), gen("User")}
	want := []string{"user", "user-1", "user-2", "claude", "claude-1", "user-3"}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("#%d = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestGFMIDs_SkipsExistingSuffix(t *testing.T) {
	ids := newGFMIDs()
	gen := func(s string) string { return string(ids.Generate([]byte(s), ast.KindHeading)) }

	// A heading whose own ID looks like a generated one must not be reused.
	got := []string{gen("a"), gen("a-2"), gen("a"), gen("a"), gen("a")}
	want := []string{"a", "a-2", "a-1", "a-3", "a-4"}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("#%d = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestGFMIDs_AllUnique(t *testing.T) {
	ids := newGFMIDs()
	seen := map[string]bool{}
	for i := 0; i < 5000; i++ {
		id := string(ids.Generate([]byte("User"), ast.KindHeading))
		if seen[id] {
			t.Fatalf("duplicate ID %q at %d", id, i)
		}
		seen[id] = true
	}
}

// A conversation log has hundreds of identical headings; heading ID
// generation must stay linear in their number.
func BenchmarkConvert_DuplicateHeadings(b *testing.B) {
	var sb strings.Builder
	for i := 0; i < 10000; i++ {
		fmt.Fprintf(&sb, "## User\n\nquestion %d\n\n## Claude\n\nanswer %d\n\n", i, i)
	}
	src := []byte(sb.String())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := Convert(src, ""); err != nil {
			b.Fatal(err)
		}
	}
}
