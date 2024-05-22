// Copyright 2022 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package traverse

import (
	"context"

	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/scmlogger"
	"golang.org/x/sync/errgroup"
)

// Repos returns the full repository list, traversing and
// combining paginated responses if necessary.
func Repos(ctx context.Context, client *scm.Client, additional scm.AdditionalInfo) ([]*scm.Repository, error) {
	list := []*scm.Repository{}
	opts := scm.ListOptions{Size: 100, Meta: additional}
	for {
		result, meta, err := client.Repositories.List(ctx, opts)
		if err != nil {
			return nil, err
		}
		list = addNonNil(list, result)
		opts.Page = meta.Page.Next
		opts.URL = meta.Page.NextURL

		if opts.Page == 0 && opts.URL == "" {
			break
		}
	}
	return list, nil
}

// ReposV2 same as Repos but uses errgroup to fetch repos in parallel
func ReposV2(ctx context.Context, client *scm.Client, opts scm.ListOptions) ([]*scm.Repository, error) {
	list := []*scm.Repository{}

	result, meta, err := client.Repositories.List(ctx, opts)
	if err != nil {
		return nil, err
	}
	list = addNonNil(list, result)
	if meta.Page.Next == 0 && meta.Page.NextURL == "" {
		return list, nil
	}
	maxPage := meta.Page.Last
	if opts.MaxPage != 0 && maxPage > opts.MaxPage {
		maxPage = opts.MaxPage
	}
	errGroup, ectx := errgroup.WithContext(ctx)
	for i := meta.Page.Next; i <= maxPage; i++ {
		opts := scm.ListOptions{Size: opts.Size, Page: i, Meta: opts.Meta}
		errGroup.Go(func() error {
			scmlogger.GetLogger().Log("Checking the page %d", opts.Page)
			result, _, err := client.Repositories.List(ectx, opts)
			if err != nil {
				return err
			}
			list = addNonNil(list, result)
			return nil
		})
	}
	return list, errGroup.Wait()
}

func addNonNil(list []*scm.Repository, result []*scm.Repository) []*scm.Repository {
	for _, src := range result {
		if src != nil {
			list = append(list, src)
		}
	}
	return list
}
