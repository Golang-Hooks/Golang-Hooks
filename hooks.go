package hooks

func CreateHooks() Core {
	actions := Hooks{}
	actions.Hooks = make(map[string]Handlers)
	filters := Hooks{}
	filters.Hooks = make(map[string]Handlers)

	rv := Core{}

	rv.AddAction = createAddHook(&rv, &actions)
	rv.DoAction = createRunHook(&rv, &actions, false)
	rv.AddFilter = createAddHook(&rv, &filters)
	rv.ApplyFilters = createRunHook(&rv, &filters, true)
	rv.CurrentAction = createCurrentHook(&rv, &actions)
	rv.CurrentFilter = createCurrentHook(&rv, &filters)
	rv.DidAction = createDidHook(&rv, &actions)
	rv.DidFilter = createDidHook(&rv, &filters)
	rv.DoingAction = createDoingHook(&rv, &actions)
	rv.DoingFilter = createDoingHook(&rv, &filters)
	rv.HasAction = createHasHook(&rv, &actions)
	rv.HasFilter = createHasHook(&rv, &filters)
	rv.RemoveAction = createRemoveHook(&rv, &actions, false)
	rv.RemoveFilter = createRemoveHook(&rv, &filters, false)
	rv.RemoveAllActions = createRemoveHook(&rv, &actions, true)
	rv.RemoveAllFilters = createRemoveHook(&rv, &filters, true)
	rv.Actions = actions
	rv.Filters = filters

	return rv
}
