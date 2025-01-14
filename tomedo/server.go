package tomedo

import "net/url"

const DemoServerURL = "http://allgemeinmedizin.demo.tomedo.org:8080/tomedo_live/"

type Server struct {
	Scheme string
	Addr   string
	Port   string
	Path   string
}

func DefaultServer() *Server {
	u, err := url.Parse(DemoServerURL)
	if err != nil {
		panic(err)
	}
	return &Server{
		Scheme: u.Scheme,
		Addr:   u.Hostname(),
		Port:   u.Port(),
		Path:   u.Path,
	}
}

func (s *Server) URL() *url.URL {
	return &url.URL{
		Scheme: s.Scheme,
		Host:   s.Addr + ":" + s.Port,
		Path:   s.Path,
	}
}

func (s *Server) String() string {
	return s.URL().String()
}
