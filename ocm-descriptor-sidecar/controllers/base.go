/*
  OCM-DESCRIPTOR-SIDECAR
  Copyright © 2022-2024 EVIDEN

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.

  This work has received funding from the European Union's HORIZON research
  and innovation programme under grant agreement No. 101070177.
*/

package controllers

import (
	"icos/server/ocm-descriptor-sidecar/utils/logs"
	"time"

	"github.com/gorilla/mux"
)

type Server struct {
	Router *mux.Router
}

func (server *Server) Init() {
	// server.Router = mux.NewRouter()
	// server.initializeRoutes()
}

func (server *Server) Run() {
	logs.Logger.Println("Starting to Schedule")
	ticker := time.NewTicker(15 * time.Second) // TODO parametrize
	quit := make(chan struct{})
	// go func() {
	for {
		select {
		case <-ticker.C:
			// do stuff
			status, err := Schedule()
			if err != nil {
				logs.Logger.Println("ERROR " + err.Error())
			} else {
				logs.Logger.Println("Status of the execution: " + status)
			}
		case <-quit:
			ticker.Stop()
			return
		}
	}
	// }()
	// after stopping server
	// logs.Logger.Println("Closing connections ...")

	// var shutdownTimeout = flag.Duration("shutdown-timeout", 10*time.Second, "shutdown timeout (5s,5m,5h) before connections are cancelled")
	// _, cancel := context.WithTimeout(context.Background(), *shutdownTimeout)
	// defer cancel()
}
