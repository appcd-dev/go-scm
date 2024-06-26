// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gitea

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/drone/go-scm/scm"

	"github.com/google/go-cmp/cmp"
	"github.com/h2non/gock"
)

func TestOrgFind(t *testing.T) {
	defer gock.Off()

	gock.New("https://try.gitea.io").
		Get("/api/v1/orgs/gogits").
		Reply(200).
		Type("application/json").
		File("testdata/organization.json")

	client, _ := New("https://try.gitea.io")
	got, _, err := client.Organizations.Find(context.Background(), "gogits")
	if err != nil {
		t.Error(err)
	}

	want := new(scm.Organization)
	raw, _ := os.ReadFile("testdata/organization.json.golden")
	json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestOrganizationFindMembership(t *testing.T) {
	defer gock.Off()

	gock.New("https://try.gitea.io").
		Get("/api/v1/orgs/gogits/members/jcitizen").
		Reply(204)

	gock.New("https://try.gitea.io").
		Get("/api/v1/users/jcitizen/orgs/gogits/permissions").
		Reply(200).
		Type("application/json").
		File("testdata/permissions.json")

	client, _ := New("https://try.gitea.io")
	got, _, err := client.Organizations.FindMembership(context.Background(), "gogits", "jcitizen")
	if err != nil {
		t.Error(err)
	}

	want := new(scm.Membership)
	raw, _ := os.ReadFile("testdata/membership.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestOrgList(t *testing.T) {
	defer gock.Off()

	gock.New("https://try.gitea.io").
		Get("/api/v1/user/orgs").
		Reply(200).
		Type("application/json").
		File("testdata/organizations.json")

	client, _ := New("https://try.gitea.io")
	got, _, err := client.Organizations.List(context.Background(), scm.ListOptions{})
	if err != nil {
		t.Error(err)
	}

	want := []*scm.Organization{}
	raw, _ := os.ReadFile("testdata/organizations.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}
