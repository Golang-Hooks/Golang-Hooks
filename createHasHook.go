package hooks

// Returns a function which, when invoked, will return whether a hook exists or not.
func createHasHook(core *Core, hooks *Hooks) func(string) bool {
	return func(hookName string) bool {
		if _, ok := hooks.Hooks[hookName]; ok {
			return true
		}
		return false
	}
}
