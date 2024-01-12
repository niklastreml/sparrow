// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package checks

import (
	"context"
	"github.com/caas-team/sparrow/pkg/api"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/prometheus/client_golang/prometheus"
	"sync"
)

// Ensure, that CheckMock does implement Check.
// If this is not the case, regenerate this file with moq.
var _ Check = &CheckMock{}

// CheckMock is a mock implementation of Check.
//
//	func TestSomethingThatUsesCheck(t *testing.T) {
//
//		// make and configure a mocked Check
//		mockedCheck := &CheckMock{
//			DeregisterHandlerFunc: func(ctx context.Context, router *api.RoutingTree)  {
//				panic("mock out the DeregisterHandler method")
//			},
//			GetMetricCollectorsFunc: func() []prometheus.Collector {
//				panic("mock out the GetMetricCollectors method")
//			},
//			RegisterHandlerFunc: func(ctx context.Context, router *api.RoutingTree)  {
//				panic("mock out the RegisterHandler method")
//			},
//			RunFunc: func(ctx context.Context) error {
//				panic("mock out the Run method")
//			},
//			SchemaFunc: func() (*openapi3.SchemaRef, error) {
//				panic("mock out the Schema method")
//			},
//			SetConfigFunc: func(ctx context.Context, config any) error {
//				panic("mock out the SetConfig method")
//			},
//			ShutdownFunc: func(ctx context.Context) error {
//				panic("mock out the Shutdown method")
//			},
//			StartupFunc: func(ctx context.Context, cResult chan<- Result) error {
//				panic("mock out the Startup method")
//			},
//		}
//
//		// use mockedCheck in code that requires Check
//		// and then make assertions.
//
//	}
type CheckMock struct {
	// DeregisterHandlerFunc mocks the DeregisterHandler method.
	DeregisterHandlerFunc func(ctx context.Context, router *api.RoutingTree)

	// GetMetricCollectorsFunc mocks the GetMetricCollectors method.
	GetMetricCollectorsFunc func() []prometheus.Collector

	// RegisterHandlerFunc mocks the RegisterHandler method.
	RegisterHandlerFunc func(ctx context.Context, router *api.RoutingTree)

	// RunFunc mocks the Run method.
	RunFunc func(ctx context.Context) error

	// SchemaFunc mocks the Schema method.
	SchemaFunc func() (*openapi3.SchemaRef, error)

	// SetConfigFunc mocks the SetConfig method.
	SetConfigFunc func(ctx context.Context, config any) error

	// ShutdownFunc mocks the Shutdown method.
	ShutdownFunc func(ctx context.Context) error

	// StartupFunc mocks the Startup method.
	StartupFunc func(ctx context.Context, cResult chan<- Result) error

	// calls tracks calls to the methods.
	calls struct {
		// DeregisterHandler holds details about calls to the DeregisterHandler method.
		DeregisterHandler []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Router is the router argument value.
			Router *api.RoutingTree
		}
		// GetMetricCollectors holds details about calls to the GetMetricCollectors method.
		GetMetricCollectors []struct {
		}
		// RegisterHandler holds details about calls to the RegisterHandler method.
		RegisterHandler []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Router is the router argument value.
			Router *api.RoutingTree
		}
		// Run holds details about calls to the Run method.
		Run []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
		}
		// Schema holds details about calls to the Schema method.
		Schema []struct {
		}
		// SetConfig holds details about calls to the SetConfig method.
		SetConfig []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Config is the config argument value.
			Config any
		}
		// Shutdown holds details about calls to the Shutdown method.
		Shutdown []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
		}
		// Startup holds details about calls to the Startup method.
		Startup []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// CResult is the cResult argument value.
			CResult chan<- Result
		}
	}
	lockDeregisterHandler   sync.RWMutex
	lockGetMetricCollectors sync.RWMutex
	lockRegisterHandler     sync.RWMutex
	lockRun                 sync.RWMutex
	lockSchema              sync.RWMutex
	lockSetConfig           sync.RWMutex
	lockShutdown            sync.RWMutex
	lockStartup             sync.RWMutex
}

