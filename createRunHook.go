package hooks

// Returns a function which, when invoked, will execute all callbacks
// registered to a hook of the specified type, optionally returning the final
// value of the call chain.
func createRunHook(core *Core, hooks *Hooks, returnFirstArg bool) func(string, ...interface{}) interface{} {
	return func(hookName string, args ...interface{}) interface{} {

		// Increase Runs by 1
		if entry, ok := hooks.Hooks[hookName]; ok {
			entry.Runs++
			hooks.Hooks[hookName] = entry
		} else {
			hooks.Hooks[hookName] = Handlers{
				Handlers: []Handler{},
				Runs: 1,
			}
		}

		if len(hooks.Hooks[hookName].Handlers) == 0 {
			if returnFirstArg {
				return args[0]
			}
			return nil
		}

		hookInfo := HookInfo{
			Name:         hookName,
			CurrentIndex: 0,
		}

		// append hookInfo to the end of the slice
		hooks.Current = append(hooks.Current, &hookInfo)

		var result interface{}

		for hookInfo.CurrentIndex < len(hooks.Hooks[hookName].Handlers) {
			handler := hooks.Hooks[hookName].Handlers[hookInfo.CurrentIndex]
			result = handler.Callback(args...)
			if returnFirstArg {
				args[0] = result
			}
			hookInfo.CurrentIndex++
		}

		// Remove the last element
		hooks.Current = hooks.Current[:len(hooks.Current)-1]

		if returnFirstArg {
			return args[0]
		}

		return nil
	}
}
