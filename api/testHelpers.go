package api

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"strings"
)

const (
	configPath = "../config"
)

func testerSetupIssue(isPullRequest, assigned, assignedMulti bool) Issue {
	var i Issue
	if isPullRequest {
		var prLinkStr = "https://www.IsPr.com"
		var prLinks = PullRequestLinks{
			URL:      &prLinkStr,
			HTMLURL:  &prLinkStr,
			DiffURL:  &prLinkStr,
			PatchURL: &prLinkStr,
		}
		i.PullRequestLinks = &prLinks
	}
	if assigned {
		var a = "testUser"
		var assignee = User{
			Login: &a,
		}
		i.Assignee = &assignee
	}
	if assignedMulti {
		tt1 := "tt1"
		tt2 := "tt2"
		var a1 = User{Login: &tt1}
		var a2 = User{Login: &tt2}
		var assignees []*User
		assignees = append(assignees, &a1)
		assignees = append(assignees, &a2)
		i.Assignees = assignees
	}
	return i
}

func testerSetupBlankIssueAssigneesArray() Issue {
	x := make([]*User, 0)
	return Issue{Assignees: x}
}

func testerSetupIssueArray() []Issue {
	var i []Issue
	var noPrString = "https://www.notPr.com"
	log.Println("--testerSetupIssueArray()--")
	npr := 0
	for x := 0; x < 10; {
		var h = Issue{URL: &noPrString}
		t := int64(x)
		h.ID = &t
		i = append(i, h)
		x++
		npr = x
	}
	log.Println(npr)
	var prLinkStr = "https://www.IsPr.com"
	var prLinks = PullRequestLinks{
		URL:      &prLinkStr,
		HTMLURL:  &prLinkStr,
		DiffURL:  &prLinkStr,
		PatchURL: &prLinkStr,
	}
	pr := 0
	for x := 0; x < 10; {
		var h = Issue{PullRequestLinks: &prLinks}
		t := int64(x)
		h.ID = &t
		i = append(i, h)
		x++
		pr = x
	}
	log.Println(pr)
	log.Printf("len(i): %v", len(i))
	return i
}

func testerMixedUpIssuePr() []Issue {
	var i []Issue
	var noPrString = "https://www.notPr.com"
	var prLinkStr = "https://www.IsPr.com"
	var prLinks = PullRequestLinks{
		URL:      &prLinkStr,
		HTMLURL:  &prLinkStr,
		DiffURL:  &prLinkStr,
		PatchURL: &prLinkStr,
	}
	npr := 0
	pr := 0
	for x := 0; x < 60; {
		var h Issue
		t := int64(x)
		h.ID = &t
		if x%2 == 0 || x%5 == 0 {
			h.PullRequestLinks = &prLinks
			pr++
		} else {
			h.URL = &noPrString
			npr++
		}
		i = append(i, h)
		x++
	}
	log.Println("--testerMixedUpIssuePr()--")
	log.Printf("pr count: %v\n", pr)
	log.Printf("normal count: %v\n", npr)
	log.Printf("len(i): %v", len(i))
	return i
}

func testerEvery3Pr() []Issue {
	var i []Issue
	var noPrString = "https://www.notPr.com"
	var prLinkStr = "https://www.IsPr.com"
	var prLinks = PullRequestLinks{
		URL:      &prLinkStr,
		HTMLURL:  &prLinkStr,
		DiffURL:  &prLinkStr,
		PatchURL: &prLinkStr,
	}
	npr := 0
	pr := 0
	for x := 0; x < 60; {
		var h Issue
		t := int64(x)
		h.ID = &t
		if x%3 == 0 {
			h.PullRequestLinks = &prLinks
			pr++
		} else {
			h.URL = &noPrString
			npr++
		}
		i = append(i, h)
		x++
	}
	log.Println("--testerEvery3Pr()--")
	log.Printf("pr count: %v\n", pr)
	log.Printf("normal count: %v\n", npr)
	log.Printf("len(i): %v", len(i))
	return i
}

func testerSetupJsonRaw(jsonData string) json.RawMessage {
	body, err := ioutil.ReadAll(strings.NewReader(jsonData))
	if err != nil {
		log.Println(err)
	}
	return body
}

var byteTrending io.ReadCloser

