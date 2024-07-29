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
	"fmt"
	"icos/server/ocm-descriptor-sidecar/models"
	"icos/server/ocm-descriptor-sidecar/utils/logs"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
)

var (
	deployManagerURL   = os.Getenv("DEPLOY_MANAGER_URL")
	lighthouseBaseURL  = os.Getenv("LIGHTHOUSE_BASE_URL")
	apiV3              = "/api/v3"
	matchmackerBaseURL = os.Getenv("MATCHMAKING_URL")
)

func Schedule() (execStatus string, err error) {
	logs.Logger.Println("Scheduling Started")

	// ------------------------- trigger the execution of the jobs -------------------------
	reqExecution, err := http.NewRequest("GET", deployManagerURL+"/execute", http.NoBody)
	if err != nil {
		logs.Logger.Println("ERROR " + err.Error())
		return
	}
	// get token from keycloak
	requester := models.KeycloakTokenRequester{}
	token, err := models.FetchKeycloakToken(requester)

	reqExecution.Header.Add("Authorization", "Bearer "+token.AccessToken)
	// debug
	logs.Logger.Println("Execution Request to send to Deployment Manager: ")
	b, err := httputil.DumpRequest(reqExecution, true)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(b))
	// do request
	client := &http.Client{}
	respExecution, err := client.Do(reqExecution)
	if err != nil {
		logs.Logger.Println("ERROR " + err.Error())
		return
	}
	defer respExecution.Body.Close()

	logs.Logger.Println("Execution Response Info: ")
	b2, err := httputil.DumpResponse(respExecution, true)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(b2))

	// ------------------------- trigger the sync of the resources -------------------------
	// update status of all deployed resources into JM periodically
	reqSync, err := http.NewRequest("GET", deployManagerURL+"/resource/sync", http.NoBody)
	if err != nil {
		logs.Logger.Println("ERROR " + err.Error())
		return
	}
	reqSync.Header.Add("Authorization", "Bearer "+token.AccessToken)
	// debug
	logs.Logger.Println("Sync Request to send to Deployment Manager: ")
	b3, err := httputil.DumpRequest(reqSync, true)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(b3))
	// do request
	respSync, err := client.Do(reqSync)
	if err != nil {
		logs.Logger.Println("ERROR " + err.Error())
		return
	}
	defer respExecution.Body.Close()

	logs.Logger.Println("Sync Response Info: ")
	b4, err := httputil.DumpResponse(respSync, true)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(b4))

	return respExecution.Status, err

}
