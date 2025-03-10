// Code generated by mockery v2.46.3. DO NOT EDIT.

package github

import (
	context "context"

	v69github "github.com/google/go-github/v69/github"
	mock "github.com/stretchr/testify/mock"
)

// MockGitHubIssueAdapter is an autogenerated mock type for the GitHubIssueAdapter type
type MockGitHubIssueAdapter struct {
	mock.Mock
}

type MockGitHubIssueAdapter_Expecter struct {
	mock *mock.Mock
}

func (_m *MockGitHubIssueAdapter) EXPECT() *MockGitHubIssueAdapter_Expecter {
	return &MockGitHubIssueAdapter_Expecter{mock: &_m.Mock}
}

// CreateIssueComment provides a mock function with given fields: ctx, owner, repo, number, comment
func (_m *MockGitHubIssueAdapter) CreateIssueComment(ctx context.Context, owner string, repo string, number int, comment string) (*v69github.IssueComment, error) {
	ret := _m.Called(ctx, owner, repo, number, comment)

	if len(ret) == 0 {
		panic("no return value specified for CreateIssueComment")
	}

	var r0 *v69github.IssueComment
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, int, string) (*v69github.IssueComment, error)); ok {
		return rf(ctx, owner, repo, number, comment)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, int, string) *v69github.IssueComment); ok {
		r0 = rf(ctx, owner, repo, number, comment)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v69github.IssueComment)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, int, string) error); ok {
		r1 = rf(ctx, owner, repo, number, comment)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockGitHubIssueAdapter_CreateIssueComment_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateIssueComment'
type MockGitHubIssueAdapter_CreateIssueComment_Call struct {
	*mock.Call
}

// CreateIssueComment is a helper method to define mock.On call
//   - ctx context.Context
//   - owner string
//   - repo string
//   - number int
//   - comment string
func (_e *MockGitHubIssueAdapter_Expecter) CreateIssueComment(ctx interface{}, owner interface{}, repo interface{}, number interface{}, comment interface{}) *MockGitHubIssueAdapter_CreateIssueComment_Call {
	return &MockGitHubIssueAdapter_CreateIssueComment_Call{Call: _e.mock.On("CreateIssueComment", ctx, owner, repo, number, comment)}
}

func (_c *MockGitHubIssueAdapter_CreateIssueComment_Call) Run(run func(ctx context.Context, owner string, repo string, number int, comment string)) *MockGitHubIssueAdapter_CreateIssueComment_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(int), args[4].(string))
	})
	return _c
}

func (_c *MockGitHubIssueAdapter_CreateIssueComment_Call) Return(_a0 *v69github.IssueComment, _a1 error) *MockGitHubIssueAdapter_CreateIssueComment_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockGitHubIssueAdapter_CreateIssueComment_Call) RunAndReturn(run func(context.Context, string, string, int, string) (*v69github.IssueComment, error)) *MockGitHubIssueAdapter_CreateIssueComment_Call {
	_c.Call.Return(run)
	return _c
}

// ListIssueComments provides a mock function with given fields: ctx, owner, repo, number
func (_m *MockGitHubIssueAdapter) ListIssueComments(ctx context.Context, owner string, repo string, number int) ([]*v69github.IssueComment, error) {
	ret := _m.Called(ctx, owner, repo, number)

	if len(ret) == 0 {
		panic("no return value specified for ListIssueComments")
	}

	var r0 []*v69github.IssueComment
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, int) ([]*v69github.IssueComment, error)); ok {
		return rf(ctx, owner, repo, number)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, int) []*v69github.IssueComment); ok {
		r0 = rf(ctx, owner, repo, number)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*v69github.IssueComment)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, int) error); ok {
		r1 = rf(ctx, owner, repo, number)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockGitHubIssueAdapter_ListIssueComments_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListIssueComments'
type MockGitHubIssueAdapter_ListIssueComments_Call struct {
	*mock.Call
}

// ListIssueComments is a helper method to define mock.On call
//   - ctx context.Context
//   - owner string
//   - repo string
//   - number int
func (_e *MockGitHubIssueAdapter_Expecter) ListIssueComments(ctx interface{}, owner interface{}, repo interface{}, number interface{}) *MockGitHubIssueAdapter_ListIssueComments_Call {
	return &MockGitHubIssueAdapter_ListIssueComments_Call{Call: _e.mock.On("ListIssueComments", ctx, owner, repo, number)}
}

func (_c *MockGitHubIssueAdapter_ListIssueComments_Call) Run(run func(ctx context.Context, owner string, repo string, number int)) *MockGitHubIssueAdapter_ListIssueComments_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(int))
	})
	return _c
}

func (_c *MockGitHubIssueAdapter_ListIssueComments_Call) Return(_a0 []*v69github.IssueComment, _a1 error) *MockGitHubIssueAdapter_ListIssueComments_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockGitHubIssueAdapter_ListIssueComments_Call) RunAndReturn(run func(context.Context, string, string, int) ([]*v69github.IssueComment, error)) *MockGitHubIssueAdapter_ListIssueComments_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateIssueComment provides a mock function with given fields: ctx, owner, repo, commentId, comment
func (_m *MockGitHubIssueAdapter) UpdateIssueComment(ctx context.Context, owner string, repo string, commentId int, comment string) (*v69github.IssueComment, error) {
	ret := _m.Called(ctx, owner, repo, commentId, comment)

	if len(ret) == 0 {
		panic("no return value specified for UpdateIssueComment")
	}

	var r0 *v69github.IssueComment
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, int, string) (*v69github.IssueComment, error)); ok {
		return rf(ctx, owner, repo, commentId, comment)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, int, string) *v69github.IssueComment); ok {
		r0 = rf(ctx, owner, repo, commentId, comment)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v69github.IssueComment)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, int, string) error); ok {
		r1 = rf(ctx, owner, repo, commentId, comment)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockGitHubIssueAdapter_UpdateIssueComment_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateIssueComment'
type MockGitHubIssueAdapter_UpdateIssueComment_Call struct {
	*mock.Call
}

// UpdateIssueComment is a helper method to define mock.On call
//   - ctx context.Context
//   - owner string
//   - repo string
//   - commentId int
//   - comment string
func (_e *MockGitHubIssueAdapter_Expecter) UpdateIssueComment(ctx interface{}, owner interface{}, repo interface{}, commentId interface{}, comment interface{}) *MockGitHubIssueAdapter_UpdateIssueComment_Call {
	return &MockGitHubIssueAdapter_UpdateIssueComment_Call{Call: _e.mock.On("UpdateIssueComment", ctx, owner, repo, commentId, comment)}
}

func (_c *MockGitHubIssueAdapter_UpdateIssueComment_Call) Run(run func(ctx context.Context, owner string, repo string, commentId int, comment string)) *MockGitHubIssueAdapter_UpdateIssueComment_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(int), args[4].(string))
	})
	return _c
}

func (_c *MockGitHubIssueAdapter_UpdateIssueComment_Call) Return(_a0 *v69github.IssueComment, _a1 error) *MockGitHubIssueAdapter_UpdateIssueComment_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockGitHubIssueAdapter_UpdateIssueComment_Call) RunAndReturn(run func(context.Context, string, string, int, string) (*v69github.IssueComment, error)) *MockGitHubIssueAdapter_UpdateIssueComment_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockGitHubIssueAdapter creates a new instance of MockGitHubIssueAdapter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockGitHubIssueAdapter(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockGitHubIssueAdapter {
	mock := &MockGitHubIssueAdapter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
