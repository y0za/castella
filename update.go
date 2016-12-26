package main

// Update represents file updated lines data
type Update struct {
	Name  string   `json:"name"`
	Lines []string `json:"lines"`
}