// DeregisterHandler calls DeregisterHandlerFunc.
func (mock *CheckMock) DeregisterHandler(ctx context.Context, router *api.RoutingTree) {
	if mock.DeregisterHandlerFunc == nil {
		panic("CheckMock.DeregisterHandlerFunc: method is nil but Check.DeregisterHandler was just called")
	}
	callInfo := struct {
		Ctx    context.Context
		Router *api.RoutingTree
	}{
		Ctx:    ctx,
		Router: router,
	}
	mock.lockDeregisterHandler.Lock()
	mock.calls.DeregisterHandler = append(mock.calls.DeregisterHandler, callInfo)
	mock.lockDeregisterHandler.Unlock()
	mock.DeregisterHandlerFunc(ctx, router)
}

// DeregisterHandlerCalls gets all the calls that were made to DeregisterHandler.
// Check the length with:
//
//	len(mockedCheck.DeregisterHandlerCalls())
func (mock *CheckMock) DeregisterHandlerCalls() []struct {
	Ctx    context.Context
	Router *api.RoutingTree
} {
	var calls []struct {
		Ctx    context.Context
		Router *api.RoutingTree
	}
	mock.lockDeregisterHandler.RLock()
	calls = mock.calls.DeregisterHandler
	mock.lockDeregisterHandler.RUnlock()
	return calls
}

// GetMetricCollectors calls GetMetricCollectorsFunc.
func (mock *CheckMock) GetMetricCollectors() []prometheus.Collector {
	if mock.GetMetricCollectorsFunc == nil {
		panic("CheckMock.GetMetricCollectorsFunc: method is nil but Check.GetMetricCollectors was just called")
	}
	callInfo := struct {
	}{}
	mock.lockGetMetricCollectors.Lock()
	mock.calls.GetMetricCollectors = append(mock.calls.GetMetricCollectors, callInfo)
	mock.lockGetMetricCollectors.Unlock()
	return mock.GetMetricCollectorsFunc()
}

// GetMetricCollectorsCalls gets all the calls that were made to GetMetricCollectors.
// Check the length with:
//
//	len(mockedCheck.GetMetricCollectorsCalls())
func (mock *CheckMock) GetMetricCollectorsCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockGetMetricCollectors.RLock()
	calls = mock.calls.GetMetricCollectors
	mock.lockGetMetricCollectors.RUnlock()
	return calls
}

// RegisterHandler calls RegisterHandlerFunc.
func (mock *CheckMock) RegisterHandler(ctx context.Context, router *api.RoutingTree) {
	if mock.RegisterHandlerFunc == nil {
		panic("CheckMock.RegisterHandlerFunc: method is nil but Check.RegisterHandler was just called")
	}
	callInfo := struct {
		Ctx    context.Context
		Router *api.RoutingTree
	}{
		Ctx:    ctx,
		Router: router,
	}
	mock.lockRegisterHandler.Lock()
	mock.calls.RegisterHandler = append(mock.calls.RegisterHandler, callInfo)
	mock.lockRegisterHandler.Unlock()
	mock.RegisterHandlerFunc(ctx, router)
}

// RegisterHandlerCalls gets all the calls that were made to RegisterHandler.
// Check the length with:
//
//	len(mockedCheck.RegisterHandlerCalls())
func (mock *CheckMock) RegisterHandlerCalls() []struct {
	Ctx    context.Context
	Router *api.RoutingTree
} {
	var calls []struct {
		Ctx    context.Context
		Router *api.RoutingTree
	}
	mock.lockRegisterHandler.RLock()
	calls = mock.calls.RegisterHandler
	mock.lockRegisterHandler.RUnlock()
	return calls
}

// Run calls RunFunc.
func (mock *CheckMock) Run(ctx context.Context) error {
	if mock.RunFunc == nil {
		panic("CheckMock.RunFunc: method is nil but Check.Run was just called")
	}
	callInfo := struct {
		Ctx context.Context
	}{
		Ctx: ctx,
	}
	mock.lockRun.Lock()
	mock.calls.Run = append(mock.calls.Run, callInfo)
	mock.lockRun.Unlock()
	return mock.RunFunc(ctx)
}

