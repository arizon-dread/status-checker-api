package businesslayer

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arizon-dread/status-checker-api/models"
)

type Test struct {
	id          int
	name        string
	server      *httptest.Server
	response    *throwaway
	expectedErr error
	body        string
	alertBody   string
	system      *models.Systemstatus
}

type fakeUpdateSystem struct {
	s   *models.Systemstatus
	err error
}

func (f *fakeUpdateSystem) setup(s *models.Systemstatus) error {

	return nil
}

type throwaway struct {
	test string
}

func TestSendStatus(t *testing.T) {

	tcs := []Test{
		{
			id:   1,
			name: "happypath",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})),
			response: &throwaway{
				test: "test",
			},
			expectedErr: nil,
			body:        "message body",
			alertBody:   "oh shizzle",
			system: &models.Systemstatus{
				ID:               1,
				Name:             "google",
				CallUrl:          "https://google.com",
				HttpMethod:       "GET",
				Message:          "Message",
				ResponseMatch:    "<html>",
				AlertBody:        "payload={\"text\": \"Google is down\"}",
				AlertHasBeenSent: false,
				AlertUrl:         "",
			},
		}, {
			id:   2,
			name: "emptyModel",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})),
			response: &throwaway{
				test: "test",
			},
			expectedErr: nil,
			body:        "message body",
			alertBody:   "oh shizzle",
			system:      &models.Systemstatus{},
		},
		{
			id:   3,
			name: "emptyBody",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})),
			response: &throwaway{
				test: "test",
			},
			expectedErr: nil,
			body:        "",
			alertBody:   "oh shizzle",
			system: &models.Systemstatus{
				ID:               1,
				Name:             "google",
				CallUrl:          "https://google.com",
				HttpMethod:       "GET",
				Message:          "Message",
				ResponseMatch:    "<html>",
				AlertBody:        "payload={\"text\": \"Google is down\"}",
				AlertHasBeenSent: false,
				AlertUrl:         "",
			},
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			tc.system.AlertUrl = tc.server.URL
			defer tc.server.Close()
			upd := &fakeUpdateSystem{s: tc.system, err: nil}
			UpdateSystem = upd.setup
			err := sendStatus(tc.system, tc.body)

			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("test FAILED, expected err %v, but was %v\n", tc.expectedErr, err)
			} else {
				t.Logf("test PASSED, err should be nil and was %v\n", err)
			}
		})

	}

}
