package orchestrator

import (
	"GoMapEnum/src/utils"
	"reflect"
	"strings"
	"sync"
)

func (orchestrator *Orchestrator) UserEnum(optionsModules Options) []string {
	optionsInterface := reflect.ValueOf(optionsModules).Interface()
	options := optionsModules.GetBaseOptions()
	options.Users = utils.GetStringOrFile(options.Users)
	options.UsernameList = strings.Split(options.Users, "\n")
	mux := &sync.Mutex{}
	var wg sync.WaitGroup
	var validusers []string
	queue := make(chan string)
	if orchestrator.PreActionUserEnum != nil {
		orchestrator.PreActionUserEnum(&optionsInterface)
	}
	for i := 0; i < options.Thread; i++ {

		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for username := range queue {
				options.Log.Verbose("Testing " + username)
				if orchestrator.CheckBeforeEnumFunc != nil {
					if !orchestrator.CheckBeforeEnumFunc(&optionsInterface, username) {
						continue
					}
				}
				if orchestrator.UserEnumFunc(&optionsInterface, username) {
					mux.Lock()
					validusers = append(validusers, username)
					mux.Unlock()
				}
			}
		}(i)
	}

	// Trim usernames and send them to the pool of workers
	for _, username := range options.UsernameList {
		username = strings.ToValidUTF8(username, "")
		username = strings.Trim(username, "\r")
		username = strings.Trim(username, "\n")
		if username == "" {
			continue
		}
		queue <- username
	}

	close(queue)
	wg.Wait()
	if orchestrator.PostActionUserEnum != nil {
		orchestrator.PostActionUserEnum(&optionsInterface)
	}
	return validusers
}
