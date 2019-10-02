package middleware

import (
	"context"

	"github.com/justinas/alice"
	"tracfox.io/tracfox/pkg/tracfox/middleware/logger"
	"tracfox.io/tracfox/pkg/tracfox/middleware/recovery"
	"tracfox.io/tracfox/pkg/tracfox/middleware/version"
)

// NewDefaultChain .
func NewDefaultChain(ctx context.Context) (alice.Chain, error) {
	chain := alice.New()
	recoveryMiddleware, err := recovery.New(ctx)
	if err != nil {
		return chain, err
	}
	chain = chain.Append(recoveryMiddleware)
	loggerMiddleware, err := logger.New(ctx)
	if err != nil {
		return chain, err
	}
	chain = chain.Append(loggerMiddleware)
	versionMiddleware, err := version.New(ctx)
	if err != nil {
		return chain, err
	}
	return chain.Append(versionMiddleware), nil
}
