package hooks_test

import (
	"fmt"
	"reflect"
	"testing"

	hooks "github.com/Golang-Hooks/Golang-Hooks"
)

var h hooks.Core

var actionValue string

type Arg1 struct {
	a int
}

type Arg2 struct {
	b int
}

func filterA(i ...interface{}) interface{} {
	arg1 := i[0]
	if p, ok := arg1.(string); ok {
		return p + "a"
	}
	return nil
}

func filterB(i ...interface{}) interface{} {
	arg1 := i[0]
	if p, ok := arg1.(string); ok {
		return p + "b"
	}
	return nil
}

func filterC(i ...interface{}) interface{} {
	arg1 := i[0]
	if p, ok := arg1.(string); ok {
		return p + "c"
	}
	return nil
}

func filterCRemovesSelf(i ...interface{}) interface{} {
	h.RemoveFilter("test.filter", "my_callback_filter_c_removes_self")
	arg1 := i[0]
	if p, ok := arg1.(string); ok {
		return p + "b"
	}
	return nil
}

func filterRemovesB(i ...interface{}) interface{} {
	h.RemoveFilter("test.filter", "my_callback_filter_b")
	arg1 := i[0]
	if p, ok := arg1.(string); ok {
		return p
	}
	return nil
}

func filterRemovesC(i ...interface{}) interface{} {
	h.RemoveFilter("test.filter", "my_callback_filter_c")
	arg1 := i[0]
	if p, ok := arg1.(string); ok {
		return p
	}
	return nil
}

func actionA(i ...interface{}) interface{} {
	actionValue += "a"
	return nil
}

func actionB(i ...interface{}) interface{} {
	actionValue += "b"
	return nil
}

func actionC(i ...interface{}) interface{} {
	actionValue += "c"
	return nil
}

// Almost the same as setupSuite, but this one is for single test instead of collection of tests
func setupTest(tb testing.TB) func(tb testing.TB) {
	// Setup
	h = hooks.CreateHooks()
	actionValue = ""

	return func(tb testing.TB) {
		// Teardown
		actionValue = ""
		h = hooks.CreateHooks()
	}
}

// Run a filter with no callbacks
func TestFilterNoCallbacks(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	expected := 42
	v := h.ApplyFilters("test.filter", 42)
	if v != expected {
		t.Errorf("Expected %d to be equal to %d", v, expected)
	}
}

// Add and remove a filter
func TestAddRemoveFilter(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	expected := "42"
	h.AddFilter("test.filter", "my_callback", filterA, 10)
	removeAllFilters := h.RemoveAllFilters("test.filter", "my_callback")
	removeAllFiltersExpected := 1
	if removeAllFilters != removeAllFiltersExpected {
		t.Errorf("Expected %d to be equal to %d", removeAllFilters, removeAllFiltersExpected)
	}

	v := h.ApplyFilters("test.filter", "test")
	expected = "test"
	if v != expected {
		t.Errorf("Expected %s to be equal to %s", v, expected)
	}

	v2 := h.RemoveAllFilters("test.filter", "my_callback")
	expected2 := 0
	if v2 != expected2 {
		t.Errorf("Expected %d to be equal to %d", v2, expected2)
	}
}

// Add a filter and run it
func TestAddFilterRun(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	expected := "testa"
	h.AddFilter("test.filter", "my_callback", filterA, 10)
	v := h.ApplyFilters("test.filter", "test")
	if v != expected {
		t.Errorf("Expected %s to be equal to %s", v, expected)
	}
}

// Add 2 filters in a row and run them
func TestAdd2FiltersRun(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	expected := "testab"
	h.AddFilter("test.filter", "my_callback", filterA, 10)
	h.AddFilter("test.filter", "my_callback", filterB, 10)
	v := h.ApplyFilters("test.filter", "test")
	if v != expected {
		t.Errorf("Expected %s to be equal to %s", v, expected)
	}
}

// Remove a non-existent filter
func TestRemoveNonExistentFilter(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	expected := 0
	v := h.RemoveFilter("test.filter", "my_callback")
	if v != expected {
		t.Errorf("Expected %d to be equal to %d", v, expected)
	}

	v = h.RemoveAllFilters("test.filter", "my_callback")
	if v != expected {
		t.Errorf("Expected %d to be equal to %d", v, expected)
	}
}

