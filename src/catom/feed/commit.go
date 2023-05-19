package feed

import (
	"errors"
	"io"
	"net/url"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/johncgriffin/overflow"
)

type commitSt struct {
	isInitial                  bool
	commit, prevCommit         object.Commit
	MdFileInfo, prevMdFileInfo []MdFileInfo
}

type MdFileInfo struct {
	name, title, urlpath, hash string
	author, email              string
	commitid                   string
	commitime                  time.Time
}

func (self *MdFileInfo) Name() string {
	return self.name
}

func (self *MdFileInfo) Title() string {
	return self.title
}

func (self *MdFileInfo) UrlPath() string {
	return self.urlpath
}

func (self *MdFileInfo) Hash() string {
	return self.hash
}

func (self *MdFileInfo) Author() string {
	return self.author
}

func (self *MdFileInfo) Email() string {
	return self.email
}

func (self *MdFileInfo) CommitId() string {
	return self.commitid
}

func (self *MdFileInfo) CommitTime() time.Time {
	return self.commitime
}

func newCatomEntryID(host string, commitTime time.Time, commitID, fileHash string) (string, error) {
	form := "2006-01-02"
	tag := strings.ToLower("tag:" + host + ":" + commitTime.Format(form) + ":" + fileHash + ":" + commitID)

	_, err := url.ParseRequestURI(tag)
	if err != nil {
		return "", err
	}

	return tag, nil
}

func Feed4git(dotgitpath, urlroot string, root []string, now time.Time, diff int64) ([]MdFileInfo, error) {
	commitList, err := feed4gitInit(dotgitpath, now, diff)
	if err != nil {
		return nil, err
	}

	var entries []MdFileInfo
	if !isEmptyCommitList(commitList) {
		entries, err = createFeedData(commitList, root, urlroot)
		if err != nil {
			return nil, err
		}
	}

	return entries, nil
}

func isEmptyCommitList(c []commitSt) bool {
	return len(c) == 0
}

func feed4gitInit(dotgitpath string, now time.Time, diff int64) ([]commitSt, error) {
	cIter, err := getIter(dotgitpath)
	if err != nil {
		return nil, err
	}

	var commitList []commitSt
	commit, err := cIter.Next()
	if err == io.EOF {
		return nil, errors.New("No commit error")
	}
	if err != nil {
		return nil, err
	}
	for {
		ok, err := isTargetCommit(commit, now, diff)
		if err != nil {
			return nil, err
		}
		if !ok {
			break
		}
		var cm commitSt
		cm.commit = *commit
		prevCommit, err := cIter.Next()
		if err == io.EOF {
			cm.isInitial = true
			commitList = append(commitList, cm)
			break
		}
		if err != nil {
			return nil, err
		}

		cm.isInitial = false
		cm.prevCommit = *prevCommit

		commitList = append(commitList, cm)
		commit = prevCommit
	}

	return commitList, nil
}

func getIter(dotgitpath string) (object.CommitIter, error) {
	r, err := git.PlainOpen(dotgitpath)
	if err != nil {
		return nil, err
	}

	ref, err := r.Head()
	if err != nil {
		return nil, err
	}

	cIter, err := r.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return nil, err
	}

	return cIter, nil
}

func isTargetCommit(commit *object.Commit, now time.Time, base int64) (bool, error) {
	nowdiff := now.Sub(commit.Author.When).Hours()
	if nowdiff < 0 {
		return false, errors.New("Commit time is future")
	}
	diff, ok := overflow.Sub64(base, int64(nowdiff))
	if !ok {
		return false, errors.New("Var diff causes Sub calc is Overflow")
	}

	if diff < 0 {
		return false, nil
	} else {
		return true, nil
	}
}

func isMarkdown(file string) bool {
	switch filepath.Ext(file) {
	case ".md", ".markdown":
		return true
	default:
		return false
	}
}

