package server

import (
	"log"
	"time"

	"github.com/hinst/go-common"
	"github.com/hinst/hinst-website/server/sitemap"
)

type siteMapSubmitter struct {
	db          *database
	siteMapPath string
	client      *GoogleIndexingClient
}

func (me *siteMapSubmitter) run() {
	var siteMap = common.DecodeXml(readBytesFile(me.siteMapPath), new(sitemap.XmlSitemap))
	var siteMapItems = siteMap.URL
	var pingedEarlierCount = 0
	var pingedNowCount = 0
	for _, siteMapItem := range siteMapItems {
		var url = siteMapItem.Loc
		var record = me.db.getUrlPing(url)
		if record == nil {
			me.db.insertUrlPing(url)
			record = me.db.getUrlPing(url)
		}
		if record.GooglePingedAt == nil {
			var ok = me.getClient().updateUrl(url)
			log.Printf("Pinged Google URL: %v, ok: %v", url, ok)
			if ok {
				me.db.updateUrlPingGoogle(url, time.Now())
				pingedNowCount++
			} else {
				log.Printf("Rate limit reached")
				break
			}
		} else {
			pingedEarlierCount++
		}
	}
	log.Printf("Total URLs in sitemap: %v, pinged now: %v, pinged earlier: %v",
		len(siteMapItems), pingedNowCount, pingedEarlierCount)
}

func (me *siteMapSubmitter) getClient() *GoogleIndexingClient {
	if me.client == nil {
		me.client = &GoogleIndexingClient{}
		me.client.connect()
	}
	return me.client
}
