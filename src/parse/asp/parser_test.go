package asp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseBasic(t *testing.T) {
	statements, err := newParser().parse("src/parse/asp/test_data/basic.build")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(statements))
	assert.NotNil(t, statements[0].FuncDef)
	assert.Equal(t, "test", statements[0].FuncDef.Name)
	assert.Equal(t, 1, len(statements[0].FuncDef.Arguments))
	assert.Equal(t, "x", statements[0].FuncDef.Arguments[0].Name)
	assert.Equal(t, 1, len(statements[0].FuncDef.Statements))
	assert.Equal(t, "pass", statements[0].FuncDef.Statements[0].Pass)
}

func TestParseDefaultArguments(t *testing.T) {
	statements, err := newParser().parse("src/parse/asp/test_data/default_arguments.build")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(statements))
	assert.NotNil(t, statements[0].FuncDef)
	assert.Equal(t, "test", statements[0].FuncDef.Name)
	assert.Equal(t, 3, len(statements[0].FuncDef.Arguments))
	assert.Equal(t, 1, len(statements[0].FuncDef.Statements))
	assert.Equal(t, "pass", statements[0].FuncDef.Statements[0].Pass)

	args := statements[0].FuncDef.Arguments
	assert.Equal(t, "name", args[0].Name)
	assert.Equal(t, "\"name\"", args[0].Value.String)
	assert.Equal(t, "timeout", args[1].Name)
	assert.Equal(t, 10, args[1].Value.Int.Int)
	assert.Equal(t, "args", args[2].Name)
	assert.Equal(t, "None", args[2].Value.Bool)
}

func TestParseFunctionCalls(t *testing.T) {
	statements, err := newParser().parse("src/parse/asp/test_data/function_call.build")
	assert.NoError(t, err)
	assert.Equal(t, 5, len(statements))

	assert.NotNil(t, statements[0].Ident.Action.Call)
	assert.Equal(t, "package", statements[0].Ident.Name)
	assert.Equal(t, 0, len(statements[0].Ident.Action.Call.Arguments))

	assert.NotNil(t, statements[2].Ident.Action.Call)
	assert.Equal(t, "package", statements[2].Ident.Name)
	assert.Equal(t, 1, len(statements[2].Ident.Action.Call.Arguments))
	arg := statements[2].Ident.Action.Call.Arguments[0]
	assert.Equal(t, "default_visibility", arg.Expr.Ident.Name)
	assert.NotNil(t, arg.Value)
	assert.Equal(t, 1, len(arg.Value.List.Values))
	assert.Equal(t, "\"PUBLIC\"", arg.Value.List.Values[0].String)

	assert.NotNil(t, statements[3].Ident.Action.Call)
	assert.Equal(t, "python_library", statements[3].Ident.Name)
	assert.Equal(t, 2, len(statements[3].Ident.Action.Call.Arguments))
	args := statements[3].Ident.Action.Call.Arguments
	assert.Equal(t, "name", args[0].Expr.Ident.Name)
	assert.Equal(t, "\"lib\"", args[0].Value.String)
	assert.Equal(t, "srcs", args[1].Expr.Ident.Name)
	assert.NotNil(t, args[1].Value.List)
	assert.Equal(t, 2, len(args[1].Value.List.Values))
	assert.Equal(t, "\"lib1.py\"", args[1].Value.List.Values[0].String)
	assert.Equal(t, "\"lib2.py\"", args[1].Value.List.Values[1].String)

	assert.NotNil(t, statements[4].Ident.Action.Call)
	assert.Equal(t, "subinclude", statements[4].Ident.Name)
	assert.NotNil(t, statements[4].Ident.Action.Call)
	assert.Equal(t, 1, len(statements[4].Ident.Action.Call.Arguments))
	assert.Equal(t, "\"//build_defs:version\"", statements[4].Ident.Action.Call.Arguments[0].Expr.String)
}

func TestParseAssignments(t *testing.T) {
	statements, err := newParser().parse("src/parse/asp/test_data/assignments.build")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(statements))

	assert.NotNil(t, statements[0].Ident.Action.Assign)
	assert.Equal(t, "x", statements[0].Ident.Name)
	ass := statements[0].Ident.Action.Assign.Dict
	assert.NotNil(t, ass)
	assert.Equal(t, 3, len(ass.Items))
	assert.Equal(t, "\"mickey\"", ass.Items[0].Key)
	assert.Equal(t, 3, ass.Items[0].Value.Int.Int)
	assert.Equal(t, "\"donald\"", ass.Items[1].Key)
	assert.Equal(t, "\"sora\"", ass.Items[1].Value.String)
	assert.Equal(t, "\"goofy\"", ass.Items[2].Key)
	assert.Equal(t, "riku", ass.Items[2].Value.Ident.Name)
}

