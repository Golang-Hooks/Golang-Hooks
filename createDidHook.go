package hooks

// Returns a function which, when invoked, will return the number of times a hook has been called.
func createDidHook(core *Core, hooks *Hooks) func(string) int {
	return func(hookName string) int {
		if v, ok := hooks.Hooks[hookName]; ok {
			return v.Runs
		}
		return 0
	}
}
