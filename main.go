// Copyright 2018 eel Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"github.com/gingerxman/eel"
	"github.com/gingerxman/eel/cron"
	"github.com/gingerxman/ginger-finance/business/user"
	_ "github.com/gingerxman/ginger-finance/cron"
	_ "github.com/gingerxman/ginger-finance/middleware"
	_ "github.com/gingerxman/ginger-finance/models"
	_ "github.com/gingerxman/ginger-finance/routers"
)

func main() {
	if eel.GetServiceMode() == eel.SERVICE_MODE_CRON{
		endRunning := make(chan bool, 1)
		cron.StartCronTasks()
		defer cron.StopCronTasks()
		<- endRunning
	}else{
		eel.Runtime.NewBusinessContext = user.NewBusinessContext
		eel.RunService()
	}
}

