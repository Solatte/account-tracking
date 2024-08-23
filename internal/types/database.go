package types

type MySQLFilter struct {
	Query  []MySQLQuery `json:"query"`
	Limit  int          `json:"limit"`
	Offset int          `json:"offset"`
}

type MySQLQuery struct {
	Column string `json:"column"`
	Op     string `json:"op"`
	Query  string `json:"query"`
}
