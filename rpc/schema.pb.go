// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        v5.29.1
// source: rpc/schema.proto

package rpc

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Generic entitities inside of dataq
type Email struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	From        string `protobuf:"bytes,1,opt,name=from,proto3" json:"from,omitempty"`
	To          string `protobuf:"bytes,2,opt,name=to,proto3" json:"to,omitempty"`
	Subject     string `protobuf:"bytes,3,opt,name=subject,proto3" json:"subject,omitempty"`
	Body        string `protobuf:"bytes,4,opt,name=body,proto3" json:"body,omitempty"`
	Date        string `protobuf:"bytes,5,opt,name=date,proto3" json:"date,omitempty"`
	MessageId   string `protobuf:"bytes,6,opt,name=message_id,json=messageId,proto3" json:"message_id,omitempty"`
	ThreadId    string `protobuf:"bytes,7,opt,name=thread_id,json=threadId,proto3" json:"thread_id,omitempty"`
	InReplyTo   string `protobuf:"bytes,8,opt,name=in_reply_to,json=inReplyTo,proto3" json:"in_reply_to,omitempty"`
	References  string `protobuf:"bytes,9,opt,name=references,proto3" json:"references,omitempty"`
	Cc          string `protobuf:"bytes,10,opt,name=cc,proto3" json:"cc,omitempty"`
	Bcc         string `protobuf:"bytes,11,opt,name=bcc,proto3" json:"bcc,omitempty"`
	Attachments string `protobuf:"bytes,12,opt,name=attachments,proto3" json:"attachments,omitempty"`
	MimeType    string `protobuf:"bytes,13,opt,name=mime_type,json=mimeType,proto3" json:"mime_type,omitempty"`
	ContentType string `protobuf:"bytes,14,opt,name=content_type,json=contentType,proto3" json:"content_type,omitempty"`
	Content     string `protobuf:"bytes,15,opt,name=content,proto3" json:"content,omitempty"`
	Html        string `protobuf:"bytes,16,opt,name=html,proto3" json:"html,omitempty"`
	Text        string `protobuf:"bytes,17,opt,name=text,proto3" json:"text,omitempty"`
}

func (x *Email) Reset() {
	*x = Email{}
	mi := &file_rpc_schema_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Email) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Email) ProtoMessage() {}

func (x *Email) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_schema_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Email.ProtoReflect.Descriptor instead.
func (*Email) Descriptor() ([]byte, []int) {
	return file_rpc_schema_proto_rawDescGZIP(), []int{0}
}

func (x *Email) GetFrom() string {
	if x != nil {
		return x.From
	}
	return ""
}

func (x *Email) GetTo() string {
	if x != nil {
		return x.To
	}
	return ""
}

func (x *Email) GetSubject() string {
	if x != nil {
		return x.Subject
	}
	return ""
}

func (x *Email) GetBody() string {
	if x != nil {
		return x.Body
	}
	return ""
}

func (x *Email) GetDate() string {
	if x != nil {
		return x.Date
	}
	return ""
}

func (x *Email) GetMessageId() string {
	if x != nil {
		return x.MessageId
	}
	return ""
}

func (x *Email) GetThreadId() string {
	if x != nil {
		return x.ThreadId
	}
	return ""
}

func (x *Email) GetInReplyTo() string {
	if x != nil {
		return x.InReplyTo
	}
	return ""
}

func (x *Email) GetReferences() string {
	if x != nil {
		return x.References
	}
	return ""
}

func (x *Email) GetCc() string {
	if x != nil {
		return x.Cc
	}
	return ""
}

func (x *Email) GetBcc() string {
	if x != nil {
		return x.Bcc
	}
	return ""
}

func (x *Email) GetAttachments() string {
	if x != nil {
		return x.Attachments
	}
	return ""
}

func (x *Email) GetMimeType() string {
	if x != nil {
		return x.MimeType
	}
	return ""
}

func (x *Email) GetContentType() string {
	if x != nil {
		return x.ContentType
	}
	return ""
}

