package namespace

type kind int8

type fileTree struct {
	entries map[string] kind
}

func newFileTree(root string, rootNamespaceName string) *fileTree {

}

