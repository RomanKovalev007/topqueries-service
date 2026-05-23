package domain

type TopEntry struct {
    Query string `json:"query"`
    Count int64  `json:"count"`
}
