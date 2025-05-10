package model8

type CustomMessagesFinding struct {
	Template    string `json:"template,omitempty"`
	Type        string `json:"type,omitempty"`
	Info        string `json:"info,omitempty"`
	Description string `json:"description,omitempty"`
	Found       string `json:"found,omitempty"`
}

type CustomMessagesSeverity struct {
	Severity     string                  `json:"severity,omitempty"`
	Per_severity []CustomMessagesFinding `json:"per_severity,omitempty"`
}

type CustomMessagesPort struct {
	Port     string                   `json:"port,omitempty"`
	Per_port []CustomMessagesSeverity `json:"per_port,omitempty"`
}

type CustomMessagesHost struct {
	Host     string               `json:"host,omitempty"`
	Per_host []CustomMessagesPort `json:"per_host,omitempty"`
}
