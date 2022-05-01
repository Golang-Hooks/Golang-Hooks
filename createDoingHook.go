package hooks

// createDoingHook Returns a function which, when invoked, will return whether a hook is currently being executed.
func createDoingHook(core *Core, hooks *Hooks) func(string) bool {
	return func(hookName string) bool {
		if len(hooks.Current) > 0 {
			// If the hookName was not passed
			// or if current hook is the same as the hook we're looking for
			if  hookName == "" || hooks.Current[0].Name == hookName {
				return true
			}
		}

		return false
	}
}
