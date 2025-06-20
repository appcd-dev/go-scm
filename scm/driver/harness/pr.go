// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package harness

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/drone/go-scm/scm/driver/internal/null"

	"github.com/drone/go-scm/scm"
)

type pullService struct {
	client *wrapper
}

func (s *pullService) Find(ctx context.Context, repo string, index int) (*scm.PullRequest, *scm.Response, error) {
	harnessURI := buildHarnessURI(s.client.account, s.client.organization, s.client.project, repo)
	repoId, queryParams, err := getRepoAndQueryParams(harnessURI)
	if err != nil {
		return nil, nil, err
	}
	path := fmt.Sprintf("api/v1/repos/%s/pullreq/%d?%s", repoId, index, queryParams)
	out := new(pr)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertPullRequest(out), res, err

}

func (s *pullService) FindComment(context.Context, string, int, int) (*scm.Comment, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *pullService) List(ctx context.Context, repo string, opts scm.PullRequestListOptions) ([]*scm.PullRequest, *scm.Response, error) {
	harnessURI := buildHarnessURI(s.client.account, s.client.organization, s.client.project, repo)
	repoId, queryParams, err := getRepoAndQueryParams(harnessURI)
	if err != nil {
		return nil, nil, err
	}
	path := fmt.Sprintf("api/v1/repos/%s/pullreq?%s&%s", repoId, encodePullRequestListOptions(opts), queryParams)
	out := []*pr{}
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertPullRequestList(out), res, err
}

func (s *pullService) ListComments(context.Context, string, int, scm.ListOptions) ([]*scm.Comment, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *pullService) ListCommits(ctx context.Context, repo string, index int, opts scm.ListOptions) ([]*scm.Commit, *scm.Response, error) {
	harnessURI := buildHarnessURI(s.client.account, s.client.organization, s.client.project, repo)
	repoId, queryParams, err := getRepoAndQueryParams(harnessURI)
	if err != nil {
		return nil, nil, err
	}
	path := fmt.Sprintf("api/v1/repos/%s/pullreq/%d/commits?%s&%s", repoId, index, encodeListOptions(opts), queryParams)
	out := []*commit{}
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertCommits(out), res, err
}

func (s *pullService) ListChanges(ctx context.Context, repo string, number int, _ scm.ListOptions) ([]*scm.Change, *scm.Response, error) {
	harnessURI := buildHarnessURI(s.client.account, s.client.organization, s.client.project, repo)
	repoId, queryParams, err := getRepoAndQueryParams(harnessURI)
	if err != nil {
		return nil, nil, err
	}
	path := fmt.Sprintf("api/v1/repos/%s/pullreq/%d/diff?%s", repoId, number, queryParams)
	out := []*fileDiff{}
	res, err := s.client.do(ctx, "POST", path, nil, &out)
	return convertFileDiffs(out), res, err
}

func (s *pullService) Create(ctx context.Context, repo string, input *scm.PullRequestInput) (*scm.PullRequest, *scm.Response, error) {
	harnessURI := buildHarnessURI(s.client.account, s.client.organization, s.client.project, repo)
	repoId, queryParams, err := getRepoAndQueryParams(harnessURI)
	if err != nil {
		return nil, nil, err
	}
	path := fmt.Sprintf("api/v1/repos/%s/pullreq?%s", repoId, queryParams)
	in := &prInput{
		Title:        input.Title,
		Description:  input.Body,
		SourceBranch: input.Source,
		TargetBranch: input.Target,
	}
	out := new(pr)
	res, err := s.client.do(ctx, "POST", path, in, out)
	return convertPullRequest(out), res, err
}

func (s *pullService) CreateComment(ctx context.Context, repo string, prNumber int, input *scm.CommentInput) (*scm.Comment, *scm.Response, error) {
	harnessURI := buildHarnessURI(s.client.account, s.client.organization, s.client.project, repo)
	repoId, queryParams, err := getRepoAndQueryParams(harnessURI)
	if err != nil {
		return nil, nil, err
	}
	path := fmt.Sprintf("api/v1/repos/%s/pullreq/%d/comments?%s", repoId, prNumber, queryParams)
	in := &prComment{
		Text: input.Body,
	}
	out := new(prCommentResponse)
	res, err := s.client.do(ctx, "POST", path, in, out)
	return convertComment(out), res, err
}

