// Code generated by protoc-gen-dataq-index. DO NOT EDIT.

package rpc

func (m *OAuth2) SchemaKind() string {
	return "OAuth2"
}

func (m *OAuth2) SchemaMetadata() map[string]interface{} {
	metadata := make(map[string]interface{})

	if m.Config != nil {
		metadata["config"] = m.Config
	}
	if m.Token != nil {
		metadata["token"] = m.Token
	}
	return metadata
}
