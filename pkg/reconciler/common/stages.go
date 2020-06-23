/*
Copyright 2020 The Knative Authors

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

package common

import (
	"context"

	mf "github.com/manifestival/manifestival"
	"knative.dev/operator/pkg/apis/operator/v1alpha1"
	"knative.dev/pkg/logging"
)

// Stage represents a step in the reconcile process
type Stage func(context.Context, *mf.Manifest, v1alpha1.KComponent) error

// Stages are a list of steps
type Stages []Stage

// Execute each stage in sequence until one returns an error
func (stages Stages) Execute(ctx context.Context, manifest *mf.Manifest, instance v1alpha1.KComponent) error {
	for _, stage := range stages {
		if err := stage(ctx, manifest, instance); err != nil {
			return err
		}
	}
	return nil
}

// NoOp does nothing
func NoOp(context.Context, *mf.Manifest, v1alpha1.KComponent) error {
	return nil
}

// AppendTarget mutates the passed manifest by appending one
// appropriate for the passed KComponent
func AppendTarget(ctx context.Context, manifest *mf.Manifest, instance v1alpha1.KComponent) error {
	m, err := TargetManifest(instance)
	if err != nil {
		return err
	}
	*manifest = manifest.Append(m)
	return nil
}

// AppendInstalled mutates the passed manifest by appending one
// appropriate for the passed KComponent, which may not be the one
// corresponding to status.version
func AppendInstalled(ctx context.Context, manifest *mf.Manifest, instance v1alpha1.KComponent) error {
	logger := logging.FromContext(ctx)
	m, err := InstalledManifest(instance)
	if err != nil {
		// TODO: return the oldest instead of the latest?
		logger.Error("Unable to fetch installed manifest, trying target", err)
		m, err = TargetManifest(instance)
	}
	if err != nil {
		return err
	}
	*manifest = manifest.Append(m)
	return nil
}
