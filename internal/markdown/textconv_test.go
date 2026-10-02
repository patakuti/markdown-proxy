package markdown

import (
	"strings"
	"testing"
)

func TestConvertText(t *testing.T) {
	out := string(ConvertText([]byte("first\n<b>second</b> & more\nthird")))

	for _, want := range []string{
		`<pre class="text-file"><code>`,
		`<span id="L1" class="source-line">first</span>`,
		`<span id="L2" class="source-line">&lt;b&gt;second&lt;/b&gt; &amp; more</span>`,
		`<span id="L3" class="source-line">third</span></code></pre>`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	}
	if strings.Contains(out, "<b>") {
		t.Errorf("text content was not escaped:\n%s", out)
	}
}

func TestConvertText_Empty(t *testing.T) {
	out := string(ConvertText(nil))
	if !strings.Contains(out, `id="L1"`) {
		t.Errorf("empty input should still produce one line: %s", out)
	}
}
