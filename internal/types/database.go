package types

type MySQLFilter struct {
	Query []MySQLQuery
}

type MySQLQuery struct {
	Column string `json:"column"`
	Op     string `json:"op"`
	Query  string `json:"query"`
}
