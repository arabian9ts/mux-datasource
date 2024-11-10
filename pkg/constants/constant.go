package constants

var MetricIDs = []string{
	"aggregate_startup_time",
	"downscale_percentage",
	"exits_before_video_start",
	"live_stream_latency",
	"max_downscale_percentage",
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
	"video_startup_time",
	"viewer_experience_score",
	"views",
	"weighted_average_bitrate",
	"video_startup_failure_percentage",
	"playback_business_exception_percentage",
	"video_startup_business_exception_percentage",
}

var Measurements = []string{
	"95th",
	"median",
	"avg",
	"count",
	"sum",
}
