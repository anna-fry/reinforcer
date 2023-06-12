package retryable_test

import (
	"bytes"
	"github.com/anna-fry/reinforcer/internal/generator/method"
	"github.com/anna-fry/reinforcer/internal/generator/retryable"
	rtypes "github.com/anna-fry/reinforcer/internal/types"
	"github.com/stretchr/testify/require"
	"go/token"
	"go/types"
	"testing"
)

func TestRetryable_Statement(t *testing.T) {
	errVar := types.NewVar(token.NoPos, nil, "", rtypes.ErrType)
	ctxVar := types.NewVar(token.NoPos, nil, "ctx", rtypes.ContextType)

	tests := []struct {
		name       string
		methodName string
		signature  *types.Signature
		want       string
		wantErr    bool
	}{
		{
			name:       "Function returns error",
			methodName: "MyFunction",
			signature:  types.NewSignature(nil, types.NewTuple(), types.NewTuple(errVar), false),
			want: `func (r *Resilient) MyFunction() error {
	var nonRetryableErr error
	err := r.run(context.Background(), ResilientMethods.MyFunction, func(_ context.Context) error {
		var err error
		err = r.delegate.MyFunction()
		if r.errorPredicate(ResilientMethods.MyFunction, err) {
			return err
		}
		nonRetryableErr = err
		return nil
	})
	if nonRetryableErr != nil {
		return nonRetryableErr
	}
	return err
}`,
			wantErr: false,
		},
		{
			name:       "Function returns string and error",
			methodName: "MyFunction",
			signature:  types.NewSignature(nil, types.NewTuple(), types.NewTuple(types.NewVar(token.NoPos, nil, "", types.Typ[types.String]), errVar), false),
			want: `func (r *Resilient) MyFunction() (string, error) {
	var nonRetryableErr error
	var r0 string
	err := r.run(context.Background(), ResilientMethods.MyFunction, func(_ context.Context) error {
		var err error
		r0, err = r.delegate.MyFunction()
		if r.errorPredicate(ResilientMethods.MyFunction, err) {
			return err
		}
		nonRetryableErr = err
		return nil
	})
	if nonRetryableErr != nil {
		return r0, nonRetryableErr
	}
	return r0, err
}`,
			wantErr: false,
		},
		{
			name:       "Function passes arguments",
			methodName: "MyFunction",
			signature: types.NewSignature(nil, types.NewTuple(
				ctxVar,
				types.NewVar(token.NoPos, nil, "myArg", types.Typ[types.String]),
			), types.NewTuple(types.NewVar(token.NoPos, nil, "", types.Typ[types.String]), errVar), false),
			want: `func (r *Resilient) MyFunction(ctx context.Context, arg1 string) (string, error) {
	var nonRetryableErr error
	var r0 string
	err := r.run(ctx, ResilientMethods.MyFunction, func(ctx context.Context) error {
		var err error
		r0, err = r.delegate.MyFunction(ctx, arg1)
		if r.errorPredicate(ResilientMethods.MyFunction, err) {
			return err
		}
		nonRetryableErr = err
		return nil
	})
	if nonRetryableErr != nil {
		return r0, nonRetryableErr
	}
	return r0, err
}`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := method.ParseMethod(tt.methodName, tt.signature)
			require.NoError(t, err)
			ret := retryable.NewRetryable(m, "Resilient", "r")
			buf := &bytes.Buffer{}
			s, err := ret.Statement()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NotNil(t, s)
				require.NoError(t, err)
				renderErr := s.Render(buf)
				require.NoError(t, renderErr)

				got := buf.String()
				require.Equal(t, tt.want, got)
			}
		})
	}

	t.Run("Function does not return error", func(t *testing.T) {
		require.Panics(t, func() {
			m, err := method.ParseMethod("Fn", types.NewSignature(nil, types.NewTuple(), types.NewTuple(), false))
			require.NoError(t, err)
			retryable.NewRetryable(m, "Resilient", "r")
		})
	})
}
