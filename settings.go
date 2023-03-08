package main

var supportedExtensions = map[string][]string{
	"aac":  []string{"", ""},
	"ai":   []string{"", ""},
	"bmp":  []string{"bi-image", ""},
	"cs":   []string{"", ""},
	"css":  []string{"", ""},
	"csv":  []string{"", ""},
	"doc":  []string{"", ""},
	"docx": []string{"", ""},
	/* "exe":    []string{"", ""}, */
	"gif":    []string{"bi-image", ""},
	"heic":   []string{"bi-image", ""},
	"html":   []string{"", ""},
	"java":   []string{"", ""},
	"jpg":    []string{"bi-image", ""},
	"js":     []string{"", ""},
	"json":   []string{"", ""},
	"jsx":    []string{"", ""},
	"key":    []string{"", ""},
	"m4p":    []string{"", ""},
	"md":     []string{"", ""},
	"mdx":    []string{"", ""},
	"mov":    []string{"", ""},
	"mp3":    []string{"bi-file-earmark-music-fill", ""},
	"mp4":    []string{"bi-film", ""},
	"otf":    []string{"", ""},
	"pdf":    []string{"bi-file-pdf-fill", "red"},
	"php":    []string{"", ""},
	"png":    []string{"bi-image", ""},
	"ppt":    []string{"", ""},
	"pptx":   []string{"", ""},
	"psd":    []string{"", ""},
	"py":     []string{"", ""},
	"raw":    []string{"bi-image", ""},
	"rb":     []string{"", ""},
	"sass":   []string{"", ""},
	"scss":   []string{"", ""},
	"sh":     []string{"", ""},
	"sql":    []string{"", ""},
	"svg":    []string{"bi-image", ""},
	"tiff":   []string{"bi-image", ""},
	"tsx":    []string{"", ""},
	"ttf":    []string{"", ""},
	"txt":    []string{"", ""},
	"wav":    []string{"", ""},
	"woff":   []string{"", ""},
	"xls":    []string{"", ""},
	"xlsx":   []string{"", ""},
	"xml":    []string{"", ""},
	"yml":    []string{"", ""},
	"zip":    []string{"bi-file-earmark-zip-fill", "orange"},
	"gz":     []string{"bi-file-earmark-zip-fill", "orange"},
	"tar":    []string{"bi-file-earmark-zip-fill", "orange"},
	"epub":   []string{"bi-book", ""},
	"bundle": []string{"bi-terminal-fill", ""},
	"run":    []string{"bi-terminal-fill", ""},
	"exe":    []string{"bi-terminal-fill", ""},
	"rpm":    []string{"bi-terminal-fill", ""},
}