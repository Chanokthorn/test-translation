package translation

import (
	"fmt"
	"strconv"
)

type CollectorTree interface {
	CollectTranslationItemsFromRoot(root any, whitelist any) []TranslationItem
}

type collectorTree struct {
}

func NewCollectorTree() CollectorTree {
	return &collectorTree{}
}

func (c *collectorTree) CollectTranslationItemsFromRoot(root any, whitelist any) []TranslationItem {
	dataNode := CreateTreeFromMap(root)
	whitelistNode := CreateMatchingConditionTree(whitelist)
	return CollectTranslationItemsFromRoot(dataNode, whitelistNode)
}

type Node struct {
	Value    any
	Children map[string]*Node
}

func CreateTreeFromMap(root any) *Node {
	return createTreeFromMapFunc(root)
}

func createTreeFromMapFunc(root any) *Node {
	if root == nil {
		return nil
	}

	node := &Node{
		Value:    root,
		Children: make(map[string]*Node),
	}

	switch typedValue := root.(type) {
	case string:
		// base case - is leaf node
	case map[string]any:
		// common case - is map
		for k := range typedValue {
			childNode := createTreeFromMapFunc(typedValue[k])
			if childNode != nil {
				node.Children[k] = childNode
			}
		}
	case []any:
		// common case - is array
		for i := range typedValue {
			childNode := createTreeFromMapFunc(typedValue[i])
			if childNode != nil {
				node.Children[strconv.Itoa(i)] = childNode
			}
		}
	default:
		// unsupported type
		return nil
	}

	return node
}

type MatchingNode struct {
	Children map[string]*MatchingNode
}

func CreateMatchingConditionTree(root any) *MatchingNode {
	node := &MatchingNode{
		Children: make(map[string]*MatchingNode),
	}

	switch typedValue := root.(type) {
	case map[string]any:
		// common case - is map
		for k, v := range typedValue {
			node.Children[k] = CreateMatchingConditionTree(v)
		}
	case bool:
		// base case - is leaf node
	default:
		// unsupported type
		return nil
	}

	return node
}

func CollectTranslationItemsFromRoot(root *Node, whitelist *MatchingNode) []TranslationItem {
	if root == nil || whitelist == nil {
		return nil
	}

	return CollectTranslationItems(root, whitelist, nil, nil, "")
}

// get node by data and whitelist,
func CollectTranslationItems(dataNode *Node, whitelistNode *MatchingNode, container any, key any, path string) []TranslationItem {
	// base case
	if dataNode == nil || whitelistNode == nil {
		return nil
	}

	translationItems := []TranslationItem{}
	switch typedValue := dataNode.Value.(type) {
	case string:
		// base case - is leaf node
		translationItems = append(translationItems, TranslationItem{
			Container: container,
			Key:       key,
			Value:     typedValue,
			Path:      path,
		})
	case map[string]any:
		for k, v := range dataNode.Children {
			whitelistChild, ok := whitelistNode.Children[k]
			if !ok {
				continue
			}

			childPath := fmt.Sprintf("%s.%s", path, k)
			// see if we can use typedValue as container instead
			translationItems = append(translationItems, CollectTranslationItems(v, whitelistChild, dataNode.Value, k, childPath)...)
		}
	case []any:
		for i_string, v := range dataNode.Children {
			i, err := strconv.Atoi(i_string)
			if err != nil {
				continue
			}
			childPath := fmt.Sprintf("%s[%d]", path, i)
			translationItems = append(translationItems, CollectTranslationItems(v, whitelistNode, dataNode.Value, i, childPath)...)
		}
	default:
		// unsupported type
		return nil
	}

	return translationItems
}
