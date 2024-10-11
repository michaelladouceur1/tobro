// Package http_server provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.3.0 DO NOT EDIT.
package http_server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// Defines values for SetupPinRequestMode.
const (
	Input  SetupPinRequestMode = "input"
	Output SetupPinRequestMode = "output"
)

// AnalogWritePinRequest defines model for AnalogWritePinRequest.
type AnalogWritePinRequest struct {
	PinNumber int `json:"pinNumber"`
	Value     int `json:"value"`
}

// AnalogWritePinResponse defines model for AnalogWritePinResponse.
type AnalogWritePinResponse struct {
	PinNumber *int `json:"pinNumber,omitempty"`
	Value     *int `json:"value,omitempty"`
}

// BoardsResponse defines model for BoardsResponse.
type BoardsResponse struct {
	Boards []string `json:"boards"`
}

// CircuitResponse defines model for CircuitResponse.
type CircuitResponse struct {
	Board string        `json:"board"`
	Id    int           `json:"id"`
	Name  string        `json:"name"`
	Pins  []PinResponse `json:"pins"`
}

// ConnectRequest defines model for ConnectRequest.
type ConnectRequest struct {
	Port string `json:"port"`
}

// ConnectResponse defines model for ConnectResponse.
type ConnectResponse struct {
	Port      *string `json:"port,omitempty"`
	Timestamp *int    `json:"timestamp,omitempty"`
}

// CreateCircuitRequest defines model for CreateCircuitRequest.
type CreateCircuitRequest struct {
	Board string `json:"board"`
	Name  string `json:"name"`
}

// DigitalWritePinRequest defines model for DigitalWritePinRequest.
type DigitalWritePinRequest struct {
	PinNumber int `json:"pinNumber"`
	Value     int `json:"value"`
}

// DigitalWritePinResponse defines model for DigitalWritePinResponse.
type DigitalWritePinResponse struct {
	PinNumber int `json:"pinNumber"`
	Value     int `json:"value"`
}

// ErrorResponse defines model for ErrorResponse.
type ErrorResponse struct {
	Message *string `json:"message,omitempty"`
}

// PinResponse defines model for PinResponse.
type PinResponse struct {
	AnalogRead   bool   `json:"analogRead"`
	AnalogWrite  bool   `json:"analogWrite"`
	DigitalRead  bool   `json:"digitalRead"`
	DigitalWrite bool   `json:"digitalWrite"`
	Max          int    `json:"max"`
	Min          int    `json:"min"`
	Mode         int    `json:"mode"`
	PinNumber    int    `json:"pinNumber"`
	Type         string `json:"type"`
}

// SaveCircuitRequest defines model for SaveCircuitRequest.
type SaveCircuitRequest struct {
	Id int `json:"id"`
}

// SetupPinRequest defines model for SetupPinRequest.
type SetupPinRequest struct {
	Mode      SetupPinRequestMode `json:"mode"`
	PinNumber int                 `json:"pinNumber"`
}

// SetupPinRequestMode defines model for SetupPinRequest.Mode.
type SetupPinRequestMode string

// SetupPinResponse defines model for SetupPinResponse.
type SetupPinResponse struct {
	Mode      string `json:"mode"`
	PinNumber int    `json:"pinNumber"`
}

// SketchAPI defines model for SketchAPI.
type SketchAPI struct {
	Id    int             `json:"id"`
	Name  string          `json:"name"`
	Steps []SketchStepAPI `json:"steps"`
}

// SketchStepAPI defines model for SketchStepAPI.
type SketchStepAPI struct {
	Action    string `json:"action"`
	End       int    `json:"end"`
	Id        int    `json:"id"`
	PinNumber int    `json:"pinNumber"`
	Start     int    `json:"start"`
}

// PostAnalogWritePinJSONRequestBody defines body for PostAnalogWritePin for application/json ContentType.
type PostAnalogWritePinJSONRequestBody = AnalogWritePinRequest

// PostCircuitJSONRequestBody defines body for PostCircuit for application/json ContentType.
type PostCircuitJSONRequestBody = CreateCircuitRequest

// PostConnectJSONRequestBody defines body for PostConnect for application/json ContentType.
type PostConnectJSONRequestBody = ConnectRequest

// PostDigitalWritePinJSONRequestBody defines body for PostDigitalWritePin for application/json ContentType.
type PostDigitalWritePinJSONRequestBody = DigitalWritePinRequest

// PostSaveCircuitJSONRequestBody defines body for PostSaveCircuit for application/json ContentType.
type PostSaveCircuitJSONRequestBody = SaveCircuitRequest

