package braze

import (
	"context"
	"net/http"
)

const (
	messagingMessagesSendPath         = "/messages/send"
	messagingTransactionalSendPath    = "/transactional/v1/campaigns/%s/send"
	messagingCampaignsTriggerSendPath = "/campaigns/trigger/send"
)

var (
	ApplePushMessageInterruptionLevelPassive       ApplePushMessageInterruptionLevel = "passive"
	ApplePushMessageInterruptionLevelActive        ApplePushMessageInterruptionLevel = "active"
	ApplePushMessageInterruptionLevelTimeSensitive ApplePushMessageInterruptionLevel = "time-sensitive"
	ApplePushMessageInterruptionLevelCritical      ApplePushMessageInterruptionLevel = "critical"

	ApplePushMessageFileTypeAIF ApplePushMessageFileType = "aif"
	ApplePushMessageFileTypeGIF ApplePushMessageFileType = "gif"
	ApplePushMessageFileTypeJPG ApplePushMessageFileType = "jpg"
	ApplePushMessageFileTypeM4A ApplePushMessageFileType = "m4a"
	ApplePushMessageFileTypeMP3 ApplePushMessageFileType = "mp3"
	ApplePushMessageFileTypeMP4 ApplePushMessageFileType = "mp4"
	ApplePushMessageFileTypePNG ApplePushMessageFileType = "png"
	ApplePushMessageFileTypeWAV ApplePushMessageFileType = "wav"
)

type MessagingEndpoint interface {
	SendMessages(context.Context, *SendMessagesRequest) (*Response, error)
	TriggerCampaign(context.Context, *TriggerCampaignRequest) (*Response, error)
}

var _ MessagingEndpoint = (*MessagingService)(nil)

type MessagingService struct {
	client *Client
}

