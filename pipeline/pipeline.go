package pipeline

type Pipeline struct {
	repo   string
	build  []string
	deploy []string
}
