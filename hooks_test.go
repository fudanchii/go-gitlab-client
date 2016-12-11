package gogitlab

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHook(t *testing.T) {
	ts, gitlab := Stub("stubs/hooks/show.json")
	hook, err := gitlab.ProjectHook("1", 2)

	assert.Equal(t, err, nil)
	assert.IsType(t, new(Hook), hook)
	assert.Equal(t, hook.URL, "http://example.com/hook")
	defer ts.Close()
}

func TestParsePushHook(t *testing.T) {
	stub, _ := ioutil.ReadFile("stubs/hooks/push.json")
	p, err := ParseHook([]byte(stub))

	assert.Equal(t, err, nil)
	assert.IsType(t, new(HookPayload), p)
	assert.Equal(t, p.After, "da1560886d4f094c3e6c9ef40349f7d38b5d27d7")
	assert.Equal(t, p.Repository.URL, "git@localhost:diaspora.git")
	assert.Equal(t, len(p.Commits), 2)
	assert.Equal(t, p.Commits[0].Author.Email, "jordi@softcatala.org")
	assert.Equal(t, p.Commits[1].Id, "da1560886d4f094c3e6c9ef40349f7d38b5d27d7")
	assert.Equal(t, p.Branch(), "master")
	assert.Equal(t, p.Head().Message, "fixed readme")
}

func TestParseIssueHook(t *testing.T) {
	stub, _ := ioutil.ReadFile("stubs/hooks/issue.json")
	p, err := ParseHook([]byte(stub))

	assert.Equal(t, err, nil)
	assert.Equal(t, p.ObjectKind, "issue")
	assert.Equal(t, p.ObjectAttributes.Id, 301)
}

func TestParseMergeRequestHook(t *testing.T) {
	stub, _ := ioutil.ReadFile("stubs/hooks/merge_request.json")
	p, err := ParseHook([]byte(stub))

	assert.Equal(t, err, nil)
	assert.Equal(t, p.ObjectKind, "merge_request")
	assert.Equal(t, p.ObjectAttributes.TargetBranch, "master")
	assert.Equal(t, p.ObjectAttributes.SourceProjectId, p.ObjectAttributes.TargetProjectId)
}

func TestBuildHookQuery(t *testing.T) {
	t_ := true
	f_ := false
	h := Hook{
		ID:                    123214,
		URL:                   "https://ci.example/hooks/test",
		CreatedAt:             "",
		PushEvents:            &t_,
		IssuesEvents:          &t_,
		MergeRequestsEvents:   &t_,
		NoteEvents:            &f_,
		BuildEvents:           &f_,
		EnableSSLVerification: &t_,
	}
	q := buildHookQuery(&h)
	assert.Equal(t, "", q)
}
