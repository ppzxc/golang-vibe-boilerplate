// Package pagination provides cursor-based pagination utilities.
// Implements RFC 8288 Link header and Total-Count header per ppzxc RESTful Guidelines.
package pagination

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

const (
	// DefaultPageSize is the default number of items per page.
	DefaultPageSize = 20
	// MaxPageSize is the maximum allowed page size.
	MaxPageSize = 100
	// HeaderTotalCount is the response header for total item count.
	HeaderTotalCount = "Total-Count"
)

// Request holds parsed pagination query parameters.
type Request struct {
	PageToken string
	PageSize  int
}

// ParseRequest parses pagination parameters from the request URL.
// Returns 400 Bad Request error if pageSize is invalid.
func ParseRequest(r *http.Request) (Request, error) {
	q := r.URL.Query()

	pageSize := DefaultPageSize
	if ps := q.Get("pageSize"); ps != "" {
		n, err := strconv.Atoi(ps)
		if err != nil || n < 1 {
			return Request{}, fmt.Errorf("pageSize must be a positive integer")
		}
		if n > MaxPageSize {
			n = MaxPageSize
		}
		pageSize = n
	}

	return Request{
		PageToken: q.Get("pageToken"),
		PageSize:  pageSize,
	}, nil
}

// EncodeToken encodes a cursor value as an opaque page token.
// Clients must not parse or construct tokens.
func EncodeToken(cursor string) string {
	return base64.RawURLEncoding.EncodeToString([]byte(cursor))
}

// DecodeToken decodes a page token back to a cursor value.
func DecodeToken(token string) (string, error) {
	b, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return "", fmt.Errorf("invalid page token: %w", err)
	}
	return string(b), nil
}

// Response holds pagination metadata for a response.
type Response struct {
	TotalCount int64
	NextToken  string // empty if no next page
	PrevToken  string // empty if no previous page
	BaseURL    string // base URL for Link header construction
}

// SetHeaders sets the Total-Count and Link headers on the response.
func (p *Response) SetHeaders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(HeaderTotalCount, strconv.FormatInt(p.TotalCount, 10))

	baseURL := p.baseURL(r)
	links := buildLinkHeader(baseURL, r.URL.Query(), p.NextToken, p.PrevToken)
	if links != "" {
		w.Header().Set("Link", links)
	}
}

func (p *Response) baseURL(r *http.Request) string {
	if p.BaseURL != "" {
		return p.BaseURL
	}
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s%s", scheme, r.Host, r.URL.Path)
}

func buildLinkHeader(baseURL string, q url.Values, nextToken, prevToken string) string {
	var links []string

	if nextToken != "" {
		nextQ := cloneQuery(q)
		nextQ.Set("pageToken", nextToken)
		links = append(links, fmt.Sprintf(`<%s?%s>; rel="next"`, baseURL, nextQ.Encode()))
	}

	if prevToken != "" {
		prevQ := cloneQuery(q)
		prevQ.Set("pageToken", prevToken)
		links = append(links, fmt.Sprintf(`<%s?%s>; rel="prev"`, baseURL, prevQ.Encode()))
	}

	if len(links) == 0 {
		return ""
	}

	result := links[0]
	for _, l := range links[1:] {
		result += ", " + l
	}
	return result
}

func cloneQuery(q url.Values) url.Values {
	clone := make(url.Values, len(q))
	for k, v := range q {
		clone[k] = append([]string(nil), v...)
	}
	return clone
}
