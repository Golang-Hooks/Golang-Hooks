package hooks

// Returns a function which, when invoked, will remove a specified hook or all
// hooks by the given name.
func createRemoveHook(core *Core, hooks *Hooks, removeAll bool) func(string, string) int {
	return func(hookName string, namespace string) int {
		handlersRemoved := 0

		if entry, ok := hooks.Hooks[hookName]; ok {

			if removeAll {
				handlersRemoved = len(hooks.Hooks[hookName].Handlers)
				entry.Handlers = []Handler{}
				hooks.Hooks[hookName] = entry
			} else {
				for i, handler := range hooks.Hooks[hookName].Handlers {
					if handler.Namespace == namespace {
						handlersRemoved++
						entry.Handlers = append(hooks.Hooks[hookName].Handlers[:i], hooks.Hooks[hookName].Handlers[i+1:]...)
						hooks.Hooks[hookName] = entry

						// This callback may also be part of a hook that is
						// currently executing.  If the callback we're removing
						// comes after the current callback, there's no problem;
						// otherwise we need to decrease the execution index of any
						// other runs by 1 to account for the removed element.
						for _, hookInfo := range hooks.Current {
							if hookInfo.Name == hookName && hookInfo.CurrentIndex >= i {
								hookInfo.CurrentIndex--
							}
						}
					}
				}
			}
		} else {
			return handlersRemoved
		}

		if hookName != "HookRemoved" {
			core.DoAction("HookRemoved", hookName, namespace)
		}

		return handlersRemoved
	}
}
