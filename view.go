package couch

import (
	"encoding/json"
)

// CouchDB Design Document (not yet public)
type design struct {
	Doc
	Views map[string]view `json:"views"`
	// There a more elements to a design document, they will be added when they are implemented
}

// CouchDB View (not yet public)
type view struct {
	Map    string `json:"map,omitempty"`
	Reduce string `json:"reduce,omitempty"`
}

// Container for ViewResultRows
type ViewResult struct {
	Offset uint64
	Rows   []ViewResultRow
}

// A single view result
type ViewResultRow struct {
	ID    string
	Key   interface{}
	Value interface{}
	Doc   json.RawMessage
}

func (r *ViewResultRow) ValueInt() int {
	num, _ := r.Value.(float64)
	return int(num)
}

// Checks if a view really exists
func (db *Database) HasView(designID, viewID string) bool {
	ok, _ := checkHead(db.viewURL(designID, viewID))
	return ok
}

// Query a view with options, see http://docs.couchdb.org/en/latest/api/ddoc/views.html#db-design-design-doc-view-view-name
func (db *Database) Query(designID, viewID string, options map[string]interface{}) (*ViewResult, error) {
	result := &ViewResult{}
	url := db.viewURL(designID, viewID) + urlEncode(options)
	_, err := Do(db.server.httpClient, url, "GET", db.Cred(), nil, &result)

	return result, err
}

func (db *Database) RetrieveAll() (*ViewResult, error) {
	result := &ViewResult{}
	options := make(map[string]interface{})
	options[OPTIONS_INCLUDE_DOCS] = true
	url := db.allDocsURL() + urlEncode(options)
	_, err := Do(db.server.httpClient, url, "GET", db.Cred(), nil, &result)
	return result, err
}

// Create a new design document (not yet public)
func newDesign() *design {
	d := &design{}
	d.Views = make(map[string]view)
	return d
}

// Get the complete url to a view of a design document
func (db *Database) viewURL(designID string, viewID string) string {
	return db.URL() + "/_design/" + designID + "/_view/" + viewID
}

// Get the complete url to the _all_docs view
func (db *Database) allDocsURL() string {
	return db.URL() + "/_all_docs"
}
