package translation

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateTree(t *testing.T) {
	t.Run("TestCreateTree", func(t *testing.T) {
		data, err := os.ReadFile("small_data.json")
		if err != nil {
			t.Fatalf("failed to read data.json: %v", err)
		}

		var root any
		if err := json.Unmarshal(data, &root); err != nil {
			t.Fatalf("failed to unmarshal JSON: %v", err)
		}

		node := CreateTreeFromMap(root)
		require.NotNil(t, node)

		var traverse func(node *Node)
		traverse = func(node *Node) {
			if node == nil {
				return
			}
			t.Logf("Node: %v", node.Value)
			for _, child := range node.Children {
				traverse(child)
			}
		}

		traverse(node)
	})
}

func TestCreateMatchingConditionTree(t *testing.T) {
	t.Run("TestCreateMatchingConditionTree", func(t *testing.T) {
		data, err := os.ReadFile("whitelist.json")
		if err != nil {
			t.Fatalf("failed to read data.json: %v", err)
		}

		var root any
		if err := json.Unmarshal(data, &root); err != nil {
			t.Fatalf("failed to unmarshal JSON: %v", err)
		}

		node := CreateMatchingConditionTree(root)
		require.NotNil(t, node)
	})
}

func TestCollectTranslationItemsFromRoot(t *testing.T) {
	t.Run("TestCollectTranslationItemsFromRoot", func(t *testing.T) {
		dataFile, err := os.ReadFile("small_data.json")
		if err != nil {
			t.Fatalf("failed to read data.json: %v", err)
		}

		var dataJson any
		if err := json.Unmarshal(dataFile, &dataJson); err != nil {
			t.Fatalf("failed to unmarshal JSON: %v", err)
		}

		whitelistFile, err := os.ReadFile("whitelist.json")
		if err != nil {
			t.Fatalf("failed to read data.json: %v", err)
		}

		var whitelistJson any
		if err := json.Unmarshal(whitelistFile, &whitelistJson); err != nil {
			t.Fatalf("failed to unmarshal JSON: %v", err)
		}

		dataNode := CreateTreeFromMap(dataJson)
		whitelistNode := CreateMatchingConditionTree(whitelistJson)

		require.NotNil(t, dataNode)
		require.NotNil(t, whitelistNode)

		translationItems := CollectTranslationItemsFromRoot(dataNode, whitelistNode)
		require.NotNil(t, translationItems)

		for _, item := range translationItems {
			t.Logf(item.Path, item.Value)
		}
	})
}