func TestForStatement(t *testing.T) {
	statements, err := newParser().parse("src/parse/asp/test_data/for_statement.build")
	assert.NoError(t, err)
	assert.Equal(t, 2, len(statements))

	assert.NotNil(t, statements[0].Ident.Action.Assign)
	assert.Equal(t, "LANGUAGES", statements[0].Ident.Name)
	assert.Equal(t, 2, len(statements[0].Ident.Action.Assign.List.Values))

	assert.NotNil(t, statements[1].For)
	assert.Equal(t, []string{"language"}, statements[1].For.Names)
	assert.Equal(t, "LANGUAGES", statements[1].For.Expr.Ident.Name)
	assert.Equal(t, 2, len(statements[1].For.Statements))
}

func TestOperators(t *testing.T) {
	statements, err := newParser().parse("src/parse/asp/test_data/operators.build")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(statements))

	assert.NotNil(t, statements[0].Ident.Action.Call)
	assert.Equal(t, "genrule", statements[0].Ident.Name)
	assert.Equal(t, 2, len(statements[0].Ident.Action.Call.Arguments))

	arg := statements[0].Ident.Action.Call.Arguments[1]
	assert.Equal(t, "srcs", arg.Expr.Ident.Name)
	assert.NotNil(t, arg.Value.List)
	assert.Equal(t, 1, len(arg.Value.List.Values))
	assert.Equal(t, "\"//something:test_go\"", arg.Value.List.Values[0].String)
	assert.NotNil(t, arg.Value.Op)
	assert.Equal(t, Add, arg.Value.Op.Op)
	call := arg.Value.Op.Expr.Ident.Action[0].Call
	assert.Equal(t, "glob", arg.Value.Op.Expr.Ident.Name)
	assert.NotNil(t, call)
	assert.Equal(t, 1, len(call.Arguments))
	assert.NotNil(t, call.Arguments[0].Expr.List)
	assert.Equal(t, 1, len(call.Arguments[0].Expr.List.Values))
	assert.Equal(t, "\"*.go\"", call.Arguments[0].Expr.List.Values[0].String)
}

func TestIndexing(t *testing.T) {
	statements, err := newParser().parse("src/parse/asp/test_data/indexing.build")
	assert.NoError(t, err)
	assert.Equal(t, 6, len(statements))

	assert.Equal(t, "x", statements[0].Ident.Name)
	assert.NotNil(t, statements[0].Ident.Action.Assign)
	assert.Equal(t, "\"test\"", statements[0].Ident.Action.Assign.String)

	assert.Equal(t, "y", statements[1].Ident.Name)
	assert.NotNil(t, statements[1].Ident.Action.Assign)
	assert.Equal(t, "x", statements[1].Ident.Action.Assign.Ident.Name)
	assert.NotNil(t, statements[1].Ident.Action.Assign.Slice)
	assert.Equal(t, 2, statements[1].Ident.Action.Assign.Slice.Start.Int.Int)
	assert.Equal(t, "", statements[1].Ident.Action.Assign.Slice.Colon)
	assert.Nil(t, statements[1].Ident.Action.Assign.Slice.End)

	assert.Equal(t, "z", statements[2].Ident.Name)
	assert.NotNil(t, statements[2].Ident.Action.Assign)
	assert.Equal(t, "x", statements[2].Ident.Action.Assign.Ident.Name)
	assert.NotNil(t, statements[2].Ident.Action.Assign.Slice)
	assert.Equal(t, 1, statements[2].Ident.Action.Assign.Slice.Start.Int.Int)
	assert.Equal(t, ":", statements[2].Ident.Action.Assign.Slice.Colon)
	assert.Equal(t, -1, statements[2].Ident.Action.Assign.Slice.End.Int.Int)

	assert.Equal(t, "a", statements[3].Ident.Name)
	assert.NotNil(t, statements[3].Ident.Action.Assign)
	assert.Equal(t, "x", statements[3].Ident.Action.Assign.Ident.Name)
	assert.NotNil(t, statements[3].Ident.Action.Assign.Slice)
	assert.Equal(t, 2, statements[3].Ident.Action.Assign.Slice.Start.Int.Int)
	assert.Equal(t, ":", statements[3].Ident.Action.Assign.Slice.Colon)
	assert.Nil(t, statements[3].Ident.Action.Assign.Slice.End)

	assert.Equal(t, "b", statements[4].Ident.Name)
	assert.NotNil(t, statements[4].Ident.Action.Assign)
	assert.Equal(t, "x", statements[4].Ident.Action.Assign.Ident.Name)
	assert.NotNil(t, statements[4].Ident.Action.Assign.Slice)
	assert.Nil(t, statements[4].Ident.Action.Assign.Slice.Start)
	assert.Equal(t, ":", statements[4].Ident.Action.Assign.Slice.Colon)
	assert.Equal(t, 2, statements[4].Ident.Action.Assign.Slice.End.Int.Int)

	assert.Equal(t, "c", statements[5].Ident.Name)
	assert.NotNil(t, statements[5].Ident.Action.Assign)
	assert.Equal(t, "x", statements[5].Ident.Action.Assign.Ident.Name)
	assert.NotNil(t, statements[5].Ident.Action.Assign.Slice)
	assert.Equal(t, "y", statements[5].Ident.Action.Assign.Slice.Start.Ident.Name)
	assert.Equal(t, "", statements[5].Ident.Action.Assign.Slice.Colon)
	assert.Nil(t, statements[5].Ident.Action.Assign.Slice.End)
}

