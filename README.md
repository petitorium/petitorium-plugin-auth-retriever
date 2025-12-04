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
go build -buildmode=plugin -o auth-retriever.so .
```

3. Copy the plugin to your Petitorium plugins directory:

```bash
cp auth-retriever.so ~/.config/petitorium/plugins/
```

## Configuration

Add the plugin to your Petitorium configuration file (`~/.config/petitorium/config.yaml`):

```yaml
plugins:
  enabled:
    - auth-retriever
  config:
    auth-retriever:
      auth_url_pattern: 'login' # Optional: URL pattern for auth endpoints (default: "login")
      token_path: 'token' # Optional: JSON path to extract token (default: "token")
      logging_enabled: false # Optional: enable logging (default: false)
      log_file: '~/auth-retriever.log' # Optional: log file path (default: "auth-retriever.log")
```

### Configuration Options

| Option             | Type    | Default                | Description                                      |
| ------------------ | ------- | ---------------------- | ------------------------------------------------ |
| `auth_url_pattern` | string  | `"login"`              | URL pattern to identify authentication endpoints |
| `token_path`       | string  | `"token"`              | JSON path for token extraction from responses    |
| `logging_enabled`  | boolean | `false`                | Enable/disable logging                           |
| `log_file`         | string  | `"auth-retriever.log"` | Path to the log file                             |

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
go build -buildmode=plugin -o auth-retriever.so .

# Run tests (if any)
go test ./...
```

### Plugin Architecture

The plugin follows Petitorium's plugin architecture:

1. **Main Plugin**: Implements the `Plugin` interface and provides hook functions
2. **Hook Functions**: Process responses at the PostReceive stage
3. **Configuration**: Accepts configuration through the HookContext
4. **Environment Integration**: Stores tokens in environment variables for flexible usage

## Troubleshooting

### Plugin Not Loading

1. Verify the plugin file exists: `ls -la ~/.config/petitorium/plugins/auth-retriever.so`
2. Check Petitorium logs for plugin loading errors
3. Ensure the plugin is enabled in your configuration

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
- [Petitorium Auth Injector Plugin](https://github.com/petitorium/petitorium-plugin-auth-injector) - Authentication injection plugin
- [Petitorium Request Logger Plugin](https://github.com/petitorium/petitorium-plugin-request-logger) - Request/response logging plugin

## Support

For issues and questions:

- Create an issue in this repository
- Check the [Petitorium documentation](https://github.com/petitorium/petitorium)
- Review existing issues for similar problems
