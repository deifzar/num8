package model8

type OptionsScan8 struct {
	// Target  []string           `json:"target" binding:"required,dive,http_url"`
	T       []string           `json:"t"`                            // -t -templates
	TURL    []string           `json:"turl" binding:"dive,http_url"` // -turl, -template-url -- binding:"required,http_url"` --
	W       []string           `json:"w"`                            // -w -workflows
	WURL    []string           `json:"wurl" binding:"dive,http_url"` // -wurl, -workflow-url -- binding:"required,http_url"` --
	Filters []TemplateFilters8 `json:"filter"`                       // binding:"dive"`
}

type PostOptionsScan8 struct {
	Options OptionsScan8 `json:"options" binding:"required"`
}
