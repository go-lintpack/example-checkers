package checkers

import (
	"go/ast"

	"github.com/go-lintpack/lintpack"
)

func init() {
	lintpack.AddChecker(&panicNilChecker{})
}

type panicNilChecker struct {
	lintpack.CheckerBase
}

func (c *panicNilChecker) InitDocumentation(d *lintpack.CheckerDoc) {
	d.Summary = "Detects panic(nil) calls"
	d.Details = "Such panic calls are hard to handle during recover."
	d.Before = `panic(nil)`
	d.After = `panic("something meaningful")`
}

func (c *panicNilChecker) VisitExpr(expr ast.Expr) {
	call, ok := expr.(*ast.CallExpr)
	if !ok {
		return
	}
	fn, ok := call.Fun.(*ast.Ident)
	if !ok || fn.Name != "panic" {
		return
	}
	arg, ok := call.Args[0].(*ast.Ident)
	if !ok || arg.Name != "nil" {
		return
	}
	c.warn(expr)
}

func (c *panicNilChecker) warn(cause ast.Node) {
	c.Ctx.Warn(cause, "panic(nil) calls are discouraged")
}
