package registry

const (
	MODE_HUB   = "hub"
	MODE_SPOKE = "spoke"
)

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

func isHub(mode string) bool {
	return mode == MODE_HUB
}

func isSpoke(mode string) bool {
	return mode == MODE_SPOKE
}
