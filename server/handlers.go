package server

import (
	"encoding/json"
	"fmt"
	"github.com/geeckmc/eltrade-cc300-driver/cmd"
	"github.com/imdario/mergo"
	"github.com/xeipuuv/gojsonschema"
	"io/ioutil"
	"net/http"
)

func (s *server) createBill() {
	if s.r.Method != http.MethodPost {
		s.w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	input, err := ioutil.ReadAll(s.r.Body)
	if err != nil {
		s.w.WriteHeader(http.StatusBadRequest)
		return
	}
	inputStruct, err := s.billSchema.Validate(gojsonschema.NewStringLoader(string(input)))
	if handleError(err, s) {
		return
	}
	if !inputStruct.Valid() {
		s.w.WriteHeader(http.StatusBadRequest)
		var errors []string
		for _, error := range inputStruct.Errors() {
			errors = append(errors, error.Description())
		}
		b, _ := json.Marshal(map[string]string{"errors": fmt.Sprint(errors)})
		s.w.Write(b)
		return
	}

	res, err := cmd.CreateBill(s.dev, input)
	s.w.WriteHeader(http.StatusOK)
	resJson, _ := json.Marshal(map[string]string{"qr_code": res})
	s.w.Write(resJson)

}

func (s *server) Info() {
	if s.r.Method != http.MethodGet {
		s.w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	s.w.WriteHeader(http.StatusOK)

	r, err := cmd.GetDeviceState(s.dev)
	if err != nil {
		handleError(err, s)
		return
	}
	deviceInfo := cmd.DeviceInfo{}
	mergo.Merge(&deviceInfo, r, mergo.WithOverride)
	r, err = cmd.GetTaxPayerInfo(s.dev)
	handleError(err, s)
	if err != nil {
		handleError(err, s)
		return
	}
	mergo.Merge(&deviceInfo, r, mergo.WithOverride)
	r, err = cmd.GetTaxServerState(s.dev)
	if err != nil {
		handleError(err, s)
		return
	}
	mergo.Merge(&deviceInfo, r, mergo.WithOverride)
	deviceInfoBytes, _ := json.Marshal(deviceInfo)
	fmt.Printf("%v", deviceInfo)
	s.w.Write(deviceInfoBytes)
}

func handleError(err error, s *server) bool {
	if err != nil {
		s.w.WriteHeader(http.StatusInternalServerError)
		s.w.Write([]byte(fmt.Sprintf("{error_message: %s}", err.Error())))
		return true
	}
	return false
}
