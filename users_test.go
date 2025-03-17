package braze_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/dietdoctor/go-braze"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestServer(t *testing.T, pattern string, handler func(http.ResponseWriter, *http.Request)) (*httptest.Server, *braze.Client) {
	mux := http.NewServeMux()
	mux.HandleFunc(pattern, handler)

	srv := httptest.NewServer(mux)
	url, err := url.Parse(srv.URL)
	require.NoError(t, err)

	client, err := braze.NewClient(braze.APIKey("key"), braze.BaseURL(url))
	require.NoError(t, err)

	return srv, client
}

func TestUsersServiceTrack(t *testing.T) {
	srv, client := createTestServer(t, "/users/track", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer key", r.Header.Get("Authorization"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"attributes_processed":1,"message":"success","errors":[{"type":"'external_id' or 'braze_id' or 'user_alias' is required","input_array":"attributes","index":1}]}`))
	})
	defer srv.Close()

	attr := &braze.UserAttributes{
		ExternalID:     braze.String("users/3ceeeb51-dae1-49d5-bbba-23affdd82dde"),
		FirstName:      braze.String("Vaidas"),
		LastName:       braze.String("Test"),
		Email:          braze.String("vaidas@dietdoctor.com"),
		EmailSubscribe: &braze.AttributeSubscribeSubscribed,
		Gender:         &braze.AttributeGenderMale,
		Timezone:       braze.String("GMT"),
	}

	attr.AddAttributes(
		braze.BoolAttribute("is_user", true),
		braze.DateAttribute("seen_time", time.Now()),
		braze.ModifyStringSliceAttribute("tags", map[braze.SliceAttributeAction][]string{
			braze.SliceAttributeActionAdd:    {"user"},
			braze.SliceAttributeActionRemove: {"foo"},
		}),
	)

	resp, err := client.Users().Track(context.Background(), &braze.UsersTrackRequest{
		Attributes: []*braze.UserAttributes{
			attr,
			// This is expected to return a minor error.
			{
				Gender: &braze.AttributeGenderMale,
			},
		},
	})

	assert.NoError(t, err)
	assert.Len(t, resp.Errors, 1)
	assert.Equal(t, "success", resp.Message)
}

func TestUsersServiceTrackInternalServerError(t *testing.T) {
	srv, client := createTestServer(t, "/users/track", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer key", r.Header.Get("Authorization"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusInternalServerError)
	})
	defer srv.Close()

	resp, err := client.Users().Track(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestUsersServiceTrackCustomAttributes(t *testing.T) {
	srv, client := createTestServer(t, "/users/track", func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		assert.Equal(t, []byte(`{"attributes":[{"external_id":"123","testing":true}]}`), b)

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{}`))
	})
	defer srv.Close()

	attr := &braze.UserAttributes{
		ExternalID: braze.String("123"),
	}
	attr.AddAttributes(braze.BoolAttribute("testing", true))

	resp, err := client.Users().Track(context.Background(), &braze.UsersTrackRequest{
		Attributes: []*braze.UserAttributes{attr},
	})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestUsersServiceTrackGenericCustomAttributes(t *testing.T) {
	srv, client := createTestServer(t, "/users/track", func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		assert.Equal(t, []byte(`{"attributes":[{"external_id":"123","testing":true}]}`), b)

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{}`))
	})
	defer srv.Close()

	attr := &braze.UserAttributes{
		ExternalID: braze.String("123"),
	}
	attr.AddAttributes(braze.Attribute("testing", true))

	resp, err := client.Users().Track(context.Background(), &braze.UsersTrackRequest{
		Attributes: []*braze.UserAttributes{attr},
	})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestUsersServiceExportIdsUserExists(t *testing.T) {
	srv, client := createTestServer(t, "/users/export/ids", func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		assert.Equal(t, []byte(`{"external_ids":["123"]}`), b)

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"users":[{"external_id":"123"}]}`))
	})
	defer srv.Close()

	resp, err := client.Users().ExportIds(context.Background(), &braze.UsersExportIdsRequest{
		ExternalIDs: []string{"123"},
	})
	assert.NoError(t, err)

	expected := &braze.UserExportResponse{
		Users: []braze.ExportedUser{{ExternalID: "123"}},
	}
	assert.Equal(t, expected, resp)
}

func TestUsersServiceExportIdsUserDoesntExist(t *testing.T) {
	srv, client := createTestServer(t, "/users/export/ids", func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		assert.Equal(t, []byte(`{"external_ids":["123"]}`), b)

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"users":[],"invalid_user_ids":["123"]}`))
	})
	defer srv.Close()

	resp, err := client.Users().ExportIds(context.Background(), &braze.UsersExportIdsRequest{
		ExternalIDs: []string{"123"},
	})
	assert.NoError(t, err)

	expected := &braze.UserExportResponse{
		Users:          []braze.ExportedUser{},
		InvalidUserIds: []string{"123"},
	}
	assert.Equal(t, expected, resp)
}
