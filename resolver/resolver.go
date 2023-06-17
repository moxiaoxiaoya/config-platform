package resolver

// TODO: why have to use builder
type Builder interface {
	Build(target string, cc ClientConn) (Discovery, error)
}

type Discovery interface {
	Discover()
	Close()
}

type Register interface {
	Register()
}

type ClientConn interface {
	UpdateState(addr string) error
}
