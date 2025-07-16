package gitops

import (
	"fmt"
	"github.com/arcalyx/gitver/internal/constants"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"sort"
	"time"
)

const (
	EMPTY = ""
)

var g *GitOps

type GitOps struct {
	path          string
	name          string
	email         string
	repository    *git.Repository
	worktree      *git.Worktree
	commitMessage string
	tagMessage    string
}

type Tag struct {
	Name    string
	Message string
	Hash    plumbing.Hash
}

func init() {
	g = New()
}

func New() *GitOps {
	g := new(GitOps)
	g.commitMessage = constants.CommitMessage
	g.tagMessage = constants.TagMessage
	return g
}

func SetRepositoryPath(path string) {
	g.SetRepositoryPath(path)
}

func (g *GitOps) SetRepositoryPath(path string) {
	g.path = path
}

func ReadRepository() error {
	return g.ReadRepository()
}

func (g *GitOps) ReadRepository() error {

	r, err := git.PlainOpen(g.path)
	if err != nil {
		// RepositoryNotExists
		return err
	}

	cfg, err := r.ConfigScoped(config.GlobalScope)
	if err != nil {
		return err
	}

	user := cfg.User
	if user.Name == "" && user.Email == "" {
		return nil
	}

	w, err := r.Worktree()
	if err != nil {
		return err
	}

	g.name = user.Name
	g.email = user.Email
	g.repository = r
	g.worktree = w

	return nil
}

func GetStatus() (git.Status, error) {
	return g.GetStatus()
}

func (g *GitOps) GetStatus() (git.Status, error) {
	return g.worktree.Status()
}

func IsCleanRepo() (bool, error) {
	return g.IsCleanRepo()
}

func (g *GitOps) IsCleanRepo() (bool, error) {
	status, err := g.GetStatus()
	if err != nil {
		return false, err
	}
	return status.IsClean(), nil
}

func Add() (plumbing.Hash, error) {
	return g.Add()
}

func (g *GitOps) Add() (plumbing.Hash, error) {
	return g.worktree.Add(".")
}

func Commit(message string, amend bool) error {
	return g.Commit(message, amend)
}

func (g *GitOps) Commit(message string, amend bool) error {
	commitOptions := &git.CommitOptions{
		Author: &object.Signature{
			Name:  g.name,
			Email: g.email,
			When:  time.Now(),
		},
	}

	// Wenn amend wahr ist, verwenden Sie die Amend-Option
	if amend {
		headRef, err := g.repository.Head()
		if err != nil {
			return err
		}
		commitOptions.Parents = []plumbing.Hash{headRef.Hash()}
	}

	_, err := g.worktree.Commit(message, commitOptions)
	if err != nil {
		return err
	}
	return nil
}

func CreateTag(tag string, message string) error {
	return g.CreateTag(tag, message)
}

func (g *GitOps) CreateTag(tag string, message string) error {

	headRef, err := g.repository.Head()
	if err != nil {
		return err
	}

	headCommit, err := g.repository.CommitObject(headRef.Hash())
	if err != nil {
		return err
	}

	// Create the tag
	_, err = g.repository.CreateTag(tag, headCommit.Hash, &git.CreateTagOptions{
		Tagger: &object.Signature{
			Name:  g.name,
			Email: g.email,
			When:  time.Now(),
		},
		Message: message,
	})
	if err != nil {
		return err
	}

	return nil
}

func GetLastTag() (string, error) {
	return g.GetLatestTag()
}

func (g *GitOps) GetLatestTag() (string, error) {

	tagRefs, err := g.repository.Tags()
	if err != nil {
		return "", err
	}

	var tags []struct {
		Name string
		When time.Time
	}

	err = tagRefs.ForEach(func(t *plumbing.Reference) error {
		obj, err := g.repository.TagObject(t.Hash())
		if err != nil {
			// Es könnte ein leichtes Tag sein, also versuchen Sie, den Commit direkt zu bekommen
			commit, err := g.repository.CommitObject(t.Hash())
			if err != nil {
				return nil // Ignorieren von Fehlern, die durch leichte Tags verursacht werden
			}
			tags = append(tags, struct {
				Name string
				When time.Time
			}{t.Name().Short(), commit.Committer.When})
			return nil
		}
		tags = append(tags, struct {
			Name string
			When time.Time
		}{t.Name().Short(), obj.Tagger.When})
		return nil
	})
	if err != nil {
		return "", err
	}

	// Sortieren Sie die Tags nach Datum
	sort.Slice(tags, func(i, j int) bool {
		return tags[i].When.After(tags[j].When)
	})

	if len(tags) > 0 {
		return tags[0].Name, nil
	}

	return "", nil
}

func Push() error {
	return g.Push()
}

