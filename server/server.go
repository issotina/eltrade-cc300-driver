package server

import (
	"fmt"
	eltrade "github.com/geeckmc/eltrade-cc300-driver/lib"
	"github.com/xeipuuv/gojsonschema"
	"net/http"
)

const (
	PORT      = ":38917"
	home_page = "https://41devs.com"
)

type Status string

func (s Status) str() []byte {
	return []byte(fmt.Sprintf("{\"status\": \"%s\"}", s))
}

const (
	DeviceNotConnected Status = "DeviceNotConnected"
	Ready              Status = "Ready"
)

type server struct {
	w          http.ResponseWriter
	r          *http.Request
	dev        *eltrade.Device
	billSchema *gojsonschema.Schema
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.r = r
	s.w = w
	s.w.Header().Set("Content-Type", "application/json")
	dev, err := eltrade.Open()
	if err != nil {
		println("error is not nil", err.Error())
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write(DeviceNotConnected.str())
		return
	}
	s.dev = dev
	defer dev.Close()
	switch r.URL.Path {
	case "/bill":
		s.createBill()
		return
	case "/info":
		s.Info()
		return
	case "/check":
		s.Check()
		return
	default:
		http.Redirect(w, r, home_page, http.StatusTemporaryRedirect)
		return
	}
}

func (s *server) Check() {
	s.w.WriteHeader(http.StatusOK)
	var code Status
	if s.dev == nil || !s.dev.IsOpen {
		s.w.WriteHeader(http.StatusServiceUnavailable)
		code = DeviceNotConnected

	} else {
		code = Ready
	}

	s.w.Write(code.str())

}

func Serve() *http.Server {
	schemaString := gojsonschema.NewStringLoader(JsonSchema())
	sl := gojsonschema.NewSchemaLoader()
	server := &server{}
	server.billSchema, _ = sl.Compile(schemaString)
	httpSrv := http.Server{Addr: PORT, Handler: server}
	err := httpSrv.ListenAndServe()
	if err != nil {
		eltrade.Logger.Errorf("%s", err)
	}
	return &httpSrv

}
