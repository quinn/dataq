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
	if len(m.Configs) > 0 {
		metadata["configs"] = m.Configs
	}
	if m.Oauth != nil {
		metadata["oauth"] = m.Oauth
	}
	if len(m.Extracts) > 0 {
		metadata["extracts"] = m.Extracts
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
	if m.Label != "" {
		metadata["label"] = m.Label
	}
	return metadata
}

func (m *ExtractRequest) SchemaKind() string {
	return "ExtractRequest"
}

func (m *ExtractRequest) SchemaMetadata() map[string]interface{} {
	metadata := make(map[string]interface{})

	if m.Oauth != nil {
		metadata["oauth"] = m.Oauth
	}
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
	if len(m.Transforms) > 0 {
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
	if len(m.Extracts) > 0 {
		metadata["extracts"] = m.Extracts
	}
	if len(m.Permanodes) > 0 {
		metadata["permanodes"] = m.Permanodes
	}
	return metadata
}