func createFeedData(commitList []commitSt, root []string, urlroot string) ([]MdFileInfo, error) {
	var feedEntries []MdFileInfo
	for i := range commitList {
		files, err := commitList[i].commit.Files()
		if err != nil {
			return nil, err
		}
		author := commitList[i].commit.Author.Name
		email := commitList[i].commit.Author.Email
		id := commitList[i].commit.ID().String()
		cmtime := commitList[i].commit.Author.When
		commitList[i].MdFileInfo, err = setMarkdownInfo(files, author, email, id, urlroot, cmtime, root)
		if err != nil {
			return nil, err
		}

		if commitList[i].isInitial {
			for _, file := range commitList[i].MdFileInfo {
				feedEntries = append(feedEntries, file)
			}
			break
		}

		files, err = commitList[i].prevCommit.Files()
		if err != nil {
			return nil, err
		}
		commitList[i].prevMdFileInfo, err = setMarkdownInfo(files, "", "", "", "", time.Time{}, root)
		if err != nil {
			return nil, err
		}

		for _, file := range commitList[i].MdFileInfo {
			isCHash := isContainsHash(file, commitList[i].prevMdFileInfo)
			isCFileName := isContainsFileName(file, commitList[i].prevMdFileInfo)
			switch {
			case isCHash && isCFileName:
				//None
				continue
			case isCHash && !isCFileName:
				//R100
				feedEntries = append(feedEntries, file)
			case !isCHash && isCFileName:
				//M
				feedEntries = append(feedEntries, file)
			case !isCHash && !isCFileName:
				//A or R??
				feedEntries = append(feedEntries, file)
			}
		}
	}

	return feedEntries, nil
}

func isContainsHash(file MdFileInfo, prevFiles []MdFileInfo) bool {
	for _, prevfile := range prevFiles {
		if isSameHash(file, prevfile) {
			return true
		}
	}
	return false
}

func isContainsFileName(file MdFileInfo, prevFiles []MdFileInfo) bool {
	for _, prevfile := range prevFiles {
		if isSameFileName(file, prevfile) {
			return true
		}
	}
	return false
}

func isSameHash(file, prevFile MdFileInfo) bool {
	if file.hash == prevFile.hash {
		return true
	} else {
		return false
	}
}

func isSameFileName(file, prevFile MdFileInfo) bool {
	if file.name == prevFile.name {
		return true
	} else {
		return false
	}
}

func getMdTitle(lines []string) (string, error) {
	for _, line := range lines {
		if !strings.Contains(line, "# ") {
			continue
		}

		title := strings.SplitN(line, "# ", 2)
		if len(title) != 2 {
			continue
		}

		return title[1], nil
	}

	return "", errors.New("No title.")
}

func setMarkdownInfo(files *object.FileIter, name, email, id, urlroot string, cmtime time.Time, root []string) ([]MdFileInfo, error) {
	mdinfos := []MdFileInfo{}
	for {
		file, err := files.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		urlpath, err := chkContainsFileOnRoot(file.Name, root)
		if err != nil {
			continue
		}

		if isMarkdown(file.Name) && file.Mode.IsFile() {
			lines, err := file.Lines()
			if err != nil {
				return nil, err
			}

			title, err := getMdTitle(lines)
			if err != nil {
				title = "No Title"
			}

			conurl := urlroot + "/" + urlpath
			conurl = path.Clean(conurl)
			_, err = url.ParseRequestURI(conurl)
			if err != nil {
				return nil, err
			}

			mdinfos = append(mdinfos, MdFileInfo{
				file.Name,
				title,
				conurl,
				file.ID().String(),
				name,
				email,
				id,
				cmtime,
			})
		}
	}

	return mdinfos, nil
}

func chkContainsFileOnRoot(file string, root []string) (string, error) {
	for _, d := range root {
		split_path := strings.SplitN(filepath.Clean("/"+file), filepath.Clean("/"+d)+"/", 2)
		if len(split_path) == 2 && split_path[0] == "" {
			return split_path[1], nil
		}
	}

	return "", errors.New("Don't match dir")
}
