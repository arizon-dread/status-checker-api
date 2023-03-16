package businesslayer

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/arizon-dread/status-checker-api/models"
)

// TODO: Set up mock for:
// datalayer.GetSystemStatus(id)
// checkCert(&system)
// getTlsConfigWithClientCert(system)
// getContentType(system.CallBody)
// handleResponse(&system, resp, err)
// sendAlert(&system, message)
// datalayer.SaveSystemStatus(&system)
// type fakeUpdateSystem struct {
// 	s   *models.Systemstatus
// 	err error
// }

// func (f *fakeUpdateSystem) setup(s *models.Systemstatus) error {

// 	return nil
// }

type fakeDLGetSystemStatus struct {
	id int
	s  *models.Systemstatus
}

func (fgss *fakeDLGetSystemStatus) setup(id int) (models.Systemstatus, error) {
	return models.Systemstatus{
		ID:                 id,
		Name:               "google",
		CallStatus:         "OK",
		CallUrl:            "https://dummyjson.com/products/1",
		CallBody:           "",
		HttpMethod:         "GET",
		CertStatus:         "OK",
		CertExpirationDays: 20,
		Message:            "Message",
		ResponseMatch:      "<html>",
		AlertBody:          "google is down",
		AlertUrl:           "",
		AlertEmail:         "hi@hello.com",
		AlertHasBeenSent:   false,
		Status:             "OK",
		LastOKTime:         time.Now(),
	}, nil
}

func TestGetSystemStatus(t *testing.T) {
	type args struct {
		id int
	}
	tests := []struct {
		name    string
		server  *httptest.Server
		args    args
		want    models.Systemstatus
		wantErr bool
	}{
		{
			name: "happy path",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})),
			args: args{1},
			want: models.Systemstatus{
				ID:                 1,
				Name:               "google",
				CallStatus:         "OK",
				CallUrl:            "https://dummyjson.com/products/1",
				CallBody:           "",
				HttpMethod:         "GET",
				CertStatus:         "OK",
				CertExpirationDays: 20,
				Message:            "Message",
				ResponseMatch:      "<html>",
				AlertBody:          "google is down",
				AlertUrl:           "",
				AlertEmail:         "hi@hello.com",
				AlertHasBeenSent:   false,
				Status:             "OK",
				LastOKTime:         time.Now(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt.want.AlertUrl = tt.server.URL
		t.Run(tt.name, func(t *testing.T) {
			dlgss := &fakeDLGetSystemStatus{id: 1}
			dlGetSystemStatus = dlgss.setup
			got, err := GetSystemStatus(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSystemStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSystemStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}
