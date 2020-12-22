/*
Copyright 2020 The Kubernetes Authors.

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

package app

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

func TestValidateAuthFlow(t *testing.T) {
	type FlagResult struct {
		Name  string
		Flow  string
		Error error
	}
	tests := []FlagResult{
		{Name: "validate gcr auth flow", Flow: gcrAuthFlow},
		{Name: "validate docker-cfg auth flow option", Flow: dockerConfigAuthFlow},
		{Name: "validate docker-cfg-url auth flow option", Flow: dockerConfigURLAuthFlow},
		{Name: "bad auth flow option", Flow: "bad-flow", Error: &AuthFlowFlagError{flagValue: "bad-flow"}},
		{Name: "empty auth flow option", Flow: "", Error: &AuthFlowFlagError{flagValue: ""}},
		{Name: "case-sensitive auth flow", Flow: "Gcrauthflow", Error: &AuthFlowFlagError{flagValue: "Gcrauthflow"}},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			err := validateFlags(&CredentialOptions{AuthFlow: tc.Flow})
			if tc.Error != nil {
				if err == nil {
					t.Fatalf("with flow %q did not get expected error %q", tc.Flow, err)
				}
				if !errors.Is(err, tc.Error) {
					t.Fatalf("with flow %q got unexpected error type %q (expected %q)", tc.Flow, reflect.TypeOf(err), reflect.TypeOf(tc.Error))
				}
				return
			}
			if err != nil {
				t.Fatalf("with flow %q unexpected error %q", tc.Flow, err)
			}
		})
	}
}

func TestProviderFromFlow(t *testing.T) {
	type ProviderResult struct {
		Name  string
		Flow  string
		Type  string
		Error error
	}
	tests := []ProviderResult{
		{Name: "gcr auth provider selection", Flow: gcrAuthFlow, Type: "ContainerRegistryProvider"},
		{Name: "docker-cfg auth provider selection", Flow: dockerConfigAuthFlow, Type: "DockerConfigKeyProvider"},
		{Name: "docker-cfg-url auth provider selection", Flow: dockerConfigURLAuthFlow, Type: "DockerConfigURLKeyProvider"},
		{Name: "non-existent auth provider request", Flow: "bad-flow", Type: "", Error: &AuthFlowTypeError{requestedFlow: "bad-flow"}},
		{Name: "empty auth provider request", Flow: "", Type: "", Error: &AuthFlowTypeError{requestedFlow: ""}},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			provider, err := providerFromFlow(tc.Flow)
			if tc.Error != nil {
				if err == nil {
					t.Fatalf("with flow %q did not get expected error %q", tc.Flow, err)
				}
				if !errors.Is(err, tc.Error) {
					t.Fatalf("with flow %q got unexpected error type %q (expected %q)", tc.Flow, reflect.TypeOf(err), reflect.TypeOf(tc.Error))
				}
				return
			}
			if err != nil {
				t.Fatalf("with flow %q unexpected error %q", tc.Flow, err)
			}
			providerType := reflect.TypeOf(provider).String()
			if providerType != "*gcpcredential."+tc.Type {
				t.Errorf("with flow %q unexpected provider type %q", tc.Flow, providerType)
			}
		})
	}
}

func TestFlagError(t *testing.T) {
	badFlow := "bad-flow"
	err := validateFlags(&CredentialOptions{AuthFlow: badFlow})
	differentFlagValueErr := &AuthFlowFlagError{flagValue: "other-bad-flow"}
	if !errors.Is(err, differentFlagValueErr) {
		t.Fatalf("errors.Is should return true for different flagValues")
	}
	if !strings.Contains(err.Error(), badFlow) {
		t.Fatalf("Flow %q missing from error message of AuthFlowFlagError", badFlow)
	}
}

func TestFlowError(t *testing.T) {
	badProviderRequest := "bad-provider"
	_, err := providerFromFlow(badProviderRequest)
	differentFlagValueErr := &AuthFlowTypeError{requestedFlow: "other-bad-provider"}
	if !errors.Is(err, differentFlagValueErr) {
		t.Fatalf("errors.Is should return true for different requestedFlows")
	}
	if !strings.Contains(err.Error(), badProviderRequest) {
		t.Fatalf("Flow %q missing from error message of AuthFlowTypeError", badProviderRequest)
	}
}
