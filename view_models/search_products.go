package view_models

const SearchProductsLimit = 100

type searchProducts struct {
	Context          string
	Scope            string
	SearchProperties []string
	Query            map[string]string
	Digests          map[string][]string
	DigestsTitles    map[string]string
	Products         []listProduct
	Limit            int
	Total            int
	Constrained      bool
	ConstrainedPath  string
	Path             string
}

func NewSearchProducts(scope string, constrained bool, path string) *searchProducts {
	return &searchProducts{
		Scope:            scope,
		Context:          "filter-products",
		SearchProperties: SearchProperties,
		Query:            make(map[string]string),
		DigestsTitles:    digestTitles,
		Limit:            SearchProductsLimit,
		Constrained:      constrained,
		Path:             path,
	}
}
