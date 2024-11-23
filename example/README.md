# DataQ Example

This example demonstrates how to use DataQ to collect and organize data using plugins.

## Directory Structure

```
example/
├── config.yaml     # Plugin configuration file
├── data/          # Sample data directory
│   ├── document1.txt
│   └── notes.docx
```

## Configuration

The `config.yaml` file specifies which plugins to use and their configuration:

- `enabled`: Whether the plugin should be active
- `config`: Plugin-specific configuration parameters

### File Scanner Plugin Configuration

```yaml
plugins:
  - id: filescan
    enabled: true
    config:
      root_path: "./data"  # Directory to scan
      include_patterns:    # File patterns to include
        - "*.txt"
        - "*.pdf"
        - "*.doc*"
      exclude_patterns:    # File patterns to exclude
        - "*.tmp"
        - "~$*"
```

## Running the Example

1. Make sure you have the DataQ binary in your PATH
2. Run the example:
   ```bash
   dataq -config config.yaml
   ```

The file scanner plugin will scan the `data` directory and collect information about the files matching the specified patterns.
