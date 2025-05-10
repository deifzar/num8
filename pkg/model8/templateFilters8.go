package model8

type TemplateFilters8 struct {
	Severity             string   `json:"severity" binding:"lowercase,oneof=info low medium high critical,omitempty"`                                              // binding:"required,lowercase,oneof=info low medium high critical,omitempty"`                                                         // filter by severities (accepts CSV values of info, low, medium, high, critical)
	ExcludeSeverities    string   `json:"excludeseverity" binding:"lowercase,oneof=info low medium high critical,omitempty"`                                       // binding:"required,lowercase,oneof=info low medium high critical,omitempty"`                                           // filter by excluding severities (accepts CSV values of info, low, medium, high, critical)
	ProtocolTypes        string   `json:"protocoltype" binding:"lowercase,oneof=dns file http headless workflow websocket whois code javascript,omitempty"`        // binding:"required,lowercase,oneof=dns file http headless workflow websocket whois code javascript,omitempty"`               // filter by protocol types
	ExcludeProtocolTypes string   `json:"excludeprotocoltype" binding:"lowercase,oneof=dns file http headless workflow websocket whois code javascript,omitempty"` // binding:"required,lowercase,oneof=dns file http headless workflow websocket whois code javascript,omitempty"` // filter by excluding protocol types
	Authors              []string `json:"authors"`                                                                                                                 // fiter by author
	Tags                 []string `json:"tags"`                                                                                                                    // filter by tags present in template
	ExcludeTags          []string `json:"excludetags"`                                                                                                             // filter by excluding tags present in template
	IncludeTags          []string `json:"includetags"`                                                                                                             // filter by including tags present in template
	IDs                  []string `json:"ids"`                                                                                                                     // filter by template IDs
	ExcludeIDs           []string `json:"excludeids"`                                                                                                              // filter by excluding template IDs
	TemplateCondition    []string `json:"templatecondition"`                                                                                                       // DSL condition/ expression
}
