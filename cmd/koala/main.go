/*
Copyright 2019 The Koala Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/golang/glog"
	"github.com/shimcdn/koala/cmd/koala/app"
	"github.com/shimcdn/koala/internal/server"
)

func init() {
	log.SetFlags(log.Flags() | log.Lshortfile | log.Ldate)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	stopChan := server.SetupSignalHandler()
	cmd := app.NewKoalaCommand(stopChan)
	if err := cmd.Execute(); err != nil {
		glog.Fatalln(err.Error())
	}
}
