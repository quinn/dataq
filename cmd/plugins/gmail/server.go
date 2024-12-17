package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	pb "go.quinn.io/dataq/rpc"
	"google.golang.org/api/gmail/v1"
)

type server struct {
	pb.UnimplementedDataQPluginServer
	plugin *GmailPlugin
}

func NewServer(p *GmailPlugin) *server {
	return &server{plugin: p}
}

func (s *server) Extract(ctx context.Context, req *pb.ExtractRequest) (*pb.ExtractResponse, error) {
	srv, err := s.plugin.getClient(ctx)
	if err != nil {
		return nil, err
	}

	switch req.Kind {
	case "initial", "next_page":
		return s.handlePageExtract(ctx, req, srv)
	case "get_message":
		return s.handleMessageExtract(ctx, req, srv)
	default:
		return nil, fmt.Errorf("unknown extract kind: %s", req.Kind)
	}
}

func (s *server) Transform(ctx context.Context, req *pb.TransformRequest) (*pb.TransformResponse, error) {
	switch req.Kind {
	case "page":
		return s.handlePageTransform(ctx, req)
	case "message":
		return s.handleMessageTransform(ctx, req)
	default:
		return nil, fmt.Errorf("unknown transform kind: %s", req.Kind)
	}
}

func (s *server) handlePageExtract(ctx context.Context, req *pb.ExtractRequest, srv *gmail.Service) (*pb.ExtractResponse, error) {
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

	resp := &pb.ExtractResponse{
		Kind:        "page",
		RequestHash: req.Hash,
		Transforms: []*pb.ExtractResponse_Transform{{
			Kind: "page",
		}},
	}

	// If there's a next page, add an extract for it
	if r.NextPageToken != "" {
		resp.Transforms = append(resp.Transforms, &pb.ExtractResponse_Transform{
			Kind: "next_page",
			Metadata: map[string]string{
				"next_page_token": r.NextPageToken,
			},
		})
	}

	return resp, nil
}

func (s *server) handleMessageExtract(ctx context.Context, req *pb.ExtractRequest, srv *gmail.Service) (*pb.ExtractResponse, error) {
	messageID, ok := req.Metadata["message_id"]
	if !ok || messageID == "" {
		return nil, fmt.Errorf("get_message request requires message_id in metadata")
	}

	msg, err := srv.Users.Messages.Get("me", messageID).Do()
	if err != nil {
		return nil, fmt.Errorf("error getting message %s: %v", messageID, err)
	}

	// Store the raw response
	rawJSON, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("error marshaling message: %v", err)
	}

	return &pb.ExtractResponse{
		Kind:        "message",
		RequestHash: req.Hash,
		Transforms: []*pb.ExtractResponse_Transform{{
			Kind: "message",
		}},
	}, nil
}

func (s *server) handlePageTransform(ctx context.Context, req *pb.TransformRequest) (*pb.TransformResponse, error) {
	var pageData gmail.ListMessagesResponse
	if err := json.Unmarshal([]byte(req.Hash), &pageData); err != nil {
		return nil, fmt.Errorf("error unmarshaling page data: %v", err)
	}

	resp := &pb.TransformResponse{
		Kind:        "page",
		RequestHash: req.Hash,
	}

	// Create extract requests for each message
	for _, msg := range pageData.Messages {
		resp.Extracts = append(resp.Extracts, &pb.TransformResponse_Extract{
			Kind: "get_message",
			Metadata: map[string]string{
				"message_id": msg.Id,
			},
		})
	}

	return resp, nil
}

func (s *server) handleMessageTransform(ctx context.Context, req *pb.TransformRequest) (*pb.TransformResponse, error) {
	var msgData gmail.Message
	if err := json.Unmarshal([]byte(req.Hash), &msgData); err != nil {
		return nil, fmt.Errorf("error unmarshaling message data: %v", err)
	}

	// Extract email fields from Gmail message
	email := extractEmailFromMessage(&msgData)

	resp := &pb.TransformResponse{
		Kind:        "message",
		RequestHash: req.Hash,
		Permanodes: []*pb.TransformResponse_Permanode{{
			Kind: "email",
			Key:  msgData.Id,
			Payload: &pb.TransformResponse_Permanode_Email{
				Email: email,
			},
		}},
	}

	return resp, nil
}

func extractEmailFromMessage(msg *gmail.Message) *pb.Email {
	email := &pb.Email{
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
