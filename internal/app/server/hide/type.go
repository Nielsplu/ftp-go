package hide

type HiddenFileCollection struct {
	path  string
	files map[string]struct{}
}

func ReadFromFile(path string) HiddenFileCollection {
	return HiddenFileCollection {
		path: path,
		files: readFrom(path),
	}
}

func (hf HiddenFileCollection) save() {
	writeTo(hf.path, hf.files)
}

func (hf HiddenFileCollection) IsPathHidden(path string) bool {
	//path = "data/"+path
	for i, char := range path {
		if char == '/' {
			_, inclued := hf.files[path[:i]]
			if inclued {
				return true
			}
		}
	}

	_, inclued := hf.files[path]
	return inclued
}

func (hf HiddenFileCollection) HidePath(path string) {
	
	hf.files[path] = struct{}{}

	hf.save()
}

func (hf HiddenFileCollection) RevealPath(path string) {

	delete(hf.files, path) 
	
	hf.save()
}