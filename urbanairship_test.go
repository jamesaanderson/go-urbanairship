package urbanairship

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	c := NewClient(nil, "", "")

	if c.BaseURL.String() != baseURL {
		t.Errorf("NewClient BaseURL = %v, want %v", c.BaseURL.String(), baseURL)
	}
	if c.MimeType != mimeType {
		t.Errorf("NewClient MimeType = %v, want %v", c.MimeType, mimeType)
	}
}

func TestNewRequest(t *testing.T) {
	c := NewClient(nil, "", "")

	inURL, outURL := "/foo", baseURL+"/foo"
	req, _ := c.NewRequest("GET", inURL, nil)
	if req.URL.String() != outURL {
		t.Errorf("NewRequest(%v) URL = %v, want %v", inURL, req.URL, outURL)
	}

	mimeType := req.Header.Get("Accept")
	if c.MimeType != mimeType {
		t.Errorf("NewRequest() MimeType = %v, want %v", mimeType, c.MimeType)
	}
}
