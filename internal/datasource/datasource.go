package datasource

type ReadDataSourceResult struct {
	TeamNames []string
	Members   []string
	Teams     map[string][]string
}

type DataSource interface {
	Read(string) (*ReadDataSourceResult, error)
}
