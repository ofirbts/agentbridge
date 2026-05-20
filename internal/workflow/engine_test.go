package workflow

import (
	"reflect"
	"testing"
)

func TestRunTaskSuccess(t *testing.T) {
	engine, err := NewEngine(Config{Mode: "normal", StorePath: t.TempDir()})
	if err != nil {
		t.Fatal(err)
	}
	record, err := engine.RunTask("find top AI infra startups")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if record.Status != "success" {
		t.Fatalf("expected success, got %s", record.Status)
	}
	if len(record.Steps) != 4 {
		t.Fatalf("expected 4 steps, got %d", len(record.Steps))
	}
}

func TestDeterministicSameInputsSameOutputs(t *testing.T) {
	task := "find top AI infra startups"
	dir := t.TempDir()

	first, err := NewEngine(Config{Mode: "deterministic", StorePath: dir})
	if err != nil {
		t.Fatal(err)
	}
	r1, err := first.RunTask(task)
	if err != nil {
		t.Fatalf("first run: %v", err)
	}
	second, err := NewEngine(Config{Mode: "deterministic", StorePath: dir})
	if err != nil {
		t.Fatal(err)
	}
	r2, err := second.RunTask(task)
	if err != nil {
		t.Fatalf("second run: %v", err)
	}
	if r1.ID != r2.ID {
		t.Fatalf("run ids differ: %s vs %s", r1.ID, r2.ID)
	}
	if !reflect.DeepEqual(r1.Result, r2.Result) {
		t.Fatalf("results differ:\nfirst: %#v\nsecond: %#v", r1.Result, r2.Result)
	}
}

func TestExplainPlanDeterministic(t *testing.T) {
	plan := ExplainPlan("task", "deterministic", "mock", "", "")
	if plan.Mode != "deterministic" {
		t.Fatalf("expected deterministic mode")
	}
	if len(plan.Steps) != 4 {
		t.Fatalf("expected 4 steps")
	}
}

func TestEngineWithMCPEndpoint(t *testing.T) {
	engine, err := NewEngine(Config{
		Mode:         "normal",
		StorePath:    t.TempDir(),
		MCPEndpoint:  "http://localhost:8080/mcp",
		MCPTransport: "stub",
	})
	if err != nil {
		t.Fatal(err)
	}
	record, err := engine.RunTask("mcp task")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if record.MCP == "" {
		t.Fatal("expected mcp endpoint on result")
	}
}

func TestNewEngineInvalidMCPTransport(t *testing.T) {
	_, err := NewEngine(Config{
		MCPEndpoint:  "http://localhost:8080/mcp",
		MCPTransport: "invalid",
	})
	if err == nil {
		t.Fatal("expected error for invalid mcp transport")
	}
}
