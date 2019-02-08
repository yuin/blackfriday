package blackfriday

import (
	"bytes"
	"github.com/yuin/gluamapper"
	"github.com/yuin/gopher-lua"
	"github.com/yuin/gopher-lua/ast"
	"github.com/yuin/gopher-lua/parse"
)

func function(p *Markdown, data []byte, offset int) (int, *Node) {
	data = data[offset:]
	if len(data) > 2 && data[1] == '{' {
		_, err := parse.Parse(bytes.NewReader(data[2:]), "<bytes>")
		if err == nil {
			return 0, nil
		}
		pe := err.(*parse.Error)
		col := 2 + pe.Pos.Column - 1
		if data[col] == '}' && len(data) >= col && data[col+1] == '}' {
			script := data[2:col]
			stmts, err := parse.Parse(bytes.NewReader(script), "<bytes>")
			if err != nil || len(stmts) != 1 {
				return 0, nil
			}
			funccall, ok := stmts[0].(*ast.FuncCallStmt)
			if !ok {
				return 0, nil
			}
			funcexpr, ok := funccall.Expr.(*ast.FuncCallExpr)
			if !ok {
				return 0, nil
			}
			identexpr, ok := funcexpr.Func.(*ast.IdentExpr)
			if !ok {
				return 0, nil
			}
			name := identexpr.Value
			l := lua.NewState(lua.Options{
				CallStackSize: 10,
				RegistrySize:  128,
			})
			l.SetGlobal(name, l.NewFunction(func(l *lua.LState) int {
				tbl := l.NewTable()
				top := l.GetTop()
				for i := 1; i <= top; i++ {
					tbl.Append(l.Get(i))
				}
				l.Push(tbl)
				return 1
			}))
			if err := l.DoString("return " + string(script)); err != nil {
				return 0, nil
			}
			args := gluamapper.ToGoValue(l.Get(1), gluamapper.Option{
				NameFunc: gluamapper.Id,
			})
			function := NewNode(Function)
			function.Literal = script
			function.FunctionData.Name = name
			function.FunctionData.Arguments = args.([]interface{})
			return col + 2, function
		}
	}
	return 0, nil
}