// Add 3 filters with different priorities and run them
func TestAdd3FiltersRun(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	expected := "testbca"
	h.AddFilter("test.filter", "my_callback", filterA, 10)
	h.AddFilter("test.filter", "my_callback", filterB, 2)
	h.AddFilter("test.filter", "my_callback", filterC, 8)
	v := h.ApplyFilters("test.filter", "test")
	if v != expected {
		t.Errorf("Expected %s to be equal to %s", v, expected)
	}
}

// Filters with the same and different priorities
func TestAdd3FiltersSamePriorityRun(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	callbacks := make(map[string]func(...interface{}) interface{})

	for p := 1; p < 5; p++ {
		pv := p
		s := []string{"a", "b", "c", "d"}
		for _, v := range s {
			sv := v
			callbacks[fmt.Sprintf("fn_%d%s", pv, sv)] = func(i ...interface{}) interface{} {
				arg1 := i[0]
				if par, ok := arg1.([]string); ok {
					return append(par, fmt.Sprintf("%d%s", pv, sv))
				}
				return nil
			}
		}
	}

	h.AddFilter("test_order", "my_callback_fn_3a", callbacks["fn_3a"], 3)
	h.AddFilter("test_order", "my_callback_fn_3b", callbacks["fn_3b"], 3)
	h.AddFilter("test_order", "my_callback_fn_3c", callbacks["fn_3c"], 3)
	h.AddFilter("test_order", "my_callback_fn_2a", callbacks["fn_2a"], 2)
	h.AddFilter("test_order", "my_callback_fn_2b", callbacks["fn_2b"], 2)
	h.AddFilter("test_order", "my_callback_fn_2c", callbacks["fn_2c"], 2)

	v := h.ApplyFilters("test_order", []string{})

	expected := []string{"2a", "2b", "2c", "3a", "3b", "3c"}
	if !reflect.DeepEqual(v, expected) {
		t.Errorf("Expected %s to be equal to %s", v, expected)
	}

	h.RemoveFilter("test_order", "my_callback_fn_2b")
	h.RemoveFilter("test_order", "my_callback_fn_3a")

	v = h.ApplyFilters("test_order", []string{})

	expected = []string{"2a", "2c", "3b", "3c"}
	if !reflect.DeepEqual(v, expected) {
		t.Errorf("Expected %s to be equal to %s", v, expected)
	}

	h.AddFilter("test_order", "my_callback_fn_4a", callbacks["fn_4a"], 4)
	h.AddFilter("test_order", "my_callback_fn_4b", callbacks["fn_4b"], 4)
	h.AddFilter("test_order", "my_callback_fn_1a", callbacks["fn_1a"], 1)
	h.AddFilter("test_order", "my_callback_fn_4c", callbacks["fn_4c"], 4)
	h.AddFilter("test_order", "my_callback_fn_1b", callbacks["fn_1b"], 1)
	h.AddFilter("test_order", "my_callback_fn_3d", callbacks["fn_3d"], 3)
	h.AddFilter("test_order", "my_callback_fn_4d", callbacks["fn_4d"], 4)
	h.AddFilter("test_order", "my_callback_fn_1c", callbacks["fn_1c"], 1)
	h.AddFilter("test_order", "my_callback_fn_2d", callbacks["fn_2d"], 2)
	h.AddFilter("test_order", "my_callback_fn_1d", callbacks["fn_1d"], 1)

	v = h.ApplyFilters("test_order", []string{})

	expected = []string{"1a", "1b", "1c", "1d", "2a", "2c", "2d", "3b", "3c", "3d", "4a", "4b", "4c", "4d"}
	if !reflect.DeepEqual(v, expected) {
		t.Errorf("Expected %s to be equal to %s", v, expected)
	}
}

// Add and remove an action
func TestAddRemoveAction(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	expected := 1
	h.AddAction("test.action", "my_callback", actionA, 10)
	v := h.RemoveAllActions("test.action", "my_callback")
	if v != expected {
		t.Errorf("Expected %d to be equal to %d", v, expected)
	}

	v2 := h.DoAction("test.action")
	if v2 != nil {
		t.Errorf("Expected %p to be equal to %p", v2, interface{}(nil))
	}

	if actionValue != "" {
		t.Errorf("Expected %s to be equal to %s", actionValue, "")
	}
}

