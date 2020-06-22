// Code generated by mocker. DO NOT EDIT.
// github.com/travisjeffery/mocker
// Source: confluent_home.go

package mock

import (
	sync "sync"
)

// MockConfluentHome is a mock of ConfluentHome interface
type MockConfluentHome struct {
	lockFindFile sync.Mutex
	FindFileFunc func(pattern string) ([]string, error)

	lockGetConfig sync.Mutex
	GetConfigFunc func(service string) ([]byte, error)

	lockGetConnectPluginPath sync.Mutex
	GetConnectPluginPathFunc func() (string, error)

	lockGetConnectorConfigFile sync.Mutex
	GetConnectorConfigFileFunc func(connector string) (string, error)

	lockGetScriptFile sync.Mutex
	GetScriptFileFunc func(service string) (string, error)

	lockGetKafkaScriptFile sync.Mutex
	GetKafkaScriptFileFunc func(mode, format string) (string, error)

	lockGetACLCLIFile sync.Mutex
	GetACLCLIFileFunc func() (string, error)

	lockGetVersion sync.Mutex
	GetVersionFunc func(service string) (string, error)

	lockIsConfluentPlatform sync.Mutex
	IsConfluentPlatformFunc func() (bool, error)

	calls struct {
		FindFile []struct {
			Pattern string
		}
		GetConfig []struct {
			Service string
		}
		GetConnectPluginPath []struct {
		}
		GetConnectorConfigFile []struct {
			Connector string
		}
		GetScriptFile []struct {
			Service string
		}
		GetKafkaScriptFile []struct {
			Mode   string
			Format string
		}
		GetACLCLIFile []struct {
		}
		GetVersion []struct {
			Service string
		}
		IsConfluentPlatform []struct {
		}
	}
}

// FindFile mocks base method by wrapping the associated func.
func (m *MockConfluentHome) FindFile(pattern string) ([]string, error) {
	m.lockFindFile.Lock()
	defer m.lockFindFile.Unlock()

	if m.FindFileFunc == nil {
		panic("mocker: MockConfluentHome.FindFileFunc is nil but MockConfluentHome.FindFile was called.")
	}

	call := struct {
		Pattern string
	}{
		Pattern: pattern,
	}

	m.calls.FindFile = append(m.calls.FindFile, call)

	return m.FindFileFunc(pattern)
}

// FindFileCalled returns true if FindFile was called at least once.
func (m *MockConfluentHome) FindFileCalled() bool {
	m.lockFindFile.Lock()
	defer m.lockFindFile.Unlock()

	return len(m.calls.FindFile) > 0
}

// FindFileCalls returns the calls made to FindFile.
func (m *MockConfluentHome) FindFileCalls() []struct {
	Pattern string
} {
	m.lockFindFile.Lock()
	defer m.lockFindFile.Unlock()

	return m.calls.FindFile
}

// GetConfig mocks base method by wrapping the associated func.
func (m *MockConfluentHome) GetConfig(service string) ([]byte, error) {
	m.lockGetConfig.Lock()
	defer m.lockGetConfig.Unlock()

	if m.GetConfigFunc == nil {
		panic("mocker: MockConfluentHome.GetConfigFunc is nil but MockConfluentHome.GetConfig was called.")
	}

	call := struct {
		Service string
	}{
		Service: service,
	}

	m.calls.GetConfig = append(m.calls.GetConfig, call)

	return m.GetConfigFunc(service)
}

// GetConfigCalled returns true if GetConfig was called at least once.
func (m *MockConfluentHome) GetConfigCalled() bool {
	m.lockGetConfig.Lock()
	defer m.lockGetConfig.Unlock()

	return len(m.calls.GetConfig) > 0
}

// GetConfigCalls returns the calls made to GetConfig.
func (m *MockConfluentHome) GetConfigCalls() []struct {
	Service string
} {
	m.lockGetConfig.Lock()
	defer m.lockGetConfig.Unlock()

	return m.calls.GetConfig
}

// GetConnectPluginPath mocks base method by wrapping the associated func.
func (m *MockConfluentHome) GetConnectPluginPath() (string, error) {
	m.lockGetConnectPluginPath.Lock()
	defer m.lockGetConnectPluginPath.Unlock()

	if m.GetConnectPluginPathFunc == nil {
		panic("mocker: MockConfluentHome.GetConnectPluginPathFunc is nil but MockConfluentHome.GetConnectPluginPath was called.")
	}

	call := struct {
	}{}

	m.calls.GetConnectPluginPath = append(m.calls.GetConnectPluginPath, call)

	return m.GetConnectPluginPathFunc()
}

// GetConnectPluginPathCalled returns true if GetConnectPluginPath was called at least once.
func (m *MockConfluentHome) GetConnectPluginPathCalled() bool {
	m.lockGetConnectPluginPath.Lock()
	defer m.lockGetConnectPluginPath.Unlock()

	return len(m.calls.GetConnectPluginPath) > 0
}