const (
	testGithubNewIssueResp = "{\"url\":\"https://api.github.com/repos/Jacobbrewer1/botter/issues/117\",\"repository_url\":\"https://api.github.com/repos/Jacobbrewer1/botter\",\"labels_url\":\"https://api.github.com/repos/Jacobbrewer1/botter/issues/117/labels{/name}\",\"comments_url\":\"https://api.github.com/repos/Jacobbrewer1/botter/issues/117/comments\",\"events_url\":\"https://api.github.com/repos/Jacobbrewer1/botter/issues/117/events\",\"html_url\":\"https://github.com/Jacobbrewer1/botter/issues/117\",\"id\":1146002481,\"node_id\":\"I_kwDOGvSFO85ETpwx\",\"number\":117,\"title\":\"test\",\"user\":{\"login\":\"Jacobbrewer1\",\"id\":75381570,\"node_id\":\"MDQ6VXNlcjc1MzgxNTcw\",\"avatar_url\":\"https://avatars.githubusercontent.com/u/75381570?v=4\",\"gravatar_id\":\"\",\"url\":\"https://api.github.com/users/Jacobbrewer1\",\"html_url\":\"https://github.com/Jacobbrewer1\",\"followers_url\":\"https://api.github.com/users/Jacobbrewer1/followers\",\"following_url\":\"https://api.github.com/users/Jacobbrewer1/following{/other_user}\",\"gists_url\":\"https://api.github.com/users/Jacobbrewer1/gists{/gist_id}\",\"starred_url\":\"https://api.github.com/users/Jacobbrewer1/starred{/owner}{/repo}\",\"subscriptions_url\":\"https://api.github.com/users/Jacobbrewer1/subscriptions\",\"organizations_url\":\"https://api.github.com/users/Jacobbrewer1/orgs\",\"repos_url\":\"https://api.github.com/users/Jacobbrewer1/repos\",\"events_url\":\"https://api.github.com/users/Jacobbrewer1/events{/privacy}\",\"received_events_url\":\"https://api.github.com/users/Jacobbrewer1/received_events\",\"type\":\"User\",\"site_admin\":false},\"labels\":[],\"state\":\"open\",\"locked\":false,\"assignee\":{\"login\":\"Jacobbrewer1\",\"id\":75381570,\"node_id\":\"MDQ6VXNlcjc1MzgxNTcw\",\"avatar_url\":\"https://avatars.githubusercontent.com/u/75381570?v=4\",\"gravatar_id\":\"\",\"url\":\"https://api.github.com/users/Jacobbrewer1\",\"html_url\":\"https://github.com/Jacobbrewer1\",\"followers_url\":\"https://api.github.com/users/Jacobbrewer1/followers\",\"following_url\":\"https://api.github.com/users/Jacobbrewer1/following{/other_user}\",\"gists_url\":\"https://api.github.com/users/Jacobbrewer1/gists{/gist_id}\",\"starred_url\":\"https://api.github.com/users/Jacobbrewer1/starred{/owner}{/repo}\",\"subscriptions_url\":\"https://api.github.com/users/Jacobbrewer1/subscriptions\",\"organizations_url\":\"https://api.github.com/users/Jacobbrewer1/orgs\",\"repos_url\":\"https://api.github.com/users/Jacobbrewer1/repos\",\"events_url\":\"https://api.github.com/users/Jacobbrewer1/events{/privacy}\",\"received_events_url\":\"https://api.github.com/users/Jacobbrewer1/received_events\",\"type\":\"User\",\"site_admin\":false},\"assignees\":[{\"login\":\"Jacobbrewer1\",\"id\":75381570,\"node_id\":\"MDQ6VXNlcjc1MzgxNTcw\",\"avatar_url\":\"https://avatars.githubusercontent.com/u/75381570?v=4\",\"gravatar_id\":\"\",\"url\":\"https://api.github.com/users/Jacobbrewer1\",\"html_url\":\"https://github.com/Jacobbrewer1\",\"followers_url\":\"https://api.github.com/users/Jacobbrewer1/followers\",\"following_url\":\"https://api.github.com/users/Jacobbrewer1/following{/other_user}\",\"gists_url\":\"https://api.github.com/users/Jacobbrewer1/gists{/gist_id}\",\"starred_url\":\"https://api.github.com/users/Jacobbrewer1/starred{/owner}{/repo}\",\"subscriptions_url\":\"https://api.github.com/users/Jacobbrewer1/subscriptions\",\"organizations_url\":\"https://api.github.com/users/Jacobbrewer1/orgs\",\"repos_url\":\"https://api.github.com/users/Jacobbrewer1/repos\",\"events_url\":\"https://api.github.com/users/Jacobbrewer1/events{/privacy}\",\"received_events_url\":\"https://api.github.com/users/Jacobbrewer1/received_events\",\"type\":\"User\",\"site_admin\":false}],\"milestone\":null,\"comments\":0,\"created_at\":\"2022-02-21T16:22:16Z\",\"updated_at\":\"2022-02-21T16:22:16Z\",\"closed_at\":null,\"author_association\":\"OWNER\",\"active_lock_reason\":null,\"body\":\"body\",\"closed_by\":null,\"reactions\":{\"url\":\"https://api.github.com/repos/Jacobbrewer1/botter/issues/117/reactions\",\"total_count\":0,\"+1\":0,\"-1\":0,\"laugh\":0,\"hooray\":0,\"confused\":0,\"heart\":0,\"rocket\":0,\"eyes\":0},\"timeline_url\":\"https://api.github.com/repos/Jacobbrewer1/botter/issues/117/timeline\",\"performed_via_github_app\":null}"
)
