package git

func (g *GitVersioning) GetMostRecentTag() (string, error) {
	return g.getMostRecentTag()
}
