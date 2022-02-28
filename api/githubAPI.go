package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
)

var (
	GithubApiToken   string
	ConfiguredLabels []*NewLabel
)

func GetLabels(owner, repo string) ([]Label, error) {
	rawJson, err := requestGithub(fmt.Sprintf("repos/%v/%v/labels", owner, repo))
	if err != nil {
		return nil, err
	}
	return decodeMultiLabels(rawJson)
}

func UpdateLabel(repo, owner, labelName string, l *NewLabel) (Label, error) {
	rawJson, err := postLabel(fmt.Sprintf("%v/%v/labels/%v", owner, repo, labelName), l)
	if err != nil {
		return Label{ID: nil}, err
	}
	return decodeSingleLabel(rawJson)
}

func postLabel(path string, l *NewLabel) (json.RawMessage, error) {
	labelData := NewLabel{
		//Name:        l.Name,
		NewName:     l.NewName,
		Color:       l.Color,
		Description: l.Description,
	}
	body, err := json.Marshal(labelData)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("https://api.github.com/repos/%v", path), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("token %v", GithubApiToken))
	req.Header.Set("accept", "application/vnd.github.v3+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func decodeSingleLabel(rawJson json.RawMessage) (Label, error) {
	var l Label
	err := json.Unmarshal(rawJson, &l)
	return l, err
}

func decodeMultiLabels(rawJson json.RawMessage) ([]Label, error) {
	var l []Label
	err := json.Unmarshal(rawJson, &l)
	return l, err
}

func GetUser(logon string) (User, error) {
	rawJson, err := requestGithub(fmt.Sprintf("users/%v", logon))
	if err != nil {
		return User{ID: nil}, err
	}
	return decodeSingleUser(rawJson)
}

func GetRepoCollaborators(repo string) ([]User, error) {
	rawJson, err := requestGithub(fmt.Sprintf("repos/Jacobbrewer1/%v/collaborators", repo))
	if err != nil {
		return nil, err
	}
	return decodeMultiUser(rawJson)
}

func GetBranches(repo string) ([]Branch, error) {
	rawJson, err := requestGithub(fmt.Sprintf("repos/Jacobbrewer1/%v/branches", repo))
	if err != nil {
		return nil, err
	}
	return decodeMultiBranch(rawJson)
}

func decodeSingleBranch(rawJson json.RawMessage) (Branch, error) {
	var b Branch
	err := json.Unmarshal(rawJson, &b)
	return b, err
}

func decodeMultiBranch(rawJson json.RawMessage) ([]Branch, error) {
	var b []Branch
	err := json.Unmarshal(rawJson, &b)
	return b, err
}

func getRandomRepoCollaborator(repo string) (User, error) {
	rawJson, err := requestGithub(fmt.Sprintf("repos/Jacobbrewer1/%v/collaborators", repo))
	if err != nil {
		return User{ID: nil}, err
	}
	users, err := decodeMultiUser(rawJson)
	return users[rand.Intn(len(users))], err
}

func decodeSingleUser(rawJson json.RawMessage) (User, error) {
	var u User
	err := json.Unmarshal(rawJson, &u)
	return u, err
}

func decodeMultiUser(rawJson json.RawMessage) ([]User, error) {
	var u []User
	err := json.Unmarshal(rawJson, &u)
	return u, err
}

func removePullRequests(i []Issue) ([]Issue, error) {
	if len(i) == 1 && i[0].IsPullRequest() {
		log.Println("Only issue found is a pr")
		return nil, nil
	}
	for z, x := range i {
		if x.IsPullRequest() {
			//log.Println("pr found in issues array")
			if z >= len(i) {
				break
			}
			i = removeIssueSlice(i, z)
		}
	}
	return i, nil
}

func decodeMultiIssues(jsonRaw json.RawMessage) ([]Issue, error) {
	var issues []Issue
	err := json.Unmarshal(jsonRaw, &issues)
	return issues, err
}

func decodeSingleIssue(jsonRaw json.RawMessage) (Issue, error) {
	var i Issue
	err := json.Unmarshal(jsonRaw, &i)
	return i, err
}

func removeIssueSlice(i []Issue, x int) []Issue {
	if x < 0 {
		return i
	}
	if i[len(i)-1].IsPullRequest() {
		for j := len(i) - 2; j >= x; {
			if !i[j].IsPullRequest() {
				i[x] = i[j]
				return i[:j]
			}
			j--
		}
		return i[:x]
	}
	i[x] = i[len(i)-1]
	return i[:len(i)-1]
}

func GetRepos() ([]Repository, error) {
	jsonRaw, err := requestGithub("user/repos")
	if err != nil {
		return nil, err
	}
	return decodeRepos(jsonRaw)
}

func decodeRepos(jsonRaw json.RawMessage) ([]Repository, error) {
	var repositories []Repository
	err := json.Unmarshal(jsonRaw, &repositories)
	return repositories, err
}

func requestGithub(endpoint string) (json.RawMessage, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.github.com/%v", endpoint), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("token %v", GithubApiToken))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
