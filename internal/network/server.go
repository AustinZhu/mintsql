package network

const (
	HOST     = "localhost"
	PROTOCOL = "tcp"
)

type Server struct {
	Host string
	Port uint16
	//ConnPool
}

func New(host string, port uint16, poolSize int) {
	panic("Not Implemented")
}

func (s *Server) Run() {
	panic("Not Implemented")
}
