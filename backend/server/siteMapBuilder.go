package server

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hinst/go-common"
	"github.com/hinst/hinst-website/server/sitemap"
)

type siteMapBuilder struct {
	webPath            string
	newFilesPath       string
	newFilesPathPrefix string
	oldFilesPath       string
	oldSiteMap         *sitemap.XmlSitemap
	items              []*sitemap.SitemapItem
}

func (me *siteMapBuilder) run() {
	me.loadOldSitemap()
	me.newFilesPathPrefix = me.newFilesPath
	me.newFilesPathPrefix = strings.TrimPrefix(me.newFilesPathPrefix, "./") + "/"
	common.AssertError(filepath.WalkDir(me.newFilesPath, me.createItem))
	var options = sitemap.Options{PrettyOutput: true, WithXMLHeader: true, Validate: true}
	var siteMap = sitemap.NewSitemap(me.items, &options)
	var siteMapText = common.AssertResultError(siteMap.ToXMLString())
	writeTextFile(me.newFilesPath+"/sitemap.xml", siteMapText)
}

func (me *siteMapBuilder) loadOldSitemap() {
	var oldSiteMapPath = me.oldFilesPath + "/sitemap.xml"
	if !checkFileExists(oldSiteMapPath) {
		return
	}
	var text = readBytesFile(me.oldFilesPath + "/sitemap.xml")
	me.oldSiteMap = &sitemap.XmlSitemap{}
	readXml(text, me.oldSiteMap)
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
	var haveChange = !checkTextFilesEqual(oldFilePath, newFilePath)
	var url = me.webPath + "/" + relativePath
	var item = &sitemap.SitemapItem{
		Loc:        url,
		Priority:   me.getDefaultPriority(),
		ChangeFreq: me.getDefaultChangeFreq(),
	}
	if haveChange {
		item.LastMod = time.Now()
	} else {
		var previousLastMod = me.findPreviousLastMod(url)
		if previousLastMod != nil {
			item.LastMod = *previousLastMod
		} else {
			item.LastMod = time.Now()
		}
	}
	me.items = append(me.items, item)
	return nil
}

func (siteMapBuilder) getDefaultPriority() float32 {
	return 0.5
}

func (siteMapBuilder) getDefaultChangeFreq() string {
	return "yearly"
}

func (me *siteMapBuilder) findPreviousLastMod(url string) *time.Time {
	if nil == me.oldSiteMap {
		return nil
	}
	for _, item := range me.oldSiteMap.URL {
		if item.Loc == url {
			var previousTime = common.AssertResultError(time.Parse(time.DateOnly, item.LastMod)).UTC()
			return &previousTime
		}
	}
	return nil
}
