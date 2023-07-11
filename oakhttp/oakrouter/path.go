package oakrouter

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/dkotik/oakacs/oakhttp"
)

type MatchError struct {
	cause error
}

func (m *MatchError) Error() string {
	return "route did not match: " + m.cause.Error()
}

func (m *MatchError) Unwrap() error {
	return m.cause
}

func (m *MatchError) HTTPStatusCode() int {
	return http.StatusNotFound
}

var (
	ErrPathEnd        = &MatchError{cause: errors.New("no more path segments left")}
	ErrPathNotEnd     = &MatchError{cause: errors.New("path has more available segments")}
	ErrTrailingSlash  = &MatchError{cause: errors.New("trailing slash")}
	ErrDoubleSlash    = &MatchError{cause: errors.New("two consequitive slashes in path")}
	ErrPrefixMismatch = &MatchError{cause: errors.New("required prefix is does not match request URL path")}
)

type Matcher struct {
	source []byte
	slices [][]byte
	cursor int
}

func NewMatcher(p string) *Matcher {
	if p == "" {
		return &Matcher{}
	}
	source := []byte(p)
	cursor := 0
	if source[0] == '/' {
		cursor++
	}
	return &Matcher{source: source, cursor: cursor}
}

func (p *Matcher) MatchBytes() ([]byte, error) {
	if last := len(p.source) - 1; p.cursor >= last {
		// panic("trailing")
		if p.source[last] == '/' {
			return nil, ErrTrailingSlash
		}
		return nil, ErrPathEnd
	}

	var (
		i int
		b byte
	)

	// fmt.Println(string(p.source[p.cursor:]))
	for i, b = range p.source[p.cursor:] {
		if b == '/' {
			if i == 0 {
				// if p.cursor == len(p.source)-1 {
				// 	return nil, ErrTrailingSlash // TODO: test.
				// }
				return nil, ErrDoubleSlash
			}
			result := p.source[p.cursor : p.cursor+i]
			p.cursor += i + 1
			return result, nil
		}
	}
	p.cursor += i
	return p.source[p.cursor-i : p.cursor+1], nil
}

func (p *Matcher) MatchString() (string, error) {
	b, err := p.MatchBytes()
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (p *Matcher) MatchInt() (int, error) {
	b, err := p.MatchBytes()
	if err != nil {
		return 0, err
	}
	i, err := strconv.Atoi(string(b))
	if err != nil {
		return 0, err
	}
	return i, nil
}

func (p *Matcher) MatchUint() (uint, error) {
	i, err := p.MatchInt()
	if err != nil {
		return 0, err
	}
	return uint(i), nil
}

func (p *Matcher) MatchSlug() (string, error) {
	return "", errors.New("unimplemented")
}

func (p *Matcher) Skip() (err error) {
	_, err = p.MatchBytes()
	return err
}

func (p *Matcher) MatchEnd() error {
	_, err := p.MatchBytes()
	if err == nil {
		return ErrPathNotEnd
	}
	if errors.Is(err, ErrPathEnd) {
		return nil
	}
	return err
}

func TailTrailingSlashRedirectOrNotFound(w http.ResponseWriter, r *http.Request, tail string) error {
	for _, character := range tail {
		if character != '/' {
			return oakhttp.NewNotFoundError(r.URL.Path)
		}
	}

	URL := r.URL.String()
	if URL == "" {
		return errors.New("trailing slash redirect received request with empty URL")
	}
	http.Redirect(w, r, URL[:len(URL)-1], http.StatusTemporaryRedirect)
	return nil
}