// Add an action and run it
func TestAddActionRun(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	expected := "a"
	h.AddAction("test.action", "my_callback", actionA, 10)
	h.DoAction("test.action")
	if actionValue != expected {
		t.Errorf("Expected %s to be equal to %s", actionValue, expected)
	}
}

// Add 2 actions in a row and then run them
func TestAdd2ActionsRun(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	expected := "ab"
	h.AddAction("test.action", "my_callback", actionA, 10)
	h.AddAction("test.action", "my_callback", actionB, 10)
	h.DoAction("test.action")
	if actionValue != expected {
		t.Errorf("Expected %s to be equal to %s", actionValue, expected)
	}
}

// Add 3 actions with different priorities and run them
func TestAdd3ActionsRun(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	expected := "bca"
	h.AddAction("test.action", "my_callback", actionA, 10)
	h.AddAction("test.action", "my_callback", actionB, 2)
	h.AddAction("test.action", "my_callback", actionC, 8)
	h.DoAction("test.action")
	if actionValue != expected {
		t.Errorf("Expected %s to be equal to %s", actionValue, expected)
	}
}

// Pass in two arguments to an action
func TestAddActionArgs(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	arg1 := Arg1{a: 1}
	arg2 := Arg2{b: 1}

	h.AddAction("test.action", "my_callback", func(i ...interface{}) interface{} {
		if p, ok := i[0].(Arg1); ok {
			expected := 1
			if p.a != expected {
				t.Errorf("Expected %d to be equal to %d", p.a, expected)
			}
		} else {
			t.Errorf("Expected %d to be equal to %d", p, Arg1{})
		}
		if p, ok := i[1].(Arg2); ok {
			expected := 1
			if p.b != expected {
				t.Errorf("Expected %d to be equal to %d", p.b, expected)
			}
		} else {
			t.Errorf("Expected %d to be equal to %d", p, Arg2{})
		}
		return nil
	}, 10)

	h.DoAction("test.action", arg1, arg2)
}

// Fire action multiple times
func TestAddActionFire(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	act1 := func(i ...interface{}) interface{} {
		if true != true {
			t.Errorf("Expected %t to be equal to %t", true, true)
		}
		return nil
	}

	h.AddAction("test.action", "my_callback", act1, 10)
	h.DoAction("test.action")
	h.DoAction("test.action")
}

// Add a filter before the one currently executing
func TestAddFilterBefore(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	h.AddFilter("test.filter", "my_callback", func(i ...interface{}) interface{} {
		h.AddFilter("test.filter", "my_callback", func(j ...interface{}) interface{} {
			if p, ok := j[0].(string); ok {
				return p + "a"
			}
			return nil
		}, 1)

		if p, ok := i[0].(string); ok {
			return p + "b"
		}

		return nil
	}, 2)

	expected := "test_b"
	v := h.ApplyFilters("test.filter", "test_")
	if v != expected {
		t.Errorf("Expected %s to be equal to %s", v, expected)
	}
}

// Add a filter after the one currently executing
func TestAddFilterAfter(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	h.AddFilter("test.filter", "my_callback", func(i ...interface{}) interface{} {
		h.AddFilter("test.filter", "my_callback", func(j ...interface{}) interface{} {
			if p, ok := j[0].(string); ok {
				return p + "b"
			}
			return nil
		}, 2)

		if p, ok := i[0].(string); ok {
			return p + "a"
		}

		return nil
	}, 1)

	expected := "test_ab"
	v := h.ApplyFilters("test.filter", "test_")
	if v != expected {
		t.Errorf("Expected %s to be equal to %s", v, expected)
	}
}

// Add a filter immediately after the one currently executing
func TestAddFilterImmediatelyAfter(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	h.AddFilter("test.filter", "my_callback", func(i ...interface{}) interface{} {
		h.AddFilter("test.filter", "my_callback", func(j ...interface{}) interface{} {
			if p, ok := j[0].(string); ok {
				return p + "b"
			}
			return nil
		}, 1)

		if p, ok := i[0].(string); ok {
			return p + "a"
		}

		return nil
	}, 1)

	expected := "test_ab"
	v := h.ApplyFilters("test.filter", "test_")
	if v != expected {
		t.Errorf("Expected %s to be equal to %s", v, expected)
	}
}

