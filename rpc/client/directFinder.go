package client

type directFinder struct {
	address string
}

func (finder *directFinder) Init() error {
	return nil
}

func (finder *directFinder) Get() (string, error) {
	return finder.address, nil
}

func newDirectFinder(conf *Direct) *directFinder {
	return &directFinder{conf.Address}
}
