package pub

type Publish interface {
  PublishData(<-chan bool) (chan<- []byte)
}
