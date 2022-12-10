package notifier

import (
	"context"
	"fmt"
	"github.com/fluxcd/pkg/runtime/events"
	"github.com/hashicorp/go-retryablehttp"
	"net/url"
)

type Zoom struct {
	URL      string
	Token    string
	ProxyUrl string
}

type ZoomPayload struct {
	IsMarkdownSupport bool        `json:"is_markdown_support"`
	Content           ZoomContent `json:"content"`
}

type ZoomSubHead struct {
	Text string `json:"text"`
}
type ZoomHead struct {
	Text    string      `json:"text"`
	SubHead ZoomSubHead `json:"sub_head"`
}
type ZoomBody struct {
	Type string `json:"type"`
	Text string `json:"text"`
}
type ZoomContent struct {
	Head ZoomHead   `json:"head"`
	Body []ZoomBody `json:"body"`
}

func NewZoom(hookUrl string, proxyUrl string, token string) (*Zoom, error) {
	_, err := url.ParseRequestURI(hookUrl)
	if err != nil {
		return nil, fmt.Errorf("invalid Zoom hook URL %s: '%w'", hookUrl, err)
	}
	return &Zoom{
		URL:      hookUrl,
		Token:    token,
		ProxyUrl: proxyUrl,
	}, nil
}

func (zoom *Zoom) Post(context context.Context, event events.Event) error {
	if isCommitStatus(event.Metadata, "update") {
		return nil
	}

	payload := ZoomPayload{}

	//color := "green"
	//
	//if event.Severity == events.EventSeverityError {
	//	color = "red"
	//}

	zFields := make([]ZoomBody, 0, len(event.Metadata))

	for k, v := range event.Metadata {
		zFields = append(zFields, ZoomBody{k, v})
	}

	err := postMessage(context, zoom.URL, zoom.ProxyUrl, nil, payload, func(request *retryablehttp.Request) {
		if zoom.Token != "" {
			request.Header.Add("Authorization", "Bearer "+zoom.Token)
		}
	})

	if err != nil {
		return fmt.Errorf("postMessage failed: %w", err)
	}
	return nil

}