func (s *MessagingService) SendMessages(ctx context.Context, r *SendMessagesRequest) (*Response, error) {
	req, err := s.client.http.newRequest(http.MethodPost, messagingMessagesSendPath, r)
	if err != nil {
		return nil, err
	}

	var res Response
	if err := s.client.http.do(ctx, req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (s *MessagingService) TriggerCampaign(ctx context.Context, r *TriggerCampaignRequest) (*Response, error) {
	req, err := s.client.http.newRequest(http.MethodPost, messagingCampaignsTriggerSendPath, r)
	if err != nil {
		return nil, err
	}

	var res Response
	if err := s.client.http.do(ctx, req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

type SendMessagesRequest struct {
	Messages *Messages
}

type Messages struct {
	AndroidPush *AndroidPushMessage
	ApplePush   *ApplePushMessage
	Email       *EmailMessage
}

// https://www.braze.com/docs/api/objects_filters/messaging/android_object/
type AndroidPushMessage struct {
	Alert                      string                         `json:"alert"`
	Title                      string                         `json:"title"`
	Extra                      *string                        `json:"extra,omitempty"`
	MessageVariationID         *string                        `json:"message_variation_id,omitempty"`
	NotificationChannelID      *string                        `json:"notification_channel_id,omitempty"`
	Priority                   *int                           `json:"priority,omitempty"`
	SendToSync                 *bool                          `json:"send_to_sync,omitempty"`
	CollapseKey                *string                        `json:"collapse_key,omitempty"`
	Sound                      *string                        `json:"sound,omitempty"`
	CustomURI                  *string                        `json:"custom_uri,omitempty"`
	SummaryText                *string                        `json:"summary_text,omitempty"`
	TimeToLive                 *int                           `json:"time_to_live,omitempty"`
	NotificationID             *int                           `json:"notification_id,omitempty"`
	PushIconImageURL           *string                        `json:"push_icon_image_url,omitempty"`
	AccentColor                *int                           `json:"accent_color,omitempty"`
	SendToMostRecentDeviceOnly *bool                          `json:"send_to_most_recent_device_only,omitempty"`
	Buttons                    []*AndroidPushActionButton     `json:"buttons,omitempty"`
	ConversationData           []*AndroidPushConversationData `json:"conversation_data,omitempty"`
}

type AndroidPushActionButton struct {
	Text       string  `json:"text"`
	Action     *string `json:"action,omitempty"`
	URI        *string `json:"uri,omitempty"`
	UseWebview *string `json:"use_webview,omitempty"`
}

type AndroidPushConversationData struct {
	ShortcutID    string                            `json:"shortcut_id"`
	ReplyPersonID string                            `json:"reply_person_id"`
	Messages      []*AndroidPushConversationMessage `json:"messages"`
	Persons       []*AndroidPushConversationPerson  `json:"persons"`
}

type AndroidPushConversationMessage struct {
	Text      string `json:"text"`
	Timestamp int    `json:"timestamp"`
	PersonID  string `json:"person_id"`
}

type AndroidPushConversationPerson struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// https://www.braze.com/docs/api/objects_filters/messaging/apple_object/
type ApplePushMessage struct {
	Badge                      *int                               `json:"badge,omitempty"`
	Alert                      *ApplePushAlert                    `json:"alert,omitempty"`
	Sound                      *string                            `json:"sound,omitempty"`
	Extra                      *string                            `json:"extra,omitempty"`
	ContentAvailable           *bool                              `json:"content-available,omitempty"`
	InterruptionLevel          *ApplePushMessageInterruptionLevel `json:"interruption_level,omitempty"`
	RelevanceScore             *float64                           `json:"relevance_score,omitempty"`
	Expiry                     *string                            `json:"expiry,omitempty"`
	CustomURI                  *string                            `json:"custom_uri,omitempty"`
	MessageVariationID         *string                            `json:"message_variation_id,omitempty"`
	NotificationGroupThreadID  *string                            `json:"notification_group_thread_id,omitempty"`
	AssetURL                   *string                            `json:"asset_url,omitempty"`
	AssetFileType              *ApplePushMessageFileType          `json:"asset_file_type,omitempty"`
	CollapseID                 *string                            `json:"collapse_id,omitempty"`
	MutableContent             *bool                              `json:"mutable_content,omitempty"`
	SendToMostRecentDeviceOnly *bool                              `json:"send_to_most_recent_device_only,omitempty"`
	Category                   *string                            `json:"category,omitempty"`
	Buttons                    []*ApplePushActionButton           `json:"buttons,omitempty"`
}

type ApplePushAlert struct {
	Body         string    `json:"body"`
	Title        *string   `json:"title,omitempty"`
	TitleLocKey  *string   `json:"title_loc_key,omitempty"`
	TitleLocArgs []string  `json:"title_loc_args,omitempty"`
	ActionLocKey *string   `json:"action_loc_key,omitempty"`
	LocKey       *string   `json:"loc_key,omitempty"`
	LocArgs      []*string `json:"loc_args,omitempty"`
}

type (
	ApplePushMessageInterruptionLevel string
	ApplePushMessageFileType          string
)

type ApplePushActionButton struct {
	ActionID   string `json:"action_id"`
	Action     string `json:"action"`
	URI        string `json:"uri"`
	UseWebview *bool  `json:"use_webview,omitempty"`
}

// https://www.braze.com/docs/api/objects_filters/messaging/email_object/
type EmailMessage struct {
	AppID              string                    `json:"app_id"`
	Subject            *string                   `json:"subject,omitempty"`
	From               string                    `json:"from"`
	ReplyTo            *string                   `json:"reply_to,omitempty"`
	BCC                *string                   `json:"bcc,omitempty"`
	Body               *string                   `json:"body,omitempty"`
	PlaintextBody      *string                   `json:"plaintext_body,omitempty"`
	PreHeader          *string                   `json:"preheader,omitempty"`
	EmailTemplateID    *string                   `json:"email_template_id,omitempty"`
	MessageVariationID *string                   `json:"message_variation_id,omitempty"`
	Extras             map[string]interface{}    `json:"extras,omitempty"`
	Headers            map[string]interface{}    `json:"headers,omitempty"`
	ShouldInlineCSS    *bool                     `json:"should_inline_css,omitempty"`
	Attachments        []*EmailMessageAttachment `json:"attachments,omitempty"`
}

type EmailMessageAttachment struct {
	FileName string `json:"file_name"`
	URL      string `json:"url"`
}

type TriggerCampaignRequest struct {
	CampaignID        string                 `json:"campaign_id,omitempty"`
	SendID            *string                `json:"send_id,omitempty"`
	TriggerProperties map[string]interface{} `json:"trigger_properties,omitempty"`
	Broadcast         *bool                  `json:"broadcast,omitempty"`
	Recipients        []*Recipient           `json:"recipients,omitempty"`
	// Audience TODO
}

type Recipient struct {
	UserAlias             *UserAlias             `json:"user_alias,omitempty"`
	ExternalUserID        *string                `json:"external_user_id,omitempty"`
	TriggerProperties     map[string]interface{} `json:"trigger_properties,omitempty"`
	CanvasEntryProperties map[string]interface{} `json:"canvas_entry_properties,omitempty"`
	SendToExistingOnly    *bool                  `json:"send_to_existing_only,omitempty"`
}
