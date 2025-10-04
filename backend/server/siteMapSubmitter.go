package server

import (
	"log"

	"github.com/hinst/hinst-website/server/sitemap"
)

type siteMapSubmitter struct {
	db          *database
	siteMapPath string
}

func (me *siteMapSubmitter) run() {
	var siteMap = &sitemap.XmlSitemap{}
	readXml(readBytesFile(me.siteMapPath), siteMap)
	var client = &GoogleIndexingClient{}
	client.connect()
	for _, item := range siteMap.URL {
		var url = item.Loc
		var ok = client.updateUrl(url)
		log.Printf("URL: %v, ok: %v\n", url, ok)
		break
	}
}
