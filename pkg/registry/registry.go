package registry

//type FileServer
type Registry struct {
	Mode string
}

//Constructor NewFileServer
func NewRegistry(mode string) *Registry {
	return &Registry{
		Mode: mode,
	}
}
