package rest

// ResDetails contains the client facing messages to an action
type ResDetails struct {
	Status  string   `json:"status"`
	Message string   `json:"message"`
	Errors  []string `json:"errors"`
}
