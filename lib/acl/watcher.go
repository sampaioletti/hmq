package acl

import (
	"github.com/fsnotify/fsnotify"
	"go.uber.org/zap"
)

var (
	watchList = []string{"./conf"}
)

func NewWatcher(c *ACLConfig) *Watcher {
	return &Watcher{c}
}

type Watcher struct {
	config *ACLConfig
}

func (b *Watcher) handleFsEvent(event fsnotify.Event) error {
	switch event.Name {
	case b.config.File:
		if event.Op&fsnotify.Write == fsnotify.Write ||
			event.Op&fsnotify.Create == fsnotify.Create {
			log.Info("text:handling acl config change event:", zap.String("filename", event.Name))
			aclconfig, err := AclConfigLoad(event.Name)
			if err != nil {
				log.Error("aclconfig change failed, load acl conf error: ", zap.Error(err))
				return err
			}
			if b.config == nil {
				b.config = aclconfig
			} else {
				*b.config = *aclconfig
			}

		}
	}
	return nil
}

func (b *Watcher) StartAclWatcher(c *ACLConfig) {
	go func() {
		wch, e := fsnotify.NewWatcher()
		if e != nil {
			log.Error("start monitor acl config file error,", zap.Error(e))
			return
		}
		defer wch.Close()

		for _, i := range watchList {
			if err := wch.Add(i); err != nil {
				log.Error("start monitor acl config file error,", zap.Error(err))
				return
			}
		}
		log.Info("watching acl config file change...")
		for {
			select {
			case evt := <-wch.Events:
				b.handleFsEvent(evt)
			case err := <-wch.Errors:
				log.Error("error:", zap.Error(err))
			}
		}
	}()
}
