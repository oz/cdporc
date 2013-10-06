package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
)

type Quote struct {
	Class     string `json:"json_class"`
	Id        int    `json:"id"`
	Topic     string `json:"topic"`
	Body      string `json:"body"`
	CreatedAt string `json:"created_at"`
	State     string `json:"state"`
}

type Quotes struct {
	Page              int     `json:"page"`
	PageCount         int     `json:"page_count"`
	TotalEntriesCount int     `json:"total_entries_count"`
	Entries           []Quote `json:"entries"`
	Server            *url.URL
}

func NewCoteDePorc(server *url.URL) *Quotes {
	return &Quotes{Server: server}
}

// Get all quotes at this path
func (q *Quotes) GetAll(path string) error {
	currentPage := 1
	stopAt := currentPage

	// TODO: after fetching the first page, if the PageCount is huge, having a
	//       pool of fetchers would be nice. :)
	for currentPage <= stopAt {
		if err := q.getPage(path, currentPage); err != nil {
			return err
		}
		if q.PageCount > stopAt {
			stopAt = q.PageCount
		}
		currentPage += 1
	}
	sort.Sort(q)

	return nil
}

// Delete a quote
func (q *Quotes) Delete(id string) error {
	resp, err := q.setupRequest("DELETE", "/quotes/"+id)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// Confirm a quote
func (q *Quotes) Confirm(id string) error {
	resp, err := q.setupRequest("PUT", "/quotes/"+id+"/confirm")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// Get a Random quote
func (q *Quotes) Random() (Quote, error) {
	quote := Quote{}
	resp, err := q.setupRequest("GET", "/quotes/random")
	if err != nil {
		return quote, err
	}
	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return quote, err
	}
	json.Unmarshal(buf, &quote)

	return quote, nil
}

func (q *Quotes) getPage(path string, page int) error {
	if page > 1 {
		path = fmt.Sprintf("%s?page=%d", path, page)
	}
	if err := q.getJSON(path); err != nil {
		return err
	}
	return nil
}

func (q *Quotes) getJSON(path string) error {
	body, err := q.getPath(path)
	if err != nil {
		return err
	}

	quotes := Quotes{}
	json.Unmarshal(body, &quotes)

	// Append fetched quotes, and update internal state
	q.Entries = append(q.Entries, quotes.Entries...)
	q.Page = quotes.Page
	q.PageCount = quotes.PageCount
	q.TotalEntriesCount = quotes.TotalEntriesCount
	return nil
}

// Read "all" the data for the requested path on the server.
func (q *Quotes) getPath(path string) ([]byte, error) {
	resp, err := q.setupRequest("GET", path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

// Setup a basic HTTP client with the given method and path, and return
// a reference to the http.Response reply.
func (q *Quotes) setupRequest(method string, relativePath string) (*http.Response, error) {
	path := q.absoluteUrl(relativePath)
	client := &http.Client{}

	req, err := http.NewRequest(method, path, nil)
	if err != nil {
		return nil, err
	}
	q.setupAuth(req)

	resp, err := client.Do(req)
	if err != nil {
		return resp, nil
	}
	return resp, nil
}

// Setup HTTP basic auth on a http.Request, if configured
func (q *Quotes) setupAuth(req *http.Request) {
	password, passwordIsSet := q.Server.User.Password()

	if passwordIsSet {
		username := q.Server.User.Username()
		req.SetBasicAuth(username, password)
	}
}

func (q *Quotes) absoluteUrl(path string) string {
	return q.Server.String() + path
}

// Sortable interface
func (q *Quotes) Len() int {
	return len(q.Entries)
}

func (q *Quotes) Swap(i, j int) {
	q.Entries[i], q.Entries[j] = q.Entries[j], q.Entries[i]
}

func (q *Quotes) Less(i, j int) bool {
	return q.Entries[i].Id < q.Entries[j].Id
}
