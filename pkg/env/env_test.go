package env

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	t.Run("should load .env file if path is provided", func(t *testing.T) {
		err := os.WriteFile(".env.test", []byte("TEST_KEY=TEST_VALUE"), 0o644)
		assert.NoError(t, err)
		defer func() { _ = os.Remove(".env.test") }()

		LoadConfig(".env.test")

		value := os.Getenv("TEST_KEY")
		assert.Equal(t, "TEST_VALUE", value)
	})

	t.Run("should use system environment variables if .env file is not found", func(t *testing.T) {
		_ = os.Setenv("SYSTEM_KEY", "SYSTEM_VALUE")
		defer func() { _ = os.Unsetenv("SYSTEM_KEY") }()

		LoadConfig("nonexistent.env")

		value := os.Getenv("SYSTEM_KEY")
		assert.Equal(t, "SYSTEM_VALUE", value)
	})
}

func TestGetEnv(t *testing.T) {
	t.Run("should return the environment variable value if it exists", func(t *testing.T) {
		_ = os.Setenv("EXISTING_KEY", "EXISTING_VALUE")
		defer func() { _ = os.Unsetenv("EXISTING_KEY") }()

		value := GetEnv("EXISTING_KEY", "DEFAULT_VALUE")
		assert.Equal(t, "EXISTING_VALUE", value)
	})

	t.Run(
		"should return the fallback value if the environment variable does not exist",
		func(t *testing.T) {
			value := GetEnv("NON_EXISTING_KEY", "DEFAULT_VALUE")
			assert.Equal(t, "DEFAULT_VALUE", value)
		},
	)
}
