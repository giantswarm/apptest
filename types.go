package apptest

type index struct {
	Entries map[string][]entry `json:"entries"`
}

type entry struct {
	AppVersion string   `json:"appVersion"`
	Created    string   `json:"created"`
	Name       string   `json:"name"`
	Urls       []string `json:"urls"`
	Version    string   `json:"version"`
}
