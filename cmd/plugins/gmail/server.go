package main

import (
	"context"
	"encoding/json"
	"fmt"

	"go.quinn.io/dataq/rpc"
	"google.golang.org/api/gmail/v1"
)

type server struct {
	rpc.UnimplementedDataQPluginServer
	p *GmailPlugin
	// service *gmail.Service
}

func NewServer(p *GmailPlugin) *server {
	return &server{
		p: p,
	}
}

func (s *server) Install(ctx context.Context, req *rpc.InstallRequest) (*rpc.InstallResponse, error) {
	return &rpc.InstallResponse{
		PluginId: "gmail",
		Oauth: &rpc.OAuth2{
			Config: &rpc.OAuth2_Config{
				Endpoint: &rpc.OAuth2_Endpoint{
					AuthUrl:  "https://accounts.google.com/o/oauth2/auth",
					TokenUrl: "https://oauth2.googleapis.com/token",
				},
				Scopes: []string{
					"https://mail.google.com/",
				},
			},
		},
		Extracts: []*rpc.InstallResponse_Extract{
			{
				Kind:        "initial",
				Label:       "Initial",
				Description: "Get initial page of messages",
			},
			{
				Kind:        "next_page",
				Label:       "Next Page",
				Description: "Get next page of messages",
				Configs: []*rpc.PluginConfig{
					{
						Key:   "next_page_token",
						Label: "Next Page Token",
					},
				},
			},
			{
				Kind:        "get_message",
				Label:       "Get Message",
				Description: "Get a specific message",
				Configs: []*rpc.PluginConfig{
					{
						Key:   "message_id",
						Label: "Message ID",
					},
				},
			},
		},
	}, nil
}

func (s *server) Extract(ctx context.Context, req *rpc.ExtractRequest) (*rpc.ExtractResponse, error) {
	var resp *rpc.ExtractResponse
	var err error
	switch req.Kind {
	case "initial", "next_page":
		resp, err = s.handlePageExtract(ctx, req)
	case "get_message":
		resp, err = s.handleMessageExtract(ctx, req)
	default:
		return nil, fmt.Errorf("unknown extract kind: %s", req.Kind)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to run handler for %s: %v", req.Kind, err)
	}

	return resp, nil
}

func (s *server) Transform(ctx context.Context, req *rpc.TransformRequest) (*rpc.TransformResponse, error) {
	switch req.Kind {
	case "page":
		return s.handlePageTransform(ctx, req)
	case "message":
		return s.handleMessageTransform(ctx, req)
	default:
		return nil, fmt.Errorf("unknown transform kind: %s", req.Kind)
	}
}

func (s *server) handlePageExtract(_ context.Context, req *rpc.ExtractRequest) (*rpc.ExtractResponse, error) {
	srv, err := s.p.getClient(context.Background(), req.Oauth)
	if err != nil {
		return nil, fmt.Errorf("failed to get client: %v", err)
	}
	gmailReq := srv.Users.Messages.List("me").MaxResults(100)

	if req.Kind == "next_page" {
		pageToken, ok := req.Metadata["next_page_token"]
		if !ok || pageToken == "" {
			return nil, fmt.Errorf("next_page request requires next_page_token in metadata")
		}
		gmailReq.PageToken(pageToken)
	}

	r, err := gmailReq.Do()
	if err != nil {
		return nil, fmt.Errorf("error listing messages: %v", err)
	}

	if r.Messages == nil {
		return nil, fmt.Errorf("no messages found")
	}

	// Store the raw response
	rawJSON, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("error marshaling response: %v", err)
	}

	resp := &rpc.ExtractResponse{
		Kind: "page",
		Data: &rpc.ExtractResponse_Content{Content: rawJSON},
		Transforms: []*rpc.ExtractResponse_Transform{{
			Kind: "page",
		}},
	}

	return resp, nil
}

