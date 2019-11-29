package ec

import (
	"encoding/json"
	"net/http"
	"strings"

	std "github.com/PKUJohnson/solar/std"
)

// TrackInfoHeader includes the device info
type TrackInfoHeader struct {
	AppID          string `json:"appId"`
	AppVersion     string `json:"appVersion"`
	ConnectionType string `json:"connectionType"`
	Locale         string `json:"locale"`
	Carrier        string `json:"carrier"`
	DeviceBrand    string `json:"deviceBrand"`
	DeviceModel    string `json:"deviceModel"`
	OSName         string `json:"osName"`
	OSVersion      string `json:"osVersion"`
	Resolution     string `json:"resolution"`
	Orientation    string `json:"orientation"`
}

// GetTrackInfo extracts the device track info from header
func GetTrackInfo(req *http.Request) *TrackInfoHeader {
	trackInfoVal := req.Header.Get("X-Track-Info")
	if trackInfoVal == "" {
		return nil
	}

	header := &TrackInfoHeader{}
	if err := json.Unmarshal([]byte(trackInfoVal), header); err != nil {
		std.LogErrorc("json", err, "fail to extract device tracking info")
		return nil
	}
	return header
}

// ResolveAppInfo resolves appID into package id and app type
func ResolveAppInfo(appID string) (string, string) {
	if appID != "" {
		idx := strings.LastIndex(appID, ".")
		if idx > -1 {
			return appID[0:idx], appID[idx+1:]
		}
	}
	return "", ""
}
