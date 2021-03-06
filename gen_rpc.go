package main

import (
	"go/types"
	"strings"
)

type RpcServiceMeta struct {
	// i.e. projectv1
	RpcPackage string
	// i.e. AccountServiceServer
	Name string
	// i.e. AccountServicePathPrefix
	EndpointPath string
	// i.e. AccountServiceRpc
	EntityServiceRpc string
	// i.e. newAccountUseCase
	NewEntityUseCase string
}

const rpcTemplate = `
// Code generated by generator, DO NOT EDIT.
package wire

{{- if gt (len .RpcServices) 0 }}
import (
	"context"
	"errors"
	"strings"
	"{{.Project}}/domain"
	v4 "github.com/labstack/echo/v4"
	twirp "github.com/twitchtv/twirp"
	rpc "{{.Project}}/internal/api/gen/{{.RpcPackage}}"
	handler "{{.Project}}/internal/handler"
)
{{- end }}

func ServeRpcServices() {
{{- range .RpcServices }}
	serve{{.EntityServiceRpc}}()
{{- end }}
}

{{ range .RpcServices }}
func serve{{.EntityServiceRpc}}() {

	prefix := strings.Replace(rpc.{{.EndpointPath}}, "/twirp", ServicePathPrefix, 1)

	handler.E.Group(prefix).POST("*", func(c v4.Context) error {
		rpc.{{.Name}}(
			rpc.{{.EntityServiceRpc}}{UseCase: {{.NewEntityUseCase}}(handler.GetAppContext(c))},
			twirp.WithServerPathPrefix(ServicePathPrefix),
			twirp.WithServerHooks(newEchoHook(c)),
		).ServeHTTP(c.Response().Writer, c.Request())

		return nil
	})
}
{{ end }}

// Twirp hook to log errors to stdout in the service
func newEchoHook(c v4.Context) *twirp.ServerHooks {
    return &twirp.ServerHooks{
        Error: func(ctx context.Context, twerr twirp.Error) context.Context {
			originalErr := errors.Unwrap(twerr)
			var internal string
			if domainErr, ok := originalErr.(domain.Error); ok {
				internal = domainErr.Internal
			}
			handler.WriteErrorLog(c, originalErr, internal)
            return ctx
        },
    }
}
`

const (
	genFileName             = "gen_rpc_router.go"
	echoQual                = "github.com/labstack/echo/v4"
	twirpQual               = "github.com/twitchtv/twirp"
	serviceServerSuffix     = "ServiceServer"
	newPrefix               = "New"
	servicePathPrefixSuffix = "ServicePathPrefix"
	serviceRpcPrefix        = "ServiceRpc"
)

func genRpc() {
	// Inspect package and use type checker to infer imported types
	pkg := loadPackage(meta.PackagePath)
	paths := strings.Split(meta.PackagePath, "/")
	rpcPkg := paths[len(paths)-1]

	if len(rpcPkg) == 0 {
		rpcPkg = pkg.Types.Name()
	}

	meta.RpcPackage = rpcPkg
	// Lookup the given source type name in the package declarations
	for _, rpcServiceName := range pkg.Types.Scope().Names() {
		// Assert suffix
		if !(strings.HasPrefix(rpcServiceName, newPrefix) && strings.HasSuffix(rpcServiceName, serviceServerSuffix)) {
			continue
		}

		obj := pkg.Types.Scope().Lookup(rpcServiceName)

		// Assert a declared function type
		if _, ok := obj.(*types.Func); !ok {
			continue
		}

		// Assert exportable type
		if !obj.Exported() {
			continue
		}

		// Build Rpc instance and binding
		rpcMeta := RpcServiceMeta{
			Name:       rpcServiceName,
			RpcPackage: rpcPkg,
		}

		// i.e. `NewAccountServiceServer` becomes `Account`
		var entity = strings.Replace(
			strings.Replace(rpcServiceName, newPrefix, "", 1),
			serviceServerSuffix, "", 1,
		)

		// i.e. AccountServicePathPrefix
		rpcMeta.EndpointPath = entity + servicePathPrefixSuffix
		// i.e. AccountServiceRpc
		rpcMeta.EntityServiceRpc = entity + serviceRpcPrefix
		// i.e. newAccountUseCase
		rpcMeta.NewEntityUseCase = "new" + entity + "UseCase"

		meta.RpcServices = append(meta.RpcServices, rpcMeta)
	}
}
