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

type Drawable interface {
	Y() float64
	X() float64
	GetFrame() int
	FlipX() bool
	TextureRow() int
}

type Drawables []Drawable

func (s Drawables) Len() int      { return len(s) }
func (s Drawables) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type ByY struct{ Drawables }

func (s ByY) Less(i, j int) bool {
	if s.Drawables[i].Y() == s.Drawables[j].Y() {
		return s.Drawables[i].X() < s.Drawables[j].X()
	}
	return s.Drawables[i].Y() < s.Drawables[j].Y()
}
