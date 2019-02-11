package db

import (
	"errors"

	"github.com/nanobox-io/golang-scribble"
	log "github.com/sirupsen/logrus"
)

type Db struct {
	scribble *scribble.Driver
}

type Opts func(*Db) error

func NewDb(path string, opts ...Opts) (*Db, error) {
	d, err := scribble.New(path, nil)
	if err != nil {
		return nil, errors.New("couldn't create db, " + err.Error())
	}

	db := Db{
		scribble: d,
	}

	for _, opt := range opts {
		if opt != nil {
			err := opt(&db)
			if err != nil {
				log.Error("couldn't set options: " + err.Error())
			}
		}
	}

	return &db, nil
}

type settingsKey string

const (
	SettingsCountdownDate    settingsKey = "countdownDate"
	SettingsCountdownEnabled settingsKey = "countdownEnabled"
)

func SettingsWrite(db *Db, key settingsKey, val string) {
	if err := db.scribble.Write("settings", string(key), val); err != nil {
		log.Errorf("could not save setting %s:%s", key, val)
	}
}

func SettingsRead(db *Db, key settingsKey) string {
	var val string
	_ = db.scribble.Read("settings", string(key), &val)
	return val
}
