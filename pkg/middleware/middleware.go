/*
Copyright 2019 The Koala Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package middleware

import (
	"context"

	"github.com/justinas/alice"
	"github.com/shimcdn/koala/pkg/middleware/logger"
	"github.com/shimcdn/koala/pkg/middleware/recovery"
	"github.com/shimcdn/koala/pkg/middleware/version"
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
