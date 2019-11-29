package data

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/PKUJohnson/solar/gateway/helper"
	std "github.com/PKUJohnson/solar/std"
	"github.com/PKUJohnson/solar/std/ec"
	"github.com/PKUJohnson/solar/std/ec/kafka"
)

// Payload wraps the event message.
type Payload struct {
	Timestamp int64             `json:"timestamp"`
	Event     string            `json:"event"`
	Extra     map[string]string `json:"extra"`
}

// UserEventPayload wraps the user event message.
type UserEventPayload struct {
	AppType        string            `json:"app_type"`
	AppName        string            `json:"app_name"`
	AppID          string            `json:"app_id"`
	Backend        string            `json:"backend"`
	AppVersion     string            `json:"app_version"`
	Time           int64             `json:"time"`
	IP             string            `json:"ip"`
	Environment    string            `json:"environment"`
	Referer        string            `json:"referer"`
	UID            int64             `json:"uid"`
	Event          string            `json:"event"`
	OSName         string            `json:"os_name"`
	OSVersion      string            `json:"os_version"`
	Resolution     string            `json:"resolution"`
	ConnectionType string            `json:"connection_type"`
	DeviceBrand    string            `json:"device_brand"`
	DeviceModel    string            `json:"device_model"`
	DeviceID       string            `json:"device_id"`
	LegacyDeviceID string            `json:"legacy_device_id"`
	Carrier        string            `json:"carrier"`
	Locale         string            `json:"locale"`
	Orientation    string            `json:"orientation"`
	ResourceType   string            `json:"resource_type"`
	ResourceID     string            `json:"resource_id"`
	Extra          map[string]string `json:"extra"`
}

var c ec.EventCollector
var enabled bool

// Init initializes the event collector.
func Init(config *ec.Config) {
	enabled = config.Enabled

	if enabled && c == nil {
		c, _ = kafka.NewCollectorWithConfig(config)
	}
}

// Collector returns the event collector.
func Collector() ec.EventCollector {
	return c
}

// CollectUserEvent send user events to kafka.
func CollectUserEvent(ctx echo.Context, event string, resourceType string, resourceID string, extra map[string]string) error {
	if !enabled {
		return nil
	}

	if c == nil {
		std.LogErrorc("data_collector", nil, "event collector need to be initialized!")
		return nil
	}

	if extra == nil {
		extra = make(map[string]string)
	}

	payload := UserEventPayload{
		Backend:        "wscn",
		Time:           time.Now().UnixNano() / 1e6,
		IP:             helper.GetClientIp(ctx),
		Environment:    curEnv(),
		Referer:        ctx.Request().Referer(),
		UID:            helper.GetUserId(ctx),
		Event:          event,
		DeviceID:       helper.GetTaotieDeviceID(ctx.Request()),
		LegacyDeviceID: helper.GetDeviceID(ctx.Request()),
		ResourceID:     resourceID,
		Extra:          extra,
	}

	trackHeader := ec.GetTrackInfo(ctx.Request())
	if trackHeader != nil {
		payload.AppID = trackHeader.AppID
		payload.AppVersion = trackHeader.AppVersion
		payload.ConnectionType = trackHeader.ConnectionType
		payload.Locale = trackHeader.Locale
		payload.Carrier = trackHeader.Carrier
		payload.DeviceBrand = trackHeader.DeviceBrand
		payload.DeviceModel = trackHeader.DeviceModel
		payload.OSName = trackHeader.OSName
		payload.OSVersion = trackHeader.OSVersion
		payload.Resolution = trackHeader.Resolution
		payload.Orientation = trackHeader.Orientation
		payload.AppName, payload.AppType = ec.ResolveAppInfo(trackHeader.AppID)
	}

	p, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	if event == "" {
		return nil
	}

	event = "Lina" + strings.ToUpper(event[0:1]) + event[1:]
	if err := c.CollectByTopic(event, &ec.Event{Body: p}); err != nil {
		std.LogErrorc("event_collector", err, "fail to collect user event")
		return err
	}
	return nil
}

func curEnv() string {
	env := strings.ToLower(os.Getenv("CONFIGOR_ENV"))
	if strings.Index(env, "prod") > -1 {
		return "prod"
	} else if strings.Index(env, "stage") > -1 {
		return "stage"
	}
	return "test"
}

// Collect send events to kafka.
func Collect(ctx echo.Context, event string, extra map[string]string) error {
	if !enabled {
		return nil
	}

	if c == nil {
		panic("event collector need to be initialized!")
	}

	if extra == nil {
		extra = make(map[string]string)
	}

	if ctx != nil {
		extra["user_id"] = strconv.FormatInt(helper.GetUserId(ctx), 10)
		extra["device_id"] = helper.GetDeviceID(ctx.Request())
		extra["device_type"] = helper.GetDeviceType(ctx)
	}

	p, err := json.Marshal(Payload{
		Timestamp: time.Now().Unix(),
		Event:     event,
		Extra:     extra,
	})
	if err != nil {
		return err
	}

	if err := c.Collect(&ec.Event{Body: p}); err != nil {
		std.LogErrorc("event_collector", err, "fail to collect user event")
		return err
	}
	return nil
}

// Close the kafka collector.
func Close() error {
	if c == nil {
		return nil
	}

	if e := c.Close(); e != nil {
		std.LogErrorc("event_collector", e, "fail to close user event collector")
		return e
	}
	return nil
}
