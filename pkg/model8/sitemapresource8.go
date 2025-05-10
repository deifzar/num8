package model8

import (
	"crypto/sha256"
	"net/url"
	"strings"
)

type SitemapResource8 struct {
	Endpoint   string   `json:"endpoint,omitempty"`
	Method     string   `json:"method,omitempty"`
	Parameters []string `json:"parameters,omitempty"`
	RawRequest string   `json:"rawrequest,omitempty"`
	Hash       []byte   `json:"hash,omitempty"`
}

func NewModel8Resource8() Model8SitemapResource8Interface {
	return &SitemapResource8{}
}

func (r8 *SitemapResource8) SetHash() {
	var endpoint string
	p := strings.Join(r8.Parameters[:], ",")
	u, _ := url.Parse(r8.Endpoint)
	endpoint = u.Scheme + "://" + u.Host + u.Path
	h := sha256.New()
	s := endpoint + r8.Method + p
	h.Write([]byte(s))
	hash := h.Sum(nil)
	r8.Hash = hash
}

func (r8 *SitemapResource8) ReturnParammsOneline() string {
	return strings.Join(r8.Parameters[:], ",")
}
