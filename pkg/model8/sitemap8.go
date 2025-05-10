package model8

type Sitemap8 struct {
	Sitemap []*SitemapResource8 `json:"sitemap,omitempty"`
}

func (s *Sitemap8) AddResource8(r8 *SitemapResource8) []*SitemapResource8 {
	s.Sitemap = append(s.Sitemap, r8)
	return s.Sitemap
}

func (s *Sitemap8) SetHash() {
	for _, r8 := range s.Sitemap {
		r8.SetHash()
	}
}

func (s *Sitemap8) Uniq() {
	seen := make(map[string]bool)
	result := []*SitemapResource8{}
	for _, r8 := range s.Sitemap {
		val := string(r8.Hash)
		if _, ok := seen[val]; !ok {
			seen[val] = true
			result = append(result, r8)
		}
	}
	s.Sitemap = result
}