// Remove specific action callback
func TestRemoveActionCallback(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	h.AddAction("test.action", "my_callback_action_a", actionA, 10)
	h.AddAction("test.action", "my_callback_action_b", actionB, 2)
	h.AddAction("test.action", "my_callback_action_b", actionC, 8)

	expected := 1
	ra := h.RemoveAction("test.action", "my_callback_action_b")
	if ra != expected {
		t.Errorf("Expected %d to be equal to %d", ra, expected)
	}

	h.DoAction("test.action")
	expected2 := "ca"
	if actionValue != expected2 {
		t.Errorf("Expected %s to be equal to %s", actionValue, expected2)
	}
}

// Remove all action callbacks
func TestRemoveActionAll(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	h.AddAction("test.action", "my_callback_action_a", actionA, 10)
	h.AddAction("test.action", "my_callback_action_b", actionB, 2)
	h.AddAction("test.action", "my_callback_action_c", actionC, 8)

	expected := 3
	ra := h.RemoveAllActions("test.action", "")
	if ra != expected {
		t.Errorf("Expected %d to be equal to %d", ra, expected)
	}

	h.DoAction("test.action")
	expected2 := ""
	if actionValue != expected2 {
		t.Errorf("Expected %s to be equal to %s", actionValue, expected2)
	}
}

// Remove specific filter callback
func TestRemoveFilterCallback(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	h.AddFilter("test.filter", "my_callback_filter_a", filterA, 10)
	h.AddFilter("test.filter", "my_callback_filter_b", filterB, 2)
	h.AddFilter("test.filter", "my_callback_filter_c", filterC, 8)

	expected := 1
	ra := h.RemoveFilter("test.filter", "my_callback_filter_b")
	if ra != expected {
		t.Errorf("Expected %d to be equal to %d", ra, expected)
	}

	v := h.ApplyFilters("test.filter", "test")
	expected2 := "testca"
	if v != expected2 {
		t.Errorf("Expected %s to be equal to %s", v, expected2)
	}
}

// Filter removes a callback that has already executed
func TestRemoveFilterCallbackAfter(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	h.AddFilter("test.filter", "my_callback_filter_a", filterA, 1)
	h.AddFilter("test.filter", "my_callback_filter_b", filterB, 3)
	h.AddFilter("test.filter", "my_callback_filter_c", filterC, 5)
	h.AddFilter("test.filter", "my_callback_filter_removes_b", filterRemovesB, 4)

	v := h.ApplyFilters("test.filter", "test")
	expected := "testabc"
	if v != expected {
		t.Errorf("Expected %s to be equal to %s", v, expected)
	}
}

// Filter removes a callback that has already executed (same priority)
func TestRemoveFilterCallbackAfter2(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	h.AddFilter("test.filter", "my_callback_filter_a", filterA, 1)
	h.AddFilter("test.filter", "my_callback_filter_b", filterB, 2)
	h.AddFilter("test.filter", "my_callback_filter_removes_b", filterRemovesB, 2)
	h.AddFilter("test.filter", "my_callback_filter_c", filterC, 4)

	v := h.ApplyFilters("test.filter", "test")
	expected := "testabc"
	if v != expected {
		t.Errorf("Expected %s to be equal to %s", v, expected)
	}
}

// Filter removes the current callback
func TestRemoveFilterCallbackCurrent(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	h.AddFilter("test.filter", "my_callback_filter_a", filterA, 1)
	h.AddFilter("test.filter", "my_callback_filter_c_removes_self", filterCRemovesSelf, 3)
	h.AddFilter("test.filter", "my_callback_filter_c", filterC, 5)

	v := h.ApplyFilters("test.filter", "test")
	expected := "testabc"
	if v != expected {
		t.Errorf("Expected %s to be equal to %s", v, expected)
	}
}

