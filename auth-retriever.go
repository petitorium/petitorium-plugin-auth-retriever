// Package main provides an example auth retriever plugin.
// To build: go build -buildmode=plugin -o auth-retriever.so .
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/go-plugin"
	"github.com/petitorium/petitorium-plugin-sdk/shared"
	"github.com/petitorium/petitorium-plugin-sdk/types"
)

// logf writes formatted log messages to the configured log file if logging is enabled
func logf(ctx *types.HookContext, format string, args ...interface{}) {
	// Check if logging is enabled in config
	if ctx.Config != nil {
		if pluginConfig, ok := ctx.Config["auth-retriever"].(map[string]interface{}); ok {
			// Check if logging is enabled, default to false
			if enabled, exists := pluginConfig["logging_enabled"].(bool); !exists || !enabled {
				return
			}

			// Get log file path, default to "auth-retriever.log"
			logFile := "auth-retriever.log"
			if filePath, exists := pluginConfig["log_file"].(string); exists && filePath != "" {
				logFile = filePath
			}

			file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return
			}
			defer file.Close()
			timestamp := time.Now().Format("2006-01-02 15:04:05")
			fmt.Fprintf(file, "["+timestamp+"] "+format+"\n", args...)
		}
	}
}

// getNestedValue extracts a value from a nested map using dot-separated path
func getNestedValue(data map[string]interface{}, path string) interface{} {
	keys := strings.Split(path, ".")
	current := data
	for i, key := range keys {
		if val, exists := current[key]; exists {
			if i == len(keys)-1 {
				return val
			}
			if nested, ok := val.(map[string]interface{}); ok {
				current = nested
			} else {
				return nil
			}
		} else {
			return nil
		}
	}
	return nil
}

// getMapKeys returns all keys from a map
func getMapKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// AuthRetriever is an example plugin that retrieves auth tokens and stores them in environment variables
type AuthRetriever struct{}

// Name returns the plugin name
func (ar *AuthRetriever) Name() string {
	return "auth-retriever"
}

// Version returns the plugin version
func (ar *AuthRetriever) Version() string {
	return "1.0.0"
}

// Description returns the plugin description
func (ar *AuthRetriever) Description() string {
	return "Retrieves authentication tokens from responses and stores them in environment variables"
}

// Hooks returns the hook types this plugin implements
func (ar *AuthRetriever) Hooks() []types.HookType {
	return []types.HookType{types.PostReceive}
}

// ExecuteHook executes a specific hook with the given context.
func (ar *AuthRetriever) ExecuteHook(hookType types.HookType, ctx *types.HookContext) (*types.HookContext, error) {
	switch hookType {
	case types.PostReceive:
		err := ar.captureAuth(ctx)
		return ctx, err
	}
	return ctx, nil
}

// captureAuth captures auth token from responses
func (ar *AuthRetriever) captureAuth(ctx *types.HookContext) error {
	// Get auth URL pattern from config, default to "login"
	authURLPattern := "login"
	if ctx.Config != nil {
		if pluginConfig, ok := ctx.Config["auth-retriever"].(map[string]interface{}); ok {
			logf(ctx, "[auth-retriever] Found plugin config: %v", pluginConfig)
			if pattern, exists := pluginConfig["auth_url_pattern"].(string); exists && pattern != "" {
				authURLPattern = pattern
			}
		} else {
			logf(ctx, "[auth-retriever] No plugin config found for 'auth-retriever', config keys: %v", getMapKeys(ctx.Config))
		}
	} else {
		logf(ctx, "[auth-retriever] Config is nil")
	}
	logf(ctx, "[auth-retriever] Using auth_url_pattern: %s", authURLPattern)

	// Check if this is an auth request (URL contains the pattern)
	if !strings.Contains(strings.ToLower(ctx.Request.URL), strings.ToLower(authURLPattern)) {
		logf(ctx, "[auth-retriever] Skipping token capture for %s %s (no match for pattern '%s')", ctx.Request.Method, ctx.Request.URL, authURLPattern)
		return nil
	}

	// Get token path from config, default to "token"
	tokenPath := "token"
	if pluginConfig, ok := ctx.Config["auth-retriever"].(map[string]interface{}); ok {
		if path, exists := pluginConfig["token_path"].(string); exists && path != "" {
			tokenPath = path
		}
	}

	// Check if response is successful
	if ctx.Response != nil {
		statusCode := ctx.Response.StatusCode
		body := ctx.Response.Body

		if statusCode == 200 && body != "" {
			logf(ctx, "[auth-retriever] Processing auth response from %s %s (status %d), body length: %d", ctx.Request.Method, ctx.Request.URL, statusCode, len(body))
			// Try to parse JSON response
			var data map[string]interface{}
			if err := json.Unmarshal([]byte(body), &data); err == nil {
				logf(ctx, "[auth-retriever] Parsed JSON response, keys: %v", getMapKeys(data))
				if tokenVal := getNestedValue(data, tokenPath); tokenVal != nil {
					if token, ok := tokenVal.(string); ok && token != "" {
						logf(ctx, "[auth-retriever] Captured auth token from %s %s", ctx.Request.Method, ctx.Request.URL)
						// Store captured token in environment
						if ctx.Environment == nil {
							ctx.Environment = make(map[string]string)
						}
						ctx.Environment["auth_token"] = token
						logf(ctx, "[auth-retriever] Stored token in environment: %s", token)
					} else {
						logf(ctx, "[auth-retriever] Token at path '%s' is not a string or empty, type: %T, value: %v", tokenPath, tokenVal, tokenVal)
					}
				} else {
					logf(ctx, "[auth-retriever] No token found at path '%s' in response, available keys: %v", tokenPath, getMapKeys(data))
				}
			} else {
				logf(ctx, "[auth-retriever] Failed to unmarshal JSON response: %v", err)
				logf(ctx, "[auth-retriever] Response body preview: %.200s", body)
			}
		} else {
			logf(ctx, "[auth-retriever] Response status not 200 or body empty: status=%d, bodyLen=%d", statusCode, len(body))
		}
	} else {
		logf(ctx, "[auth-retriever] Response is nil")
	}

	return nil
}

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.Handshake,
		Plugins: map[string]plugin.Plugin{
			"auth-retriever": &shared.PetitoriumPlugin{Impl: &AuthRetriever{}},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
