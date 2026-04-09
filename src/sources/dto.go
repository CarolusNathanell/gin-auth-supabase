package sources

type SourcesAdd struct {
	Name       string `json:"name" binding:"required"`
	Type       string `json:"type" binding:"required,oneof=MP4 RTSP Webcam Youtube Other"`
	Url        string `json:"url" binding:"required,url"`
	FpsTarget  int32  `json:"fps_target" binding:"required"`
	Resolution string `json:"resolution" binding:"required"`
	Status     *bool  `json:"status" binding:"required"`
}

type SourcesUpdate struct {
	Name       string `json:"name" binding:"required"`
	Type       string `json:"type" binding:"required,oneof=MP4 RTSP Webcam Youtube Other"`
	Url        string `json:"url" binding:"required,url"`
	FpsTarget  int32  `json:"fps_target" binding:"required"`
	Resolution string `json:"resolution" binding:"required"`
	Status     *bool  `json:"status" binding:"required"`
}

type SourcesResponse struct {
	Sources any `json:"sources"`
}
