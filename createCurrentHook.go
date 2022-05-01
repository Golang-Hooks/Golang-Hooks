package hooks

import "errors"

// Returns a function which, when invoked, will return the HookInfo of the currently running hook
// or an error if no hook is currently running.
func createCurrentHook(core *Core, hooks *Hooks) func() (HookInfo, error) {
	return func() (HookInfo, error) {
		if len(hooks.Current) == 0 {
			return HookInfo{}, errors.New("no currently running hook")
		}

		hookInfo := hooks.Current[len(hooks.Current)-1]

		return *hookInfo, nil
	}
}
