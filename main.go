package main

import (
    "strings"
    "fmt"
    "go/parser"
    "go/token"
    "go/ast"
    "math/rand"
    "strconv"
    "os"
)

type FuncTest []string

type FuncTests []FuncTest

type FileTest map[string]FuncTests

func GenerateString (length int) (str string) {
    n := rand.Intn(length)
    str = ""
    for i := 0; i < n; i++ {
        str += string(rand.Intn(128))
    }
    return
}

func GenerateFuncTests (testedFunc *ast.FuncDecl, funcTestAmt int) (funcTests FuncTests) {
    funcTests = make([]FuncTest, funcTestAmt)
    paramTypes := []string{}
    for _, param := range testedFunc.Type.Params.List {
        typeName := param.Type.(*ast.Ident).Name
        for range param.Names {
            paramTypes = append(paramTypes, typeName)
        }
    }
    paramAmt := len(paramTypes)
    for i := range funcTests {
        funcTests[i] = make(FuncTest, paramAmt)
        for j := range paramTypes {
            switch paramTypes[j] {
            case "int":
                funcTests[i][j] = strconv.Itoa(int(rand.Uint64()))
            case "string":
                funcTests[i][j] = fmt.Sprintf("%#v", GenerateString(42))
            default:
                funcTests[i][j] = "nil"
            }
        }
    }
    return
}

func GenerateFileTest (testedFile string, funcTestAmt int) (fileTest FileTest) {
    fset := token.NewFileSet()
    file, _ := parser.ParseFile(fset, testedFile, nil, 0)
    fileTest = make(FileTest)
    for _, d := range file.Decls {
        if fun, ok := d.(*ast.FuncDecl); ok {
            fileTest[fun.Name.Name] = GenerateFuncTests(fun, funcTestAmt)
        }
    }
    return
}

func GenerateTestingFile (fileTest FileTest) (testingFile string) {
    testingFile = "package test1\n\nimport \"testing\""
    for name, tests := range fileTest {
        testingFile += "\n\nfunc Test_" + name + " (t *testing.T) {\n"
        for i, test := range tests {
            testStr := strings.Join(test, ", ")
            funcStr := name + "(" + testStr + ")"
            iStr := strconv.Itoa(i)
            testingFile +=
                "    t.Run(\"" + iStr + "\", func (t *testing.T) { defer func () { if recover() != nil { t.Fail() } }(); " + funcStr + " })\n"
        }
        testingFile += "}"
    }
    return
}

func main () {
    testedFile := "test-files/test1.go"
    fileTest := GenerateFileTest(testedFile, 5)
    testingFile := GenerateTestingFile(fileTest)
    file, _ := os.OpenFile("test-files/test1_test.go", os.O_WRONLY|os.O_CREATE, 0644)
    defer file.Close()
    file.WriteString(testingFile)
}
