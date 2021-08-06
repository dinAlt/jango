package videoroom

import (
	"fmt"

	"github.com/dinalt/jango"
)

type Error struct {
	Code        int
	Description string
}

func (e *Error) Error() string {
	return fmt.Sprintf("videoroom: [%d] %s", e.Code, e.Description)
}

type pluginImpl struct{}

func (pluginImpl) Plugin() string {
	return "janus.plugin.videoroom"
}

type asyncFalseImpl struct{}

func (r *asyncFalseImpl) Async() bool {
	return false
}

var _ jango.PluginRequest = (*ListRequest)(nil)

type ListRequest struct {
	Request string `json:"request,omitempty"`
	pluginImpl
	asyncFalseImpl
}

func (r *ListRequest) Build() interface{} {
	r.Request = "list"
	return r
}

type ListResponse struct {
	Videoroom string `json:"videoroom,omitempty"`
	ErrorCode int    `json:"error_code,omitempty"`
	Error     string `json:"error,omitempty"`
	List      []struct {
		Room               interface{} `json:"room,omitempty"`
		Description        string      `json:"description,omitempty"`
		PinRequired        bool        `json:"pin_required,omitempty"`
		MaxPublishers      int         `json:"max_publishers,omitempty"`
		Bitrate            int         `json:"bitrate,omitempty"`
		FirFreq            int         `json:"fir_freq,omitempty"`
		RequirePvtid       bool        `json:"require_pvtid,omitempty"`
		RequireE2ee        bool        `json:"require_e2ee,omitempty"`
		NotifyJoining      bool        `json:"notify_joining,omitempty"`
		Audiocodec         string      `json:"audiocodec,omitempty"`
		Videocodec         string      `json:"videocodec,omitempty"`
		Record             bool        `json:"record,omitempty"`
		LockRecord         bool        `json:"lock_record,omitempty"`
		NumParticipants    int         `json:"num_participants,omitempty"`
		AudiolevelExt      bool        `json:"audiolevel_ext,omitempty"`
		AudiolevelEvent    bool        `json:"audiolevel_event,omitempty"`
		VideoorientExt     bool        `json:"videoorient_ext,omitempty"`
		PlayoutdelayExt    bool        `json:"playoutdelay_ext,omitempty"`
		TransportWideCcExt bool        `json:"transport_wide_cc_ext,omitempty"`
	} `json:"list,omitempty"`
}

func (r *ListResponse) Err() error {
	if r.ErrorCode == 0 {
		return nil
	}
	return &Error{r.ErrorCode, r.Error}
}

var _ jango.PluginRequest = (*RTPForwardRequest)(nil)

type RTPForwardRequest struct {
	Request       string      `json:"request,omitempty"`
	Room          interface{} `json:"room,omitempty"`
	PublisherID   interface{} `json:"publisher_id,omitempty"`
	Host          string      `json:"host,omitempty"`
	HostFamily    string      `json:"host_family,omitempty"`
	AudioPort     int         `json:"audio_port,omitempty"`
	AudioSsrc     int         `json:"audio_ssrc,omitempty"`
	AudioPt       int         `json:"audio_pt,omitempty"`
	AudioRtcpPort int         `json:"audio_rtcp_port,omitempty"`
	VideoPort     int         `json:"video_port,omitempty"`
	VideoSsrc     int         `json:"video_ssrc,omitempty"`
	VideoPt       int         `json:"video_pt,omitempty"`
	VideoRtcpPort int         `json:"video_rtcp_port,omitempty"`
	Simulcast     bool        `json:"simulcast,omitempty"`
	VideoPport2   int         `json:"video_pport_2,omitempty"`
	VideoSsrc2    int         `json:"video_ssrc_2,omitempty"`
	VideoPt2      int         `json:"video_pt_2,omitempty"`
	VideoPort3    int         `json:"video_port_3,omitempty"`
	VideoSsrc3    int         `json:"video_ssrc_3,omitempty"`
	VideoPt3      int         `json:"video_pt_3,omitempty"`
	DataPort      int         `json:"data_port,omitempty"`
	SrtpSuite     int         `json:"srtp_suite,omitempty"`
	SrtpCrypto    string      `json:"srtp_crypto,omitempty"`
	AdminKey      string      `json:"admin_key,omitempty"`
	Secret        string      `json:"secret,omitempty"`
	pluginImpl
	asyncFalseImpl
}

func (r *RTPForwardRequest) Build() interface{} {
	r.Request = "rtp_forward"
	return r
}

type RTPForwardResponse struct {
	Videoroom   string      `json:"videoroom,omitempty"`
	ErrorCode   int         `json:"error_code,omitempty"`
	Error       string      `json:"error,omitempty"`
	Room        interface{} `json:"room,omitempty"`
	PublisherID interface{} `json:"publisher_id,omitempty"`
	RtpStream   struct {    // nolint:stylecheck
		Host           string `json:"host,omitempty"`
		Audio          int    `json:"audio,omitempty"`
		AudioRtcp      int    `json:"audio_rtcp,omitempty"`
		AudioStreamID  int    `json:"audio_stream_id,omitempty"`
		Video          int    `json:"video,omitempty"`
		VideoRtcp      int    `json:"video_rtcp,omitempty"`
		VideoStreamID  int    `json:"video_stream_id,omitempty"`
		Video2         int    `json:"video_2,omitempty"`
		VideoStreamID2 int    `json:"video_stream_id_2,omitempty"`
		Video3         int    `json:"video_3,omitempty"`
		VideoStreamID3 int    `json:"video_stream_id_3,omitempty"`
		Data           int    `json:"data,omitempty"`
		DataStreamID   int    `json:"data_stream_id,omitempty"`
	} `json:"rtp_stream,omitempty"`
}

func (r *RTPForwardResponse) Err() error {
	if r.ErrorCode == 0 {
		return nil
	}
	return &Error{r.ErrorCode, r.Error}
}
