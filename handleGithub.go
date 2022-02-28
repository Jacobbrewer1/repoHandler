package main

import (
	"github.com/Jacobbrewer1/repoHandler/api"
	"log"
	"sync"
)

func handleGithub() error {
	repos, err := api.GetRepos()
	if err != nil {
		return err
	}
	return labelUpdate(repos)
}

func labelUpdate(repos []api.Repository) error {
	var wg sync.WaitGroup
	for _, r := range repos {
		gotLabels, err := api.GetLabels(*r.Owner.Login, *r.Name)
		if err != nil {
			return err
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
	return nil
}

func matchLabels(name string, array []*api.NewLabel) (bool, *api.NewLabel) {
	for _, a := range array {
		if name == *a.Name {
			return true, a
		}
	}
	return false, &api.NewLabel{NewName: nil}
}
