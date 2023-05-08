package git

import "github.com/go-git/go-git/v5/plumbing"

func (g *GitVersioning) GetMostRecentTag() (string, error) {
	return g.getMostRecentTag()
}

func (g *GitVersioning) BranchHead() *plumbing.Reference {
	return g.branchHead
}
