package generator_test

import (
	"fmt"
	"github.com/csueiras/reinforcer/internal/generator"
	"github.com/csueiras/reinforcer/internal/loader"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/packages/packagestest"
	"testing"
)

type input struct {
	interfaceName string
	code          string
}

func TestGenerator_Generate(t *testing.T) {
	tests := []struct {
		name                  string
		ignoreNoReturnMethods bool
		inputs                map[string]input
		outCode               *generator.Generated
		wantErr               bool
	}{
		{
			name:                  "Using aliased import",
			ignoreNoReturnMethods: false,
			inputs: map[string]input{
				"my_service.go": {
					interfaceName: "Service",
					code: `package fake

import goctx "context"

type Service interface {
	A(ctx goctx.Context) error
}
`,
				},
			},
			outCode: &generator.Generated{
				Common: `// Code generated by reinforcer, DO NOT EDIT.

package resilient

import (
	"context"
	goresilience "github.com/slok/goresilience"
)

type base struct {
	errorPredicate func(string, error) bool
	runnerFactory  runnerFactory
}
type runnerFactory interface {
	GetRunner(name string) goresilience.Runner
}

var RetryAllErrors = func(_ string, _ error) bool {
	return true
}

type Option func(*base)

func WithRetryableErrorPredicate(fn func(string, error) bool) Option {
	return func(o *base) {
		o.errorPredicate = fn
	}
}
func (b *base) run(ctx context.Context, name string, fn func(ctx context.Context) error) error {
	return b.runnerFactory.GetRunner(name).Run(ctx, fn)
}
`,
				Constants: `// Code generated by reinforcer, DO NOT EDIT.

package resilient

// GeneratedServiceMethods are the methods in GeneratedService
var GeneratedServiceMethods = struct {
	A string
}{
	A: "A",
}
`,
				Files: []*generator.GeneratedFile{
					{
						TypeName: "GeneratedService",
						Contents: `// Code generated by reinforcer, DO NOT EDIT.

package resilient

import "context"

type targetService interface {
	A(ctx context.Context) error
}
type GeneratedService struct {
	*base
	delegate targetService
}

func NewService(delegate targetService, runnerFactory runnerFactory, options ...Option) *GeneratedService {
	if delegate == nil {
		panic("provided nil delegate")
	}
	if runnerFactory == nil {
		panic("provided nil runner factory")
	}
	c := &GeneratedService{
		base: &base{
			errorPredicate: RetryAllErrors,
			runnerFactory:  runnerFactory,
		},
		delegate: delegate,
	}
	for _, o := range options {
		o(c.base)
	}
	return c
}
func (g *GeneratedService) A(ctx context.Context) error {
	var nonRetryableErr error
	err := g.run(ctx, GeneratedServiceMethods.A, func(ctx context.Context) error {
		var err error
		err = g.delegate.A(ctx)
		if g.errorPredicate(GeneratedServiceMethods.A, err) {
			return err
		}
		nonRetryableErr = err
		return nil
	})
	if nonRetryableErr != nil {
		return nonRetryableErr
	}
	return err
}
`,
					},
				},
			},
		},
		{
			name:                  "Complex",
			ignoreNoReturnMethods: false,
			inputs: map[string]input{
				"users_service.go": {
					interfaceName: "Service",
					code: `package fake

import "context"

type User struct {
	Name string
}

type Service interface {
	A()
	B(ctx context.Context)
	C(ctx context.Context, param1 int, param2 *int32, param3 *User)
	GetUserID(ctx context.Context, userID string) (string, error)
	GetUserID2(ctx context.Context, userID *string) (*User, error)
	HasVariadic(ctx context.Context, fields ...string) error
}`,
				},
			},
			outCode: &generator.Generated{
				Common: `// Code generated by reinforcer, DO NOT EDIT.

package resilient

import (
	"context"
	goresilience "github.com/slok/goresilience"
)

type base struct {
	errorPredicate func(string, error) bool
	runnerFactory  runnerFactory
}
type runnerFactory interface {
	GetRunner(name string) goresilience.Runner
}

var RetryAllErrors = func(_ string, _ error) bool {
	return true
}

type Option func(*base)

func WithRetryableErrorPredicate(fn func(string, error) bool) Option {
	return func(o *base) {
		o.errorPredicate = fn
	}
}
func (b *base) run(ctx context.Context, name string, fn func(ctx context.Context) error) error {
	return b.runnerFactory.GetRunner(name).Run(ctx, fn)
}
`,
				Constants: `// Code generated by reinforcer, DO NOT EDIT.

package resilient

// GeneratedServiceMethods are the methods in GeneratedService
var GeneratedServiceMethods = struct {
	A           string
	B           string
	C           string
	GetUserID   string
	GetUserID2  string
	HasVariadic string
}{
	A:           "A",
	B:           "B",
	C:           "C",
	GetUserID:   "GetUserID",
	GetUserID2:  "GetUserID2",
	HasVariadic: "HasVariadic",
}
`,
				Files: []*generator.GeneratedFile{
					{
						TypeName: "GeneratedService",
						Contents: `// Code generated by reinforcer, DO NOT EDIT.

package resilient

import (
	"context"
	unresilient "github.com/csueiras/fake/unresilient"
)

type targetService interface {
	A()
	B(ctx context.Context)
	C(ctx context.Context, arg1 int, arg2 *int32, arg3 *unresilient.User)
	GetUserID(ctx context.Context, arg1 string) (string, error)
	GetUserID2(ctx context.Context, arg1 *string) (*unresilient.User, error)
	HasVariadic(ctx context.Context, arg1 ...string) error
}
type GeneratedService struct {
	*base
	delegate targetService
}

func NewService(delegate targetService, runnerFactory runnerFactory, options ...Option) *GeneratedService {
	if delegate == nil {
		panic("provided nil delegate")
	}
	if runnerFactory == nil {
		panic("provided nil runner factory")
	}
	c := &GeneratedService{
		base: &base{
			errorPredicate: RetryAllErrors,
			runnerFactory:  runnerFactory,
		},
		delegate: delegate,
	}
	for _, o := range options {
		o(c.base)
	}
	return c
}
func (g *GeneratedService) A() {
	err := g.run(context.Background(), GeneratedServiceMethods.A, func(_ context.Context) error {
		g.delegate.A()
		return nil
	})
	if err != nil {
		panic(err)
	}
}
func (g *GeneratedService) B(ctx context.Context) {
	err := g.run(ctx, GeneratedServiceMethods.B, func(ctx context.Context) error {
		g.delegate.B(ctx)
		return nil
	})
	if err != nil {
		panic(err)
	}
}
func (g *GeneratedService) C(ctx context.Context, arg1 int, arg2 *int32, arg3 *unresilient.User) {
	err := g.run(ctx, GeneratedServiceMethods.C, func(ctx context.Context) error {
		g.delegate.C(ctx, arg1, arg2, arg3)
		return nil
	})
	if err != nil {
		panic(err)
	}
}
func (g *GeneratedService) GetUserID(ctx context.Context, arg1 string) (string, error) {
	var nonRetryableErr error
	var r0 string
	err := g.run(ctx, GeneratedServiceMethods.GetUserID, func(ctx context.Context) error {
		var err error
		r0, err = g.delegate.GetUserID(ctx, arg1)
		if g.errorPredicate(GeneratedServiceMethods.GetUserID, err) {
			return err
		}
		nonRetryableErr = err
		return nil
	})
	if nonRetryableErr != nil {
		return r0, nonRetryableErr
	}
	return r0, err
}
func (g *GeneratedService) GetUserID2(ctx context.Context, arg1 *string) (*unresilient.User, error) {
	var nonRetryableErr error
	var r0 *unresilient.User
	err := g.run(ctx, GeneratedServiceMethods.GetUserID2, func(ctx context.Context) error {
		var err error
		r0, err = g.delegate.GetUserID2(ctx, arg1)
		if g.errorPredicate(GeneratedServiceMethods.GetUserID2, err) {
			return err
		}
		nonRetryableErr = err
		return nil
	})
	if nonRetryableErr != nil {
		return r0, nonRetryableErr
	}
	return r0, err
}
func (g *GeneratedService) HasVariadic(ctx context.Context, arg1 ...string) error {
	var nonRetryableErr error
	err := g.run(ctx, GeneratedServiceMethods.HasVariadic, func(ctx context.Context) error {
		var err error
		err = g.delegate.HasVariadic(ctx, arg1...)
		if g.errorPredicate(GeneratedServiceMethods.HasVariadic, err) {
			return err
		}
		nonRetryableErr = err
		return nil
	})
	if nonRetryableErr != nil {
		return nonRetryableErr
	}
	return err
}
`,
					},
				},
			},
		},
		{
			name:                  "Ignore No Return Methods",
			ignoreNoReturnMethods: true,
			inputs: map[string]input{
				"users_service.go": {
					interfaceName: "Service",
					code: `package fake

import "context"

type User struct {
	Name string
}

type Service interface {
	A()
	B(ctx context.Context, userID string) (string, error)
}`,
				},
			},
			outCode: &generator.Generated{
				Common: `// Code generated by reinforcer, DO NOT EDIT.

package resilient

import (
	"context"
	goresilience "github.com/slok/goresilience"
)

type base struct {
	errorPredicate func(string, error) bool
	runnerFactory  runnerFactory
}
type runnerFactory interface {
	GetRunner(name string) goresilience.Runner
}

var RetryAllErrors = func(_ string, _ error) bool {
	return true
}

type Option func(*base)

func WithRetryableErrorPredicate(fn func(string, error) bool) Option {
	return func(o *base) {
		o.errorPredicate = fn
	}
}
func (b *base) run(ctx context.Context, name string, fn func(ctx context.Context) error) error {
	return b.runnerFactory.GetRunner(name).Run(ctx, fn)
}
`,
				Constants: `// Code generated by reinforcer, DO NOT EDIT.

package resilient

// GeneratedServiceMethods are the methods in GeneratedService
var GeneratedServiceMethods = struct {
	A string
	B string
}{
	A: "A",
	B: "B",
}
`,
				Files: []*generator.GeneratedFile{
					{
						TypeName: "GeneratedService",
						Contents: `// Code generated by reinforcer, DO NOT EDIT.

package resilient

import "context"

type targetService interface {
	A()
	B(ctx context.Context, arg1 string) (string, error)
}
type GeneratedService struct {
	*base
	delegate targetService
}

func NewService(delegate targetService, runnerFactory runnerFactory, options ...Option) *GeneratedService {
	if delegate == nil {
		panic("provided nil delegate")
	}
	if runnerFactory == nil {
		panic("provided nil runner factory")
	}
	c := &GeneratedService{
		base: &base{
			errorPredicate: RetryAllErrors,
			runnerFactory:  runnerFactory,
		},
		delegate: delegate,
	}
	for _, o := range options {
		o(c.base)
	}
	return c
}
func (g *GeneratedService) A() {
	g.delegate.A()
}
func (g *GeneratedService) B(ctx context.Context, arg1 string) (string, error) {
	var nonRetryableErr error
	var r0 string
	err := g.run(ctx, GeneratedServiceMethods.B, func(ctx context.Context) error {
		var err error
		r0, err = g.delegate.B(ctx, arg1)
		if g.errorPredicate(GeneratedServiceMethods.B, err) {
			return err
		}
		nonRetryableErr = err
		return nil
	})
	if nonRetryableErr != nil {
		return r0, nonRetryableErr
	}
	return r0, err
}
`,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ifaces := loadInterface(t, tt.inputs)
			got, err := generator.Generate(generator.Config{
				OutPkg:                "resilient",
				Files:                 ifaces,
				IgnoreNoReturnMethods: tt.ignoreNoReturnMethods,
			})

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, got)

				require.Equal(t, tt.outCode.Constants, got.Constants)
				require.Equal(t, tt.outCode.Common, got.Common)

				for fileName, genFile := range got.Files {
					require.Equal(t, tt.outCode.Files[fileName].TypeName, genFile.TypeName, "File %s type names don't match", fileName)
					require.Equal(t, tt.outCode.Files[fileName].Contents, genFile.Contents, "File %s contents doesn't match. Got:\n%s", fileName, genFile.Contents)
				}
				require.Equal(t, tt.outCode, got)
			}
		})
	}
}

func loadInterface(t *testing.T, filesCode map[string]input) []*generator.FileConfig {
	pkg := "github.com/csueiras/fake/unresilient"
	m := map[string]interface{}{}
	for fileName, in := range filesCode {
		m[fileName] = in.code
	}

	mods := []packagestest.Module{
		{
			Name:  pkg,
			Files: m,
		},
	}
	exported := packagestest.Export(t, packagestest.GOPATH, mods)
	defer exported.Cleanup()

	l := loader.NewLoader(func(cfg *packages.Config, patterns ...string) ([]*packages.Package, error) {
		exported.Config.Mode = cfg.Mode
		return packages.Load(exported.Config, patterns...)
	})

	var loadedTypes []*generator.FileConfig
	for _, in := range filesCode {
		svc, err := l.LoadOne(pkg, in.interfaceName)
		require.NoError(t, err)
		loadedTypes = append(loadedTypes, &generator.FileConfig{
			SrcTypeName:   in.interfaceName,
			OutTypeName:   fmt.Sprintf("Generated%s", in.interfaceName),
			InterfaceType: svc,
		})
	}
	return loadedTypes
}
