package statpoll

type ReqParam string

const (
	YtParamNowPlaying    ReqParam = "yt_now_playing"
	YtParamTrackDuration ReqParam = "yt_track_duration"
	YtParamElapsedTime   ReqParam = "yt_elapsed_time"
	YtParamPaused        ReqParam = "yt_paused"
	YtParamTrackIndex    ReqParam = "yt_track_index"
	YtParamQueuedTracks  ReqParam = "yt_queued_tracks"
)
