package mapper

import "image/color"
import "fmt"

type Block struct {
	Id   int16 // composed of "byte" and "add" from a chunk Section
	Data int8  // actually only 4 bits (0-15)
	// missing: light/skyligh
}

const (
	Air = iota
	Stone
	Grass
	Dirt
	Cobblestone
	Planks
	Sapling
	Bedrock
	Flowing_water
	Water
	Flowing_lava
	Lava
	Sand
	Gravel
	Gold_ore
	Iron_ore
	Coal_ore
	Log
	Leaves
	Sponge
	Glass
	Lapis_ore
	Lapis_block
	Dispenser
	Sandstone
	Noteblock
	Bed
	Golden_rail
	Detector_rail
	Sticky_piston
	Web
	Tallgrass
	Deadbush
	Piston
	Piston_head
	Wool
	Piston_extension
	Yellow_flower
	Red_flower
	Brown_mushroom
	Red_mushroom
	Gold_block
	Iron_block
	Double_stone_slab
	Stone_slab
	Brick_block
	Tnt
	Bookshelf
	Mossy_cobblestone
	Obsidian
	Torch
	Fire
	Mob_spawner
	Oak_stairs
	Chest
	Redstone_wire
	Diamond_ore
	Diamond_block
	Crafting_table
	Wheat
	Farmland
	Furnace
	Lit_furnace
	Standing_sign
	Wooden_door
	Ladder
	Rail
	Stone_stairs
	Wall_sign
	Lever
	Stone_pressure_plate
	Iron_door
	Wooden_pressure_plate
	Redstone_ore
	Lit_redstone_ore
	Unlit_redstone_torch
	Redstone_torch
	Stone_button
	Snow_layer
	Ice
	Snow
	Cactus
	Clay
	Reeds
	Jukebox
	Fence
	Pumpkin
	Netherrack
	Soul_sand
	Glowstone
	Portal
	Lit_pumpkin
	Cake
	Unpowered_repeater
	Powered_repeater
	Stained_glass
	Trapdoor
	Monster_egg
	Stonebrick
	Brown_mushroom_block
	Red_mushroom_block
	Iron_bars
	Glass_pane
	Melon_block
	Pumpkin_stem
	Melon_stem
	Vine
	Fence_gate
	Brick_stairs
	Stone_brick_stairs
	Mycelium
	Waterlily
	Nether_brick
	Nether_brick_fence
	Nether_brick_stairs
	Nether_wart
	Enchanting_table
	Brewing_stand
	Cauldron
	End_portal
	End_portal_frame
	End_stone
	Dragon_egg
	Redstone_lamp
	Lit_redstone_lamp
	Double_wooden_slab
	Wooden_slab
	Cocoa
	Sandstone_stairs
	Emerald_ore
	Ender_chest
	Tripwire_hook
	Tripwire
	Emerald_block
	Spruce_stairs
	Birch_stairs
	Jungle_stairs
	Command_block
	Beacon
	Cobblestone_wall
	Flower_pot
	Carrots
	Potatoes
	Wooden_button
	Skull
	Anvil
	Trapped_chest
	Light_weighted_pressure_plate
	Heavy_weighted_pressure_plate
	Unpowered_comparator
	Powered_comparator
	Daylight_detector
	Redstone_block
	Quartz_ore
	Hopper
	Quartz_block
	Quartz_stairs
	Activator_rail
	Dropper
	Stained_hardened_clay
	Stained_glass_pane
	Leaves2
	Log2
	Acacia_stairs
	Dark_oak_stairs
	Slime
	Barrier
	Iron_trapdoor
	Prismarine
	Sea_lantern
	Hay_block
	Carpet
	Hardened_clay
	Coal_block
	Packed_ice
	Double_plant
	Standing_banner
	Wall_banner
	Daylight_detector_inverted
	Red_sandstone
	Red_sandstone_stairs
	Double_stone_slab2
	Stone_slab2
	Spruce_fence_gate
	Birch_fence_gate
	Jungle_fence_gate
	Dark_oak_fence_gate
	Acacia_fence_gate
	Spruce_fence
	Birch_fence
	Jungle_fence
	Dark_oak_fence
	Acacia_fence
	Spruce_door
	Birch_door
	Jungle_door
	Acacia_door
	Dark_oak_door
	End_rod
	Chorus_plant
	Chorus_flower
	Purpur_block
	Purpur_pillar
	Purpur_stairs
	Purpur_double_slab
	Purpur_slab
	End_bricks
	Beetroots
	Grass_path
	End_gateway
	Repeating_command_block
	Chain_command_block
	Frosted_ice
	Magma
	Nether_wart_block
	Red_nether_brick
	Bone_block
	Structure_void
	Observer
	White_shulker_box
	Orange_shulker_box
	Magenta_shulker_box
	Light_blue_shulker_box
	Yellow_shulker_box
	Lime_shulker_box
	Pink_shulker_box
	Gray_shulker_box
	Light_gray_shulker_box
	Cyan_shulker_box
	Purple_shulker_box
	Blue_shulker_box
	Brown_shulker_box
	Green_shulker_box
	Red_shulker_box
	Black_shulker_box
	Structure_block
)

