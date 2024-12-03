package route

type route struct {
	prefix string
}

func New(prefix string) *route {
	return &route{
		prefix: prefix,
	}
}

func (api *route) Path(subpath string) string {
	return api.prefix + subpath

}
