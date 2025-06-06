// Code generated by mockery v2.53.3. DO NOT EDIT.

package mocks

import (
	context "context"

	sqs "github.com/aws/aws-sdk-go-v2/service/sqs"
	mock "github.com/stretchr/testify/mock"
)

// SqsClient is an autogenerated mock type for the SqsClient type
type SqsClient struct {
	mock.Mock
}

// DeleteMessage provides a mock function with given fields: ctx, params, optFns
func (_m *SqsClient) DeleteMessage(ctx context.Context, params *sqs.DeleteMessageInput, optFns ...func(*sqs.Options)) (*sqs.DeleteMessageOutput, error) {
	_va := make([]interface{}, len(optFns))
	for _i := range optFns {
		_va[_i] = optFns[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, params)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for DeleteMessage")
	}

	var r0 *sqs.DeleteMessageOutput
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *sqs.DeleteMessageInput, ...func(*sqs.Options)) (*sqs.DeleteMessageOutput, error)); ok {
		return rf(ctx, params, optFns...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *sqs.DeleteMessageInput, ...func(*sqs.Options)) *sqs.DeleteMessageOutput); ok {
		r0 = rf(ctx, params, optFns...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sqs.DeleteMessageOutput)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *sqs.DeleteMessageInput, ...func(*sqs.Options)) error); ok {
		r1 = rf(ctx, params, optFns...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ReceiveMessage provides a mock function with given fields: ctx, params, optFns
func (_m *SqsClient) ReceiveMessage(ctx context.Context, params *sqs.ReceiveMessageInput, optFns ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error) {
	_va := make([]interface{}, len(optFns))
	for _i := range optFns {
		_va[_i] = optFns[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, params)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for ReceiveMessage")
	}

	var r0 *sqs.ReceiveMessageOutput
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *sqs.ReceiveMessageInput, ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error)); ok {
		return rf(ctx, params, optFns...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *sqs.ReceiveMessageInput, ...func(*sqs.Options)) *sqs.ReceiveMessageOutput); ok {
		r0 = rf(ctx, params, optFns...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sqs.ReceiveMessageOutput)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *sqs.ReceiveMessageInput, ...func(*sqs.Options)) error); ok {
		r1 = rf(ctx, params, optFns...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewSqsClient creates a new instance of SqsClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewSqsClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *SqsClient {
	mock := &SqsClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
