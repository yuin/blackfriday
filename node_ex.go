package blackfriday

// Function is identifying a function node.
const Function NodeType = TableRow + 1

func init() {
	nodeTypeNames = append(nodeTypeNames, "Function")
}

// FunctionData contains fields relevant to a Function note type.
type FunctionData struct {
	Name      string
	Arguments []interface{}
}
