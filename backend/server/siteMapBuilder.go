package server

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/carlosstrand/go-sitemap"
)

type siteMapBuilder struct {
	newFilesPath string
	oldFilesPath string
}

func (me *siteMapBuilder) run() string {
	var items []*sitemap.SitemapItem
	var pathPrefix = me.newFilesPath
	pathPrefix = strings.TrimPrefix(pathPrefix, "./") + "/"
	filepath.WalkDir(me.newFilesPath, func(path string, directory os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if directory.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ".html" {
			var relativePath = path[len(pathPrefix):]
			relativePath = strings.ReplaceAll(relativePath, "\\", "/")
			var oldPath = me.oldFilesPath + "/" + relativePath
			log.Println(path, relativePath, checkFileExists(oldPath))
		}
		return nil
	})
	var siteMap = sitemap.NewSitemap(items, nil)
	return assertResultError(siteMap.ToXMLString())
}