// Filter removes a callback that has not yet executed (last)
func TestRemoveFilterCallbackLast(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	h.AddFilter("test.filter", "my_callback_filter_a", filterA, 1)
	h.AddFilter("test.filter", "my_callback_filter_b", filterB, 3)
	h.AddFilter("test.filter", "my_callback_filter_c", filterC, 5)
	h.AddFilter("test.filter", "my_callback_filter_removes_c", filterRemovesC, 4)

	v := h.ApplyFilters("test.filter", "test")
	expected := "testab"
	if v != expected {
		t.Errorf("Expected %s to be equal to %s", v, expected)
	}
}

// Filter removes a callback that has not yet executed (middle)
func TestRemoveFilterCallbackMiddle(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	h.AddFilter("test.filter", "my_callback_filter_a", filterA, 1)
	h.AddFilter("test.filter", "my_callback_filter_b", filterB, 3)
	h.AddFilter("test.filter", "my_callback_filter_c", filterC, 4)
	h.AddFilter("test.filter", "my_callback_filter_removes_b", filterRemovesB, 2)

	v := h.ApplyFilters("test.filter", "test")
	expected := "testac"
	if v != expected {
		t.Errorf("Expected %s to be equal to %s", v, expected)
	}
}

// Filter removes a callback that has not yet executed (same priority)
func TestRemoveFilterCallbackSamePriority(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	h.AddFilter("test.filter", "my_callback_filter_a", filterA, 1)
	h.AddFilter("test.filter", "my_callback_filter_removes_b", filterRemovesB, 2)
	h.AddFilter("test.filter", "my_callback_filter_b", filterB, 2)
	h.AddFilter("test.filter", "my_callback_filter_c", filterC, 4)

	v := h.ApplyFilters("test.filter", "test")
	expected := "testac"
	if v != expected {
		t.Errorf("Expected %s to be equal to %s", v, expected)
	}
}

// Remove all filter callbacks
func TestRemoveAllFilterCallbacks(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	h.AddFilter("test.filter", "my_callback_filter_a", filterA, 10)
	h.AddFilter("test.filter", "my_callback_filter_b", filterB, 2)
	h.AddFilter("test.filter", "my_callback_filter_c", filterC, 8)

	expected := 3
	ra := h.RemoveAllFilters("test.filter", "")
	if ra != expected {
		t.Errorf("Expected %d to be equal to %d", ra, expected)
	}

	v := h.ApplyFilters("test.filter", "test")
	expected2 := "test"
	if v != expected2 {
		t.Errorf("Expected %s to be equal to %s", v, expected2)
	}
}

