package main

import (
    "fmt"
    "go/parser"
    "go/token"
    "go/ast"
    "math/rand"
)

func generateString (length int) string {
    n := rand.Intn(length)
    s := ""
    for i := 0; i < n; i++ {
        s += string(rand.Intn(128))
    }
    return s
}

func generateTests (fun *ast.FuncDecl, testAmt int) (tests [][]interface{}) {
    tests = make([][]interface{}, testAmt)
    paramTypes := []string{}
    for _, param := range fun.Type.Params.List {
        typeName := param.Type.(*ast.Ident).Name
        for range param.Names {
            paramTypes = append(paramTypes, typeName)
        }
    }
    paramAmt := len(paramTypes)
    for i := range tests {
        tests[i] = make([]interface{}, paramAmt)
        for j := range paramTypes {
            switch paramTypes[j] {
            case "int":
                tests[i][j] = int(rand.Uint64())
            case "string":
                tests[i][j] = generateString(42)
            default:
                tests[i][j] = nil
            }
        }
    }
    return
}

func main () {
    fset := token.NewFileSet()
    testFile := "test-files/1.go"
    file, err := parser.ParseFile(fset, testFile, nil, 0)
    if err != nil {
        fmt.Println(err)
        return
    }
    for _, d := range file.Decls {
        if fun, ok := d.(*ast.FuncDecl); ok {
            tests := generateTests(fun, 5)
            for _, test := range tests {
                fmt.Printf("%s(%#v", fun.Name, test[0])
                for i := 1; i < len(test); i++ {
                    fmt.Printf(", %#v", test[i])
                }
                fmt.Println(")")
            }
        }
    }
}
