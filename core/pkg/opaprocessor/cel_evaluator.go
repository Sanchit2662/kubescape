package opaprocessor

import (
	"fmt"

	"github.com/google/cel-go/cel"
)

// CELEvaluator evaluates a single CEL expression against a Kubernetes resource object.
// It returns true when the expression matches (i.e. a violation is found), false otherwise.
// An error is returned if the expression fails to compile or evaluate.
func CELEvaluator(expression string, resource map[string]interface{}) (bool, error) {
	env, err := cel.NewEnv(
		cel.Variable("object", cel.MapType(cel.StringType, cel.DynType)),
	)
	if err != nil {
		return false, fmt.Errorf("cel env creation failed: %w", err)
	}

	ast, issues := env.Compile(expression)
	if issues != nil && issues.Err() != nil {
		return false, fmt.Errorf("cel compilation failed: %w", issues.Err())
	}

	prg, err := env.Program(ast)
	if err != nil {
		return false, fmt.Errorf("cel program creation failed: %w", err)
	}

	out, _, err := prg.Eval(map[string]interface{}{
		"object": resource,
	})
	if err != nil {
		return false, fmt.Errorf("cel evaluation failed: %w", err)
	}

	result, ok := out.Value().(bool)
	if !ok {
		return false, fmt.Errorf("cel expression must return a bool, got %T", out.Value())
	}

	return result, nil
}