// Test DoingAction, DidAction, HasAction.
func TestDoingActionDidActionHasAction(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	actionCalls := 0

	h.AddAction("another.action", "my_callback", func(j ...interface{}) interface{} { return nil }, 10)
	h.DoAction("another.action")

	// Verify no action is running yet.
	if h.DoingAction("test.action") {
		t.Errorf("Expected action to not be running.")
	}

	if h.DidAction("test.action") != 0 {
		t.Errorf("Expected no action run.")
	}

	if h.HasAction("test.action") {
		t.Errorf("Expected action to not exist.")
	}

	h.AddAction("test.action", "my_callback", func(j ...interface{}) interface{} {
		actionCalls++

		// Expected current action to be test.action
		if hi, _ := h.CurrentAction(); hi.Name != "test.action" {
			t.Errorf("Expected current action to be test.action.")
		}

		if !h.DoingAction("") {
			t.Errorf("Expected action to be running.")
		}

		if !h.DoingAction("test.action") {
			t.Errorf("Expected test.action action to be running.")
		}
		return nil
	}, 10)

	// Verify action added, not running yet.
	if h.DoingAction("test.action") {
		t.Errorf("Expected action to not be running.")
	}

	if h.DidAction("test.action") != 0 {
		fmt.Println(h.DidAction("test.action"))
		t.Errorf("Expected no action run.")
	}

	if !h.HasAction("test.action") {
		t.Errorf("Expected action to be exist.")
	}

	// Run action.
	h.DoAction("test.action")

	// Verify action added and running.
	if actionCalls != 1 {
		t.Errorf("Expected action to be called once.")
	}

	if h.DoingAction("test.action") {
		t.Errorf("Expected action to not be running.")
	}

	// DidAction should return 1
	if h.DidAction("test.action") != 1 {
		t.Errorf("Expected action to be called once.")
	}

	// HasAction should return true
	if !h.HasAction("test.action") {
		t.Errorf("Expected action to be exist.")
	}

	// DoingAction with empty string should return false
	if h.DoingAction("") {
		t.Errorf("Expected action to not be running.")
	}

	// No action with "notatest.action" name is running
	if h.DoingAction("notatest.action") {
		t.Errorf("Expected notatest.action action to not be running.")
	}

	if _, err := h.CurrentAction(); err == nil {
		t.Errorf("Expected no current action.")
	}

	h.DoAction("test.action")

	// Verify actionCalls
	expected := 2
	if actionCalls != expected {
		t.Errorf("Expected actionCalls to be %d.", expected)
	}

	// Verify DidAction
	expected = 2
	if h.DidAction("test.action") != expected {
		t.Errorf("Expected DidAction to be %d.", expected)
	}

	v := h.RemoveAllActions("test.action", "")
	expected = 1
	if v != expected {
		t.Errorf("Expected %d to be equal to %d", v, expected)
	}

	// Verify state is reset appropriately.
	if h.DoingAction("test.action") {
		t.Errorf("Expected action to not be running.")
	}

	expected = 2
	if v := h.DidAction("test.action"); v != expected {
		t.Errorf("Expected %d to be equal to %d", v, expected)
	}

	if !h.HasAction("test.action") {
		t.Errorf("Expected action to be exist.")
	}

	h.DoAction("another.action")

	if h.DoingAction("test.action") {
		t.Errorf("Expected action to not be running.")
	}

	// Verify an action with no handlers is still counted
	if h.DidAction("unattached.action") != 0 {
		t.Errorf("Expected unattached.action action to not be run.")
	}

	h.DoAction("unattached.action")

	if h.DoingAction("unattached.action") {
		t.Errorf("Expected unattached.action action to not be running.")
	}

	if h.DidAction("unattached.action") != 1 {
		t.Errorf("Expected unattached.action action to be run.")
	}

	h.DoAction("unattached.action")

	if h.DoingAction("unattached.action") {
		t.Errorf("Expected unattached.action action to not be running.")
	}

	expected = 2
	if h.DidAction("unattached.action") != expected {
		t.Errorf("Expected unattached.action action to be run %d times.", expected)
	}

	// Verify hasAction returns false when no matching action.
	if h.HasAction("notatest.action") {
		t.Errorf("Expected notatest.action action to not exist.")
	}
}

// Verify doingFilter, didFilter and hasFilter.
func TestDoingFilterDidFilterHasFilter(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	filterCalls := 0

	h.AddFilter("runtest.filter", "my_callback", func(arg ...interface{}) interface{} {
		filterCalls++

		if hi, _ := h.CurrentFilter(); hi.Name != "runtest.filter" {
			t.Errorf("Expected current filter to be runtest.filter.")
		}

		if !h.DoingFilter("") {
			t.Errorf("Expected filter to be running.")
		}

		if !h.DoingFilter("runtest.filter") {
			t.Errorf("Expected runtest.filter filter to be running.")
		}

		return arg[0]
	}, 10)

	// Verify filter added and running.
	test := h.ApplyFilters("runtest.filter", "someValue")

	if test != "someValue" {
		t.Errorf("Expected filter to return someValue.")
	}

	if filterCalls != 1 {
		t.Errorf("Expected filter to be called once.")
	}

	if h.DidFilter("runtest.filter") != 1 {
		t.Errorf("Expected DidFilter to be 1.")
	}

	if !h.HasFilter("runtest.filter") {
		t.Errorf("Expected filter to be exist.")
	}

	if h.HasFilter("notatest.filter") {
		t.Errorf("Expected notatest.filter filter to not exist.")
	}

	if h.DoingFilter("") {
		t.Errorf("Expected filter to not be running.")
	}

	if h.DoingFilter("runtest.filter") {
		t.Errorf("Expected runtest.filter filter to not be running.")
	}

	if h.DoingFilter("notatest.filter") {
		t.Errorf("Expected notatest.filter filter to not be running.")
	}

	if _, err := h.CurrentFilter(); err == nil {
		t.Errorf("Expected no current filter.")
	}

	expected := 1
	if v := h.RemoveAllFilters("runtest.filter", ""); v != expected {
		t.Errorf("Expected %d to be equal to %d", v, expected)
	}

	if !h.HasFilter("runtest.filter") {
		t.Errorf("Expected runtest.filter filter to not exist.")
	}

	expected = 1
	if v := h.DidFilter("runtest.filter"); v != expected {
		t.Errorf("Expected DidFilter to be %d.", expected)
	}
}

