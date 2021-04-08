package streaming

import (
	"fmt"

	"github.com/dinalt/jango"
)

type StreamType string

const (
	STRTP      StreamType = "rtp"
	STLive     StreamType = "live"
	STOndemand StreamType = "ondemand"
	STRTSP     StreamType = "rtsp"
)

type DataType string

const (
	DTString DataType = "string"
	DTBinary DataType = "binary"
)

type Error struct {
	Code        int
	Description string
}

func (e *Error) Error() string {
	return fmt.Sprintf("streaming: [%d] %s", e.Code, e.Description)
}

type pluginImpl struct{}

func (pluginImpl) Plugin() string {
	return "janus.plugin.streaming"
}

type asyncFasleImpl struct{}

func (r *asyncFasleImpl) Async() bool {
	return false
}

var _ jango.PluginRequest = (*ListRequest)(nil)

type ListRequest struct {
	Request string `json:"request,omitempty"`

	pluginImpl
	asyncFasleImpl
}

func (r *ListRequest) Build() interface{} {
	r.Request = "list"
	return r
}

type ListResponse struct {
	Streaming string `json:"streaming,omitempty"`
	ErrorCode int    `json:"error_code,omitempty"`
	Error     string `json:"error,omitempty"`

	List []struct {
		ID          interface{} `json:"id,omitempty"`
		Type        string      `json:"type,omitempty"`
		Description string      `json:"description,omitempty"`
		Metadata    string      `json:"metadata,omitempty"`
		Enabled     bool        `json:"enabled,omitempty"`
		AudioAgeMs  int         `json:"audio_age_ms,omitempty"`
		VideoAgeMs  int         `json:"video_age_ms,omitempty"`
	} `json:"list,omitempty"`
}

func (r *ListResponse) Err() error {
	if r.ErrorCode == 0 {
		return nil
	}
	return &Error{r.ErrorCode, r.Error}
}

var _ jango.PluginRequest = (*CreateRequest)(nil)

type CreateRequest struct {
	Request        string      `json:"request,omitempty"`
	AdminKey       string      `json:"admin_key,omitempty"`
	Type           StreamType  `json:"type,omitempty"`
	ID             interface{} `json:"id,omitempty"`
	Name           string      `json:"name,omitempty"`
	Description    string      `json:"description,omitempty"`
	Metadata       string      `json:"metadata,omitempty"`
	Secret         string      `json:"secret,omitempty"`
	Pin            string      `json:"pin,omitempty"`
	IsPrivate      bool        `json:"is_private,omitempty"`
	Audio          bool        `json:"audio,omitempty"`
	Video          bool        `json:"video,omitempty"`
	Data           bool        `json:"data,omitempty"`
	Permanent      bool        `json:"permanent,omitempty"`
	Filename       string      `json:"filename,omitempty"`
	Audioport      int         `json:"audioport,omitempty"`
	Audiortcpport  int         `json:"audiortcpport,omitempty"`
	Audiomcast     string      `json:"audiomcast,omitempty"`
	Audioiface     string      `json:"audioiface,omitempty"`
	Audiopt        int         `json:"audiopt,omitempty"`
	Audiortpmap    string      `json:"audiortpmap,omitempty"`
	Audiofmtp      string      `json:"audiofmtp,omitempty"`
	Audioskew      bool        `json:"audioskew,omitempty"`
	Videoport      int         `json:"videoport,omitempty"`
	Videortcpport  int         `json:"videortcpport,omitempty"`
	Videomcast     string      `json:"videomcast,omitempty"`
	Videoiface     string      `json:"videoiface,omitempty"`
	Videopt        int         `json:"videopt,omitempty"`
	Videortpmap    string      `json:"videortpmap,omitempty"`
	Videofmtp      string      `json:"videofmtp,omitempty"`
	Videobufferkf  bool        `json:"videobufferkf,omitempty"`
	Videosimulcast bool        `json:"videosimulcast,omitempty"`
	Videoport2     int         `json:"videoport2,omitempty"`
	Videoport3     int         `json:"videoport3,omitempty"`
	Videoskew      bool        `json:"videoskew,omitempty"`
	Videosvc       bool        `json:"videosvc,omitempty"`
	Collision      int         `json:"collision,omitempty"`
	Dataport       int         `json:"dataport,omitempty"`
	Dataiface      string      `json:"dataiface,omitempty"`
	Datatype       DataType    `json:"datatype,omitempty"`
	Databuffermsg  bool        `json:"databuffermsg,omitempty"`
	Threads        int         `json:"threads,omitempty"`
	Srtpsuite      int         `json:"srtpsuite,omitempty"`
	Srtpcrypto     string      `json:"srtpcrypto,omitempty"`

	E2ee bool `json:"e2ee,omitempty"`

	Url           string `json:"url,omitempty"`
	RtspUser      string `json:"rtsp_user,omitempty"`
	RtspPwd       string `json:"rtsp_pwd,omitempty"`
	RtspFailcheck bool   `json:"rtsp_failcheck,omitempty"`
	Rtspiface     string `json:"rtspiface,omitempty"`

	pluginImpl
	asyncFasleImpl
}