func (s *server) handleMessageExtract(_ context.Context, req *rpc.ExtractRequest) (*rpc.ExtractResponse, error) {
	messageID, ok := req.Metadata["message_id"]
	if !ok || messageID == "" {
		return nil, fmt.Errorf("get_message request requires message_id in metadata")
	}

	srv, err := s.p.getClient(context.Background(), req.Oauth)
	if err != nil {
		return nil, fmt.Errorf("failed to get client: %v", err)
	}

	msg, err := srv.Users.Messages.Get("me", messageID).Do()
	if err != nil {
		return nil, fmt.Errorf("error getting message: %v", err)
	}

	// Store the raw response
	rawJSON, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("error marshaling message: %v", err)
	}

	resp := &rpc.ExtractResponse{
		Kind: "message",
		Data: &rpc.ExtractResponse_Content{Content: rawJSON},
		Transforms: []*rpc.ExtractResponse_Transform{{
			Kind: "message",
		}},
	}

	return resp, nil
}

func (s *server) handlePageTransform(_ context.Context, req *rpc.TransformRequest) (*rpc.TransformResponse, error) {
	var pageData gmail.ListMessagesResponse
	if err := json.Unmarshal(req.GetContent(), &pageData); err != nil {
		return nil, fmt.Errorf("error unmarshaling page data: %v", err)
	}

	resp := &rpc.TransformResponse{
		Kind: "page",
	}

	// If there's a next page, add an extract for it
	if pageData.NextPageToken != "" {
		resp.Extracts = append(resp.Extracts, &rpc.TransformResponse_Extract{
			Kind: "next_page",
			Metadata: map[string]string{
				"next_page_token": pageData.NextPageToken,
			},
		})
	}

	// Create extract requests for each message
	for _, msg := range pageData.Messages {
		resp.Extracts = append(resp.Extracts, &rpc.TransformResponse_Extract{
			Kind: "get_message",
			Metadata: map[string]string{
				"message_id": msg.Id,
			},
		})
	}

	return resp, nil
}

func (s *server) handleMessageTransform(_ context.Context, req *rpc.TransformRequest) (*rpc.TransformResponse, error) {
	var msgData gmail.Message
	if err := json.Unmarshal(req.GetContent(), &msgData); err != nil {
		return nil, fmt.Errorf("error unmarshaling message data: %v", err)
	}

	// Extract email fields from Gmail message
	email := extractEmailFromMessage(&msgData)

	resp := &rpc.TransformResponse{
		Kind: "message",
		Permanodes: []*rpc.TransformResponse_Permanode{{
			Kind: "email",
			Key:  msgData.Id,
			Payload: &rpc.TransformResponse_Permanode_Email{
				Email: email,
			},
		}},
	}

	return resp, nil
}

func extractEmailFromMessage(msg *gmail.Message) *rpc.Email {
	email := &rpc.Email{
		MessageId: msg.Id,
		ThreadId:  msg.ThreadId,
	}

	for _, header := range msg.Payload.Headers {
		switch header.Name {
		case "From":
			email.From = header.Value
		case "To":
			email.To = header.Value
		case "Subject":
			email.Subject = header.Value
		case "Date":
			email.Date = header.Value
		case "Cc":
			email.Cc = header.Value
		case "Bcc":
			email.Bcc = header.Value
		case "In-Reply-To":
			email.InReplyTo = header.Value
		case "References":
			email.References = header.Value
		case "Content-Type":
			email.ContentType = header.Value
		}
	}

	// Handle message body
	if msg.Payload != nil {
		if len(msg.Payload.Parts) > 0 {
			for _, part := range msg.Payload.Parts {
				switch part.MimeType {
				case "text/plain":
					email.Text = part.Body.Data
				case "text/html":
					email.Html = part.Body.Data
				}
			}
		} else if msg.Payload.Body != nil {
			// Single part message
			if msg.Payload.MimeType == "text/html" {
				email.Html = msg.Payload.Body.Data
			} else {
				email.Text = msg.Payload.Body.Data
			}
		}
		email.MimeType = msg.Payload.MimeType
	}

	return email
}
