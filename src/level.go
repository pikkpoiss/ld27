// Copyright 2013 Arne Roomann-Kurrik
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"./system"
	"fmt"
)

type Level struct {
	Map        *system.TiledMap
	Camera     *Camera
	StartX     int
	StartY     int
	Goal       *Tile
	tiles      []Tile
	TileWidth  int
	TileHeight int
}

func LoadLevel(path string) (out *Level, err error) {
	var (
		tm    *system.TiledMap
		cw    float64
		ch    float64
		tiles []Tile
	)
	if tm, err = system.LoadMap(path); err != nil {
		return
	}
	cw = float64(tm.Width * tm.Tilewidth)
	ch = float64(tm.Height * tm.Tileheight)
	tiles = make([]Tile, tm.Width*tm.Height)
	out = &Level{
		Map:        tm,
		Camera:     NewCamera(0, 0, cw, ch),
		TileWidth:  tm.Tilewidth,
		TileHeight: tm.Tileheight,
		tiles:      tiles,
	}
	if err = out.parseTiles(); err != nil {
		return
	}
	if err = out.parseObjects(); err != nil {
		return
	}
	return
}

func (l *Level) Update() (err error) {
	var (
		layer *system.TiledLayer
	)
	if layer, err = l.getLayer("tilelayer", "Tiles"); err != nil {
		return
	}
	for _, a := range TILE_ANIMATIONS {
		a.Next()
	}
	for i, t := range l.tiles {
		layer.Data[i] = TILE_ANIMATIONS[t.Type].Curr()
	}
	return
}

func (l *Level) getLayer(t string, n string) (out *system.TiledLayer, err error) {
	for i, _ := range l.Map.Layers {
		if l.Map.Layers[i].Type != t && l.Map.Layers[i].Name != n {
			continue
		}
		out = &l.Map.Layers[i]
		return
	}
	err = fmt.Errorf("Could not find layer with type %v and name %v", t, n)
	return
}

func (l *Level) getTileAt(x int, y int) (t *Tile, err error) {
	var i = l.Map.Width*y + x
	if i >= len(l.tiles) || i < 0 {
		err = fmt.Errorf("No tile at (%v,%v)", x, y)
		return
	}
	t = &l.tiles[i]
	return
}

func (l *Level) getTileAtPixel(x int, y int) (t *Tile, err error) {
	x = x / l.Map.Tilewidth
	y = y / l.Map.Tileheight
	return l.getTileAt(x, y)
}

func (l *Level) iToX(i int) int {
	return i % l.Map.Width
}

func (l *Level) iToY(i int) int {
	return i / l.Map.Width
}

func (l *Level) parseTiles() (err error) {
	var (
		layer *system.TiledLayer
	)
	if layer, err = l.getLayer("tilelayer", "Tiles"); err != nil {
		return
	}
	for i, v := range layer.Data {
		l.tiles[i] = Tile{
			X:    l.iToX(i),
			Y:    l.iToY(i),
			Type: v,
		}
	}
	return
}

func (l *Level) parseObjects() (err error) {
	var (
		layer *system.TiledLayer
	)
	if layer, err = l.getLayer("objectlayer", "Objects"); err != nil {
		return
	}
	for _, obj := range layer.Objects {
		switch obj.Type {
		case "player":
			l.StartX = obj.X
			l.StartY = obj.Y
		case "goal":
			l.Goal, err = l.getTileAtPixel(obj.X, obj.Y)
		}
		if err != nil {
			return
		}
	}
	return
}

const (
	TILE_GRASS = 1 + iota
	TILE_STONE
	TILE_BRICK
)

var TILE_ANIMATIONS = map[int]*system.Animation{
	TILE_GRASS: system.Anim([]int{1}, 4),
	TILE_STONE: system.Anim([]int{2}, 4),
	TILE_BRICK: system.Anim([]int{3, 2}, 4),
}

type Tile struct {
	X    int
	Y    int
	Type int
}