func (r *CreateRequest) Build() interface{} {
	r.Request = "create"
	return r
}

type CreateResponse struct {
	Streaming string `json:"streaming,omitempty"`
	ErrorCode int    `json:"error_code,omitempty"`
	Error     string `json:"error,omitempty"`

	ID             interface{} `json:"id,omitempty"`
	Name           string      `json:"name,omitempty"`
	Description    string      `json:"description,omitempty"`
	Metadata       string      `json:"metadata,omitempty"`
	Secret         string      `json:"secret,omitempty"`
	Pin            string      `json:"pin,omitempty"`
	IsPrivate      bool        `json:"is_private,omitempty"`
	Audio          bool        `json:"audio,omitempty"`
	Video          bool        `json:"video,omitempty"`
	Data           bool        `json:"data,omitempty"`
	Permanent      bool        `json:"permanent,omitempty"`
	Filename       string      `json:"filename,omitempty"`
	Audioport      int         `json:"audioport,omitempty"`
	Audiortcpport  int         `json:"audiortcpport,omitempty"`
	Audiomcast     string      `json:"audiomcast,omitempty"`
	Audioiface     string      `json:"audioiface,omitempty"`
	Audiopt        int         `json:"audiopt,omitempty"`
	Audiortpmap    string      `json:"audiortpmap,omitempty"`
	Audiofmtp      string      `json:"audiofmtp,omitempty"`
	Audioskew      bool        `json:"audioskew,omitempty"`
	Videoport      int         `json:"videoport,omitempty"`
	Videortcpport  int         `json:"videortcpport,omitempty"`
	Videomcast     string      `json:"videomcast,omitempty"`
	Videoiface     string      `json:"videoiface,omitempty"`
	Videopt        int         `json:"videopt,omitempty"`
	Videortpmap    string      `json:"videortpmap,omitempty"`
	Videofmtp      string      `json:"videofmtp,omitempty"`
	Videobufferkf  bool        `json:"videobufferkf,omitempty"`
	Videosimulcast bool        `json:"videosimulcast,omitempty"`
	Videoport2     int         `json:"videoport2,omitempty"`
	Videoport3     int         `json:"videoport3,omitempty"`
	Videoskew      bool        `json:"videoskew,omitempty"`
	Videosvc       bool        `json:"videosvc,omitempty"`
	Collision      int         `json:"collision,omitempty"`
	Dataport       int         `json:"dataport,omitempty"`
	Dataiface      string      `json:"dataiface,omitempty"`
	Datatype       DataType    `json:"datatype,omitempty"`
	Databuffermsg  bool        `json:"databuffermsg,omitempty"`
	Threads        int         `json:"threads,omitempty"`
	Srtpsuite      int         `json:"srtpsuite,omitempty"`
	Srtpcrypto     string      `json:"srtpcrypto,omitempty"`

	E2ee bool `json:"e2ee,omitempty"`

	Url           string `json:"url,omitempty"`
	RtspUser      string `json:"rtsp_user,omitempty"`
	RtspPwd       string `json:"rtsp_pwd,omitempty"`
	RtspFailcheck bool   `json:"rtsp_failcheck,omitempty"`
	Rtspiface     string `json:"rtspiface,omitempty"`
}

func (r *CreateResponse) Err() error {
	if r.ErrorCode == 0 {
		return nil
	}
	return &Error{r.ErrorCode, r.Error}
}

var _ jango.PluginRequest = (*DestroyRequest)(nil)

type DestroyRequest struct {
	Request   string      `json:"request,omitempty"`
	ID        interface{} `json:"id,omitempty"`
	Secret    string      `json:"secret,omitempty"`
	Permanent string      `json:"permanent,omitempty"`

	pluginImpl
	asyncFasleImpl
}

func (r *DestroyRequest) Build() interface{} {
	r.Request = "destroy"
	return r
}

type DestroyResponse struct {
	Streaming string `json:"streaming,omitempty"`
	ErrorCode int    `json:"error_code,omitempty"`
	Error     string `json:"error,omitempty"`

	ID interface{} `json:"id,omitempty"`
}

func (r *DestroyResponse) Err() error {
	if r.ErrorCode == 0 {
		return nil
	}
	return &Error{r.ErrorCode, r.Error}
}
