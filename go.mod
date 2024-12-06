module github.com/TheBitDrifter/coldbrew

go 1.23.3

require (
	github.com/TheBitDrifter/blueprint v0.0.0-00010101000000-000000000000
	github.com/TheBitDrifter/table v0.0.0-20241202222727-621c10848124
	github.com/TheBitDrifter/warehouse v0.0.0-20241202220617-cb9ecc34a5c3
	github.com/hajimehoshi/ebiten/v2 v2.8.5
)

require (
	github.com/TheBitDrifter/mask v0.0.0-20241122180741-07926f8b9e86 // indirect
	github.com/TheBitDrifter/util v0.0.0-20241102212109-342f4c0a810e // indirect
	github.com/ebitengine/gomobile v0.0.0-20240911145611-4856209ac325 // indirect
	github.com/ebitengine/hideconsole v1.0.0 // indirect
	github.com/ebitengine/purego v0.8.0 // indirect
	github.com/jezek/xgb v1.1.1 // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/sys v0.25.0 // indirect
)

replace github.com/TheBitDrifter/blueprint => ../blueprint/

replace github.com/TheBitDrifter/table => ../table/

replace github.com/TheBitDrifter/warehouse => ../warehouse/
