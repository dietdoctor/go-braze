package braze

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"sync"
)

const (
	usersTrackPath       = "/users/track"
	usersCreateAliasPath = "/users/alias/new"
	usersDeletePath      = "/users/delete"
	usersIdentifyPath    = "/users/identify"
	usersMergePath       = "/users/merge"
)

var (
	_ UsersEndpoint = (*UsersService)(nil)

	AttributeSubscribeOptedIn      AttributeSubscribe = "opted_in"
	AttributeSubscribeUnsubscribed AttributeSubscribe = "unsubscribed"
	AttributeSubscribeSubscribed   AttributeSubscribe = "subscribed"
	AttributeGenderMale            AttributeGender    = "M"
	AttributeGenderFemale          AttributeGender    = "F"
	AttributeGenderOther           AttributeGender    = "O"
	AttributeGenderNotApplicable   AttributeGender    = "N"
	AttributeGenderPreferNotToSay  AttributeGender    = "P"
)

type UsersEndpoint interface {
	Track(ctx context.Context, r *UsersTrackRequest) (*Response, error)
	Delete(ctx context.Context, r *UsersDeleteRequest) (*Response, error)
	Identify(ctx context.Context, r *UsersIdentifyRequest) (*Response, error)
	CreateAlias(ctx context.Context, r *UsersCreateAliasRequest) (*Response, error)
	Merge(ctx context.Context, r *UsersMergeRequest) (*Response, error)
}

type (
	AttributeSubscribe string
	AttributeGender    string
)

type UsersService struct {
	client *Client
}

type UsersTrackRequest struct {
	Attributes []*UserAttributes `json:"attributes,omitempty"`
	Events     []*UserEvent      `json:"events,omitempty"`
	Purchases  []*UserPurchase   `json:"purchases,omitempty"`
}

type UsersDeleteRequest struct {
	ExternalIDs []string     `json:"external_ids,omitempty"`
	UserAliases []*UserAlias `json:"user_aliases,omitempty"`
	BrazeIDs    []string     `json:"braze_ids,omitempty"`
}

type UsersIdentifyRequest struct{}

type UsersCreateAliasRequest struct{}

type UsersMergeRequest struct {
	MergeUpdates []*UsersMergeUpdates `json:"merge_updates,omitempty"`
}

// https://www.braze.com/docs/api/objects_filters/user_attributes_object/
type UserAttributes struct {
	// Of the unique user identifier.
	ExternalID         *string    `json:"external_id,omitempty"`
	UserAlias          *UserAlias `json:"user_alias,omitempty"`
	BrazeID            *string    `json:"braze_id,omitempty"`
	UpdateExistingOnly *bool      `json:"update_existing_only,omitempty"`
	PushTokenImport    *bool      `json:"push_token_import,omitempty"`
	FirstName          *string    `json:"first_name,omitempty"`
	LastName           *string    `json:"last_name,omitempty"`
	Email              *string    `json:"email,omitempty"`

	// We require that country codes be passed to Braze in the ISO-3166-1 alpha-2 standard .
	Country *string `json:"country,omitempty"`

	// Available values are “opted_in” (explicitly registered to receive email
	// messages), “unsubscribed” (explicitly opted out of email messages), and
	// “subscribed” (neither opted in nor out).
	EmailSubscribe *AttributeSubscribe `json:"email_subscribe,omitempty"`

	// Set to true to disable the open tracking pixel from being added to all future emails sent to this user.
	EmailOpenTrackingDisabled *bool `json:"email_open_tracking_disabled,omitempty"`

	// Set to true to disable the click tracking for all links within a future email, sent to this user.
	EmailClickTrackingDisabled *bool              `json:"email_click_tracking_disabled,omitempty"`
	Facebook                   *AttributeFacebook `json:"facebook,omitempty"`
	Gender                     *AttributeGender   `json:"gender,omitempty"`
	HomeCity                   *string            `json:"home_city,omitempty"`

	// URL of image to be associated with user profile.
	ImageURL *string `json:"image_url,omitempty"`

	// We require that language be passed to Braze in the ISO-639-1 standard . List of accepted Languages
	Language *string `json:"language,omitempty"`

	// Date at which the user’s email was marked as spam. Appears in ISO 8601 format or in yyyy-MM-dd’T’HH:mm:ss:SSSZ format.
	MarkedEmailAsSpamAt *string             `json:"marked_email_as_spam_at,omitempty"`
	Phone               *string             `json:"phone,omitempty"`
	PushSubscribe       *AttributeSubscribe `json:"push_subscribe,omitempty"`

	// Array of objects with app_id and token string. You may optionally provide a
	// device_id for the device this token is associated with, e.g., [{"app_id":
	// App Identifier, "token": "abcd", "device_id": "optional_field_value"}]. If a
	// device_id is not provided, one will be randomly generated.
	PushTokens []*PushToken `json:"push_tokens,omitempty"`

	// Of time zone name from IANA Time Zone Database  (e.g., “America/New_York” or “Eastern Time (US & Canada)”). Only valid time zone values will be set.
	Timezone *string `json:"time_zone,omitempty"`

	// Hash containing any of id (integer), screen_name (string, Twitter handle), followers_count (integer), friends_count (integer), statuses_count (integer).
	Twitter *AttributeTwitter `json:"twitter,omitempty"`

	mu               sync.Mutex
	customAttributes map[string]any
}