// this is independat of the zany way the chunk format stores stuff
// just take the id and the data fields as parameters
func NewBlock(id int16, data int8) Block {
	// returns new struct
	return Block{id, data}
}

func (b Block) Name() string {
	// giant case statement which returns things like "minecraft:stone"
	switch b.Id {
	case Air:
		return "minecraft:air"
	case Stone:
		return "minecraft:stone"
	case Grass:
		return "minecraft:grass"
	case Dirt:
		return "minecraft:dirt"
	case Cobblestone:
		return "minecraft:cobblestone"
	case Planks:
		return "minecraft:planks"
	case Sapling:
		return "minecraft:sapling"
	case Bedrock:
		return "minecraft:bedrock"
	case Flowing_water:
		return "minecraft:flowing_water"
	case Water:
		return "minecraft:water"
	case Flowing_lava:
		return "minecraft:flowing_lava"
	case Lava:
		return "minecraft:lava"
	case Sand:
		return "minecraft:sand"
	case Gravel:
		return "minecraft:gravel"
	case Gold_ore:
		return "minecraft:gold_ore"
	case Iron_ore:
		return "minecraft:iron_ore"
	case Coal_ore:
		return "minecraft:coal_ore"
	case Log:
		return "minecraft:log"
	case Leaves:
		return "minecraft:leaves"
	case Sponge:
		return "minecraft:sponge"
	case Glass:
		return "minecraft:glass"
	case Lapis_ore:
		return "minecraft:lapis_ore"
	case Lapis_block:
		return "minecraft:lapis_block"
	case Dispenser:
		return "minecraft:dispenser"
	case Sandstone:
		return "minecraft:sandstone"
	case Noteblock:
		return "minecraft:noteblock"
	case Bed:
		return "minecraft:bed"
	case Golden_rail:
		return "minecraft:golden_rail"
	case Detector_rail:
		return "minecraft:detector_rail"
	case Sticky_piston:
		return "minecraft:sticky_piston"
	case Web:
		return "minecraft:web"
	case Tallgrass:
		return "minecraft:tallgrass"
	case Deadbush:
		return "minecraft:deadbush"
	case Piston:
		return "minecraft:piston"
	case Piston_head:
		return "minecraft:piston_head"
	case Wool:
		return "minecraft:wool"
	case Piston_extension:
		return "minecraft:piston_extension"
	case Yellow_flower:
		return "minecraft:yellow_flower"
	case Red_flower:
		return "minecraft:red_flower"
	case Brown_mushroom:
		return "minecraft:brown_mushroom"
	case Red_mushroom:
		return "minecraft:red_mushroom"
	case Gold_block:
		return "minecraft:gold_block"
	case Iron_block:
		return "minecraft:iron_block"
	case Double_stone_slab:
		return "minecraft:double_stone_slab"
	case Stone_slab:
		return "minecraft:stone_slab"
	case Brick_block:
		return "minecraft:brick_block"
	case Tnt:
		return "minecraft:tnt"
	case Bookshelf:
		return "minecraft:bookshelf"
	case Mossy_cobblestone:
		return "minecraft:mossy_cobblestone"
	case Obsidian:
		return "minecraft:obsidian"
	case Torch:
		return "minecraft:torch"
	case Fire:
		return "minecraft:fire"
	case Mob_spawner:
		return "minecraft:mob_spawner"
	case Oak_stairs:
		return "minecraft:oak_stairs"
	case Chest:
		return "minecraft:chest"
	case Redstone_wire:
		return "minecraft:redstone_wire"
	case Diamond_ore:
		return "minecraft:diamond_ore"
	case Diamond_block:
		return "minecraft:diamond_block"
	case Crafting_table:
		return "minecraft:crafting_table"
	case Wheat:
		return "minecraft:wheat"
	case Farmland:
		return "minecraft:farmland"
	case Furnace:
		return "minecraft:furnace"
	case Lit_furnace:
		return "minecraft:lit_furnace"
	case Standing_sign:
		return "minecraft:standing_sign"
	case Wooden_door:
		return "minecraft:wooden_door"
	case Ladder:
		return "minecraft:ladder"
	case Rail:
		return "minecraft:rail"
	case Stone_stairs:
		return "minecraft:stone_stairs"
	case Wall_sign:
		return "minecraft:wall_sign"
	case Lever:
		return "minecraft:lever"
	case Stone_pressure_plate:
		return "minecraft:stone_pressure_plate"
	case Iron_door:
		return "minecraft:iron_door"
	case Wooden_pressure_plate:
		return "minecraft:wooden_pressure_plate"
	case Redstone_ore:
		return "minecraft:redstone_ore"
	case Lit_redstone_ore:
		return "minecraft:lit_redstone_ore"
	case Unlit_redstone_torch:
		return "minecraft:unlit_redstone_torch"
	case Redstone_torch:
		return "minecraft:redstone_torch"
	case Stone_button:
		return "minecraft:stone_button"
	case Snow_layer:
		return "minecraft:snow_layer"
	case Ice:
		return "minecraft:ice"
	case Snow:
		return "minecraft:snow"
	case Cactus:
		return "minecraft:cactus"
	case Clay:
		return "minecraft:clay"
	case Reeds:
		return "minecraft:reeds"
	case Jukebox:
		return "minecraft:jukebox"
	case Fence:
		return "minecraft:fence"
	case Pumpkin:
		return "minecraft:pumpkin"
	case Netherrack:
		return "minecraft:netherrack"
	case Soul_sand:
		return "minecraft:soul_sand"
	case Glowstone:
		return "minecraft:glowstone"
	case Portal:
		return "minecraft:portal"
	case Lit_pumpkin:
		return "minecraft:lit_pumpkin"
	case Cake:
		return "minecraft:cake"
	case Unpowered_repeater:
		return "minecraft:unpowered_repeater"
	case Powered_repeater:
		return "minecraft:powered_repeater"
	case Stained_glass:
		return "minecraft:stained_glass"
	case Trapdoor:
		return "minecraft:trapdoor"
	case Monster_egg:
		return "minecraft:monster_egg"
	case Stonebrick:
		return "minecraft:stonebrick"
	case Brown_mushroom_block:
		return "minecraft:brown_mushroom_block"
	case Red_mushroom_block:
		return "minecraft:red_mushroom_block"
	case Iron_bars:
		return "minecraft:iron_bars"
	case Glass_pane:
		return "minecraft:glass_pane"
	case Melon_block:
		return "minecraft:melon_block"
	case Pumpkin_stem:
		return "minecraft:pumpkin_stem"
	case Melon_stem:
		return "minecraft:melon_stem"
	case Vine:
		return "minecraft:vine"
	case Fence_gate:
		return "minecraft:fence_gate"
	case Brick_stairs:
		return "minecraft:brick_stairs"
	case Stone_brick_stairs:
		return "minecraft:stone_brick_stairs"
	case Mycelium:
		return "minecraft:mycelium"
	case Waterlily:
		return "minecraft:waterlily"
	case Nether_brick:
		return "minecraft:nether_brick"
	case Nether_brick_fence:
		return "minecraft:nether_brick_fence"
	case Nether_brick_stairs:
		return "minecraft:nether_brick_stairs"
	case Nether_wart:
		return "minecraft:nether_wart"
	case Enchanting_table:
		return "minecraft:enchanting_table"
	case Brewing_stand:
		return "minecraft:brewing_stand"
	case Cauldron:
		return "minecraft:cauldron"
	case End_portal:
		return "minecraft:end_portal"
	case End_portal_frame:
		return "minecraft:end_portal_frame"
	case End_stone:
		return "minecraft:end_stone"
	case Dragon_egg:
		return "minecraft:dragon_egg"
	case Redstone_lamp:
		return "minecraft:redstone_lamp"
	case Lit_redstone_lamp:
		return "minecraft:lit_redstone_lamp"
	case Double_wooden_slab:
		return "minecraft:double_wooden_slab"
	case Wooden_slab:
		return "minecraft:wooden_slab"
	case Cocoa:
		return "minecraft:cocoa"
	case Sandstone_stairs:
		return "minecraft:sandstone_stairs"
	case Emerald_ore:
		return "minecraft:emerald_ore"
	case Ender_chest:
		return "minecraft:ender_chest"
	case Tripwire_hook:
		return "minecraft:tripwire_hook"
	case Tripwire:
		return "minecraft:tripwire"
	case Emerald_block:
		return "minecraft:emerald_block"
	case Spruce_stairs:
		return "minecraft:spruce_stairs"
	case Birch_stairs:
		return "minecraft:birch_stairs"
	case Jungle_stairs:
		return "minecraft:jungle_stairs"
	case Command_block:
		return "minecraft:command_block"
	case Beacon:
		return "minecraft:beacon"
	case Cobblestone_wall:
		return "minecraft:cobblestone_wall"
	case Flower_pot:
		return "minecraft:flower_pot"
	case Carrots:
		return "minecraft:carrots"
	case Potatoes:
		return "minecraft:potatoes"
	case Wooden_button:
		return "minecraft:wooden_button"
	case Skull:
		return "minecraft:skull"
	case Anvil:
		return "minecraft:anvil"
	case Trapped_chest:
		return "minecraft:trapped_chest"
	case Light_weighted_pressure_plate:
		return "minecraft:light_weighted_pressure_plate"
	case Heavy_weighted_pressure_plate:
		return "minecraft:heavy_weighted_pressure_plate"
	case Unpowered_comparator:
		return "minecraft:unpowered_comparator"
	case Powered_comparator:
		return "minecraft:powered_comparator"
	case Daylight_detector:
		return "minecraft:daylight_detector"
	case Redstone_block:
		return "minecraft:redstone_block"
	case Quartz_ore:
		return "minecraft:quartz_ore"
	case Hopper:
		return "minecraft:hopper"
	case Quartz_block:
		return "minecraft:quartz_block"
	case Quartz_stairs:
		return "minecraft:quartz_stairs"
	case Activator_rail:
		return "minecraft:activator_rail"
	case Dropper:
		return "minecraft:dropper"
	case Stained_hardened_clay:
		return "minecraft:stained_hardened_clay"
	case Stained_glass_pane:
		return "minecraft:stained_glass_pane"
	case Leaves2:
		return "minecraft:leaves2"
	case Log2:
		return "minecraft:log2"
	case Acacia_stairs:
		return "minecraft:acacia_stairs"
	case Dark_oak_stairs:
		return "minecraft:dark_oak_stairs"
	case Slime:
		return "minecraft:slime"
	case Barrier:
		return "minecraft:barrier"
	case Iron_trapdoor:
		return "minecraft:iron_trapdoor"
	case Prismarine:
		return "minecraft:prismarine"
	case Sea_lantern:
		return "minecraft:sea_lantern"
	case Hay_block:
		return "minecraft:hay_block"
	case Carpet:
		return "minecraft:carpet"
	case Hardened_clay:
		return "minecraft:hardened_clay"
	case Coal_block:
		return "minecraft:coal_block"
	case Packed_ice:
		return "minecraft:packed_ice"
	case Double_plant:
		return "minecraft:double_plant"
	case Standing_banner:
		return "minecraft:standing_banner"
	case Wall_banner:
		return "minecraft:wall_banner"
	case Daylight_detector_inverted:
		return "minecraft:daylight_detector_inverted"
	case Red_sandstone:
		return "minecraft:red_sandstone"
	case Red_sandstone_stairs:
		return "minecraft:red_sandstone_stairs"
	case Double_stone_slab2:
		return "minecraft:double_stone_slab2"
	case Stone_slab2:
		return "minecraft:stone_slab2"
	case Spruce_fence_gate:
		return "minecraft:spruce_fence_gate"
	case Birch_fence_gate:
		return "minecraft:birch_fence_gate"
	case Jungle_fence_gate:
		return "minecraft:jungle_fence_gate"
	case Dark_oak_fence_gate:
		return "minecraft:dark_oak_fence_gate"
	case Acacia_fence_gate:
		return "minecraft:acacia_fence_gate"
	case Spruce_fence:
		return "minecraft:spruce_fence"
	case Birch_fence:
		return "minecraft:birch_fence"
	case Jungle_fence:
		return "minecraft:jungle_fence"
	case Dark_oak_fence:
		return "minecraft:dark_oak_fence"
	case Acacia_fence:
		return "minecraft:acacia_fence"
	case Spruce_door:
		return "minecraft:spruce_door"
	case Birch_door:
		return "minecraft:birch_door"
	case Jungle_door:
		return "minecraft:jungle_door"
	case Acacia_door:
		return "minecraft:acacia_door"
	case Dark_oak_door:
		return "minecraft:dark_oak_door"
	case End_rod:
		return "minecraft:end_rod"
	case Chorus_plant:
		return "minecraft:chorus_plant"
	case Chorus_flower:
		return "minecraft:chorus_flower"
	case Purpur_block:
		return "minecraft:purpur_block"
	case Purpur_pillar:
		return "minecraft:purpur_pillar"
	case Purpur_stairs:
		return "minecraft:purpur_stairs"
	case Purpur_double_slab:
		return "minecraft:purpur_double_slab"
	case Purpur_slab:
		return "minecraft:purpur_slab"
	case End_bricks:
		return "minecraft:end_bricks"
	case Beetroots:
		return "minecraft:beetroots"
	case Grass_path:
		return "minecraft:grass_path"
	case End_gateway:
		return "minecraft:end_gateway"
	case Repeating_command_block:
		return "minecraft:repeating_command_block"
	case Chain_command_block:
		return "minecraft:chain_command_block"
	case Frosted_ice:
		return "minecraft:frosted_ice"
	case Magma:
		return "minecraft:magma"
	case Nether_wart_block:
		return "minecraft:nether_wart_block"
	case Red_nether_brick:
		return "minecraft:red_nether_brick"
	case Bone_block:
		return "minecraft:bone_block"
	case Structure_void:
		return "minecraft:structure_void"
	case Observer:
		return "minecraft:observer"
	case White_shulker_box:
		return "minecraft:white_shulker_box"
	case Orange_shulker_box:
		return "minecraft:orange_shulker_box"
	case Magenta_shulker_box:
		return "minecraft:magenta_shulker_box"
	case Light_blue_shulker_box:
		return "minecraft:light_blue_shulker_box"
	case Yellow_shulker_box:
		return "minecraft:yellow_shulker_box"
	case Lime_shulker_box:
		return "minecraft:lime_shulker_box"
	case Pink_shulker_box:
		return "minecraft:pink_shulker_box"
	case Gray_shulker_box:
		return "minecraft:gray_shulker_box"
	case Light_gray_shulker_box:
		return "minecraft:light_gray_shulker_box"
	case Cyan_shulker_box:
		return "minecraft:cyan_shulker_box"
	case Purple_shulker_box:
		return "minecraft:purple_shulker_box"
	case Blue_shulker_box:
		return "minecraft:blue_shulker_box"
	case Brown_shulker_box:
		return "minecraft:brown_shulker_box"
	case Green_shulker_box:
		return "minecraft:green_shulker_box"
	case Red_shulker_box:
		return "minecraft:red_shulker_box"
	case Black_shulker_box:
		return "minecraft:black_shulker_box"
	case Structure_block:
		return "minecraft:structure_block"
	}
	panic(fmt.Sprintf("Block.Name: unhandled id %d", b.Id))
}

