package frinkconv_server

import (
	"fmt"
	"github.com/initialed85/frinkconv-api/internal/helpers"
	"github.com/initialed85/frinkconv-api/pkg/frinkconv_repl"
	"github.com/initialed85/frinkconv-api/pkg/pool"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

type Server struct {
	serveMux *http.ServeMux
	server   *http.Server
	repls    []*frinkconv_repl.REPL
	pool     *pool.Pool
}

func New(port int, processes int) (*Server, error) {
	log.Printf("!!! port=%#+v, processes=%#+v", port, processes)

	s := Server{
		serveMux: http.NewServeMux(),
		server:   nil,
		repls:    make([]*frinkconv_repl.REPL, 0),
		pool:     pool.New(processes),
	}

	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	replErrors := make([]error, 0)

	for i := 0; i < processes; i++ {
		wg.Add(1)

		go func() {
			repl, err := frinkconv_repl.New()
			if err != nil {
				mu.Lock()
				replErrors = append(replErrors, err)
				mu.Unlock()
				return
			}

			err = s.pool.PutTimeout(repl, time.Second)
			if err != nil {
				mu.Lock()
				replErrors = append(replErrors, err)
				mu.Unlock()
				return
			}

			mu.Lock()
			s.repls = append(s.repls, repl)
			mu.Unlock()

			wg.Done()
		}()
	}

	wg.Wait()

	if len(replErrors) > 0 {
		return nil, replErrors[0]
	}

	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%v", port),
		Handler: s.serveMux,
	}

	s.serveMux.HandleFunc("/convert/", s.convert)
	s.serveMux.HandleFunc("/batch_convert/", s.batchConvert)

	errors := helpers.GetErrorChannel()

	go func() {
		err := s.server.ListenAndServe()
		if err != nil {
			errors <- err
		}
	}()

	err := helpers.WaitForError(errors, time.Millisecond*100)
	if err != nil {
		s.Close()
		return nil, err
	}

	return &s, nil
}

func (s *Server) isMethodValid(responseWriter http.ResponseWriter, request *http.Request) bool {
	if request.Method != http.MethodPost {
		responseWriter.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = responseWriter.Write(getConvertErrorResponseJSON(
			fmt.Errorf("%#+v not supported", request.Method),
		))
		return false
	}

	return true
}

func (s *Server) getBody(responseWriter http.ResponseWriter, request *http.Request) ([]byte, error) {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		_, _ = responseWriter.Write(getConvertErrorResponseJSON(err))
		return []byte{}, err
	}

	return body, nil
}

func (s *Server) getREPL() (*frinkconv_repl.REPL, error) {
	possibleREPL, err := s.pool.GetTimeout(time.Second * 5)
	if err != nil {
		return nil, err
	}

	repl, ok := possibleREPL.(*frinkconv_repl.REPL)
	if !ok {
		return nil, fmt.Errorf("failed to cast %#+v to REPL", repl)
	}

	return repl, nil
}

func (s *Server) convert(responseWriter http.ResponseWriter, request *http.Request) {
	if !s.isMethodValid(responseWriter, request) {
		return
	}

	body, err := s.getBody(responseWriter, request)
	if err != nil {
		return
	}

	defer func() {
		_ = request.Body.Close()
	}()

	convertRequest, err := getConvertRequestFromJSON(body)
	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		_, _ = responseWriter.Write(getConvertErrorResponseJSON(err))
		return
	}

	repl, err := s.getREPL()
	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		_, _ = responseWriter.Write(getConvertErrorResponseJSON(err))
	}

	defer s.pool.Put(repl)

	destinationValue, err := repl.Convert(convertRequest.SourceValue, convertRequest.SourceUnits, convertRequest.DestinationUnits)
	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		_, _ = responseWriter.Write(getConvertErrorResponseJSON(err))
		return
	}

	responseWriter.WriteHeader(http.StatusOK)
	_, _ = responseWriter.Write(getConvertSuccessResponseJSON(destinationValue))
}

func (s *Server) batchConvert(responseWriter http.ResponseWriter, request *http.Request) {
	if !s.isMethodValid(responseWriter, request) {
		return
	}

	body, err := s.getBody(responseWriter, request)
	if err != nil {
		return
	}

	defer func() {
		_ = request.Body.Close()
	}()

	convertRequests, err := getConvertRequestsFromJSON(body)
	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		_, _ = responseWriter.Write(getConvertErrorResponseJSON(err))
		return
	}

	repl, err := s.getREPL()
	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		_, _ = responseWriter.Write(getConvertErrorResponseJSON(err))
	}

	defer s.pool.Put(repl)

	convertResponses := make([]interface{}, 0)
	anyErrors := false

	for _, convertRequest := range convertRequests {
		destinationValue, err := repl.Convert(convertRequest.SourceValue, convertRequest.SourceUnits, convertRequest.DestinationUnits)
		if err != nil {
			anyErrors = true
			convertResponses = append(
				convertResponses,
				ConvertErrorResponse{
					Error: err.Error(),
				},
			)
			continue
		}

		convertResponses = append(
			convertResponses,
			ConvertSuccessResponse{
				DestinationValue: destinationValue,
			},
		)
	}

	if anyErrors {
		responseWriter.WriteHeader(http.StatusBadRequest)
	} else {
		responseWriter.WriteHeader(http.StatusOK)
	}

	_, _ = responseWriter.Write(getConvertResponsesJSON(convertResponses))
}

func (s *Server) Close() {
	_ = s.server.Close()

	for _, repl := range s.repls {
		repl.Close()
	}
}
