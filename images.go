package kubenventoree

type ImageDescription struct {
	ImageRepository string
	ImageTag        string
	ImagePullable   string
	Count           int
}

func ReadAllImages() []ImageDescription {
	return nil
}
