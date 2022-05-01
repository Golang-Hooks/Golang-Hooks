package hooks

// Returns a function which, when invoked, will add a hook.
func createAddHook(core *Core, hooks *Hooks) func(string, string, func(...interface{}) interface{}, int) {
	return func(hookName string, namespace string, callback func(...interface{}) interface{}, priority int) {
		handler := Handler{
			Namespace: namespace,
			Callback:  callback,
			Priority:  priority,
		}

		if _, ok := hooks.Hooks[hookName]; ok {
			handlers := hooks.Hooks[hookName].Handlers

			i := len(handlers)
			for ; i > 0; i-- {
				if priority >= handlers[i-1].Priority {
					break
				}
			}

			if i == len(handlers) {
				handlers = append(handlers, handler)
			} else {
				// Otherwise, insert before index.
				handlers = insert(handlers, i, handler)
			}

			if entry, ok := hooks.Hooks[hookName]; ok {
				entry.Handlers = handlers
				hooks.Hooks[hookName] = entry
			}

			if len(hooks.Current) > 0 {
				for _, hookInfo := range hooks.Current {
					if hookInfo.Name == hookName && hookInfo.CurrentIndex >= i {
						hookInfo.CurrentIndex++
					}
				}
			}
		} else {
			hooks.Hooks[hookName] = Handlers{
				Handlers: []Handler{
					handler,
				},
				Runs: 0,
			}
		}

		if hookName != "HookAdded" {
			core.DoAction("HookAdded", hookName, namespace, callback, priority)
		}
	}
}
