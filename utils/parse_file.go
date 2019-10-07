package utils

import (
	"fmt"
	"go/ast"
	"time"
)

func (config *GlobalConfig) ParseFile(astFile *ast.File) {
	fmt.Print(fmt.Sprintf(`[%s]--------`, time.Now().Format("2006-01-02 15:04:05")))
	fmt.Printf("%c[%d;%d;%dmparse file: %s%c[0m\n", 0x1B, 0 /*字体*/, 0 /*背景*/, 36 /*前景*/, config.Filepath, 0x1B)
	ast.Inspect(astFile, func(x ast.Node) bool {
		switch x.(type) {
		case *ast.ImportSpec:
			typeSpec := x.(*ast.ImportSpec)
			fmt.Print(fmt.Sprintf(`[%s]----------`, time.Now().Format("2006-01-02 15:04:05")))
			fmt.Printf("%c[%d;%d;%dmfind import: \t%s%c[0m\n", 0x1B, 0 /*字体*/, 0 /*背景*/, 32 /*前景*/, typeSpec.Path.Value, 0x1B)
			config.Import[typeSpec.Path.Value] = func() string {
				if typeSpec.Name != nil {
					return typeSpec.Name.Name
				} else {
					return ""
				}
			}()
		case *ast.ValueSpec:
			typeSpec := x.(*ast.ValueSpec)
			fmt.Print(fmt.Sprintf(`[%s]----------`, time.Now().Format("2006-01-02 15:04:05")))
			fmt.Printf("%c[%d;%d;%dmfind variable: \t%s%c[0m\n", 0x1B, 0 /*字体*/, 0 /*背景*/, 32 /*前景*/, typeSpec.Names[0].Name, 0x1B)
			for key := range typeSpec.Names {
				config.Variable[typeSpec.Names[key].Name] = ParseExpr(typeSpec.Values[key])
			}
		case *ast.TypeSpec:
			typeSpec := x.(*ast.TypeSpec)
			fmt.Println(fmt.Sprintf(`[%s]----------find struct: %s`, time.Now().Format("2006-01-02 15:04:05"), typeSpec.Name.Name))
			config.Database.Tables[ensnake(typeSpec.Name.Name)] = &[]TableField{}
		default:
			return true
		}
		return false
	})
}

func ParseExpr(expr ast.Expr) (value string) {
	switch expr.(type) {
	case *ast.Ident:
		ident := expr.(*ast.Ident)
		return ident.Name
	case *ast.CallExpr:
		ident := expr.(*ast.CallExpr)
		return ParseExpr(ident.Fun) + "()"
	case *ast.BasicLit:
		ident := expr.(*ast.BasicLit)
		return ident.Value
	case *ast.StarExpr:
		starExpr := expr.(*ast.StarExpr)
		return "*" + ParseExpr(starExpr.X)
	case *ast.SelectorExpr:
		selectorExpr := expr.(*ast.SelectorExpr)
		return ParseExpr(selectorExpr.X) + "." + selectorExpr.Sel.Name
	default:
		return value
	}
}