func (s *pullService) DeleteComment(context.Context, string, int, int) (*scm.Response, error) {
	return nil, scm.ErrNotSupported
}

func (s *pullService) Merge(ctx context.Context, repo string, index int) (*scm.Response, error) {
	return nil, scm.ErrNotSupported
}

func (s *pullService) Close(context.Context, string, int) (*scm.Response, error) {
	return nil, scm.ErrNotSupported
}

func (s *pullService) Update(ctx context.Context, repo string, number int, input *scm.PullRequestInput) (*scm.PullRequest, *scm.Response, error) {
	harnessURI := buildHarnessURI(s.client.account, s.client.organization, s.client.project, repo)
	repoId, queryParams, err := getRepoAndQueryParams(harnessURI)
	if err != nil {
		return nil, nil, err
	}
	path := fmt.Sprintf("api/v1/repos/%s/pullreq/%d?%s", repoId, number, queryParams)
	in := &prInput{}
	if input.Title != "" {
		in.Title = input.Title
	}
	if input.Body != "" {
		in.Description = input.Body
	}
	if input.Target != "" {
		in.TargetBranch = input.Target
	}
	out := new(pr)
	res, err := s.client.do(ctx, "PUT", path, in, out)
	return convertPullRequest(out), res, err
}

// native data structures
type (
	pr struct {
		Author      principal `json:"author"`
		Created     int64     `json:"created"`
		Description string    `json:"description"`
		Edited      int64     `json:"edited"`
		IsDraft     bool      `json:"is_draft"`

		MergeTargetSHA   null.String `json:"merge_target_sha"`
		MergeBaseSha     string      `json:"merge_base_sha"`
		Merged           null.Int    `json:"merged"`
		MergeMethod      null.String `json:"merge_method"`
		MergeSHA         null.String `json:"merge_sha"`
		MergeCheckStatus string      `json:"merge_check_status"`
		MergeConflicts   []string    `json:"merge_conflicts,omitempty"`
		Merger           *principal  `json:"merger"`

		Number int64 `json:"number"`

		SourceBranch string `json:"source_branch"`
		SourceRepoID int64  `json:"source_repo_id"`
		SourceSHA    string `json:"source_sha"`
		TargetBranch string `json:"target_branch"`
		TargetRepoID int64  `json:"target_repo_id"`

		State string `json:"state"`
		Stats struct {
			Commits         null.Int `json:"commits,omitempty"`
			Conversations   int      `json:"conversations,omitempty"`
			FilesChanged    null.Int `json:"files_changed,omitempty"`
			UnresolvedCount int      `json:"unresolved_count,omitempty"`
		} `json:"stats"`

		Title string `json:"title"`
	}

	reference struct {
		Repo repository `json:"repo"`
		Name string     `json:"ref"`
		Sha  string     `json:"sha"`
	}

	prInput struct {
		Description   string `json:"description"`
		IsDraft       bool   `json:"is_draft"`
		SourceBranch  string `json:"source_branch"`
		SourceRepoRef string `json:"source_repo_ref"`
		TargetBranch  string `json:"target_branch"`
		Title         string `json:"title"`
	}

	commit struct {
		Author struct {
			Identity struct {
				Email string `json:"email"`
				Name  string `json:"name"`
			} `json:"identity"`
			When time.Time `json:"when"`
		} `json:"author"`
		Committer struct {
			Identity struct {
				Email string `json:"email"`
				Name  string `json:"name"`
			} `json:"identity"`
			When time.Time `json:"when"`
		} `json:"committer"`
		Message string `json:"message"`
		Sha     string `json:"sha"`
		Title   string `json:"title"`
	}
	prComment struct {
		LineEnd         int    `json:"line_end"`
		LineEndNew      bool   `json:"line_end_new"`
		LineStart       int    `json:"line_start"`
		LineStartNew    bool   `json:"line_start_new"`
		ParentID        int    `json:"parent_id"`
		Path            string `json:"path"`
		SourceCommitSha string `json:"source_commit_sha"`
		TargetCommitSha string `json:"target_commit_sha"`
		Text            string `json:"text"`
	}
	prCommentResponse struct {
		Id        int         `json:"id"`
		Created   int64       `json:"created"`
		Updated   int64       `json:"updated"`
		Edited    int64       `json:"edited"`
		ParentId  interface{} `json:"parent_id"`
		RepoId    int         `json:"repo_id"`
		PullreqId int         `json:"pullreq_id"`
		Order     int         `json:"order"`
		SubOrder  int         `json:"sub_order"`
		Type      string      `json:"type"`
		Kind      string      `json:"kind"`
		Text      string      `json:"text"`
		Payload   struct{}    `json:"payload"`
		Metadata  interface{} `json:"metadata"`
		Author    struct {
			Id          int    `json:"id"`
			Uid         string `json:"uid"`
			DisplayName string `json:"display_name"`
			Email       string `json:"email"`
			Type        string `json:"type"`
			Created     int64  `json:"created"`
			Updated     int64  `json:"updated"`
		} `json:"author"`
	}
)