// Recursively calling a filter
func TestFilterRecursion(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	h.AddFilter("test.filter", "my_callback", func(args ...interface{}) interface{} {
		arg1 := args[0]
		if p, ok := arg1.(string); ok {
			if len(p) == 7 {
				return args[0]
			}
			return h.ApplyFilters("test.filter", p+"X")
		}

		return nil
	}, 10)

	expected := "testXXX"
	test := h.ApplyFilters("test.filter", "testXXX")
	if test != expected {
		t.Errorf("Expected %s to be equal to %s", test, expected)
	}
}

// Current filter when multiple filters are running
func TestCurretFilterWithMultipleFiltersRunning(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	h.AddFilter("test.filter1", "my_callback", func(args ...interface{}) interface{} {
		arg1 := args[0]
		if p, ok := arg1.([]string); ok {
			cf, _ := h.CurrentFilter()
			return h.ApplyFilters("test.filter2", append(p, cf.Name))
		}

		return nil
	}, 10)

	h.AddFilter("test.filter2", "my_callback", func(args ...interface{}) interface{} {
		arg1 := args[0]
		if p, ok := arg1.([]string); ok {
			cf, _ := h.CurrentFilter()
			return append(p, cf.Name)
		}

		return nil
	}, 10)

	if _, err := h.CurrentFilter(); err == nil {
		t.Errorf("Expected no current filter.")
	}

	v := h.ApplyFilters("test.filter1", []string{"test"})
	expected := []string{"test", "test.filter1", "test.filter2"}
	if !reflect.DeepEqual(v, expected) {
		t.Errorf("Expected %v to be equal to %v", v, expected)
	}

	if _, err := h.CurrentFilter(); err == nil {
		t.Errorf("Expected no current filter.")
	}
}

// Adding and removing filters with recursion
func TestAddRemoveFiltersWithRecursion(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	var removeRecurseAndAdd2 func(...interface{}) interface{}

	removeRecurseAndAdd2 = func(val ...interface{}) interface{} {
		expected := 1
		v := h.RemoveFilter("remove_and_add", "my_callback_recurse")
		if v != expected {
			t.Errorf("Expected %d to be equal to %d", v, expected)
		}

		arg1 := val[0]
		if p, ok := arg1.(string); ok {
			fv := h.ApplyFilters("remove_and_add", "")
			p2, ok2 := fv.(string)
			if !ok2 {
				return nil
			}
			p += "-" + p2 + "-"
			h.AddFilter("remove_and_add", "my_callback_recurse", removeRecurseAndAdd2, 10)
			return p + "2"
		}

		return nil
	}

	h.AddFilter("remove_and_add", "my_callback", func(i ...interface{}) interface{} {
		arg1 := i[0]
		if p, ok := arg1.(string); ok {
			return p + "1"
		}

		return nil
	}, 11)
	h.AddFilter("remove_and_add", "my_callback_recurse", removeRecurseAndAdd2, 12)
	h.AddFilter("remove_and_add", "my_callback", func(i ...interface{}) interface{} {
		arg1 := i[0]
		if p, ok := arg1.(string); ok {
			return p + "3"
		}

		return nil
	}, 13)
	h.AddFilter("remove_and_add", "my_callback", func(i ...interface{}) interface{} {
		arg1 := i[0]
		if p, ok := arg1.(string); ok {
			return p + "4"
		}

		return nil
	}, 14)

	expected := "1-134-234"
	test := h.ApplyFilters("remove_and_add", "")
	if test != expected {
		t.Errorf("Expected %s to be equal to %s", test, expected)
	}
}

