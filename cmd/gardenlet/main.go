// Copyright (c) 2018 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
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
	"context"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/gardener/gardener/cmd/gardenlet/app"
	"github.com/gardener/gardener/pkg/gardenlet/features"
)

func main() {
	features.RegisterFeatureGates()

	if err := exec.Command("which", "openvpn").Run(); err != nil {
		panic("openvpn is not installed or not executable. cannot start gardenlet.")
	}

	// Setup signal handler if running inside a Kubernetes cluster
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 2)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer func() {
		signal.Stop(c)
		cancel()
	}()
	go func() {
		<-c
		cancel()
		<-c
		os.Exit(1)
	}()

	command := app.NewCommandStartGardenlet(ctx)
	if err := command.Execute(); err != nil {
		panic(err)
	}
}