// native data structure conversion
func convertPullRequests(src []*pr) []*scm.PullRequest {
	dst := []*scm.PullRequest{}
	for _, v := range src {
		dst = append(dst, convertPullRequest(v))
	}
	return dst
}

func convertPullRequest(src *pr) *scm.PullRequest {
	return &scm.PullRequest{
		Number: int(src.Number),
		Title:  src.Title,
		Body:   src.Description,
		Sha:    src.SourceSHA,
		Source: src.SourceBranch,
		Target: src.TargetBranch,
		Merged: src.Merged.Valid,
		Author: scm.User{
			Login: src.Author.Email,
			Name:  src.Author.DisplayName,
			ID:    src.Author.UID,
			Email: src.Author.Email,
		},
		Head: scm.Reference{
			Name: src.SourceBranch,
			Path: scm.ExpandRef(src.SourceBranch, "refs/heads"),
			Sha:  src.SourceSHA,
		},
		Base: scm.Reference{
			Name: src.TargetBranch,
			Path: scm.ExpandRef(src.TargetBranch, "refs/heads"),
			Sha:  src.MergeTargetSHA.String,
		},
		Fork:    "fork",
		Ref:     fmt.Sprintf("refs/pullreq/%d/head", src.Number),
		Closed:  src.State == "closed",
		Created: time.UnixMilli(src.Created),
		Updated: time.UnixMilli(src.Edited),
	}
}

func convertCommits(src []*commit) []*scm.Commit {
	dst := []*scm.Commit{}
	for _, v := range src {
		dst = append(dst, convertCommit(v))
	}
	return dst
}

func convertCommit(src *commit) *scm.Commit {
	return &scm.Commit{
		Message: src.Message,
		Sha:     src.Sha,
		Author: scm.Signature{
			Name:  src.Author.Identity.Name,
			Email: src.Author.Identity.Email,
		},
		Committer: scm.Signature{
			Name:  src.Committer.Identity.Name,
			Email: src.Committer.Identity.Email,
		},
	}
}

func convertFileDiffs(diff []*fileDiff) []*scm.Change {
	var dst []*scm.Change
	for _, v := range diff {
		dst = append(dst, convertFileDiff(v))
	}
	return dst
}

func convertFileDiff(diff *fileDiff) *scm.Change {
	return &scm.Change{
		Path:         diff.Path,
		Added:        strings.EqualFold(diff.Status, "ADDED"),
		Renamed:      strings.EqualFold(diff.Status, "RENAMED"),
		Deleted:      strings.EqualFold(diff.Status, "DELETED"),
		Sha:          diff.SHA,
		BlobID:       "",
		PrevFilePath: diff.OldPath,
	}
}

func convertPullRequestList(from []*pr) []*scm.PullRequest {
	to := []*scm.PullRequest{}
	for _, v := range from {
		to = append(to, convertPullRequest(v))
	}
	return to
}

func convertComment(comment *prCommentResponse) *scm.Comment {
	return &scm.Comment{
		ID:   comment.Id,
		Body: comment.Text,
		Author: scm.User{
			Login:   comment.Author.Uid,
			Name:    comment.Author.DisplayName,
			ID:      strconv.Itoa(comment.Author.Id),
			Email:   comment.Author.Email,
			Created: time.UnixMilli(comment.Author.Created),
			Updated: time.UnixMilli(comment.Author.Updated),
		},
		Created: time.UnixMilli(comment.Created),
		Updated: time.UnixMilli(comment.Updated),
	}
}
