package constants

var MetricIDs = []string{
	"aggregate_startup_time",
	"downscale_percentage",
	"exits_before_video_start",
	"live_stream_latency",
	"max_downscale_percentage",
	"max_request_latency",
	"max_upscale_percentage",
	"page_load_time",
	"playback_failure_percentage",
	"playback_success_score",
	"player_startup_time",
	"playing_time",
	"rebuffer_count",
	"rebuffer_duration",
	"rebuffer_frequency",
	"rebuffer_percentage",
	"request_latency",
	"request_throughput",
	"rebuffer_score",
	"requests_for_first_preroll",
	"seek_latency",
	"startup_time_score",
	"unique_viewers",
	"upscale_percentage",
	"video_quality_score",
	"video_startup_preroll_load_time",
	"video_startup_preroll_request_time",
	"video_startup_time",
	"viewer_experience_score",
	"views",
	"weighted_average_bitrate",
	"video_startup_failure_percentage",
	"ad_attempt_count",
	"ad_break_count",
	"ad_break_error_count",
	"ad_break_error_percentage",
	"ad_error_count",
	"ad_error_percentage",
	"ad_exit_before_start_count",
	"ad_exit_before_start_percentage",
	"ad_impression_count",
	"ad_startup_error_count",
	"ad_startup_error_percentage",
	"playback_business_exception_percentage",
	"video_startup_business_exception_percentage",
	"view_content_startup_time",
	"ad_preroll_startup_time",
	"ad_watch_time",
	"view_content_watch_time",
}

var Measurements = []string{
	"95th",
	"median",
	"avg",
	"count",
	"sum",
}