// Actions preserve arguments across handlers without return value
func TestActionsPreserveArguments(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	arg1 := Arg1{a: 10}
	arg2 := Arg2{b: 20}

	h.AddAction("test.action", "my_callback1", func(i ...interface{}) interface{} {
		if p, ok := i[0].(Arg1); ok {
			expected := 10
			if p.a != expected {
				t.Errorf("Expected %d to be equal to %d", p.a, expected)
			}
		} else {
			t.Errorf("Expected %d to be equal to %d", p, Arg1{})
		}
		if p, ok := i[1].(Arg2); ok {
			expected := 20
			if p.b != expected {
				t.Errorf("Expected %d to be equal to %d", p.b, expected)
			}
		} else {
			t.Errorf("Expected %d to be equal to %d", p, Arg2{})
		}

		return nil
	}, 10)

	h.AddAction("test.action", "my_callback2", func(i ...interface{}) interface{} {
		if p, ok := i[0].(Arg1); ok {
			expected := 10
			if p.a != expected {
				t.Errorf("Expected %d to be equal to %d", p.a, expected)
			}
		} else {
			t.Errorf("Expected %d to be equal to %d", p, Arg1{})
		}
		if p, ok := i[1].(Arg2); ok {
			expected := 20
			if p.b != expected {
				t.Errorf("Expected %d to be equal to %d", p.b, expected)
			}
		} else {
			t.Errorf("Expected %d to be equal to %d", p, Arg2{})
		}

		return nil
	}, 10)

	h.DoAction("test.action", arg1, arg2)
}

// Filters pass first argument across handlers
func TestFiltersPassArguments(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	h.AddFilter("test.filter", "my_callback1", func(i ...interface{}) interface{} {
		if p, ok := i[0].(int); ok {
			return p + 1
		}
		return nil
	}, 10)

	h.AddFilter("test.filter", "my_callback2", func(i ...interface{}) interface{} {
		if p, ok := i[0].(int); ok {
			return p + 1
		}
		return nil
	}, 10)

	expected := 2
	test := h.ApplyFilters("test.filter", 0)
	if test != expected {
		t.Errorf("Expected %d to be equal to %d", test, expected)
	}
}

// Adding an action triggers a hookAdded action
func TestAddActionTriggersHookAdded(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	numberOfCalls := 0

	fn := func(i ...interface{}) interface{} {
		numberOfCalls++
		return nil
	}

	h.AddAction("HookAdded", "my_callback", fn, 10)
	h.AddAction("testAction", "my_callback2", actionA, 9)

	expected := 1
	if numberOfCalls != expected {
		t.Errorf("Expected %d to be equal to %d", numberOfCalls, expected)
	}
}

// Adding a filter triggers a hookAdded action
func TestAddFilterTriggersHookAdded(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	numberOfCalls := 0

	fn := func(i ...interface{}) interface{} {
		numberOfCalls++
		return nil
	}

	h.AddAction("HookAdded", "my_callback", fn, 10)
	h.AddFilter("testFilter", "my_callback2", filterA, 8)

	expected := 1
	if numberOfCalls != expected {
		t.Errorf("Expected %d to be equal to %d", numberOfCalls, expected)
	}
}

// Removing an action triggers a hookRemoved action
func TestRemoveActionTriggersHookRemoved(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	numberOfCalls := 0

	fn := func(i ...interface{}) interface{} {
		numberOfCalls++
		return nil
	}

	h.AddAction("HookRemoved", "my_callback", fn, 10)
	h.AddAction("testAction", "my_callback2", actionA, 9)

	h.RemoveAction("testAction", "my_callback2")

	expected := 1
	if numberOfCalls != expected {
		t.Errorf("Expected %d to be equal to %d", numberOfCalls, expected)
	}
}

// Removing a filter triggers a hookRemoved action
func TestRemoveFilterTriggersHookRemoved(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	numberOfCalls := 0

	fn := func(i ...interface{}) interface{} {
		numberOfCalls++
		return nil
	}

	h.AddAction("HookRemoved", "my_callback", fn, 10)
	h.AddFilter("testFilter", "my_callback3", filterA, 8)

	h.RemoveFilter("testFilter", "my_callback3")

	expected := 1
	if numberOfCalls != expected {
		t.Errorf("Expected %d to be equal to %d", numberOfCalls, expected)
	}
}
