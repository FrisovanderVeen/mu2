package bot

// OpusReader returns an opus frame and an error
type OpusReader interface {
	OpusFrame() ([]byte, error)
}
