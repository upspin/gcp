// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/logging/type/http_request.proto

/*
Package ltype is a generated protocol buffer package.

It is generated from these files:
	google/logging/type/http_request.proto
	google/logging/type/log_severity.proto

It has these top-level messages:
	HttpRequest
*/
package ltype

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "google.golang.org/genproto/googleapis/api/annotations"
import google_protobuf1 "github.com/golang/protobuf/ptypes/duration"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// A common proto for logging HTTP requests. Only contains semantics
// defined by the HTTP specification. Product-specific logging
// information MUST be defined in a separate message.
type HttpRequest struct {
	// The request method. Examples: `"GET"`, `"HEAD"`, `"PUT"`, `"POST"`.
	RequestMethod string `protobuf:"bytes,1,opt,name=request_method,json=requestMethod" json:"request_method,omitempty"`
	// The scheme (http, https), the host name, the path and the query
	// portion of the URL that was requested.
	// Example: `"http://example.com/some/info?color=red"`.
	RequestUrl string `protobuf:"bytes,2,opt,name=request_url,json=requestUrl" json:"request_url,omitempty"`
	// The size of the HTTP request message in bytes, including the request
	// headers and the request body.
	RequestSize int64 `protobuf:"varint,3,opt,name=request_size,json=requestSize" json:"request_size,omitempty"`
	// The response code indicating the status of response.
	// Examples: 200, 404.
	Status int32 `protobuf:"varint,4,opt,name=status" json:"status,omitempty"`
	// The size of the HTTP response message sent back to the client, in bytes,
	// including the response headers and the response body.
	ResponseSize int64 `protobuf:"varint,5,opt,name=response_size,json=responseSize" json:"response_size,omitempty"`
	// The user agent sent by the client. Example:
	// `"Mozilla/4.0 (compatible; MSIE 6.0; Windows 98; Q312461; .NET CLR 1.0.3705)"`.
	UserAgent string `protobuf:"bytes,6,opt,name=user_agent,json=userAgent" json:"user_agent,omitempty"`
	// The IP address (IPv4 or IPv6) of the client that issued the HTTP
	// request. Examples: `"192.168.1.1"`, `"FE80::0202:B3FF:FE1E:8329"`.
	RemoteIp string `protobuf:"bytes,7,opt,name=remote_ip,json=remoteIp" json:"remote_ip,omitempty"`
	// The IP address (IPv4 or IPv6) of the origin server that the request was
	// sent to.
	ServerIp string `protobuf:"bytes,13,opt,name=server_ip,json=serverIp" json:"server_ip,omitempty"`
	// The referer URL of the request, as defined in
	// [HTTP/1.1 Header Field Definitions](http://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html).
	Referer string `protobuf:"bytes,8,opt,name=referer" json:"referer,omitempty"`
	// The request processing latency on the server, from the time the request was
	// received until the response was sent.
	Latency *google_protobuf1.Duration `protobuf:"bytes,14,opt,name=latency" json:"latency,omitempty"`
	// Whether or not a cache lookup was attempted.
	CacheLookup bool `protobuf:"varint,11,opt,name=cache_lookup,json=cacheLookup" json:"cache_lookup,omitempty"`
	// Whether or not an entity was served from cache
	// (with or without validation).
	CacheHit bool `protobuf:"varint,9,opt,name=cache_hit,json=cacheHit" json:"cache_hit,omitempty"`
	// Whether or not the response was validated with the origin server before
	// being served from cache. This field is only meaningful if `cache_hit` is
	// True.
	CacheValidatedWithOriginServer bool `protobuf:"varint,10,opt,name=cache_validated_with_origin_server,json=cacheValidatedWithOriginServer" json:"cache_validated_with_origin_server,omitempty"`
	// The number of HTTP response bytes inserted into cache. Set only when a
	// cache fill was attempted.
	CacheFillBytes int64 `protobuf:"varint,12,opt,name=cache_fill_bytes,json=cacheFillBytes" json:"cache_fill_bytes,omitempty"`
	// Protocol used for the request. Examples: "HTTP/1.1", "HTTP/2", "websocket"
	Protocol string `protobuf:"bytes,15,opt,name=protocol" json:"protocol,omitempty"`
}

func (m *HttpRequest) Reset()                    { *m = HttpRequest{} }
func (m *HttpRequest) String() string            { return proto.CompactTextString(m) }
func (*HttpRequest) ProtoMessage()               {}
func (*HttpRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *HttpRequest) GetRequestMethod() string {
	if m != nil {
		return m.RequestMethod
	}
	return ""
}