// RunCalls gets all the calls that were made to Run.
// Check the length with:
//
//	len(mockedCheck.RunCalls())
func (mock *CheckMock) RunCalls() []struct {
	Ctx context.Context
} {
	var calls []struct {
		Ctx context.Context
	}
	mock.lockRun.RLock()
	calls = mock.calls.Run
	mock.lockRun.RUnlock()
	return calls
}

// Schema calls SchemaFunc.
func (mock *CheckMock) Schema() (*openapi3.SchemaRef, error) {
	if mock.SchemaFunc == nil {
		panic("CheckMock.SchemaFunc: method is nil but Check.Schema was just called")
	}
	callInfo := struct {
	}{}
	mock.lockSchema.Lock()
	mock.calls.Schema = append(mock.calls.Schema, callInfo)
	mock.lockSchema.Unlock()
	return mock.SchemaFunc()
}

// SchemaCalls gets all the calls that were made to Schema.
// Check the length with:
//
//	len(mockedCheck.SchemaCalls())
func (mock *CheckMock) SchemaCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockSchema.RLock()
	calls = mock.calls.Schema
	mock.lockSchema.RUnlock()
	return calls
}

// SetConfig calls SetConfigFunc.
func (mock *CheckMock) SetConfig(ctx context.Context, config any) error {
	if mock.SetConfigFunc == nil {
		panic("CheckMock.SetConfigFunc: method is nil but Check.SetConfig was just called")
	}
	callInfo := struct {
		Ctx    context.Context
		Config any
	}{
		Ctx:    ctx,
		Config: config,
	}
	mock.lockSetConfig.Lock()
	mock.calls.SetConfig = append(mock.calls.SetConfig, callInfo)
	mock.lockSetConfig.Unlock()
	return mock.SetConfigFunc(ctx, config)
}

// SetConfigCalls gets all the calls that were made to SetConfig.
// Check the length with:
//
//	len(mockedCheck.SetConfigCalls())
func (mock *CheckMock) SetConfigCalls() []struct {
	Ctx    context.Context
	Config any
} {
	var calls []struct {
		Ctx    context.Context
		Config any
	}
	mock.lockSetConfig.RLock()
	calls = mock.calls.SetConfig
	mock.lockSetConfig.RUnlock()
	return calls
}

// Shutdown calls ShutdownFunc.
func (mock *CheckMock) Shutdown(ctx context.Context) error {
	if mock.ShutdownFunc == nil {
		panic("CheckMock.ShutdownFunc: method is nil but Check.Shutdown was just called")
	}
	callInfo := struct {
		Ctx context.Context
	}{
		Ctx: ctx,
	}
	mock.lockShutdown.Lock()
	mock.calls.Shutdown = append(mock.calls.Shutdown, callInfo)
	mock.lockShutdown.Unlock()
	return mock.ShutdownFunc(ctx)
}

// ShutdownCalls gets all the calls that were made to Shutdown.
// Check the length with:
//
//	len(mockedCheck.ShutdownCalls())
func (mock *CheckMock) ShutdownCalls() []struct {
	Ctx context.Context
} {
	var calls []struct {
		Ctx context.Context
	}
	mock.lockShutdown.RLock()
	calls = mock.calls.Shutdown
	mock.lockShutdown.RUnlock()
	return calls
}

// Startup calls StartupFunc.
func (mock *CheckMock) Startup(ctx context.Context, cResult chan<- Result) error {
	if mock.StartupFunc == nil {
		panic("CheckMock.StartupFunc: method is nil but Check.Startup was just called")
	}
	callInfo := struct {
		Ctx     context.Context
		CResult chan<- Result
	}{
		Ctx:     ctx,
		CResult: cResult,
	}
	mock.lockStartup.Lock()
	mock.calls.Startup = append(mock.calls.Startup, callInfo)
	mock.lockStartup.Unlock()
	return mock.StartupFunc(ctx, cResult)
}

// StartupCalls gets all the calls that were made to Startup.
// Check the length with:
//
//	len(mockedCheck.StartupCalls())
func (mock *CheckMock) StartupCalls() []struct {
	Ctx     context.Context
	CResult chan<- Result
} {
	var calls []struct {
		Ctx     context.Context
		CResult chan<- Result
	}
	mock.lockStartup.RLock()
	calls = mock.calls.Startup
	mock.lockStartup.RUnlock()
	return calls
}