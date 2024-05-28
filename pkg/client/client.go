package wapi

import (
	"fmt"

	"github.com/sarthakjdev/wapi.go/internal/manager"
	"github.com/sarthakjdev/wapi.go/pkg/events"
	"github.com/sarthakjdev/wapi.go/utils"
)

// Client represents a WhatsApp client.
type Client struct {
	Media             manager.MediaManager
	Message           manager.MessageManager
	Phone             manager.PhoneNumbersManager
	webhook           manager.WebhookManager
	phoneNumberId     string
	apiAccessToken    string
	businessAccountId string
}

// ClientConfig represents the configuration options for the WhatsApp client.
type ClientConfig struct {
	PhoneNumberId     string `validate:"required"`
	ApiAccessToken    string `validate:"required"`
	BusinessAccountId string `validate:"required"`
	WebhookPath       string `validate:"required"`
	WebhookSecret     string `validate:"required"`
	WebhookServerPort int
}

// NewWapiClient creates a new instance of Client.
func New(configs ClientConfig) (*Client, error) {
	err := utils.GetValidator().Struct(configs)
	if err != nil {
		return nil, fmt.Errorf("error validating client config", err)
	}
	requester := *manager.NewRequestClient(configs.PhoneNumberId, configs.ApiAccessToken)
	eventManager := *manager.NewEventManager()
	return &Client{
		Media:             *manager.NewMediaManager(requester),
		Message:           *manager.NewMessageManager(requester),
		Phone:             *manager.NewPhoneNumbersManager(requester),
		webhook:           *manager.NewWebhook(&manager.WebhookManagerConfig{Path: configs.WebhookPath, Secret: configs.WebhookSecret, Port: configs.WebhookServerPort, EventManager: eventManager, Requester: requester}),
		phoneNumberId:     configs.PhoneNumberId,
		apiAccessToken:    configs.ApiAccessToken,
		businessAccountId: configs.BusinessAccountId,
	}, nil
}

// GetPhoneNumberId returns the phone number ID associated with the client.
func (client *Client) GetPhoneNumberId() string {
	return client.phoneNumberId
}

// SetPhoneNumberId sets the phone number ID for the client.
func (client *Client) SetPhoneNumberId(phoneNumberId string) {
	client.phoneNumberId = phoneNumberId
}

// InitiateClient initializes the client and starts listening to events from the webhook.
// It returns true if the client was successfully initiated.
func (client *Client) InitiateClient() bool {

	client.webhook.ListenToEvents()
	return true
}

// OnMessage registers a handler for a specific event type.
func (client *Client) On(eventType manager.EventType, handler func(events.BaseEvent)) {
	client.webhook.EventManager.On(eventType, handler)
}
