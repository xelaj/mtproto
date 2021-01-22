package calls

// CallConfig is UNOFFICIAL description of returned json for phone.getCallConfig method.
// this code, honestly, stolen from telegramd repository of nebula project https://git.io/JtYG0
// TODO: find official description of these parameters and add some more docs
type CallConfig struct {
	AudioFrameSize         int     `json:"audio_frame_size"`
	JitterMinDelay20       int     `json:"jitter_min_delay_20"`
	JitterMinDelay40       int     `json:"jitter_min_delay_40"`
	JitterMinDelay60       int     `json:"jitter_min_delay_60"`
	JitterMaxDelay20       int     `json:"jitter_max_delay_20"`
	JitterMaxDelay40       int     `json:"jitter_max_delay_40"`
	JitterMaxDelay60       int     `json:"jitter_max_delay_60"`
	JitterMaxSlots20       int     `json:"jitter_max_slots_20"`
	JitterMaxSlots40       int     `json:"jitter_max_slots_40"`
	JitterMaxSlots60       int     `json:"jitter_max_slots_60"`
	JitterLossesToReset    int     `json:"jitter_losses_to_reset"`
	JitterResyncThreshold  float32 `json:"jitter_resync_threshold"`
	AudioCongestionWindow  int     `json:"audio_congestion_window"`
	AudioMaxBitrate        int     `json:"audio_max_bitrate"`
	AudioMaxBitrateEdge    int     `json:"audio_max_bitrate_edge"`
	AudioMaxBitrateGprs    int     `json:"audio_max_bitrate_gprs"`
	AudioMaxBitrateSaving  int     `json:"audio_max_bitrate_saving"`
	AudioInitBitrate       int     `json:"audio_init_bitrate"`
	AudioInitBitrateEdge   int     `json:"audio_init_bitrate_edge"`
	AudioInitBitrateGrps   int     `json:"audio_init_bitrate_gprs"`
	AudioInitBitrateSaving int     `json:"audio_init_bitrate_saving"`
	AudioBitrateStepIncr   int     `json:"audio_bitrate_step_incr"`
	AudioBitrateStepDecr   int     `json:"audio_bitrate_step_decr"`
	UseSystemNs            bool    `json:"use_system_ns"`
	UseSystemAec           bool    `json:"us audioInitBitrateGrps inte_system_aec"`

	// next parameters also gotten from this method, only god knows what does they mean
	EnableVP9Decoder         bool    `json:"enable_vp9_decoder"`
	EnableH265Encoder        bool    `json:"enable_h265_encoder"`
	AdspGoodImpls            string  `json:"adsp_good_impls"`
	AudioMediumFecBitrate    int     `json:"audio_medium_fec_bitrate"`
	EnableH264Encoder        bool    `json:"enable_h264_encoder"`
	AudioMediumFecMultiplier float64 `json:"audio_medium_fec_multiplier"`
	Audio_strongFecBitrate   int     `json:"audio_strong_fec_bitrate"`
	EnableVP8Encoder         bool    `json:"enable_vp8_encoder"`
	ForceTCP                 bool    `json:"force_tcp"`
	JitterInitialDelay60     int     `json:"jitter_initial_delay_60"`
	UseIosVpioAGC            bool    `json:"use_ios_vpio_agc"`
	UseSystemAEC             bool    `json:"use_system_aec"`
	EnableVP9Encoder         bool    `json:"enable_vp9_encoder"`
	EnableH265Decoder        bool    `json:"enable_h265_decoder"`
	EnableVP8Decoder         bool    `json:"enable_vp8_decoder"`
	BadCallRating            bool    `json:"bad_call_rating"`
	EnableH264Decoder        bool    `json:"enable_h264_decoder"`
	UseTCP                   bool    `json:"use_tcp"`
}
