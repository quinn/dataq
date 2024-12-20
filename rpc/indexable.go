package rpc

func (e *ExtractRequest) SchemaKind() string {
	return "ExtractRequest"
}

func (e *ExtractRequest) SchemaMetadata() map[string]interface{} {
	return map[string]interface{}{
		"kind":      e.Kind,
		"plugin_id": e.PluginId,
		"hash":      e.Hash,
		"metadata":  e.Metadata,
	}
}
