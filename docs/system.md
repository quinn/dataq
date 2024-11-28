# DataQ System Design

## Overview
DataQ is a plugin-based data archival and organization system designed to collect, process, and store data from various sources in a unified way. The system follows a modular architecture that allows easy extension through plugins while maintaining a consistent interface for data extraction and processing.

## System Architecture

### Core Components

1. **Plugin System**
   - Plugin Registry: Manages plugin registration and lifecycle
   - Plugin Interface: Defines standard methods for data extraction
   - Plugin Runner: Executes plugins and handles their responses

2. **Data Model**
   - DataItem: Represents extracted data with metadata
   - PluginRequest/Response: Handles communication with plugins
   - Configuration: Manages system and plugin settings

3. **Plugin Communication Protocol**
   - Uses Protocol Buffers for structured data exchange
   - Standardized request/response format
   - Support for streaming data extraction

### Plugin Architecture

Each plugin implements the following interface:
- ID(): Unique identifier
- Name(): Human-readable name
- Description(): Plugin functionality description
- Configure(): Setup with plugin-specific settings
- Extract(): Data extraction implementation

## Data Flow

1. **Configuration Phase**
   - System loads configuration from YAML file
   - Plugins are initialized with their specific settings
   - Plugin registry maintains active plugin instances

2. **Extraction Phase**
   - System iterates through enabled plugins
   - Each plugin receives extraction request
   - Plugins stream DataItems back to core system
   - Data items are processed and stored as needed

3. **Plugin Communication**
   - Plugins run as separate processes
   - Communication via stdin/stdout using Protocol Buffers
   - Support for async streaming of large datasets

## Built-in Plugins

1. **File Scanner Plugin**
   - Scans local filesystem
   - Configurable file patterns
   - Extracts file metadata and content

2. **Gmail Plugin**
   - OAuth2 authentication
   - Email metadata extraction
   - Configurable email filtering

## Extensibility

The system is designed for easy extension through:
1. New plugin development
2. Custom data processors
4. Enhanced metadata extraction

## Configuration Management

- YAML-based configuration
- Plugin-specific settings

## Core Protocol

The system is built around a simple but flexible protocol defined in `proto/dataq.proto`:

1. **Plugin Definition**
   ```protobuf
   message Plugin {
     string id = 1;
     string name = 2;
     string description = 3;
     map<string, string> config = 4;
   }
   ```
   Plugins are self-describing units that can be configured via key-value pairs.

2. **Data Model**
   ```protobuf
   message DataItem {
     string plugin_id = 1;
     string source_id = 2;
     int64 timestamp = 3;
     string content_type = 4;
     bytes raw_data = 5;
     map<string, string> metadata = 6;
   }
   ```
   DataItems are the fundamental unit of data exchange, containing:
   - Source identification
   - Raw data payload
   - Extensible metadata

3. **Plugin Communication**
   ```protobuf
   message PluginRequest {
     string plugin_id = 1;
     map<string, string> config = 2;
     string operation = 3;
   }

   message PluginResponse {
     string plugin_id = 1;
     repeated DataItem items = 2;
     string error = 3;
   }
   ```
   The request/response protocol allows for:
   - Plugin configuration
   - Data extraction operations
   - Error reporting

## Data Processing Pipeline

The pipeline is designed to be generic and extensible, working with the core protocol:

1. **Plugin Operations**
   - Configure: Set up plugin with required parameters
   - Extract: Request data from the source
   - Each plugin determines its own extraction strategy

2. **Data Flow**
   - Plugins stream DataItems to the core system
   - Each DataItem is self-contained
   - Metadata enables plugin-specific features
   - Core system remains agnostic to plugin implementation details

3. **Plugin Implementation Patterns**
   Plugins can implement various data access patterns:
   - Single-shot extraction
   - Stateful extraction (via config/metadata)
   - Incremental updates
   - Streaming real-time data

4. **Processing Stages**
   - Raw data validation against content_type
   - Optional transformation based on metadata
   - Storage according to plugin_id and source_id

## Plugin Implementation Example

A plugin implementing stateful extraction (like Gmail) might:
```yaml
# Initial request
plugin_request:
  plugin_id: "gmail"
  operation: "extract"
  config: {}

# Subsequent request using metadata from previous response
plugin_request:
  plugin_id: "gmail"
  operation: "extract"
  config:
    page_token: "token_from_previous_metadata"
```

The core system remains unaware of the pagination mechanism, while the plugin can maintain its state through the generic config and metadata fields.

## Future Considerations

1. **Scalability**
   - Parallel plugin execution
   - Distributed data processing
   - Cloud storage integration

2. **Enhanced Features**
   - Real-time data monitoring
   - Advanced search capabilities
   - Data deduplication

3. **Integration**
   - API endpoints
   - Event-driven architecture
   - External service webhooks
