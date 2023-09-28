// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots.

package config

import (
	"encoding/json"
	"reflect"

	"fmt"

	"github.com/spf13/pflag"
)

// If v is a pointer, it will get its element value or the zero value of the element type.
// If v is not a pointer, it will return it as is.
func (ServerConfig) elemValueOrNil(v interface{}) interface{} {
	if t := reflect.TypeOf(v); t.Kind() == reflect.Ptr {
		if reflect.ValueOf(v).IsNil() {
			return reflect.Zero(t.Elem()).Interface()
		} else {
			return reflect.ValueOf(v).Interface()
		}
	} else if v == nil {
		return reflect.Zero(t).Interface()
	}

	return v
}

func (ServerConfig) mustJsonMarshal(v interface{}) string {
	raw, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	return string(raw)
}

func (ServerConfig) mustMarshalJSON(v json.Marshaler) string {
	raw, err := v.MarshalJSON()
	if err != nil {
		panic(err)
	}

	return string(raw)
}

// GetPFlagSet will return strongly types pflags for all fields in ServerConfig and its nested types. The format of the
// flags is json-name.json-sub-name... etc.
func (cfg ServerConfig) GetPFlagSet(prefix string) *pflag.FlagSet {
	cmdFlags := pflag.NewFlagSet("ServerConfig", pflag.ExitOnError)
	cmdFlags.Int(fmt.Sprintf("%v%v", prefix, "httpPort"), defaultServerConfig.HTTPPort, "On which http port to serve admin")
	cmdFlags.Int(fmt.Sprintf("%v%v", prefix, "grpcPort"), defaultServerConfig.GrpcPort, "deprecated")
	cmdFlags.Bool(fmt.Sprintf("%v%v", prefix, "grpcServerReflection"), defaultServerConfig.GrpcServerReflection, "deprecated")
	cmdFlags.String(fmt.Sprintf("%v%v", prefix, "kube-config"), defaultServerConfig.KubeConfig, "Path to kubernetes client config file,  default is empty,  useful for incluster config.")
	cmdFlags.String(fmt.Sprintf("%v%v", prefix, "master"), defaultServerConfig.Master, "The address of the Kubernetes API server.")
	cmdFlags.Bool(fmt.Sprintf("%v%v", prefix, "security.secure"), defaultServerConfig.Security.Secure, "")
	cmdFlags.String(fmt.Sprintf("%v%v", prefix, "security.ssl.certificateFile"), defaultServerConfig.Security.Ssl.CertificateFile, "")
	cmdFlags.String(fmt.Sprintf("%v%v", prefix, "security.ssl.keyFile"), defaultServerConfig.Security.Ssl.KeyFile, "")
	cmdFlags.Bool(fmt.Sprintf("%v%v", prefix, "security.useAuth"), defaultServerConfig.Security.UseAuth, "")
	cmdFlags.Bool(fmt.Sprintf("%v%v", prefix, "security.auditAccess"), defaultServerConfig.Security.AuditAccess, "")
	cmdFlags.Bool(fmt.Sprintf("%v%v", prefix, "security.allowCors"), defaultServerConfig.Security.AllowCors, "")
	cmdFlags.StringSlice(fmt.Sprintf("%v%v", prefix, "security.allowedOrigins"), defaultServerConfig.Security.AllowedOrigins, "")
	cmdFlags.StringSlice(fmt.Sprintf("%v%v", prefix, "security.allowedHeaders"), defaultServerConfig.Security.AllowedHeaders, "")
	cmdFlags.Int(fmt.Sprintf("%v%v", prefix, "grpc.port"), defaultServerConfig.GrpcConfig.Port, "On which grpc port to serve admin")
	cmdFlags.Bool(fmt.Sprintf("%v%v", prefix, "grpc.serverReflection"), defaultServerConfig.GrpcConfig.ServerReflection, "Enable GRPC Server Reflection")
	cmdFlags.Int(fmt.Sprintf("%v%v", prefix, "grpc.maxMessageSizeBytes"), defaultServerConfig.GrpcConfig.MaxMessageSizeBytes, "The max size in bytes for incoming gRPC messages")
	cmdFlags.String(fmt.Sprintf("%v%v", prefix, "thirdPartyConfig.flyteClient.clientId"), defaultServerConfig.DeprecatedThirdPartyConfig.FlyteClientConfig.ClientID, "public identifier for the app which handles authorization for a Flyte deployment")
	cmdFlags.String(fmt.Sprintf("%v%v", prefix, "thirdPartyConfig.flyteClient.redirectUri"), defaultServerConfig.DeprecatedThirdPartyConfig.FlyteClientConfig.RedirectURI, "This is the callback uri registered with the app which handles authorization for a Flyte deployment")
	cmdFlags.StringSlice(fmt.Sprintf("%v%v", prefix, "thirdPartyConfig.flyteClient.scopes"), defaultServerConfig.DeprecatedThirdPartyConfig.FlyteClientConfig.Scopes, "Recommended scopes for the client to request.")
	cmdFlags.String(fmt.Sprintf("%v%v", prefix, "thirdPartyConfig.flyteClient.audience"), defaultServerConfig.DeprecatedThirdPartyConfig.FlyteClientConfig.Audience, "Audience to use when initiating OAuth2 authorization requests.")
	cmdFlags.String(fmt.Sprintf("%v%v", prefix, "dataProxy.upload.maxSize"), defaultServerConfig.DataProxy.Upload.MaxSize.String(), "Maximum allowed upload size.")
	cmdFlags.String(fmt.Sprintf("%v%v", prefix, "dataProxy.upload.maxExpiresIn"), defaultServerConfig.DataProxy.Upload.MaxExpiresIn.String(), "Maximum allowed expiration duration.")
	cmdFlags.Int(fmt.Sprintf("%v%v", prefix, "dataProxy.upload.defaultFileNameLength"), defaultServerConfig.DataProxy.Upload.DefaultFileNameLength, "Default length for the generated file name if not provided in the request.")
	cmdFlags.String(fmt.Sprintf("%v%v", prefix, "dataProxy.upload.storagePrefix"), defaultServerConfig.DataProxy.Upload.StoragePrefix, "Storage prefix to use for all upload requests.")
	cmdFlags.String(fmt.Sprintf("%v%v", prefix, "dataProxy.download.maxExpiresIn"), defaultServerConfig.DataProxy.Download.MaxExpiresIn.String(), "Maximum allowed expiration duration.")
	cmdFlags.Int(fmt.Sprintf("%v%v", prefix, "readHeaderTimeoutSeconds"), defaultServerConfig.ReadHeaderTimeoutSeconds, "The amount of time allowed to read request headers.")
	cmdFlags.Int32(fmt.Sprintf("%v%v", prefix, "kubeClientConfig.qps"), defaultServerConfig.KubeClientConfig.QPS, "Max QPS to the master for requests to KubeAPI. 0 defaults to 5.")
	cmdFlags.Int(fmt.Sprintf("%v%v", prefix, "kubeClientConfig.burst"), defaultServerConfig.KubeClientConfig.Burst, "Max burst rate for throttle. 0 defaults to 10")
	cmdFlags.String(fmt.Sprintf("%v%v", prefix, "kubeClientConfig.timeout"), defaultServerConfig.KubeClientConfig.Timeout.String(), "Max duration allowed for every request to KubeAPI before giving up. 0 implies no timeout.")
	return cmdFlags
}