func (m *HttpRequest) GetRequestUrl() string {
	if m != nil {
		return m.RequestUrl
	}
	return ""
}

func (m *HttpRequest) GetRequestSize() int64 {
	if m != nil {
		return m.RequestSize
	}
	return 0
}

func (m *HttpRequest) GetStatus() int32 {
	if m != nil {
		return m.Status
	}
	return 0
}

func (m *HttpRequest) GetResponseSize() int64 {
	if m != nil {
		return m.ResponseSize
	}
	return 0
}

func (m *HttpRequest) GetUserAgent() string {
	if m != nil {
		return m.UserAgent
	}
	return ""
}

func (m *HttpRequest) GetRemoteIp() string {
	if m != nil {
		return m.RemoteIp
	}
	return ""
}

func (m *HttpRequest) GetServerIp() string {
	if m != nil {
		return m.ServerIp
	}
	return ""
}

func (m *HttpRequest) GetReferer() string {
	if m != nil {
		return m.Referer
	}
	return ""
}

func (m *HttpRequest) GetLatency() *google_protobuf1.Duration {
	if m != nil {
		return m.Latency
	}
	return nil
}

func (m *HttpRequest) GetCacheLookup() bool {
	if m != nil {
		return m.CacheLookup
	}
	return false
}

func (m *HttpRequest) GetCacheHit() bool {
	if m != nil {
		return m.CacheHit
	}
	return false
}

func (m *HttpRequest) GetCacheValidatedWithOriginServer() bool {
	if m != nil {
		return m.CacheValidatedWithOriginServer
	}
	return false
}

func (m *HttpRequest) GetCacheFillBytes() int64 {
	if m != nil {
		return m.CacheFillBytes
	}
	return 0
}

func (m *HttpRequest) GetProtocol() string {
	if m != nil {
		return m.Protocol
	}
	return ""
}

func init() {
	proto.RegisterType((*HttpRequest)(nil), "google.logging.type.HttpRequest")
}

