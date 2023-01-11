package braze_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/dietdoctor/go-braze"
	"github.com/stretchr/testify/assert"
)

func TestPreferencesServerCreateURL(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/preference_center/v1/foo/url/bar", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer key", r.Header.Get("Authorization"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"preference_center_url":"https://foo"}`))
	})

	srv := httptest.NewServer(mux)
	defer srv.Close()

	url, _ := url.Parse(srv.URL)
	client, err := braze.NewClient(braze.APIKey("key"), braze.BaseURL(url))
	assert.NoError(t, err)

	resp, err := client.PreferenceCenter.CreateURL(context.Background(), &braze.PreferenceCenterCreateURLRequest{
		PreferenceCenterID: "foo",
		UserID:             "bar",
	})

	assert.NoError(t, err)
	assert.Equal(t, "https://foo", resp.URL)
}

func TestPreferencesServerCreateURLInternalServerError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/preference_center/v1/foo/url/bar", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer key", r.Header.Get("Authorization"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusInternalServerError)
	})

	srv := httptest.NewServer(mux)
	defer srv.Close()

	url, _ := url.Parse(srv.URL)
	client, err := braze.NewClient(braze.APIKey("key"), braze.BaseURL(url))
	assert.NoError(t, err)

	resp, err := client.PreferenceCenter.CreateURL(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, resp)
}
