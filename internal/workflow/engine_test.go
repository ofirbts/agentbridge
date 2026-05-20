package workflow

import (
	"reflect"
	"testing"
)

func TestRunTaskSuccess(t *testing.T) {
	engine := NewEngine(Config{Mode: "normal", StorePath: t.TempDir()})
	result, err := engine.RunTask("find top AI infra startups")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Status != "success" {
		t.Fatalf("expected success, got %s", result.Status)
	}
	if len(result.Steps) != 4 {
		t.Fatalf("expected 4 steps, got %d", len(result.Steps))
	}
}

func TestDeterministicSameInputsSameOutputs(t *testing.T) {
	task := "find top AI infra startups"
	dir := t.TempDir()

	first, err := NewEngine(Config{Mode: "deterministic", StorePath: dir}).RunTask(task)
	if err != nil {
		t.Fatalf("first run: %v", err)
	}
	second, err := NewEngine(Config{Mode: "deterministic", StorePath: dir}).RunTask(task)
	if err != nil {
		t.Fatalf("second run: %v", err)
	}
	if first.RunID != second.RunID {
		t.Fatalf("run ids differ: %s vs %s", first.RunID, second.RunID)
	}
	if !reflect.DeepEqual(first.Result, second.Result) {
		t.Fatalf("results differ:\nfirst: %#v\nsecond: %#v", first.Result, second.Result)
	}
}

func TestExplainPlanDeterministic(t *testing.T) {
	plan := ExplainPlan("task", "deterministic", "mock", "")
	if plan.Mode != "deterministic" {
		t.Fatalf("expected deterministic mode")
	}
	if len(plan.Steps) != 4 {
		t.Fatalf("expected 4 steps")
	}
}

func TestEngineWithMCPEndpoint(t *testing.T) {
	engine := NewEngine(Config{
		Mode:        "normal",
		StorePath:   t.TempDir(),
		MCPEndpoint: "http://localhost:8080/mcp",
	})
	result, err := engine.RunTask("mcp task")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.MCP == "" {
		t.Fatal("expected mcp endpoint on result")
	}
}
