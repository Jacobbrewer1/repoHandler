package main

import (
	"github.com/Jacobbrewer1/repoHandler/api"
	"log"
	"sync"
)

var (
	waiter sync.WaitGroup
)

func handleGithub(w *sync.WaitGroup) {
	defer w.Done()
	repos, err := api.GetRepos()
	if err != nil {
		log.Println(err)
		return
	}
	if api.ConfiguredLabels != nil {
		waiter.Add(1)
		go labelUpdate(repos)
	}
	if api.AllRepoConfig != nil {
		waiter.Add(1)
		go allRepoUpdate(repos)
	}
	if api.NewLabels != nil {
		waiter.Add(1)
		go addNewLabels(repos)
	}
	waiter.Wait()
}

func addNewLabels(repos []api.Repository) {
	defer waiter.Done()
	var w sync.WaitGroup
	for _, r := range repos {
		if r.IsOrganisationsRepo() {
			continue
		}
		gotLabels, err := api.GetLabels(*r.Owner.Login, *r.Name)
		if err != nil {
			log.Println(err)
		}
		for _, l := range api.NewLabels {
			w.Add(1)
			go func(label *api.Label, repo api.Repository, existingLabels []api.Label) {
				defer w.Done()
				if checkLabel(existingLabels, *label.Name, *repo.FullName) {
					return
				}
				log.Printf("creating label %v in repo %v\n", *label.Name, *repo.FullName)
				newLabel, err := api.AddNewLabels(*label, repo)
				if err != nil {
					log.Println(err)
					return
				}
				if newLabel.Name == nil {
					return
				}
				log.Printf("label %v created in %v\n", *newLabel.Name, *repo.FullName)
			}(l, r, gotLabels)
		}
	}
	w.Wait()
}

func checkLabel(existingLabel []api.Label, newLabel, repo string) bool {
	log.Println("checking if label is already in repo")
	for _, l := range existingLabel {
		if *l.Name == newLabel {
			log.Printf("label %v exists in %v\n", *l.Name, repo)
			return true
		}
	}
	log.Printf("label %v does not exist in %v\n", newLabel, repo)
	return false
}

func allRepoUpdate(repos []api.Repository) {
	defer waiter.Done()
	var w sync.WaitGroup
	for _, r := range repos {
		if r.IsOrganisationsRepo() {
			continue
		}
		w.Add(1)
		go func(owner string, repos api.Repository) {
			defer w.Done()
			log.Printf("%v updating\n", *repos.FullName)
			newRepo, err := api.UpdateRepo(owner, *repos.Name)
			if err != nil {
				log.Println(err)
				return
			}
			if newRepo.Name == nil {
				return
			}
			log.Printf("%v : updated to config", *newRepo.FullName)
		}(*r.Owner.Login, r)
	}
	w.Wait()
}

func labelUpdate(repos []api.Repository) {
	defer waiter.Done()
	var wg sync.WaitGroup
	for _, r := range repos {
		if r.IsOrganisationsRepo() {
			continue
		}
		gotLabels, err := api.GetLabels(*r.Owner.Login, *r.Name)
		if err != nil {
			log.Println(err)
			continue
		}
		for _, l := range gotLabels {
			run, label := matchLabels(*l.Name, api.ConfiguredLabels)
			if !run {
				continue
			}
			wg.Add(1)
			go func(lName string, l *api.NewLabel, repos api.Repository) {
				defer wg.Done()
				log.Printf("%v : updating label %v\n", *r.FullName, *l.Name)
				newLabel, err := api.UpdateLabel(*repos.Name, *repos.Owner.Login, lName, l)
				if err != nil {
					log.Println(err)
					return
				}
				log.Printf("%v : label %v -> %v\n", *repos.FullName, *l.Name, *newLabel.Name)
			}(*l.Name, label, r)
		}
	}
	wg.Wait()
	return
}

func matchLabels(name string, array []*api.NewLabel) (bool, *api.NewLabel) {
	for _, a := range array {
		if name == *a.Name {
			return true, a
		}
	}
	return false, &api.NewLabel{NewName: nil}
}
