// Code generated by mockery v2.46.3. DO NOT EDIT.

package github

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockGitHubRepositoryAdapter is an autogenerated mock type for the GitHubRepositoryAdapter type
type MockGitHubRepositoryAdapter struct {
	mock.Mock
}

type MockGitHubRepositoryAdapter_Expecter struct {
	mock *mock.Mock
}

func (_m *MockGitHubRepositoryAdapter) EXPECT() *MockGitHubRepositoryAdapter_Expecter {
	return &MockGitHubRepositoryAdapter_Expecter{mock: &_m.Mock}
}

// GetFileContent provides a mock function with given fields: ctx, owner, repo, path
func (_m *MockGitHubRepositoryAdapter) GetFileContent(ctx context.Context, owner string, repo string, path string) ([]byte, error) {
	ret := _m.Called(ctx, owner, repo, path)

	if len(ret) == 0 {
		panic("no return value specified for GetFileContent")
	}

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) ([]byte, error)); ok {
		return rf(ctx, owner, repo, path)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) []byte); ok {
		r0 = rf(ctx, owner, repo, path)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, string) error); ok {
		r1 = rf(ctx, owner, repo, path)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockGitHubRepositoryAdapter_GetFileContent_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetFileContent'
type MockGitHubRepositoryAdapter_GetFileContent_Call struct {
	*mock.Call
}

// GetFileContent is a helper method to define mock.On call
//   - ctx context.Context
//   - owner string
//   - repo string
//   - path string
func (_e *MockGitHubRepositoryAdapter_Expecter) GetFileContent(ctx interface{}, owner interface{}, repo interface{}, path interface{}) *MockGitHubRepositoryAdapter_GetFileContent_Call {
	return &MockGitHubRepositoryAdapter_GetFileContent_Call{Call: _e.mock.On("GetFileContent", ctx, owner, repo, path)}
}

func (_c *MockGitHubRepositoryAdapter_GetFileContent_Call) Run(run func(ctx context.Context, owner string, repo string, path string)) *MockGitHubRepositoryAdapter_GetFileContent_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(string))
	})
	return _c
}

func (_c *MockGitHubRepositoryAdapter_GetFileContent_Call) Return(_a0 []byte, _a1 error) *MockGitHubRepositoryAdapter_GetFileContent_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockGitHubRepositoryAdapter_GetFileContent_Call) RunAndReturn(run func(context.Context, string, string, string) ([]byte, error)) *MockGitHubRepositoryAdapter_GetFileContent_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockGitHubRepositoryAdapter creates a new instance of MockGitHubRepositoryAdapter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockGitHubRepositoryAdapter(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockGitHubRepositoryAdapter {
	mock := &MockGitHubRepositoryAdapter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
