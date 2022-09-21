package stencil_app

const (
	NavUpdates  = "Updates"
	NavProducts = "Products"
)

var NavItems = []string{NavUpdates, NavProducts}

var NavIcons = map[string]string{
	NavUpdates:  "updates",
	NavProducts: "stack",
}

var NavHrefs = map[string]string{
	NavUpdates:  UpdatesPath,
	NavProducts: ProductsPath,
}
