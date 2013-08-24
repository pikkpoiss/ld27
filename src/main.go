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
	"flag"
	"log"
	"runtime"
)

func init() {
	// See https://code.google.com/p/go/issues/detail?id=3527
	runtime.LockOSThread()
}

func main() {
	var (
		err  error
		win  *system.Window
		ctrl *system.Controller
		//snd  *system.Sound
		game *Game
	)
	flag.Parse()
	if ctrl, err = system.NewController(); err != nil {
		log.Fatalf("Couldn't init Controller: %v\n", err)
	}
	defer ctrl.Terminate()
	/*
	if snd, err = system.NewSound(); err != nil {
		log.Fatalf("Couldn't init Sound: %v\n", err)
	}
	defer snd.Terminate()
	*/
	win = &system.Window{Width: 1136, Height: 640, Resize: false}
	if err = ctrl.Open(win); err != nil {
		log.Fatalf("Couldn't open Window: %v\n", err)
	}
	if game, err = NewGame(ctrl); err != nil {
		log.Fatalf("Couldn't start Game: %v\n", err)
	}
	game.Run()
	log.Printf("Exiting peacefully")
	log.Printf("%v", win)
}
