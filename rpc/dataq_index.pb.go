// Code generated by protoc-gen-dataq-index. DO NOT EDIT.

package rpc

func (m *InstallRequest) SchemaKind() string {
	return "InstallRequest"
}

func (m *InstallRequest) SchemaMetadata() map[string]interface{} {
	metadata := make(map[string]interface{})

	if m.PluginId != "" {
		metadata["plugin_id"] = m.PluginId
	}
	return metadata
}

func (m *InstallResponse) SchemaKind() string {
	return "InstallResponse"
}

func (m *InstallResponse) SchemaMetadata() map[string]interface{} {
	metadata := make(map[string]interface{})

	if m.PluginId != "" {
		metadata["plugin_id"] = m.PluginId
	}
	if m.Configs != nil {
		metadata["configs"] = m.Configs
	}
	if m.OauthConfig != nil {
		metadata["oauth_config"] = m.OauthConfig
	}
	return metadata
}

func (m *PluginConfig) SchemaKind() string {
	return "PluginConfig"
}

func (m *PluginConfig) SchemaMetadata() map[string]interface{} {
	metadata := make(map[string]interface{})

	if m.Key != "" {
		metadata["key"] = m.Key
	}
	return metadata
}

func (m *OauthConfig) SchemaKind() string {
	return "OauthConfig"
}

func (m *OauthConfig) SchemaMetadata() map[string]interface{} {
	metadata := make(map[string]interface{})

	if m.ClientId != "" {
		metadata["client_id"] = m.ClientId
	}
	if m.ClientSecret != "" {
		metadata["client_secret"] = m.ClientSecret
	}
	if m.RedirectUri != "" {
		metadata["redirect_uri"] = m.RedirectUri
	}
	if m.Scopes != "" {
		metadata["scopes"] = m.Scopes
	}
	if m.AuthUrl != "" {
		metadata["auth_url"] = m.AuthUrl
	}
	if m.TokenUrl != "" {
		metadata["token_url"] = m.TokenUrl
	}
	return metadata
}

func (m *ExtractRequest) SchemaKind() string {
	return "ExtractRequest"
}

func (m *ExtractRequest) SchemaMetadata() map[string]interface{} {
	metadata := make(map[string]interface{})

	if m.PluginId != "" {
		metadata["plugin_id"] = m.PluginId
	}
	if m.ParentHash != "" {
		metadata["parent_hash"] = m.ParentHash
	}
	if m.Kind != "" {
		metadata["kind"] = m.Kind
	}
	if m.Metadata != nil {
		metadata["metadata"] = m.Metadata
	}
	return metadata
}

func (m *ExtractResponse) SchemaKind() string {
	return "ExtractResponse"
}

func (m *ExtractResponse) SchemaMetadata() map[string]interface{} {
	metadata := make(map[string]interface{})

	if m.Kind != "" {
		metadata["kind"] = m.Kind
	}
	if m.RequestHash != "" {
		metadata["request_hash"] = m.RequestHash
	}
	if m.Transforms != nil {
		metadata["transforms"] = m.Transforms
	}
	if m.Data != nil {
		switch {
		case m.GetHash() != "":
			metadata["data_hash"] = m.GetHash()
		case m.GetContent() != nil:
			metadata["data_content"] = m.GetContent()
		}
	}
	return metadata
}

func (m *TransformRequest) SchemaKind() string {
	return "TransformRequest"
}

func (m *TransformRequest) SchemaMetadata() map[string]interface{} {
	metadata := make(map[string]interface{})

	if m.PluginId != "" {
		metadata["plugin_id"] = m.PluginId
	}
	if m.Kind != "" {
		metadata["kind"] = m.Kind
	}
	if m.Metadata != nil {
		metadata["metadata"] = m.Metadata
	}
	if m.Data != nil {
		switch {
		case m.GetHash() != "":
			metadata["data_hash"] = m.GetHash()
		case m.GetContent() != nil:
			metadata["data_content"] = m.GetContent()
		}
	}
	return metadata
}

func (m *TransformResponse) SchemaKind() string {
	return "TransformResponse"
}

func (m *TransformResponse) SchemaMetadata() map[string]interface{} {
	metadata := make(map[string]interface{})

	if m.Kind != "" {
		metadata["kind"] = m.Kind
	}
	if m.RequestHash != "" {
		metadata["request_hash"] = m.RequestHash
	}
	if m.Extracts != nil {
		metadata["extracts"] = m.Extracts
	}
	if m.Permanodes != nil {
		metadata["permanodes"] = m.Permanodes
	}
	return metadata
}

func (m *DataSource) SchemaKind() string {
	return "DataSource"
}

func (m *DataSource) SchemaMetadata() map[string]interface{} {
	metadata := make(map[string]interface{})

	if m.PermanodeHash != "" {
		metadata["permanode_hash"] = m.PermanodeHash
	}
	if m.TransformResponseHash != "" {
		metadata["transform_response_hash"] = m.TransformResponseHash
	}
	if m.Plugin != "" {
		metadata["plugin"] = m.Plugin
	}
	if m.Key != "" {
		metadata["key"] = m.Key
	}
	return metadata
}

func (m *PermanodeVersion) SchemaKind() string {
	return "PermanodeVersion"
}

func (m *PermanodeVersion) SchemaMetadata() map[string]interface{} {
	metadata := make(map[string]interface{})

	if m.PermanodeHash != "" {
		metadata["permanode_hash"] = m.PermanodeHash
	}
	if m.Timestamp != 0 {
		metadata["timestamp"] = m.Timestamp
	}
	if m.Deleted != false {
		metadata["deleted"] = m.Deleted
	}
	if m.Source != nil {
		metadata["source"] = m.Source
	}
	if m.Payload != nil {
		switch {
		case m.GetEmail() != nil:
			metadata["payload_email"] = m.GetEmail()
		case m.GetFinancialTransaction() != nil:
			metadata["payload_financial_transaction"] = m.GetFinancialTransaction()
		}
	}
	return metadata
}
