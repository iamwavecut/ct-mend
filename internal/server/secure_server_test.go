package server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/iamwavecut/ct-mend/internal/storage"
	"github.com/iamwavecut/ct-mend/tools"
)

type TLSTestSuite struct {
	suite.Suite
}

var vars map[string]string

func (s *TLSTestSuite) SetupTest() {
	vars = map[string]string{}
}

func (s *TLSTestSuite) TestClientHandlers() {
	db := storage.NewMockAdapter(s.T())
	handler := &ClientsHandler{db: db}

	for _, tc := range []struct {
		methodName string
		method     string
		path       string
		resultCode int
		handle     http.HandlerFunc
	}{
		{
			"SelectClients",
			"GET",
			"/",
			200,
			func() func(writer http.ResponseWriter, request *http.Request) {
				db.On("SelectClients").Return([]*storage.Client{{}, {}}, nil).Once()
				return handler.Select
			}(),
		},
		{
			"GetClient",
			"GET",
			"/1",
			200,
			func() func(writer http.ResponseWriter, request *http.Request) {
				db.On("GetClient", mock.AnythingOfType("int")).Return(&storage.Client{}, nil).Once()
				vars = map[string]string{"id": "1"}
				return handler.Get
			}(),
		},
		{
			"PostClient",
			"POST",
			"/",
			201,
			func() func(writer http.ResponseWriter, request *http.Request) {
				db.On("UpsertClient", mock.IsType(&storage.Client{})).Return(&storage.Client{ID: tools.IntPtr(1)}, nil).Once()
				return handler.Post
			}(),
		},
		{
			"PutClient",
			"PUT",
			"/1",
			201,
			func() func(writer http.ResponseWriter, request *http.Request) {
				db.On("UpsertClient", mock.IsType(&storage.Client{})).Return(&storage.Client{ID: tools.IntPtr(1)}, nil).Once()
				vars = map[string]string{"id": "1"}
				return handler.Put
			}(),
		},
		{
			"DeleteClient",
			"DELETE",
			"/1",
			204,
			func() func(writer http.ResponseWriter, request *http.Request) {
				db.On("DeleteClient", mock.AnythingOfType("int")).Return(nil).Once()
				vars = map[string]string{"id": "1"}
				return handler.Delete
			}(),
		},
	} {
		s.Run(tc.methodName, func() {
			resp := httptest.NewRecorder()
			req, err := http.NewRequest(tc.method, "https://about.blank/clients"+tc.path, bytes.NewBuffer([]byte(
				`{"id":1,"name":"123123","settings":{"code_scan_interval":123123}}`,
			)))
			s.Assert().NoError(err)
			if len(vars) > 0 {
				params := url.Values{}
				for key, val := range vars {
					params.Set(key, val)
				}
				req.URL.RawQuery = params.Encode()
			}

			tc.handle(resp, req)
			s.Assert().Equal(tc.resultCode, resp.Code)
		})
	}
}

func (s *TLSTestSuite) TestProjectHandlers() {
	db := storage.NewMockAdapter(s.T())
	handler := &ProjectsHanlder{db: db}

	for _, tc := range []struct {
		methodName string
		method     string
		path       string
		resultCode int
		handle     http.HandlerFunc
	}{
		{
			"SelectProjects",
			"GET",
			"/",
			200,
			func() func(writer http.ResponseWriter, request *http.Request) {
				db.On("SelectProjects").Return([]*storage.Project{{}, {}}, nil).Once()
				return handler.Select
			}(),
		},
		{
			"GetProject",
			"GET",
			"/1",
			200,
			func() func(writer http.ResponseWriter, request *http.Request) {
				db.On("GetProject", mock.AnythingOfType("int")).Return(&storage.Project{}, nil).Once()
				vars = map[string]string{"id": "1"}
				return handler.Get
			}(),
		},
		{
			"PostProject",
			"POST",
			"/",
			201,
			func() func(writer http.ResponseWriter, request *http.Request) {
				db.On("UpsertProject", mock.IsType(&storage.Project{})).Return(&storage.Project{ID: tools.IntPtr(1)}, nil).Once()
				return handler.Post
			}(),
		},
		{
			"PutProject",
			"PUT",
			"/1",
			201,
			func() func(writer http.ResponseWriter, request *http.Request) {
				db.On("UpsertProject", mock.IsType(&storage.Project{})).Return(&storage.Project{ID: tools.IntPtr(1)}, nil).Once()
				vars = map[string]string{"id": "1"}
				return handler.Put
			}(),
		},
		{
			"DeleteProject",
			"DELETE",
			"/1",
			204,
			func() func(writer http.ResponseWriter, request *http.Request) {
				db.On("DeleteProject", mock.AnythingOfType("int")).Return(nil).Once()
				vars = map[string]string{"id": "1"}
				return handler.Delete
			}(),
		},
	} {
		s.Run(tc.methodName, func() {
			resp := httptest.NewRecorder()
			req, err := http.NewRequest(tc.method, "https://about.blank/projects"+tc.path, bytes.NewBuffer([]byte(
				`{"id":1,"name":"123123","settings":{"code_scan_interval":123123}}`,
			)))
			s.Assert().NoError(err)
			if len(vars) > 0 {
				params := url.Values{}
				for key, val := range vars {
					params.Set(key, val)
				}
				req.URL.RawQuery = params.Encode()
			}

			tc.handle(resp, req)
			s.Assert().Equal(tc.resultCode, resp.Code)
		})
	}
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(TLSTestSuite))
}
