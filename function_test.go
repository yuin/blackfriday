package blackfriday

import (
	"fmt"
	"io"
	"strings"
	"testing"
)

func TestFunction(t *testing.T) {
	input := `aaaa {{ sum(1, 2, 3) }} vvvv`
	r := NewHTMLRenderer(HTMLRendererParameters{
		FunctionHandler: func(w io.Writer, name string, args []interface{}) error {
			switch name {
			case "sum":
				total := 0
				for _, v := range args {
					total += int(v.(float64))
				}
				fmt.Fprintf(w, "%d", total)
				return nil
			}
			return DefaultFunctionHandler(w, name, args)
		},
	})
	output := Run([]byte(input), WithExtensions(CommonExtensions|Functions),
		WithRenderer(r))
	if strings.TrimSpace(string(output)) != `<p>aaaa 6 vvvv</p>` {
		t.Errorf("unexpected output: %s", output)
	}
}