func (x *Email) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

func (x *Email) GetHtml() string {
	if x != nil {
		return x.Html
	}
	return ""
}

func (x *Email) GetText() string {
	if x != nil {
		return x.Text
	}
	return ""
}

type FinancialTransaction struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Date        string `protobuf:"bytes,2,opt,name=date,proto3" json:"date,omitempty"`
	Description string `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	Amount      string `protobuf:"bytes,4,opt,name=amount,proto3" json:"amount,omitempty"`
	Currency    string `protobuf:"bytes,5,opt,name=currency,proto3" json:"currency,omitempty"`
	Category    string `protobuf:"bytes,6,opt,name=category,proto3" json:"category,omitempty"`
	Account     string `protobuf:"bytes,7,opt,name=account,proto3" json:"account,omitempty"`
	Subcategory string `protobuf:"bytes,8,opt,name=subcategory,proto3" json:"subcategory,omitempty"`
	Notes       string `protobuf:"bytes,9,opt,name=notes,proto3" json:"notes,omitempty"`
	Type        string `protobuf:"bytes,10,opt,name=type,proto3" json:"type,omitempty"`
}

func (x *FinancialTransaction) Reset() {
	*x = FinancialTransaction{}
	mi := &file_rpc_schema_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *FinancialTransaction) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FinancialTransaction) ProtoMessage() {}

func (x *FinancialTransaction) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_schema_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FinancialTransaction.ProtoReflect.Descriptor instead.
func (*FinancialTransaction) Descriptor() ([]byte, []int) {
	return file_rpc_schema_proto_rawDescGZIP(), []int{1}
}

func (x *FinancialTransaction) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *FinancialTransaction) GetDate() string {
	if x != nil {
		return x.Date
	}
	return ""
}

func (x *FinancialTransaction) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *FinancialTransaction) GetAmount() string {
	if x != nil {
		return x.Amount
	}
	return ""
}

func (x *FinancialTransaction) GetCurrency() string {
	if x != nil {
		return x.Currency
	}
	return ""
}

func (x *FinancialTransaction) GetCategory() string {
	if x != nil {
		return x.Category
	}
	return ""
}

func (x *FinancialTransaction) GetAccount() string {
	if x != nil {
		return x.Account
	}
	return ""
}

func (x *FinancialTransaction) GetSubcategory() string {
	if x != nil {
		return x.Subcategory
	}
	return ""
}

func (x *FinancialTransaction) GetNotes() string {
	if x != nil {
		return x.Notes
	}
	return ""
}

func (x *FinancialTransaction) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

var File_rpc_schema_proto protoreflect.FileDescriptor

var file_rpc_schema_proto_rawDesc = []byte{
	0x0a, 0x10, 0x72, 0x70, 0x63, 0x2f, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x05, 0x64, 0x61, 0x74, 0x61, 0x71, 0x22, 0xaf, 0x03, 0x0a, 0x05, 0x45, 0x6d,
	0x61, 0x69, 0x6c, 0x12, 0x12, 0x0a, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x12, 0x0e, 0x0a, 0x02, 0x74, 0x6f, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x02, 0x74, 0x6f, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x62, 0x6a, 0x65,
	0x63, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x73, 0x75, 0x62, 0x6a, 0x65, 0x63,
	0x74, 0x12, 0x12, 0x0a, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x62, 0x6f, 0x64, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x65, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x64, 0x61, 0x74, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x49, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x74, 0x68, 0x72, 0x65,
	0x61, 0x64, 0x5f, 0x69, 0x64, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x74, 0x68, 0x72,
	0x65, 0x61, 0x64, 0x49, 0x64, 0x12, 0x1e, 0x0a, 0x0b, 0x69, 0x6e, 0x5f, 0x72, 0x65, 0x70, 0x6c,
	0x79, 0x5f, 0x74, 0x6f, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x69, 0x6e, 0x52, 0x65,
	0x70, 0x6c, 0x79, 0x54, 0x6f, 0x12, 0x1e, 0x0a, 0x0a, 0x72, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e,
	0x63, 0x65, 0x73, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x72, 0x65, 0x66, 0x65, 0x72,
	0x65, 0x6e, 0x63, 0x65, 0x73, 0x12, 0x0e, 0x0a, 0x02, 0x63, 0x63, 0x18, 0x0a, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x02, 0x63, 0x63, 0x12, 0x10, 0x0a, 0x03, 0x62, 0x63, 0x63, 0x18, 0x0b, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x62, 0x63, 0x63, 0x12, 0x20, 0x0a, 0x0b, 0x61, 0x74, 0x74, 0x61, 0x63,
	0x68, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x61, 0x74,
	0x74, 0x61, 0x63, 0x68, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x1b, 0x0a, 0x09, 0x6d, 0x69, 0x6d,
	0x65, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6d, 0x69,
	0x6d, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e,
	0x74, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x63, 0x6f,
	0x6e, 0x74, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e,
	0x74, 0x65, 0x6e, 0x74, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74,
	0x65, 0x6e, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x74, 0x6d, 0x6c, 0x18, 0x10, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x68, 0x74, 0x6d, 0x6c, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x65, 0x78, 0x74, 0x18,
	0x11, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x65, 0x78, 0x74, 0x22, 0x92, 0x02, 0x0a, 0x14,
	0x46, 0x69, 0x6e, 0x61, 0x6e, 0x63, 0x69, 0x61, 0x6c, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x64, 0x61, 0x74, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63,
	0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64,
	0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x6d,
	0x6f, 0x75, 0x6e, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x61, 0x6d, 0x6f, 0x75,
	0x6e, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x12, 0x1a,
	0x0a, 0x08, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x63,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x63, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x12, 0x20, 0x0a, 0x0b, 0x73, 0x75, 0x62, 0x63, 0x61, 0x74, 0x65, 0x67,
	0x6f, 0x72, 0x79, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x73, 0x75, 0x62, 0x63, 0x61,
	0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x6e, 0x6f, 0x74, 0x65, 0x73, 0x18,
	0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6e, 0x6f, 0x74, 0x65, 0x73, 0x12, 0x12, 0x0a, 0x04,
	0x74, 0x79, 0x70, 0x65, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65,
	0x42, 0x17, 0x5a, 0x15, 0x67, 0x6f, 0x2e, 0x71, 0x75, 0x69, 0x6e, 0x6e, 0x2e, 0x69, 0x6f, 0x2f,
	0x64, 0x61, 0x74, 0x61, 0x71, 0x2f, 0x72, 0x70, 0x63, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_rpc_schema_proto_rawDescOnce sync.Once
	file_rpc_schema_proto_rawDescData = file_rpc_schema_proto_rawDesc
)

func file_rpc_schema_proto_rawDescGZIP() []byte {
	file_rpc_schema_proto_rawDescOnce.Do(func() {
		file_rpc_schema_proto_rawDescData = protoimpl.X.CompressGZIP(file_rpc_schema_proto_rawDescData)
	})
	return file_rpc_schema_proto_rawDescData
}

var file_rpc_schema_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_rpc_schema_proto_goTypes = []any{
	(*Email)(nil),                // 0: dataq.Email
	(*FinancialTransaction)(nil), // 1: dataq.FinancialTransaction
}
var file_rpc_schema_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_rpc_schema_proto_init() }
func file_rpc_schema_proto_init() {
	if File_rpc_schema_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_rpc_schema_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_rpc_schema_proto_goTypes,
		DependencyIndexes: file_rpc_schema_proto_depIdxs,
		MessageInfos:      file_rpc_schema_proto_msgTypes,
	}.Build()
	File_rpc_schema_proto = out.File
	file_rpc_schema_proto_rawDesc = nil
	file_rpc_schema_proto_goTypes = nil
	file_rpc_schema_proto_depIdxs = nil
}