// PostSetupPinJSONRequestBody defines body for PostSetupPin for application/json ContentType.
type PostSetupPinJSONRequestBody = SetupPinRequest

// PostSketchJSONRequestBody defines body for PostSketch for application/json ContentType.
type PostSketchJSONRequestBody = SketchAPI

// PostSketchStepJSONRequestBody defines body for PostSketchStep for application/json ContentType.
type PostSketchStepJSONRequestBody = SketchStepAPI

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (POST /analog_write_pin)
	PostAnalogWritePin(w http.ResponseWriter, r *http.Request)

	// (GET /boards)
	GetBoards(w http.ResponseWriter, r *http.Request)

	// (GET /circuit)
	GetCircuit(w http.ResponseWriter, r *http.Request)

	// (POST /circuit)
	PostCircuit(w http.ResponseWriter, r *http.Request)

	// (POST /connect)
	PostConnect(w http.ResponseWriter, r *http.Request)

	// (POST /digital_write_pin)
	PostDigitalWritePin(w http.ResponseWriter, r *http.Request)

	// (POST /save_circuit)
	PostSaveCircuit(w http.ResponseWriter, r *http.Request)

	// (POST /setup_pin)
	PostSetupPin(w http.ResponseWriter, r *http.Request)

	// (GET /sketch)
	GetSketch(w http.ResponseWriter, r *http.Request)

	// (POST /sketch)
	PostSketch(w http.ResponseWriter, r *http.Request)

	// (POST /sketch/step)
	PostSketchStep(w http.ResponseWriter, r *http.Request)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// PostAnalogWritePin operation middleware
func (siw *ServerInterfaceWrapper) PostAnalogWritePin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostAnalogWritePin(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetBoards operation middleware
func (siw *ServerInterfaceWrapper) GetBoards(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetBoards(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetCircuit operation middleware
func (siw *ServerInterfaceWrapper) GetCircuit(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetCircuit(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// PostCircuit operation middleware
func (siw *ServerInterfaceWrapper) PostCircuit(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostCircuit(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// PostConnect operation middleware
func (siw *ServerInterfaceWrapper) PostConnect(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostConnect(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// PostDigitalWritePin operation middleware
func (siw *ServerInterfaceWrapper) PostDigitalWritePin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostDigitalWritePin(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// PostSaveCircuit operation middleware
func (siw *ServerInterfaceWrapper) PostSaveCircuit(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostSaveCircuit(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// PostSetupPin operation middleware
func (siw *ServerInterfaceWrapper) PostSetupPin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostSetupPin(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetSketch operation middleware
func (siw *ServerInterfaceWrapper) GetSketch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetSketch(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// PostSketch operation middleware
func (siw *ServerInterfaceWrapper) PostSketch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostSketch(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// PostSketchStep operation middleware
func (siw *ServerInterfaceWrapper) PostSketchStep(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostSketchStep(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, GorillaServerOptions{})
}

type GorillaServerOptions struct {
	BaseURL          string
	BaseRouter       *mux.Router
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r *mux.Router) http.Handler {
	return HandlerWithOptions(si, GorillaServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r *mux.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, GorillaServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options GorillaServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = mux.NewRouter()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	r.HandleFunc(options.BaseURL+"/analog_write_pin", wrapper.PostAnalogWritePin).Methods("POST")

	r.HandleFunc(options.BaseURL+"/boards", wrapper.GetBoards).Methods("GET")

	r.HandleFunc(options.BaseURL+"/circuit", wrapper.GetCircuit).Methods("GET")

	r.HandleFunc(options.BaseURL+"/circuit", wrapper.PostCircuit).Methods("POST")

	r.HandleFunc(options.BaseURL+"/connect", wrapper.PostConnect).Methods("POST")

	r.HandleFunc(options.BaseURL+"/digital_write_pin", wrapper.PostDigitalWritePin).Methods("POST")

	r.HandleFunc(options.BaseURL+"/save_circuit", wrapper.PostSaveCircuit).Methods("POST")

	r.HandleFunc(options.BaseURL+"/setup_pin", wrapper.PostSetupPin).Methods("POST")

	r.HandleFunc(options.BaseURL+"/sketch", wrapper.GetSketch).Methods("GET")

	r.HandleFunc(options.BaseURL+"/sketch", wrapper.PostSketch).Methods("POST")

	r.HandleFunc(options.BaseURL+"/sketch/step", wrapper.PostSketchStep).Methods("POST")

	return r
}
