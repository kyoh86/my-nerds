package archive

import "io"

type Extractor func(pathTo string, reader io.Reader) error
