package controllers

// FilteredLocation represents the structure for filtered locations.
type FilteredLocation struct {
    DestId   string `json:"dest_id"`
    DestType string `json:"dest_type"`
    Value    string `json:"value"`
}