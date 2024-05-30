package entities

type Prompt struct {
    Choices  []string
    Cursor   int
    Selected map[int]struct{}
}
