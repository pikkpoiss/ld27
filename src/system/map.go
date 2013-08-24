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

package system

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type TiledObject struct {
	Height     int
	Name       string
	Properties map[string]string
	Type       string
	Width      int
	X          int
	Y          int
}

type TiledLayer struct {
	Data    []int
	Height  int
	Name    string
	Opacity float32
	Type    string
	Visible bool
	Width   int
	X       int
	Y       int
	Objects []TiledObject
}

type TiledTileset struct {
	Firstgid         int
	Lastgid          int `json:-`
	Tilecount        int `json:-`
	Image            string
	Imageheight      int
	Imagewidth       int
	Margin           int
	Name             string
	Properties       map[string]string
	Spacing          int
	Tileheight       int
	Tilewidth        int
	Transparentcolor string
	Texture          *Texture `json:-`
}

type TiledMap struct {
	Height      int
	Layers      []TiledLayer
	Orientation string
	Properties  map[string]string
	Tileheight  int
	Tilesets    []TiledTileset
	Tilewidth   int
	Version     int
	Width       int
}

func LoadMap(path string) (out *TiledMap, err error) {
	var (
		f       *os.File
		tm      TiledMap
		decoder *json.Decoder
	)
	if f, err = os.Open(path); err != nil {
		return
	}
	defer f.Close()
	decoder = json.NewDecoder(f)
	if err = decoder.Decode(&tm); err != nil {
		return
	}
	for i, ts := range tm.Tilesets {
		tspath := filepath.Join(filepath.Dir(path), ts.Image)
		if tm.Tilesets[i].Texture, err = LoadTexture(tspath, IntNearest, ts.Tilewidth, ts.Tileheight); err != nil {
			return
		}
		// The following ignores spacing, but I don't use it.
		tm.Tilesets[i].Tilecount = (ts.Imagewidth / ts.Tilewidth) * (ts.Imageheight / ts.Tileheight)
		tm.Tilesets[i].Lastgid = ts.Firstgid + tm.Tilesets[i].Tilecount
	}
	out = &tm
	return
}
