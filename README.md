# Go client for Janus WebRTC server API
Visit official [website](https://janus.conf.meetecho.com/) to get more info about `Janus`.

## Restrictions
- Quite small API subset is implemented right now.
- Only [Admin API](https://janus.conf.meetecho.com/docs/admin.html) is supported.
- Only synchronous API supported.
- No new janus multistream feature support.

## Supported transports
- [x] HTTP
- [ ] WebSockets
- [x] RabbitMQ
- [ ] MQTT
- [ ] Nanomsg
- [ ] UnixSockets

## Admin API
See [related page](https://janus.conf.meetecho.com/docs/admin.html) on Janus website.

### Generic requests
- [x] info
- [x] ping

### Session-related requests
- [ ] accept_new_sessions
- [x] list_sessions
- [ ] set_session_timeout
- [ ] destroy_session

### Handle and WebRTC-related requests
- [ ] list_handles
- [ ] handle_info
- [ ] start_pcap
- [ ] stop_pcap
- [ ] start_text2pcap
- [ ] stop_text2pcap
- [x] message_plugin
- [ ] hangup_webrtc
- [ ] detach_handle


## Plugin APIs

### Videoroom
See [related page](https://janus.conf.meetecho.com/docs/videoroom.html) on Janus website.
- [ ] create
- [ ] destroy
- [ ] edit
- [ ] exists
- [x] list
- [ ] allowed
- [ ] kick
- [ ] moderate
- [ ] enable_recording
- [ ] listparticipants
- [ ] list_forwarders
- [x] rtp_forward

### Streaming
See [related page](https://janus.conf.meetecho.com/docs/streaming.html) on Janus website.
- [x] list
- [ ] info
- [x] create
- [x] destroy
- [ ] recording
- [ ] edit
- [ ] enable
- [ ] disable
