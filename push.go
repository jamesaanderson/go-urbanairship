package urbanairship

import "net/http"

type Push struct {
	Audience     interface{}  `json:"audience"`
	Notification Notification `json:"notification"`
	DeviceTypes  interface{}  `json:"device_types"`
}

type Notification struct {
	Alert string `json:"alert"`
}

func (c *Client) Push(push interface{}) (*http.Response, error) {
	u := "/api/push"
	req, err := c.NewRequest("POST", u, push)
	if err != nil {
		return nil, err
	}
	res, err := c.Do(req, nil)
	if err != nil {
		return nil, err
	}
	return res, nil
}
