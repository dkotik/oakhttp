package turnstile

type Request struct {
	// Secret key provided by Cloudflare.
	Secret string `json:"secret"`

	// Response key provided by the verification subject.
	Response string `json:"response"`

	// RemoteIP is the subject's IP address. It is optional but improves the verification accuracy.
	RemoteIP string `json:"remoteip"`
}