func (b Block) Block() string {
	// giant case statement which returns things like "Stone"
	return ""
}

func (b Block) Description() string {
	// placeholder; for things with data values, return a specific description
	// such as "Granite".  Might just return the same as Block().  for now
	// that's all it does
	return b.Block()
}

// returns the colour which should be used by this block on the terrain map
func (b Block) Colour() color.RGBA {
	// giant case statement which returns colours
	switch b.Id {
	case Stone, Cobblestone, Gravel:
		// could look at b.Data for stone type
		return color.RGBA{169, 169, 169, 255} // DarkGray
	case Grass:
		return color.RGBA{0, 128, 0, 255} // Green
	case Dirt:
		return color.RGBA{139, 69, 19, 255} // SaddleBrown
	case Sapling, Leaves, Tallgrass, Waterlily, Leaves2, Cactus:
		return color.RGBA{34, 139, 23, 255} // ForestGreen
	case Wheat:
		return color.RGBA{235, 222, 179, 255} // Wheat
	case Reeds:
		return color.RGBA{144, 238, 144, 255} // LightGreen
	case Melon_block:
		return color.RGBA{173, 255, 47, 255} // GreenYellow
	case Flowing_water, Water:
		return color.RGBA{0, 0, 255, 255} // Blue
	case Flowing_lava, Lava:
		return color.RGBA{255, 69, 0, 255} // OrangeRed
	case Sand, Sandstone:
		return color.RGBA{240, 230, 140, 255} //Khaki
	case Snow, Snow_layer:
		return color.RGBA{255, 250, 250, 244} // Snow
	}
	return color.RGBA{0, 0, 0, 255}
}
