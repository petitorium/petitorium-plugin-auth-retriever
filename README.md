# Petitorium Auth Retriever Plugin

An authentication retrieval plugin for [Petitorium](https://github.com/petitorium/petitorium) that automatically captures authentication tokens from responses and stores them in environment variables for flexible usage.

## Features

- **Token Capture**: Automatically extracts authentication tokens from login responses
- **Environment Variable Storage**: Stores captured tokens in environment variables for user-defined usage
- **Configurable Token Paths**: Flexible JSON path configuration for token extraction
- **URL Pattern Matching**: Configurable patterns to identify authentication endpoints
- **Optional Logging**: Configurable logging with custom log file locations
- **Plugin Configuration**: Supports authentication patterns and token extraction settings

## Installation

### Prerequisites

- Go 1.23 or later
- Petitorium with plugin support

### Building from Source

1. Clone this repository:

```bash
git clone https://github.com/petitorium/petitorium-plugin-auth-retriever.git
cd petitorium-plugin-auth-retriever
```

2. Build the plugin:

```bash
go build -o auth-retriever .
```

3. Copy the plugin to your Petitorium plugins directory:

```bash
cp auth-retriever ~/.config/petitorium/plugins/
```

## Configuration

Add the plugin to your Petitorium configuration file (`~/.config/petitorium/config.yaml`):

```yaml
plugins:
  enabled:
    - auth-retriever
  config:
    auth-retriever:
      auth_url_pattern: 'login' # Optional: global URL pattern for auth endpoints (default: "login")
      token_path: 'token' # Optional: global JSON path to extract token (default: "token")
      logging_enabled: false # Optional: global enable logging (default: false)
      log_file: '~/auth-retriever.log' # Optional: global log file path (default: "auth-retriever.log")
      workspaces: # Optional: workspace-specific overrides
        "My Workspace Name":
          auth_url_pattern: 'api/v1/auth'
          token_path: 'data.token'
          logging_enabled: true
          log_file: '~/my-workspace-auth.log'
```

### Configuration Options

| Option             | Type    | Default                | Description                                      |
| ------------------ | ------- | ---------------------- | ------------------------------------------------ |
| `auth_url_pattern` | string  | `"login"`              | URL pattern to identify authentication endpoints |
| `token_path`       | string  | `"token"`              | JSON path for token extraction from responses    |
| `logging_enabled`  | boolean | `false`                | Enable/disable logging                           |
| `log_file`         | string  | `"auth-retriever.log"` | Path to the log file                             |
| `workspaces`       | object  | `null`                 | Map of workspace names to specific overrides     |

## Usage

Once installed and configured, the plugin automatically:

1. **Captures authentication tokens** from login/signin responses
2. **Stores tokens in environment variables** for flexible usage by other plugins or configurations

### Token Capture Process

1. Plugin identifies authentication requests by URL pattern matching
2. On successful responses (HTTP 200), extracts token from JSON response
3. Stores extracted token in environment variables (`auth_token`) for user-defined usage

### Example Log Output (when logging enabled)

```
[2025-11-29 10:30:45] [auth-retriever] Processing auth response from POST https://api.example.com/users/login (status 200), body length: 256
[2025-11-29 10:30:45] [auth-retriever] Parsed JSON response, keys: [token user_id expires_in]
[2025-11-29 10:30:45] [auth-retriever] Captured auth token from POST https://api.example.com/users/login
[2025-11-29 10:30:45] [auth-retriever] Stored token in environment: xyz789ghi012
```

## Plugin Hooks

This plugin implements the following Petitorium hooks:

- `PostReceive`: Captures authentication tokens from responses

## Development

### Project Structure

```
petitorium-plugin-auth-retriever/
├── auth-retriever.go      # Main plugin implementation
├── go.mod                 # Go module definition
├── go.sum                 # Dependency checksums
└── README.md              # This file
```

### Building for Development

```bash
# Download dependencies
go mod tidy

# Build the plugin
go build -o auth-retriever .

# Run tests (if any)
go test ./...
```

### Plugin Architecture

The plugin follows Petitorium's plugin architecture (using HashiCorp go-plugin):

1. **SDK Imports**: Uses `github.com/petitorium/petitorium-plugin-sdk`
2. **Main Plugin**: Implements the `types.Plugin` interface and provides hook functions via `ExecuteHook`
3. **gRPC Server**: Runs a local gRPC server via `plugin.Serve()` for communication with Petitorium
4. **Environment Integration**: Stores tokens in environment variables for flexible usage

## Troubleshooting

### Plugin Not Loading

1. Verify the plugin executable exists: `ls -la ~/.config/petitorium/plugins/auth-retriever`
2. Ensure the file has execution permissions: `chmod +x ~/.config/petitorium/plugins/auth-retriever`
3. Check Petitorium logs for plugin loading errors
4. Ensure the plugin is enabled in your configuration

### Authentication Not Working

1. Check that authentication endpoints match the configured URL pattern
2. Ensure response format matches the configured token path
3. Enable logging to debug token capture process

### Token Not Being Captured

1. Verify the authentication request URL contains the configured pattern
2. Check that the response status is HTTP 200
3. Ensure the response body contains valid JSON
4. Verify the token path configuration matches the response structure
5. Enable logging to see detailed capture process

### Log File Not Created

1. Check file permissions in the target directory
2. Verify the log file path in your configuration
3. Ensure logging is enabled (`logging_enabled: true`)
4. Check that the directory exists (the plugin won't create parent directories)

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature-name`
3. Make your changes and test thoroughly
4. Commit your changes: `git commit -am 'Add new feature'`
5. Push to the branch: `git push origin feature-name`
6. Submit a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Related Projects

- [Petitorium](https://github.com/petitorium/petitorium) - The main API client application
- [Petitorium Auth Injector Plugin](https://github.com/petitorium/petitorium-plugin-auth-injector) - An authentication injection plugin for [Petitorium](https://github.com/petitorium/petitorium) that automatically injects authentication headers into HTTP requests and captures authentication tokens from responses.
- [Petitorium Request Logger Plugin](https://github.com/petitorium/petitorium-plugin-request-logger) - A comprehensive request and response logging plugin for [Petitorium](https://github.com/petitorium/petitorium) that provides detailed logging of HTTP requests and responses with support for both raw template variables and expanded environment variables.

## Support

For issues and questions:

- Create an issue in this repository
- Check the [Petitorium documentation](https://github.com/petitorium/petitorium)
- Review existing issues for similar problems
