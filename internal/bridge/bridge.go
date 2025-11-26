package bridge

type Bridge struct {
	Commands chan string
	Events   chan string
}
