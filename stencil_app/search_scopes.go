package stencil_app

import (
	"github.com/arelate/vangogh_local_data"
	"net/url"
)

const (
	ScopeNewSearch = "New search"
	ScopeOwned     = "Owned"
	ScopeWishlist  = "Wishlist"
	ScopeSale      = "Sale"
	ScopeAll       = "All"
)

var SearchScopes = []string{
	ScopeNewSearch,
	ScopeOwned,
	ScopeWishlist,
	ScopeSale,
	ScopeAll,
}

func SearchScopeQueries() map[string]string {

	scopes := make(map[string]string)

	scopes[ScopeNewSearch] = ""

	q := make(url.Values)
	q.Set(vangogh_local_data.TypesProperty, vangogh_local_data.AccountProducts.String())
	q.Set(vangogh_local_data.SortProperty, vangogh_local_data.GOGOrderDateProperty)
	q.Set(vangogh_local_data.DescendingProperty, vangogh_local_data.TrueValue)
	scopes[ScopeOwned] = q.Encode()

	q = make(url.Values)
	q.Set(vangogh_local_data.WishlistedProperty, vangogh_local_data.TrueValue)
	q.Set(vangogh_local_data.SortProperty, vangogh_local_data.GOGReleaseDateProperty)
	q.Set(vangogh_local_data.DescendingProperty, vangogh_local_data.TrueValue)
	scopes[ScopeWishlist] = q.Encode()

	q = make(url.Values)
	q.Set(vangogh_local_data.TypesProperty, vangogh_local_data.CatalogProducts.String())
	q.Set(vangogh_local_data.OwnedProperty, vangogh_local_data.FalseValue)
	q.Set(vangogh_local_data.IsDiscountedProperty, vangogh_local_data.TrueValue)
	q.Set(vangogh_local_data.SortProperty, vangogh_local_data.DiscountPercentageProperty)
	q.Set(vangogh_local_data.DescendingProperty, vangogh_local_data.TrueValue)
	scopes[ScopeSale] = q.Encode()

	q = make(url.Values)
	q.Set(vangogh_local_data.TypesProperty, vangogh_local_data.CatalogProducts.String())
	q.Set(vangogh_local_data.SortProperty, vangogh_local_data.GOGReleaseDateProperty)
	q.Set(vangogh_local_data.DescendingProperty, vangogh_local_data.TrueValue)
	scopes[ScopeAll] = q.Encode()

	return scopes
}
