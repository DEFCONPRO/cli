// Code generated by mocker. DO NOT EDIT.
// github.com/travisjeffery/mocker
// Source: internal/pkg/auth/ccloud_client_factory.go

package mock

import (
	context "context"
	sync "sync"

	github_com_confluentinc_ccloud_sdk_go_v1_public "github.com/confluentinc/ccloud-sdk-go-v1-public"
)

// CCloudClientFactory is a mock of CCloudClientFactory interface
type CCloudClientFactory struct {
	lockAnonHTTPClientFactory sync.Mutex
	AnonHTTPClientFactoryFunc func(baseURL string) *github_com_confluentinc_ccloud_sdk_go_v1_public.Client

	lockJwtHTTPClientFactory sync.Mutex
	JwtHTTPClientFactoryFunc func(ctx context.Context, jwt, baseURL string) *github_com_confluentinc_ccloud_sdk_go_v1_public.Client

	calls struct {
		AnonHTTPClientFactory []struct {
			BaseURL string
		}
		JwtHTTPClientFactory []struct {
			Ctx     context.Context
			Jwt     string
			BaseURL string
		}
	}
}

// AnonHTTPClientFactory mocks base method by wrapping the associated func.
func (m *CCloudClientFactory) AnonHTTPClientFactory(baseURL string) *github_com_confluentinc_ccloud_sdk_go_v1_public.Client {
	m.lockAnonHTTPClientFactory.Lock()
	defer m.lockAnonHTTPClientFactory.Unlock()

	if m.AnonHTTPClientFactoryFunc == nil {
		panic("mocker: CCloudClientFactory.AnonHTTPClientFactoryFunc is nil but CCloudClientFactory.AnonHTTPClientFactory was called.")
	}

	call := struct {
		BaseURL string
	}{
		BaseURL: baseURL,
	}

	m.calls.AnonHTTPClientFactory = append(m.calls.AnonHTTPClientFactory, call)

	return m.AnonHTTPClientFactoryFunc(baseURL)
}

// AnonHTTPClientFactoryCalled returns true if AnonHTTPClientFactory was called at least once.
func (m *CCloudClientFactory) AnonHTTPClientFactoryCalled() bool {
	m.lockAnonHTTPClientFactory.Lock()
	defer m.lockAnonHTTPClientFactory.Unlock()

	return len(m.calls.AnonHTTPClientFactory) > 0
}

// AnonHTTPClientFactoryCalls returns the calls made to AnonHTTPClientFactory.
func (m *CCloudClientFactory) AnonHTTPClientFactoryCalls() []struct {
	BaseURL string
} {
	m.lockAnonHTTPClientFactory.Lock()
	defer m.lockAnonHTTPClientFactory.Unlock()

	return m.calls.AnonHTTPClientFactory
}

// JwtHTTPClientFactory mocks base method by wrapping the associated func.
func (m *CCloudClientFactory) JwtHTTPClientFactory(ctx context.Context, jwt, baseURL string) *github_com_confluentinc_ccloud_sdk_go_v1_public.Client {
	m.lockJwtHTTPClientFactory.Lock()
	defer m.lockJwtHTTPClientFactory.Unlock()

	if m.JwtHTTPClientFactoryFunc == nil {
		panic("mocker: CCloudClientFactory.JwtHTTPClientFactoryFunc is nil but CCloudClientFactory.JwtHTTPClientFactory was called.")
	}

	call := struct {
		Ctx     context.Context
		Jwt     string
		BaseURL string
	}{
		Ctx:     ctx,
		Jwt:     jwt,
		BaseURL: baseURL,
	}

	m.calls.JwtHTTPClientFactory = append(m.calls.JwtHTTPClientFactory, call)

	return m.JwtHTTPClientFactoryFunc(ctx, jwt, baseURL)
}

// JwtHTTPClientFactoryCalled returns true if JwtHTTPClientFactory was called at least once.
func (m *CCloudClientFactory) JwtHTTPClientFactoryCalled() bool {
	m.lockJwtHTTPClientFactory.Lock()
	defer m.lockJwtHTTPClientFactory.Unlock()

	return len(m.calls.JwtHTTPClientFactory) > 0
}

// JwtHTTPClientFactoryCalls returns the calls made to JwtHTTPClientFactory.
func (m *CCloudClientFactory) JwtHTTPClientFactoryCalls() []struct {
	Ctx     context.Context
	Jwt     string
	BaseURL string
} {
	m.lockJwtHTTPClientFactory.Lock()
	defer m.lockJwtHTTPClientFactory.Unlock()

	return m.calls.JwtHTTPClientFactory
}

// Reset resets the calls made to the mocked methods.
func (m *CCloudClientFactory) Reset() {
	m.lockAnonHTTPClientFactory.Lock()
	m.calls.AnonHTTPClientFactory = nil
	m.lockAnonHTTPClientFactory.Unlock()
	m.lockJwtHTTPClientFactory.Lock()
	m.calls.JwtHTTPClientFactory = nil
	m.lockJwtHTTPClientFactory.Unlock()
}