func (ua *UserAttributes) AddAttributes(attrs ...CustomAttribute) {
	ua.mu.Lock()
	defer ua.mu.Unlock()

	if ua.customAttributes == nil {
		ua.customAttributes = make(map[string]any, len(attrs))
	}

	for _, a := range attrs {
		ua.customAttributes[a.key] = a.value
	}
}

func (ua *UserAttributes) GetCustomAttributes() map[string]any {
	ua.mu.Lock()
	defer ua.mu.Unlock()

	copy := make(map[string]any, len(ua.customAttributes))
	for k, v := range ua.customAttributes {
		copy[k] = v
	}

	return copy
}

// Marshall twice basically. First to get a map and then add custom attributes and marshal again.
func (attributes *UserAttributes) MarshalJSON() ([]byte, error) {
	// Create a new type to prevent MarshalJSON going into infinite loop.
	type a UserAttributes
	ua := (*a)(attributes)

	d, err := json.Marshal(&ua)
	if err != nil {
		return nil, err
	}

	m := map[string]any{}
	if err := json.Unmarshal(d, &m); err != nil {
		return nil, err
	}

	if attributes.customAttributes != nil {
		for k, v := range attributes.customAttributes {
			m[k] = v
		}
	}
	return json.Marshal(m)
}

// https://www.braze.com/docs/api/objects_filters/user_alias_object/
type UserAlias struct {
	AliasName  string `json:"alias_name"`
	AliasLabel string `json:"alias_label"`
}

type AttributeFacebook struct {
	ID         string   `json:"id"`
	Likes      []string `json:"likes"`
	NumFriends int      `json:"num_friends"`
}

type PushToken struct {
	AppID    string  `json:"app_id"`
	Token    string  `json:"token"`
	DeviceID *string `json:"device_id,omitempty"`
}

type AttributeTwitter struct {
	ID             *string `json:"id,omitempty"`
	FollowersCount *int    `json:"followers_count,omitempty"`
	FriendsCount   *int    `json:"friends_count,omitempty"`
	StatusesCount  *int    `json:"statuses_count,omitempty"`
}

// https://www.braze.com/docs/api/objects_filters/event_object/
type UserEvent struct {
	// One of "external_id" or "user_alias" or "braze_id" is required
	ExternalID *string    `json:"external_id,omitempty"`
	UserAlias  *UserAlias `json:"user_alias,omitempty"`
	BrazeID    *string    `json:"braze_id,omitempty"`
	AppID      *string    `json:"app_id,omitempty"`

	// The name of the event. Required.
	Name string `json:"name"`
	// Datetime as string in ISO 8601 or in `yyyy-MM-dd'T'HH:mm:ss:SSSZ` format). Required.
	Time string `json:"time"`

	// https://www.braze.com/docs/api/objects_filters/event_object/#event-properties-object
	Properties map[string]any `json:"properties,omitempty"`

	// Setting this flag to true will put the API in "Update Only" mode.
	// When using a "user_alias", "Update Only" mode is always true.
	UpdateExistingOnly *bool `json:"update_existing_only,omitempty"`
}

// TODO
type UserPurchase struct{}

type UsersMergeUpdates struct {
	IdentifierToMerge *UsersIdentifierToMerge `json:"identifier_to_merge,omitempty"`
	IdentifierToKeep  *UsersIdentifierToKeep  `json:"identifier_to_keep,omitempty"`
}

// Note: Braze supports more options for merging. See https://www.braze.com/docs/api/endpoints/user_data/post_users_merge/#request-parameters for more details. Update this structure if needed.
type UsersIdentifierToMerge struct {
	ExternalID *string `json:"external_id,omitempty"`
}

// Note: Braze supports more options for keeping. See https://www.braze.com/docs/api/endpoints/user_data/post_users_merge/#request-parameters for more details. Update this structure if needed.
type UsersIdentifierToKeep struct {
	ExternalID *string `json:"external_id,omitempty"`
}

func (s *UsersService) Track(ctx context.Context, r *UsersTrackRequest) (*Response, error) {
	req, err := s.client.http.newRequest(http.MethodPost, usersTrackPath, r)
	if err != nil {
		return nil, err
	}

	var res Response
	if err := s.client.http.do(ctx, req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (s *UsersService) Delete(ctx context.Context, r *UsersDeleteRequest) (*Response, error) {
	req, err := s.client.http.newRequest(http.MethodPost, usersDeletePath, r)
	if err != nil {
		return nil, err
	}

	var res Response
	if err := s.client.http.do(ctx, req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (s *UsersService) Identify(ctx context.Context, r *UsersIdentifyRequest) (*Response, error) {
	panic(errors.New("not implemented"))
}

func (s *UsersService) CreateAlias(ctx context.Context, r *UsersCreateAliasRequest) (*Response, error) {
	panic(errors.New("not implemented"))
}

func (s *UsersService) Merge(ctx context.Context, r *UsersMergeRequest) (*Response, error) {
	req, err := s.client.http.newRequest(http.MethodPost, usersMergePath, r)
	if err != nil {
		return nil, err
	}

	var res Response
	if err := s.client.http.do(ctx, req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}