func TestIfStatement(t *testing.T) {
	statements, err := newParser().parse("src/parse/asp/test_data/if_statement.build")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(statements))

	ifs := statements[0].If
	assert.NotNil(t, ifs)
	assert.Equal(t, "condition_a", ifs.Condition.Ident.Name)
	assert.Equal(t, And, ifs.Condition.Op.Op)
	assert.Equal(t, "condition_b", ifs.Condition.Op.Expr.Ident.Name)
	assert.Equal(t, 1, len(ifs.Statements))
	assert.Equal(t, "genrule", ifs.Statements[0].Ident.Name)
}

func TestDoubleUnindent(t *testing.T) {
	statements, err := newParser().parse("src/parse/asp/test_data/double_unindent.build")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(statements))

	assert.NotNil(t, statements[0].For)
	assert.Equal(t, "y", statements[0].For.Names[0])
	assert.Equal(t, "x", statements[0].For.Expr.Ident.Name)
	assert.Equal(t, 1, len(statements[0].For.Statements))

	for2 := statements[0].For.Statements[0].For
	assert.NotNil(t, for2)
	assert.Equal(t, "z", for2.Names[0])
	assert.Equal(t, "y", for2.Expr.Ident.Name)
	assert.Equal(t, 1, len(for2.Statements))
	assert.Equal(t, "genrule", for2.Statements[0].Ident.Name)
}

func TestInlineIf(t *testing.T) {
	statements, err := newParser().parse("src/parse/asp/test_data/inline_if.build")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(statements))

	assert.Equal(t, "x", statements[0].Ident.Name)
	ass := statements[0].Ident.Action.Assign
	assert.NotNil(t, ass)
	assert.NotNil(t, ass.List)
	assert.Equal(t, 1, len(ass.List.Values))
	assert.NotNil(t, ass.If)
	assert.Equal(t, "y", ass.If.Condition.Ident.Name)
	assert.Equal(t, Is, ass.If.Condition.Op.Op)
	assert.Equal(t, "None", ass.If.Condition.Op.Expr.Bool)
	assert.NotNil(t, ass.If.Else.List)
	assert.Equal(t, 1, len(ass.If.Else.List.Values))
}

func TestFunctionDef(t *testing.T) {
	statements, err := newParser().parse("src/parse/asp/test_data/function_def.build")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(statements))
	assert.Equal(t, 4, len(statements[0].FuncDef.Statements))
}

func TestComprehension(t *testing.T) {
	statements, err := newParser().parse("src/parse/asp/test_data/comprehension.build")
	assert.NoError(t, err)
	assert.Equal(t, 2, len(statements))

	assert.NotNil(t, statements[0].Ident.Action.Assign)
	assert.NotNil(t, statements[1].Ident.Action.Assign)
	assert.Equal(t, 1, len(statements[0].Ident.Action.Assign.List.Values))
	assert.NotNil(t, statements[0].Ident.Action.Assign.List.Comprehension)
	assert.NotNil(t, statements[1].Ident.Action.Assign.Dict.Comprehension)
}

func TestMethodsOnLiterals(t *testing.T) {
	statements, err := newParser().parse("src/parse/asp/test_data/literal_methods.build")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(statements))
}

func TestUnaryOp(t *testing.T) {
	statements, err := newParser().parse("src/parse/asp/test_data/unary_op.build")
	assert.NoError(t, err)
	assert.Equal(t, 3, len(statements))

	assert.NotNil(t, statements[0].Ident.Action.Assign.UnaryOp)
	assert.Equal(t, "-", statements[0].Ident.Action.Assign.UnaryOp.Op)
	assert.Equal(t, "len", statements[0].Ident.Action.Assign.UnaryOp.Expr.Ident.Name)
	assert.NotNil(t, statements[1].Ident.Action.Assign.UnaryOp)
	assert.Equal(t, "not", statements[1].Ident.Action.Assign.UnaryOp.Op)
	assert.Equal(t, "x", statements[1].Ident.Action.Assign.UnaryOp.Expr.Ident.Name)
}

