package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserClient(t *testing.T) {
	tests := []struct {
		name     string
		user     *User
		response string
		status   int
		wantErr  bool
	}{
		{
			name:     "success",
			user:     &User{ID: 1, Name: "江湖十年"},
			response: `{"id":1,"name":"江湖十年"}`,
			status:   200,
			wantErr:  false,
		},
		{
			name:     "notfound",
			user:     &User{ID: 2},
			response: "",
			status:   404,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(tt.status)
					if tt.response != "" {
						w.Write([]byte(tt.response))
					}
				}),
			)
			defer server.Close()

			client := NewUserClient(server.URL)
			user, err := client.GetUser(tt.user.ID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.user, user)
			}
		})
	}
}
