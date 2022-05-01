package hooks

type Hooks struct {
	Hooks   map[string]Handlers
	Current []*HookInfo
}

type Handler struct {
	Namespace string
	Callback  func(...interface{}) interface{}
	Priority  int
}

type Handlers struct {
	Handlers []Handler
	Runs     int
}

type HookInfo struct {
	Name         string
	CurrentIndex int
}

type Core struct {
	AddAction    func(string, string, func(...interface{}) interface{}, int)
	DoAction     func(string, ...interface{}) interface{}
	AddFilter    func(string, string, func(...interface{}) interface{}, int)
	ApplyFilters func(string, ...interface{}) interface{}
	CurrentAction func() (HookInfo, error)
	CurrentFilter func() (HookInfo, error)
	DidAction func(string) (int)
	DidFilter func(string) (int)
	DoingAction func(string) (bool)
	DoingFilter func(string) (bool)
	HasAction func(string) (bool)
	HasFilter func(string) (bool)
	RemoveAction func(string, string) (int)
	RemoveFilter func(string, string) (int)
	RemoveAllActions func(string, string) (int)
	RemoveAllFilters func(string, string) (int)
	Actions Hooks
	Filters Hooks
}