func (g *GitOps) Push() error {

	// HEAD-Referenz holen
	headRef, err := g.repository.Head()
	if err != nil {
		return err
	}

	// HEAD dereferenzieren, um den aktuellen Branch zu bekommen
	ref, err := g.repository.Reference(headRef.Name(), true)
	if err != nil {
		return err
	}

	err = g.repository.Push(&git.PushOptions{
		RemoteName: "origin",
		RefSpecs:   []config.RefSpec{config.RefSpec(ref.Name() + ":" + ref.Name())},
	})
	if err != nil {
		return err
	}

	return nil
}

func GetHeadCommit() (*object.Commit, error) {
	return g.GetHeadCommit()
}

func (g *GitOps) GetHeadCommit() (*object.Commit, error) {
	headRef, err := g.repository.Head()
	if err != nil {
		return nil, err
	}

	commit, err := g.repository.CommitObject(headRef.Hash())
	if err != nil {

	}

	return commit, nil
}

func GetCommits(tag string) ([]*object.Commit, error) {
	return g.GetCommitsBetweenTags(EMPTY, tag)
}

func GetTag(tag string) (plumbing.Hash, error) {
	return g.GetTag(tag)
}

func (g *GitOps) GetTag(tag string) (plumbing.Hash, error) {
	tagRef, err := g.repository.Tag(tag)
	if err == nil {
		// Dereferenziert das Tag-Objekt, falls es ein annotiertes Tag ist
		resolvedTag, err := g.repository.TagObject(tagRef.Hash())
		if err == nil {
			return resolvedTag.Target, nil
		} else {
			// Wenn es kein annotiertes Tag ist, sondern ein leichtgewichtiger Tag
			return tagRef.Hash(), nil
		}
	}
	return plumbing.Hash{}, nil
}

func HasTag(tag string) bool {
	hash, _ := g.GetTag(tag)
	return !hash.IsZero()
}

func GetCommitsBetweenTags(tagStart, tagEnd string) ([]*object.Commit, error) {
	return g.GetCommitsBetweenTags(tagStart, tagEnd)
}

func (g *GitOps) GetCommitsBetweenTags(startTag, endTag string) ([]*object.Commit, error) {

	startHash, err := g.GetTag(startTag)
	if err != nil {
		headRef, err := g.repository.Head()
		if err != nil {
			return nil, err
		}
		startHash = headRef.Hash()
	}

	endHash, _ := g.GetTag(endTag)
	if err != nil {
		endHash = plumbing.Hash{}
	}

	logOptions := &git.LogOptions{From: startHash}
	commitIter, err := g.repository.Log(logOptions)
	if err != nil {
		return nil, err
	}

	var commits []*object.Commit
	err = commitIter.ForEach(func(c *object.Commit) error {
		if c.Hash == endHash {
			return fmt.Errorf("BREAK")
		}
		commits = append(commits, c)
		return nil
	})

	if err != nil && err.Error() != "BREAK" {
		return nil, err
	}

	return commits, nil
}

func GetTags() ([]Tag, error) {
	return g.GetTags()
}

func (g *GitOps) GetTags() ([]Tag, error) {
	if g.repository == nil {
		if err := g.ReadRepository(); err != nil {
			return nil, err
		}
	}

	tagRefs, err := g.repository.Tags()
	if err != nil {
		return nil, err
	}

	var tags []Tag
	err = tagRefs.ForEach(func(t *plumbing.Reference) error {
		obj, err := g.repository.TagObject(t.Hash())
		if err != nil {
			// It might be a lightweight tag, so try to get the commit directly
			commit, err := g.repository.CommitObject(t.Hash())
			if err != nil {
				return nil // Ignore errors caused by lightweight tags
			}
			tags = append(tags, Tag{
				Name:    t.Name().Short(),
				Message: commit.Message,
				Hash:    t.Hash(),
			})
			return nil
		}
		tags = append(tags, Tag{
			Name:    t.Name().Short(),
			Message: obj.Message,
			Hash:    obj.Target,
		})
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Sort tags by creation time (newest first)
	sort.Slice(tags, func(i, j int) bool {
		timeI, _ := g.GetTagTime(tags[i].Name)
		timeJ, _ := g.GetTagTime(tags[j].Name)
		return timeI.After(timeJ)
	})

	return tags, nil
}

func GetTagTime(tagName string) (time.Time, error) {
	return g.GetTagTime(tagName)
}

func (g *GitOps) GetTagTime(tagName string) (time.Time, error) {
	if g.repository == nil {
		if err := g.ReadRepository(); err != nil {
			return time.Time{}, err
		}
	}

	tagRef, err := g.repository.Tag(tagName)
	if err != nil {
		return time.Time{}, err
	}

	// Try to get the tag object (for annotated tags)
	tagObj, err := g.repository.TagObject(tagRef.Hash())
	if err == nil {
		return tagObj.Tagger.When, nil
	}

	// If it's a lightweight tag, get the commit it points to
	commit, err := g.repository.CommitObject(tagRef.Hash())
	if err != nil {
		return time.Time{}, err
	}

	return commit.Committer.When, nil
}
