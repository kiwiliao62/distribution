package configuration

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

type localConfiguration struct {
	Version       Version `yaml:"version"`
	Log           *Log    `yaml:"log"`
	Notifications []Notif `yaml:"notifications,omitempty"`
}

type Notif struct {
	Name string `yaml:"name"`
}

var expectedConfig = localConfiguration{
	Version: "0.1",
	Log: &Log{
		Formatter: "json",
	},
	Notifications: []Notif{
		{Name: "foo"},
		{Name: "bar"},
		{Name: "car"},
	},
}

const testConfig = `version: "0.1"
log:
  formatter: "text"
notifications:
  - name: "foo"
  - name: "bar"
  - name: "car"`

func TestParserOverwriteIninitializedPoiner(t *testing.T) {
	config := localConfiguration{}

	t.Setenv("REGISTRY_LOG_FORMATTER", "json")

	p := NewParser("registry", []VersionedParseInfo{
		{
			Version: "0.1",
			ParseAs: reflect.TypeOf(config),
			ConversionFunc: func(c interface{}) (interface{}, error) {
				return c, nil
			},
		},
	})

	err := p.Parse([]byte(testConfig), &config)
	require.NoError(t, err)
	require.Equal(t, expectedConfig, config)
}

const testConfig2 = `version: "0.1"
log:
  formatter: "text"
notifications:
  - name: "val1"
  - name: "val2"
  - name: "car"`

func TestParseOverwriteUnininitializedPoiner(t *testing.T) {
	config := localConfiguration{}

	t.Setenv("REGISTRY_LOG_FORMATTER", "json")

	// override only first two notificationsvalues
	// in the tetConfig: leave the last value unchanged.
	t.Setenv("REGISTRY_NOTIFICATIONS_0_NAME", "foo")
	t.Setenv("REGISTRY_NOTIFICATIONS_1_NAME", "bar")

	p := NewParser("registry", []VersionedParseInfo{
		{
			Version: "0.1",
			ParseAs: reflect.TypeOf(config),
			ConversionFunc: func(c interface{}) (interface{}, error) {
				return c, nil
			},
		},
	})

	err := p.Parse([]byte(testConfig2), &config)
	require.NoError(t, err)
	require.Equal(t, expectedConfig, config)
}
