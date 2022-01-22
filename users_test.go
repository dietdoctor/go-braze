package braze_test

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/tkuchiki/faketime"

	"github.com/dietdoctor/go-braze"
	"github.com/ferdypruis/iso4217"
	"github.com/stretchr/testify/assert"
)

func TestUsersServiceTrack(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/users/track", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer key", r.Header.Get("Authorization"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"attributes_processed":1,"message":"success","errors":[{"type":"'external_id' or 'braze_id' or 'user_alias' is required","input_array":"attributes","index":1}]}`))
	})

	srv := httptest.NewServer(mux)
	defer srv.Close()

	url, _ := url.Parse(srv.URL)
	client, err := braze.NewClient(braze.APIKey("key"), braze.BaseURL(url))
	assert.NoError(t, err)

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

	resp, err := client.Users.Track(context.Background(), &braze.UsersTrackRequest{
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
	mux := http.NewServeMux()
	mux.HandleFunc("/users/track", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer key", r.Header.Get("Authorization"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusInternalServerError)
	})

	srv := httptest.NewServer(mux)
	defer srv.Close()

	url, _ := url.Parse(srv.URL)
	client, err := braze.NewClient(braze.APIKey("key"), braze.BaseURL(url))
	assert.NoError(t, err)

	resp, err := client.Users.Track(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestUsersServiceTrackCustomAttributes(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/users/track", func(w http.ResponseWriter, r *http.Request) {
		b, _ := ioutil.ReadAll(r.Body)
		assert.Equal(t, []byte(`{"attributes":[{"external_id":"123","testing":true}]}`), b)

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{}`))
	})

	srv := httptest.NewServer(mux)
	defer srv.Close()

	url, _ := url.Parse(srv.URL)
	client, err := braze.NewClient(braze.APIKey("key"), braze.BaseURL(url))
	assert.NoError(t, err)

	attr := &braze.UserAttributes{
		ExternalID: braze.String("123"),
	}
	attr.AddAttributes(braze.BoolAttribute("testing", true))

	resp, err := client.Users.Track(context.Background(), &braze.UsersTrackRequest{
		Attributes: []*braze.UserAttributes{attr},
	})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestUsersServiceTrackPurchase(t *testing.T) {
	f := faketime.NewFaketime(2022, time.January, 22, 0, 0, 0, 0, time.UTC)
	defer f.Undo()
	f.Do()

	mux := http.NewServeMux()
	mux.HandleFunc("/users/track", func(w http.ResponseWriter, r *http.Request) {
		b, _ := ioutil.ReadAll(r.Body)
		assert.Equal(t, []byte(`{"purchases":[{"external_id":"123","app_id":"diet_doctor","product_id":"subscription","currency":"USD","price":11.49,"time":"2022-01-22T00:00:00Z","properties":{"trial":true,"trial_end":"2022-02-22T00:00:00Z"}}]}`), b)

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{}`))
	})

	srv := httptest.NewServer(mux)
	defer srv.Close()

	url, _ := url.Parse(srv.URL)
	client, err := braze.NewClient(braze.APIKey("key"), braze.BaseURL(url))
	assert.NoError(t, err)

	attr := &braze.UserPurchase{
		ExternalID: braze.String("123"),
		AppID:      braze.String("diet_doctor"),
		ProductID:  *braze.String("subscription"),
		Currency:   iso4217.USD.Alpha(),
		Price:      *braze.Float64(11.49),
		Time:       time.Now().Format(time.RFC3339),
		Properties: map[string]interface{}{
			"trial": true,
			// Trial period - 1 month
			"trial_end": time.Now().AddDate(0, 1, 0).Format(time.RFC3339),
		},
	}

	resp, err := client.Users.Track(context.Background(), &braze.UsersTrackRequest{
		Purchases: []*braze.UserPurchase{attr},
	})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestUsersServiceTrackPurchaseQuantity(t *testing.T) {
	f := faketime.NewFaketime(2022, time.January, 22, 0, 0, 0, 0, time.UTC)
	defer f.Undo()
	f.Do()

	mux := http.NewServeMux()
	mux.HandleFunc("/users/track", func(w http.ResponseWriter, r *http.Request) {
		b, _ := ioutil.ReadAll(r.Body)
		assert.Equal(t, []byte(`{"purchases":[{"external_id":"123","app_id":"diet_doctor","product_id":"subscription","currency":"USD","price":11.49,"quantity":1,"time":"2022-01-22T00:00:00Z","properties":{"trial":true,"trial_end":"2022-02-22T00:00:00Z"}}]}`), b)

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{}`))
	})

	srv := httptest.NewServer(mux)
	defer srv.Close()

	url, _ := url.Parse(srv.URL)
	client, err := braze.NewClient(braze.APIKey("key"), braze.BaseURL(url))
	assert.NoError(t, err)

	attr := &braze.UserPurchase{
		ExternalID: braze.String("123"),
		AppID:      braze.String("diet_doctor"),
		ProductID:  *braze.String("subscription"),
		Currency:   iso4217.USD.Alpha(),
		Price:      *braze.Float64(11.49),
		Time:       time.Now().Format(time.RFC3339),
		Properties: map[string]interface{}{
			"trial": true,
			// Trial period - 1 month
			"trial_end": time.Now().AddDate(0, 1, 0).Format(time.RFC3339),
		},
	}
	// Quantity - default 1
	err = attr.SetQuantity(1)

	assert.NoError(t, err)

	resp, err := client.Users.Track(context.Background(), &braze.UsersTrackRequest{
		Purchases: []*braze.UserPurchase{attr},
	})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestUsersServiceTrackPurchaseQuantityOutOfRange(t *testing.T) {

	attr := &braze.UserPurchase{
		ExternalID: braze.String("123"),
		AppID:      braze.String("diet_doctor"),
		ProductID:  *braze.String("subscription"),
		Currency:   iso4217.USD.Alpha(),
		Price:      *braze.Float64(11.49),
		Time:       time.Now().Format(time.RFC3339),
		Properties: map[string]interface{}{
			"trial": true,
			// Trial period - 1 month
			"trial_end": time.Now().AddDate(0, 1, 0).Format(time.RFC3339),
		},
	}
	// Quantity - range 1 to 100
	err := attr.SetQuantity(101)

	assert.Error(t, err)
}
