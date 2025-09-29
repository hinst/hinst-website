package server

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hinst/hinst-website/server/sitemap"
)

type siteMapBuilder struct {
	webPath            string
	newFilesPath       string
	newFilesPathPrefix string
	oldFilesPath       string
	items              []*sitemap.SitemapItem
}

func (me *siteMapBuilder) run() {
	me.newFilesPathPrefix = me.newFilesPath
	me.newFilesPathPrefix = strings.TrimPrefix(me.newFilesPathPrefix, "./") + "/"
	filepath.WalkDir(me.newFilesPath, me.createItem)
	var options = sitemap.Options{PrettyOutput: true, WithXMLHeader: true, Validate: true}
	var siteMap = sitemap.NewSitemap(me.items, &options)
	var siteMapText = assertResultError(siteMap.ToXMLString())
	log.Println(me.newFilesPath + "/sitemap.xml")
	writeTextFile(me.newFilesPath+"/sitemap.xml", siteMapText)
}

func (me *siteMapBuilder) loadOldSitemap() {
}

func (me *siteMapBuilder) createItem(newFilePath string, directory os.DirEntry, err error) error {
	if err != nil {
		return err
	}
	if directory.IsDir() {
		return nil
	}
	if filepath.Ext(newFilePath) != ".html" {
		return nil
	}
	var relativePath = newFilePath[len(me.newFilesPathPrefix):]
	relativePath = strings.ReplaceAll(relativePath, "\\", "/")
	var oldFilePath = me.oldFilesPath + "/" + relativePath
	var haveChange = !checkFilesEqual(oldFilePath, newFilePath)
	var url = me.webPath + "/" + relativePath
	log.Println(url)
	var item = &sitemap.SitemapItem{
		Loc:        url,
		Priority:   me.getDefaultPriority(),
		ChangeFreq: me.getDefaultChangeFrequency(),
	}
	if haveChange {
		item.LastMod = time.Now()
	}
	me.items = append(me.items, item)
	return nil
}

func (siteMapBuilder) getDefaultPriority() float32 {
	return 0.5
}

func (siteMapBuilder) getDefaultChangeFrequency() string {
	return "yearly"
}
