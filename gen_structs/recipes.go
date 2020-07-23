package gen_structs

import (
	"bytes"
	"encoding/json"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"net"
	"os"
)

var recipesFile *os.File
var recipesDone = false

func setupRecipes(generator *Generator) {
	recipesFile = generator.AddFile("recipes.json")
	generator.AddListener(recipes)
}

func recipes(header packet.Header, payload []byte, src, dst net.Addr) {
	if !recipesDone && header.PacketID == packet.IDCraftingData {
		ppacket := packet.CraftingData{}
		if err := ppacket.Unmarshal(bytes.NewBuffer(payload)); err != nil {
			panic(err.Error())
		}
		encoder := json.NewEncoder(recipesFile)
		mmap := getMap(ppacket)
		if err := encoder.Encode(mmap); err != nil {
			panic(err.Error())
		}
		if err := recipesFile.Sync(); err != nil {
			panic(err.Error())
		}
		recipesDone = true
	}
}

func getMap(ppacket packet.CraftingData) map[string]interface{} {
	var shapeless, shaped, furnace, furnaceAux, multi, shulkerBox, shapelessChemistry, shapedChemistry []protocol.Recipe
	for _, recipe := range ppacket.Recipes {
		switch recipe.(type) {
		case *protocol.ShapelessRecipe:
			shapeless = append(shapeless, recipe)
		case *protocol.ShapedRecipe:
			shaped = append(shaped, recipe)
		case *protocol.FurnaceRecipe:
			furnace = append(furnace, recipe)
		case *protocol.FurnaceDataRecipe:
			furnaceAux = append(furnaceAux, recipe)
		case *protocol.MultiRecipe:
			multi = append(multi, recipe)
		case *protocol.ShulkerBoxRecipe:
			shulkerBox = append(shulkerBox, recipe)
		case *protocol.ShapelessChemistryRecipe:
			shapelessChemistry = append(shapelessChemistry, recipe)
		case *protocol.ShapedChemistryRecipe:
			shapedChemistry = append(shapedChemistry, recipe)
		}
	}
	return map[string]interface{}{
		"recipes": map[string]interface{}{
			"shapeless":          shapeless,
			"shaped":             shaped,
			"furnace":            furnace,
			"furnaceAux":         furnaceAux,
			"multi":              multi,
			"shulkerBox":         shulkerBox,
			"shapelessChemistry": shapelessChemistry,
			"shapedChemistry":    shapedChemistry,
		},
		"potion": ppacket.PotionRecipes,
	}
}