func TestAugmentedAssignment(t *testing.T) {
	statements, err := newParser().parse("src/parse/asp/test_data/aug_assign.build")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(statements))
	assert.NotNil(t, statements[0].Ident.Action.AugAssign)
}

func TestElseStatement(t *testing.T) {
	statements, err := newParser().parse("src/parse/asp/test_data/else.build")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(statements))
	assert.NotNil(t, statements[0].If)
	assert.Equal(t, 1, len(statements[0].If.Statements))
	assert.Equal(t, 2, len(statements[0].If.Elif))
	assert.Equal(t, 1, len(statements[0].If.Elif[0].Statements))
	assert.NotNil(t, statements[0].If.Elif[0].Condition)
	assert.Equal(t, 1, len(statements[0].If.Elif[1].Statements))
	assert.NotNil(t, statements[0].If.Elif[1].Condition)
	assert.Equal(t, 1, len(statements[0].If.ElseStatements))
}

func TestDestructuringAssignment(t *testing.T) {
	statements, err := newParser().parse("src/parse/asp/test_data/destructuring_assign.build")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(statements))
	assert.NotNil(t, statements[0].Ident)
	assert.Equal(t, "x", statements[0].Ident.Name)
	assert.NotNil(t, statements[0].Ident.Unpack)
	assert.Equal(t, 1, len(statements[0].Ident.Unpack.Names))
	assert.Equal(t, "y", statements[0].Ident.Unpack.Names[0])
	assert.NotNil(t, statements[0].Ident.Unpack.Expr)
	assert.Equal(t, "something", statements[0].Ident.Unpack.Expr.Ident.Name)
}

func TestMultipleActions(t *testing.T) {
	statements, err := newParser().parse("src/parse/asp/test_data/multiple_action.build")
	assert.NoError(t, err)
	assert.Equal(t, 3, len(statements))
	assert.NotNil(t, statements[0].Ident.Action.Assign)
	assert.Equal(t, "y", statements[0].Ident.Action.Assign.Ident.Name)
}

func TestAssert(t *testing.T) {
	statements, err := newParser().parse("src/parse/asp/test_data/assert.build")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(statements))
}

func TestOptimise(t *testing.T) {
	p := newParser()
	statements, err := p.parse("src/parse/asp/test_data/optimise.build")
	assert.NoError(t, err)
	assert.Equal(t, 4, len(statements))
	statements = p.optimise(statements)
	assert.Equal(t, 3, len(statements))
	assert.NotNil(t, statements[0].FuncDef)
	assert.Equal(t, 0, len(statements[0].FuncDef.Statements))
	assert.NotNil(t, statements[1].FuncDef)
	assert.Equal(t, 0, len(statements[1].FuncDef.Statements))
	assert.NotNil(t, statements[2].FuncDef)
	assert.Equal(t, 1, len(statements[2].FuncDef.Statements))
	ident := statements[2].FuncDef.Statements[0].Ident
	assert.NotNil(t, ident)
	assert.Equal(t, "l", ident.Name)
	assert.NotNil(t, ident.Action.AugAssign)
	assert.NotNil(t, ident.Action.AugAssign.List)
}

func TestMultilineStringQuotes(t *testing.T) {
	statements, err := newParser().parse("src/parse/asp/test_data/multiline_string_quotes.build")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(statements))
	assert.NotNil(t, statements[0].Ident)
	assert.NotNil(t, statements[0].Ident.Action)
	assert.NotNil(t, statements[0].Ident.Action.Assign)
	expected := `"
#include "UnitTest++/UnitTest++.h"
"`
	assert.Equal(t, expected, statements[0].Ident.Action.Assign.String)
}

func TestExample1(t *testing.T) {
	// These tests are specific examples that turned out to fail.
	_, err := newParser().parse("src/parse/asp/test_data/example_1.build")
	assert.NoError(t, err)
}

func TestExample2(t *testing.T) {
	_, err := newParser().parse("src/parse/asp/test_data/example_2.build")
	assert.NoError(t, err)
}

func TestExample3(t *testing.T) {
	_, err := newParser().parse("src/parse/asp/test_data/example_3.build")
	assert.NoError(t, err)
}

func TestExample4(t *testing.T) {
	_, err := newParser().parse("src/parse/asp/test_data/example_4.build")
	assert.NoError(t, err)
}

func TestExample5(t *testing.T) {
	_, err := newParser().parse("src/parse/asp/test_data/example_5.build")
	assert.NoError(t, err)
}

func TestExample6(t *testing.T) {
	_, err := newParser().parse("src/parse/asp/test_data/example_6.build")
	assert.NoError(t, err)
}
