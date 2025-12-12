package envget

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"time"
)

// GetValueFromFileOrEnv reads <ENV>_FILE or <ENV> and converts it to T,
// falling back to defaultValue on missing or invalid values.
//
// Priority:
//  1. <ENV>_FILE - Path to file containing the value
//  2. <ENV> - Direct environment variable
//  3. defaultValue - Fallback value if neither exist or parsing fails
//
// For file reading, whitespace is trimmed from both the path and content.
func GetValueFromFileOrEnv[T any](envVariable string, defaultValue T) T {
	if pathf := strings.TrimSpace(os.Getenv(envVariable + "_FILE")); pathf != "" {
		if bf, errf := os.ReadFile(pathf); errf == nil {
			return convertEnvValue(string(bf), defaultValue)
		}
	}

	if path := strings.TrimSpace(os.Getenv(envVariable)); path != "" {
		if b, err := os.ReadFile(path); err == nil {
			return convertEnvValue(string(b), defaultValue)
		}
	}	
	if raw, ok := os.LookupEnv(envVariable); ok {
		return convertEnvValue(raw, defaultValue)
	}

	return defaultValue
}

// convertEnvValue converts a raw string value to type T based on the
// type of defaultValue. It handles primitive types natively and falls
// back to JSON unmarshaling for complex types.
func convertEnvValue[T any](raw string, defaultValue T) T {
	raw = strings.TrimSpace(raw)

	switch any(defaultValue).(type) {
	case string:
		return any(raw).(T)

	case int:
		if v, err := strconv.Atoi(raw); err == nil {
			return any(v).(T)
		}

	case int64:
		if v, err := strconv.ParseInt(raw, 10, 64); err == nil {
			return any(v).(T)
		}

	case float64:
		if v, err := strconv.ParseFloat(raw, 64); err == nil {
			return any(v).(T)
		}

	case bool:
		if v, err := strconv.ParseBool(raw); err == nil {
			return any(v).(T)
		}

	case time.Duration:
		if v, err := time.ParseDuration(raw); err == nil {
			return any(v).(T)
		}

	default:
		var v T
		if err := json.Unmarshal([]byte(raw), &v); err == nil {
			return v
		}
	}

	return defaultValue
}