// GetConnectPluginPathCalls returns the calls made to GetConnectPluginPath.
func (m *MockConfluentHome) GetConnectPluginPathCalls() []struct {
} {
	m.lockGetConnectPluginPath.Lock()
	defer m.lockGetConnectPluginPath.Unlock()

	return m.calls.GetConnectPluginPath
}

// GetConnectorConfigFile mocks base method by wrapping the associated func.
func (m *MockConfluentHome) GetConnectorConfigFile(connector string) (string, error) {
	m.lockGetConnectorConfigFile.Lock()
	defer m.lockGetConnectorConfigFile.Unlock()

	if m.GetConnectorConfigFileFunc == nil {
		panic("mocker: MockConfluentHome.GetConnectorConfigFileFunc is nil but MockConfluentHome.GetConnectorConfigFile was called.")
	}

	call := struct {
		Connector string
	}{
		Connector: connector,
	}

	m.calls.GetConnectorConfigFile = append(m.calls.GetConnectorConfigFile, call)

	return m.GetConnectorConfigFileFunc(connector)
}

// GetConnectorConfigFileCalled returns true if GetConnectorConfigFile was called at least once.
func (m *MockConfluentHome) GetConnectorConfigFileCalled() bool {
	m.lockGetConnectorConfigFile.Lock()
	defer m.lockGetConnectorConfigFile.Unlock()

	return len(m.calls.GetConnectorConfigFile) > 0
}

// GetConnectorConfigFileCalls returns the calls made to GetConnectorConfigFile.
func (m *MockConfluentHome) GetConnectorConfigFileCalls() []struct {
	Connector string
} {
	m.lockGetConnectorConfigFile.Lock()
	defer m.lockGetConnectorConfigFile.Unlock()

	return m.calls.GetConnectorConfigFile
}

// GetScriptFile mocks base method by wrapping the associated func.
func (m *MockConfluentHome) GetScriptFile(service string) (string, error) {
	m.lockGetScriptFile.Lock()
	defer m.lockGetScriptFile.Unlock()

	if m.GetScriptFileFunc == nil {
		panic("mocker: MockConfluentHome.GetScriptFileFunc is nil but MockConfluentHome.GetScriptFile was called.")
	}

	call := struct {
		Service string
	}{
		Service: service,
	}

	m.calls.GetScriptFile = append(m.calls.GetScriptFile, call)

	return m.GetScriptFileFunc(service)
}

// GetScriptFileCalled returns true if GetScriptFile was called at least once.
func (m *MockConfluentHome) GetScriptFileCalled() bool {
	m.lockGetScriptFile.Lock()
	defer m.lockGetScriptFile.Unlock()

	return len(m.calls.GetScriptFile) > 0
}

// GetScriptFileCalls returns the calls made to GetScriptFile.
func (m *MockConfluentHome) GetScriptFileCalls() []struct {
	Service string
} {
	m.lockGetScriptFile.Lock()
	defer m.lockGetScriptFile.Unlock()

	return m.calls.GetScriptFile
}

// GetKafkaScriptFile mocks base method by wrapping the associated func.
func (m *MockConfluentHome) GetKafkaScriptFile(mode, format string) (string, error) {
	m.lockGetKafkaScriptFile.Lock()
	defer m.lockGetKafkaScriptFile.Unlock()

	if m.GetKafkaScriptFileFunc == nil {
		panic("mocker: MockConfluentHome.GetKafkaScriptFileFunc is nil but MockConfluentHome.GetKafkaScriptFile was called.")
	}

	call := struct {
		Mode   string
		Format string
	}{
		Mode:   mode,
		Format: format,
	}

	m.calls.GetKafkaScriptFile = append(m.calls.GetKafkaScriptFile, call)

	return m.GetKafkaScriptFileFunc(mode, format)
}

// GetKafkaScriptFileCalled returns true if GetKafkaScriptFile was called at least once.
func (m *MockConfluentHome) GetKafkaScriptFileCalled() bool {
	m.lockGetKafkaScriptFile.Lock()
	defer m.lockGetKafkaScriptFile.Unlock()

	return len(m.calls.GetKafkaScriptFile) > 0
}

// GetKafkaScriptFileCalls returns the calls made to GetKafkaScriptFile.
func (m *MockConfluentHome) GetKafkaScriptFileCalls() []struct {
	Mode   string
	Format string
} {
	m.lockGetKafkaScriptFile.Lock()
	defer m.lockGetKafkaScriptFile.Unlock()

	return m.calls.GetKafkaScriptFile
}