func init() { proto.RegisterFile("google/logging/type/http_request.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 499 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x92, 0xcb, 0x6f, 0x13, 0x31,
	0x10, 0xc6, 0xb5, 0x7d, 0x25, 0x71, 0x1e, 0x54, 0x46, 0x02, 0x37, 0x40, 0x09, 0x45, 0xa0, 0x3d,
	0xed, 0x4a, 0xf4, 0x82, 0xc4, 0x89, 0x80, 0xa0, 0x45, 0x45, 0x54, 0x5b, 0x1e, 0x12, 0x97, 0x95,
	0x93, 0x4c, 0xbc, 0x16, 0xce, 0xda, 0xd8, 0xde, 0xa2, 0xf4, 0xca, 0x7f, 0xc3, 0x85, 0x7f, 0x11,
	0xed, 0xd8, 0x2b, 0x81, 0xc4, 0x25, 0xd2, 0x7c, 0xdf, 0xef, 0x9b, 0x71, 0x66, 0x87, 0x3c, 0x15,
	0x5a, 0x0b, 0x05, 0xb9, 0xd2, 0x42, 0xc8, 0x5a, 0xe4, 0x7e, 0x6b, 0x20, 0xaf, 0xbc, 0x37, 0xa5,
	0x85, 0xef, 0x0d, 0x38, 0x9f, 0x19, 0xab, 0xbd, 0xa6, 0xb7, 0x03, 0x97, 0x45, 0x2e, 0x6b, 0xb9,
	0xe9, 0xfd, 0x18, 0xe6, 0x46, 0xe6, 0xbc, 0xae, 0xb5, 0xe7, 0x5e, 0xea, 0xda, 0x85, 0xc8, 0xf4,
	0x38, 0xba, 0x58, 0x2d, 0x9a, 0x75, 0xbe, 0x6a, 0x2c, 0x02, 0xc1, 0x3f, 0xf9, 0xbd, 0x47, 0x86,
	0x67, 0xde, 0x9b, 0x22, 0x0c, 0xa2, 0x4f, 0xc8, 0x24, 0xce, 0x2c, 0x37, 0xe0, 0x2b, 0xbd, 0x62,
	0xc9, 0x2c, 0x49, 0x07, 0xc5, 0x38, 0xaa, 0xef, 0x51, 0xa4, 0x0f, 0xc9, 0xb0, 0xc3, 0x1a, 0xab,
	0xd8, 0x0e, 0x32, 0x24, 0x4a, 0x9f, 0xac, 0xa2, 0x8f, 0xc8, 0xa8, 0x03, 0x9c, 0xbc, 0x01, 0xb6,
	0x3b, 0x4b, 0xd2, 0xdd, 0xa2, 0x0b, 0x5d, 0xc9, 0x1b, 0xa0, 0x77, 0xc8, 0x81, 0xf3, 0xdc, 0x37,
	0x8e, 0xed, 0xcd, 0x92, 0x74, 0xbf, 0x88, 0x15, 0x7d, 0x4c, 0xc6, 0x16, 0x9c, 0xd1, 0xb5, 0x83,
	0x90, 0xdd, 0xc7, 0xec, 0xa8, 0x13, 0x31, 0xfc, 0x80, 0x90, 0xc6, 0x81, 0x2d, 0xb9, 0x80, 0xda,
	0xb3, 0x03, 0x9c, 0x3f, 0x68, 0x95, 0x97, 0xad, 0x40, 0xef, 0x91, 0x81, 0x85, 0x8d, 0xf6, 0x50,
	0x4a, 0xc3, 0x7a, 0xe8, 0xf6, 0x83, 0x70, 0x6e, 0x5a, 0xd3, 0x81, 0xbd, 0x06, 0xdb, 0x9a, 0xe3,
	0x60, 0x06, 0xe1, 0xdc, 0x50, 0x46, 0x7a, 0x16, 0xd6, 0x60, 0xc1, 0xb2, 0x3e, 0x5a, 0x5d, 0x49,
	0x4f, 0x49, 0x4f, 0x71, 0x0f, 0xf5, 0x72, 0xcb, 0x26, 0xb3, 0x24, 0x1d, 0x3e, 0x3b, 0xca, 0xe2,
	0xf7, 0xe8, 0x96, 0x9b, 0xbd, 0x8e, 0xcb, 0x2d, 0x3a, 0xb2, 0xdd, 0xc3, 0x92, 0x2f, 0x2b, 0x28,
	0x95, 0xd6, 0xdf, 0x1a, 0xc3, 0x86, 0xb3, 0x24, 0xed, 0x17, 0x43, 0xd4, 0x2e, 0x50, 0x6a, 0x9f,
	0x13, 0x90, 0x4a, 0x7a, 0x36, 0x40, 0xbf, 0x8f, 0xc2, 0x99, 0xf4, 0xf4, 0x1d, 0x39, 0x09, 0xe6,
	0x35, 0x57, 0x72, 0xc5, 0x3d, 0xac, 0xca, 0x1f, 0xd2, 0x57, 0xa5, 0xb6, 0x52, 0xc8, 0xba, 0x0c,
	0xcf, 0x66, 0x04, 0x53, 0xc7, 0x48, 0x7e, 0xee, 0xc0, 0x2f, 0xd2, 0x57, 0x1f, 0x10, 0xbb, 0x42,
	0x8a, 0xa6, 0xe4, 0x30, 0xf4, 0x5a, 0x4b, 0xa5, 0xca, 0xc5, 0xd6, 0x83, 0x63, 0x23, 0xdc, 0xed,
	0x04, 0xf5, 0x37, 0x52, 0xa9, 0x79, 0xab, 0xd2, 0x29, 0xe9, 0xe3, 0x7f, 0x5a, 0x6a, 0xc5, 0x6e,
	0x85, 0x05, 0x75, 0xf5, 0xfc, 0x67, 0x42, 0xee, 0x2e, 0xf5, 0x26, 0xfb, 0xcf, 0x2d, 0xce, 0x0f,
	0xff, 0x3a, 0xa5, 0xcb, 0x36, 0x70, 0x99, 0x7c, 0x7d, 0x1e, 0x41, 0xa1, 0x15, 0xaf, 0x45, 0xa6,
	0xad, 0xc8, 0x05, 0xd4, 0xd8, 0x2e, 0x0f, 0x16, 0x37, 0xd2, 0xfd, 0x73, 0xfb, 0x2f, 0x54, 0xfb,
	0xfb, 0x6b, 0xe7, 0xe8, 0x6d, 0x88, 0xbe, 0x52, 0xba, 0x59, 0x65, 0x17, 0x71, 0xd2, 0xc7, 0xad,
	0x81, 0xc5, 0x01, 0x36, 0x38, 0xfd, 0x13, 0x00, 0x00, 0xff, 0xff, 0x09, 0x49, 0xe6, 0xb8, 0x3b,
	0x03, 0x00, 0x00,
}
