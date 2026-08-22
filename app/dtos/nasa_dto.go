package dtos

// ApodDTO — NASA APOD response DTO
type ApodDTO struct {
	Title       string `json:"title"`
	Explanation string `json:"explanation"`
	URL         string `json:"url"`
	HdURL       string `json:"hdUrl"`
	MediaType   string `json:"mediaType"`
	Date        string `json:"date"`
}
