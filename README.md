# Golang Hooks

An Event Manager, Plugin System, Middleware Manager, Extendability System, Plugin API for Go apps. Have you ever imagined if your Go app has WordPress extendability features? The goal is to make the powerful & simple plugin API of WordPress available in Go!


## Installation

> Go 1.18+ is required.

```bash
go get github.com/Golang-Hooks/Golang-Hooks@v1.0.1
```

## Examples

> The test file `hooks_test.go` covers a wide range of examples.

### Actions Example

```go
package main

import (
	"github.com/Golang-Hooks/Golang-Hooks"
)

func main() {
	h := hooks.CreateHooks()

	h.AddAction("test", "vendor/plugin/function", func(i ...interface{}) interface{} {
		println("just doing something...")
		return nil
	}, 10)

	h.AddAction("test", "vendor/plugin/function", func(i ...interface{}) interface{} {
		arg1 := i[0]
		arg2 := i[1]
		if p, ok := arg1.(int); ok {
			println("arg1", p)
		}
		if p, ok := arg2.(string); ok {
			println("arg2", p)
		}
		return nil
	}, 9)

	h.DoAction("test", 33, "awesome")
}
```

### Filters Example

```go
package main

import (
	"fmt"

	"github.com/Golang-Hooks/Golang-Hooks"
)

func main() {
	h := hooks.CreateHooks()

	h.AddFilter("MyFilter", "vendor/plugin/function", func(i ...interface{}) interface{} {
		arg1 := i[0]
		if p, ok := arg1.(int); ok {
			println("f1", p)
			return p + 1
		}
		return nil
	}, 10)

	h.AddFilter("MyFilter", "vendor/plugin/function", func(i ...interface{}) interface{} {
		arg1 := i[0]
		arg2 := i[1]
		p1, ok1 := arg1.(int)
		p2, ok2 := arg2.(int)
		if !ok1 || !ok2 {
			println("not ok")
			return nil
		}
		println("f2", p1, p2)
		return p1 + p2 + 1
	}, 9)

	v := h.ApplyFilters("MyFilter", 3, 2)
	fmt.Println(v)
}
```

## API Usage

- `CreateHooks()`
- `AddAction("HookName", "namespace", callback, priority)`
- `AddFilter("HookName", "namespace", callback, priority)`
- `RemoveAction("HookName", "namespace")`
- `RemoveFilter("HookName", "namespace")`
- `RemoveAllActions("HookName", "")`
- `RemoveAllFilters("HookName", "")`
- `DoAction("HookName", arg1, arg2, moreArgs, finalArg)`
- `ApplyFilters("HookName", content, arg1, arg2, moreArgs, finalArg)`
- `DoingAction("HookName")`
- `DoingFilter("HookName")`
- `DidAction("HookName")`
- `DidFilter("HookName")`
- `HasAction("HookName")`
- `HasFilter("HookName")`
- `Actions`
- `Filters`

> The namespace is a unique string used to identify the callback, the best practice to make it in the form `vendor/plugin/function`

### Events on action/filter add or remove

Whenever an action or filter is added or removed, a matching `HookAdded` or `HookRemoved` action is triggered.

- `HookAdded` action is triggered when `AddFilter()` or `AddAction()` method is called, passing values for `HookName`, `functionName`, `callback` and `priority`.
- `HookRemoved` action is triggered when `RemoveFilter()` or `RemoveAction()` method is called, passing values for `HookName` and `functionName`.