// GetACLCLIFile mocks base method by wrapping the associated func.
func (m *MockConfluentHome) GetACLCLIFile() (string, error) {
	m.lockGetACLCLIFile.Lock()
	defer m.lockGetACLCLIFile.Unlock()

	if m.GetACLCLIFileFunc == nil {
		panic("mocker: MockConfluentHome.GetACLCLIFileFunc is nil but MockConfluentHome.GetACLCLIFile was called.")
	}

	call := struct {
	}{}

	m.calls.GetACLCLIFile = append(m.calls.GetACLCLIFile, call)

	return m.GetACLCLIFileFunc()
}

// GetACLCLIFileCalled returns true if GetACLCLIFile was called at least once.
func (m *MockConfluentHome) GetACLCLIFileCalled() bool {
	m.lockGetACLCLIFile.Lock()
	defer m.lockGetACLCLIFile.Unlock()

	return len(m.calls.GetACLCLIFile) > 0
}

// GetACLCLIFileCalls returns the calls made to GetACLCLIFile.
func (m *MockConfluentHome) GetACLCLIFileCalls() []struct {
} {
	m.lockGetACLCLIFile.Lock()
	defer m.lockGetACLCLIFile.Unlock()

	return m.calls.GetACLCLIFile
}

// GetVersion mocks base method by wrapping the associated func.
func (m *MockConfluentHome) GetVersion(service string) (string, error) {
	m.lockGetVersion.Lock()
	defer m.lockGetVersion.Unlock()

	if m.GetVersionFunc == nil {
		panic("mocker: MockConfluentHome.GetVersionFunc is nil but MockConfluentHome.GetVersion was called.")
	}

	call := struct {
		Service string
	}{
		Service: service,
	}

	m.calls.GetVersion = append(m.calls.GetVersion, call)

	return m.GetVersionFunc(service)
}

// GetVersionCalled returns true if GetVersion was called at least once.
func (m *MockConfluentHome) GetVersionCalled() bool {
	m.lockGetVersion.Lock()
	defer m.lockGetVersion.Unlock()

	return len(m.calls.GetVersion) > 0
}

// GetVersionCalls returns the calls made to GetVersion.
func (m *MockConfluentHome) GetVersionCalls() []struct {
	Service string
} {
	m.lockGetVersion.Lock()
	defer m.lockGetVersion.Unlock()

	return m.calls.GetVersion
}

// IsConfluentPlatform mocks base method by wrapping the associated func.
func (m *MockConfluentHome) IsConfluentPlatform() (bool, error) {
	m.lockIsConfluentPlatform.Lock()
	defer m.lockIsConfluentPlatform.Unlock()

	if m.IsConfluentPlatformFunc == nil {
		panic("mocker: MockConfluentHome.IsConfluentPlatformFunc is nil but MockConfluentHome.IsConfluentPlatform was called.")
	}

	call := struct {
	}{}

	m.calls.IsConfluentPlatform = append(m.calls.IsConfluentPlatform, call)

	return m.IsConfluentPlatformFunc()
}

// IsConfluentPlatformCalled returns true if IsConfluentPlatform was called at least once.
func (m *MockConfluentHome) IsConfluentPlatformCalled() bool {
	m.lockIsConfluentPlatform.Lock()
	defer m.lockIsConfluentPlatform.Unlock()

	return len(m.calls.IsConfluentPlatform) > 0
}

// IsConfluentPlatformCalls returns the calls made to IsConfluentPlatform.
func (m *MockConfluentHome) IsConfluentPlatformCalls() []struct {
} {
	m.lockIsConfluentPlatform.Lock()
	defer m.lockIsConfluentPlatform.Unlock()

	return m.calls.IsConfluentPlatform
}

// Reset resets the calls made to the mocked methods.
func (m *MockConfluentHome) Reset() {
	m.lockFindFile.Lock()
	m.calls.FindFile = nil
	m.lockFindFile.Unlock()
	m.lockGetConfig.Lock()
	m.calls.GetConfig = nil
	m.lockGetConfig.Unlock()
	m.lockGetConnectPluginPath.Lock()
	m.calls.GetConnectPluginPath = nil
	m.lockGetConnectPluginPath.Unlock()
	m.lockGetConnectorConfigFile.Lock()
	m.calls.GetConnectorConfigFile = nil
	m.lockGetConnectorConfigFile.Unlock()
	m.lockGetScriptFile.Lock()
	m.calls.GetScriptFile = nil
	m.lockGetScriptFile.Unlock()
	m.lockGetKafkaScriptFile.Lock()
	m.calls.GetKafkaScriptFile = nil
	m.lockGetKafkaScriptFile.Unlock()
	m.lockGetACLCLIFile.Lock()
	m.calls.GetACLCLIFile = nil
	m.lockGetACLCLIFile.Unlock()
	m.lockGetVersion.Lock()
	m.calls.GetVersion = nil
	m.lockGetVersion.Unlock()
	m.lockIsConfluentPlatform.Lock()
	m.calls.IsConfluentPlatform = nil
	m.lockIsConfluentPlatform.Unlock()
}