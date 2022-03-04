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
	waiter.Wait()
}

func allRepoUpdate(repos []api.Repository) {
	defer waiter.Done()
	var w sync.WaitGroup
	for _, r := range repos {
		if r.IsOrganisationsRepo() {
			continue
		}
		w.Add(1)
		go func(o, n string) {
			defer w.Done()
			newRepo, err := api.UpdateRepo(o, n)
			if err != nil {
				log.Println(err)
				return
			}
			if newRepo.Name == nil {
				return
			}
			log.Printf("repo %v updated as per config", *newRepo.Name)
		}(*r.Owner.Login, *r.Name)
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
			go func(lName string, l *api.NewLabel) {
				defer wg.Done()
				log.Printf("Updating label %v in %v\n", *l.Name, *r.Name)
				newLabel, err := api.UpdateLabel(*r.Name, *r.Owner.Login, lName, l)
				if err != nil {
					log.Println(err)
					return
				}
				log.Printf("%v -> %v\n", *l.Name, *newLabel.Name)
			}(*l.Name, label)
			wg.Wait()
		}
	}
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
