// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"github.com/rs/cors"
	"gopkg.in/yaml.v2"
	"i18n-manager"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"i18n-manager/restapi/operations"
	"i18n-manager/restapi/operations/operation"
	"i18n-manager/restapi/operations/query"
)

//go:generate swagger generate server --target .. --name I18n --spec ../swagger.yaml

func configureFlags(api *operations.I18nAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.I18nAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.MultipartformConsumer = runtime.DiscardConsumer

	api.JSONProducer = runtime.JSONProducer()

	api.TxtProducer = runtime.TextProducer()

	api.QueryGetQueryLanguageStatusHandler = query.GetQueryLanguageStatusHandlerFunc(func(params query.GetQueryLanguageStatusParams) middleware.Responder {
		i18n := manager.Query(params.Language, params.Status, params.Languages)
		accept := params.HTTPRequest.Header.Get("Accept")
		if strings.Contains(accept, "application/json") {
			return query.NewGetQueryLanguageStatusOK().WithPayload(i18n.ToApiModel())
		} else if bytes, err := yaml.Marshal(i18n); nil != err {
			return query.NewGetQueryLanguageStatusBadRequest().WithPayload(err.Error())
		} else {
			return middleware.ResponderFunc(func(w http.ResponseWriter, p runtime.Producer) {
				w.Header().Set("Content-type", "text/yaml")
				w.WriteHeader(200)
				w.Write(bytes)
			})
		}
	})
	api.OperationSaveOrUpdateKeyHandler = operation.SaveOrUpdateKeyHandlerFunc(func(params operation.SaveOrUpdateKeyParams) middleware.Responder {
		if err := manager.SaveOrUpdate(params.Body); nil != err {
			return operation.NewSaveOrUpdateKeyBadRequest().WithPayload(err.Error())
		} else {
			return operation.NewSaveOrUpdateKeyOK()
		}

	})
	api.OperationUploadHandler = operation.UploadHandlerFunc(func(params operation.UploadParams) middleware.Responder {
		if bytes, err := ioutil.ReadAll(params.File); nil != err {
			return operation.NewUploadBadRequest().WithPayload(err.Error())
		} else {
			i18n := &i18n_manager.I18N{}
			if err := yaml.Unmarshal(bytes, i18n); nil != err {
				return operation.NewUploadBadRequest().WithPayload(err.Error())
			}
			if err := manager.Update(i18n); nil != err {
				return operation.NewUploadBadRequest().WithPayload(err.Error())
			}
		}

		return operation.NewUploadOK()
	})

	api.ServerShutdown = func() {
		if err := store.Shutdown(); nil != err {
			api.Logger(err.Error())
		}
	}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

var manager *i18n_manager.Manager
var store i18n_manager.Store

func init() {
	var initError error = nil
	if cwd, err := os.Getwd(); nil == err {
		storePath := path.Join(cwd, "data")
		if store, err = i18n_manager.CreateBadgerStore(storePath); nil == err {
			manager, err = i18n_manager.CreateManager(store)
		}

		initError = err
	}

	if nil != initError {
		panic(initError)
	}
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	options := cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"HEAD", "GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}
	return cors.New(options).Handler(handler)
}
