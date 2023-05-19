package feed

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	atom "github.com/denisbrodbeck/atomfeed"
	"github.com/naoina/toml"
)

type FeedConfig struct {
	Dotgit  string
	Diff    int64
	Root    []string
	Urlroot string
	Proto   string
	Host    string
	Feedurl string
	Outpath string
	Outfile string
	Feed    *feed
	Author  *author
}

type feed struct {
	Feedid   string
	Title    string
	Subtitle string
}

type author struct {
	Name  string
	Email string
	Url   string
}

const (
	http  = "http"
	https = "https"
)

func NewConfig(envPath string) (*FeedConfig, error) {
	env, err := os.Open(envPath)
	if err != nil {
		return nil, err
	}
	defer env.Close()

	var cfg FeedConfig
	err = toml.NewDecoder(env).Decode(&cfg)
	if err != nil {
		return nil, err
	}

	if !filepath.IsAbs(cfg.Dotgit) {
		return nil, errors.New("dotgitpath isn't Abs")
	}
	if filepath.Base(cfg.Dotgit) != ".git" {
		return nil, errors.New("dotigitpath don't match .git")
	}
	splitpath := strings.Split(cfg.Dotgit, ".git")
	if len(splitpath) > 2 {
		return nil, errors.New("Include only one .git directory.")
	}
	if !filepath.IsAbs(cfg.Outpath) {
		return nil, errors.New(fmt.Sprintf("Don't Abs: %s", cfg.Outpath))
	}
	if filepath.Ext(cfg.Outfile) != ".atom" {
		return nil, errors.New(fmt.Sprintf("Don't Ext .atom: %s", filepath.Base(cfg.Outfile)))
	}

	if cfg.Feed == nil {
		return nil, errors.New("*FeedConfig.Feed is nil")
	}
	if cfg.Proto != http && cfg.Proto != https {
		return nil, errors.New("Scheme Error [http or https]")
	}

	return &cfg, nil
}

func (self *FeedConfig) NewAtom(now time.Time, gEntries []MdFileInfo) (*atom.Feed, error) {
	feedID := atom.NewID(self.Feed.Feedid)
	var author *atom.Person
	if self.Author != nil {
		author = atom.NewPerson(self.Author.Name, self.Author.Email, self.Author.Url)
	} else {
		author = atom.NewPerson("Anonymous", "", "")
	}
	title := self.Feed.Title
	subtitle := self.Feed.Subtitle
	url := &url.URL{
		Scheme: self.Proto,
		Host:   self.Host,
	}
	baseURL := url.String()
	feedURL := self.Feedurl
	updated := now

	var entries []atom.Entry
	for _, gEntry := range gEntries {
		tag, err := newCatomEntryID(self.Host, gEntry.CommitTime(), gEntry.CommitId(), gEntry.Hash())
		if err != nil {
			return nil, err
		}
		entryID := atom.NewID(tag)
		commitTime := gEntry.CommitTime()
		entryTitle := gEntry.Title()
		url.Path = gEntry.UrlPath()
		updated := gEntry.CommitTime()

		entries = append(entries, atom.NewEntry(
			atom.NewEntryID(entryID, commitTime),
			entryTitle,
			url.String(),
			nil,
			updated, time.Time{},
			nil, nil, nil,
		))
	}

	feed := atom.NewFeed(feedID, author, title, subtitle, baseURL, feedURL, updated, entries)
	if err := feed.Verify(); err != nil {
		return nil, err
	}

	return &feed, nil
}
