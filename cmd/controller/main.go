/*
Copyright 2019 The Knative Authors

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
	"context"
	"log"
	"net/http"

	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/injection/sharedmain"
	"knative.dev/pkg/signals"

	"github.com/mattmoor/http01-solver/pkg/challenger"
	"github.com/mattmoor/http01-solver/pkg/reconciler/certificate"
)

func main() {
	// Uncomment this to use the Let's Encrypt Staging environment.
	// ordermanager.Endpoint = ordermanager.Staging

	ctx := signals.NewContext()

	chlr, err := challenger.New(ctx)
	if err != nil {
		log.Fatalf("Error creating challenger: %v", err)
	}

	go http.ListenAndServe(":8080", chlr)

	sharedmain.MainWithContext(
		ctx,
		"controller",
		func(ctx context.Context, cmw configmap.Watcher) *controller.Impl {
			return certificate.NewController(ctx, cmw, chlr)
		},
	)
}
