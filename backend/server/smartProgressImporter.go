package server

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/hinst/go-gophers"
	"github.com/hinst/hinst-website/server/base"
	"github.com/hinst/hinst-website/server/db_objects"
	"github.com/hinst/hinst-website/server/html_util"
	"github.com/hinst/hinst-website/server/smart_progress"
	"golang.org/x/net/html"
)

type smartProgressImporter struct {
	goalIds  []string
	database *database
}

func (me *smartProgressImporter) run() {
	for _, goalId := range me.goalIds {
		me.syncGoal(goalId)
	}
}

func (me *smartProgressImporter) syncGoal(goalId string) {
	var goalInfo = me.readGoalInfo(goalId)
	me.saveGoalInfo(goalInfo)
	me.syncPosts(goalId)
}

func (me *smartProgressImporter) syncPosts(goalId string) {
	var posts = me.readAllPosts(goalId)
	var newCount = 0
	for _, post := range posts {
		var isNew = me.savePost(post)
		var comments = me.readComments(post.Id)
		me.saveComments(post, comments)
		if isNew {
			newCount += 1
		}
		var age = time.Since(me.parseDateTime(post.Date))
		if isNew || age < 30*24*time.Hour {
			var images = me.readImages(post)
			me.saveImages(post, images)
		}
	}
	log.Printf("Sync complete: goal=%s posts=%d new=%d", goalId, len(posts), newCount)
}

func (me *smartProgressImporter) parseDateTime(text string) (result time.Time) {
	var err error
	result, err = time.ParseInLocation("2006-01-02 15:04:05", text, time.UTC)
	if err != nil {
		panic("Cannot parse date time: \"" + text + "\"")
	}
	return
}

func (me *smartProgressImporter) saveComments(post smart_progress.Post, comments []smart_progress.Comment) {
	var parentDateTime = me.parseDateTime(post.Date).UTC().Unix()
	var goalId = gophers.GetInt64FromString(post.ObjId)
	for _, comment := range comments {
		var row = db_objects.GoalPostCommentRow{
			GoalId:              goalId,
			ParentDateTime:      parentDateTime,
			DateTime:            me.parseDateTime(comment.Date).UTC().Unix(),
			SmartProgressUserId: gophers.GetInt64FromStringOptional(comment.UserId),
			Username:            comment.Username,
			Text:                convertHtmlToMarkdown(me.unpackRedirects(comment.Msg)),
		}
		me.database.saveGoalPostComment(row)
	}
}

func (me *smartProgressImporter) unpackRedirects(htmlText string) (result string) {
	var document = html_util.ParseHtmlFragment(htmlText)
	html_util.Walk(document, func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "a" {
			var hrefAttr = html_util.Attr(node, "href")
			if hrefAttr == nil {
				return
			}
			var href = hrefAttr.Val
			if strings.HasPrefix(href, smart_progress.RedirectPrefix) {
				href = strings.TrimSuffix(strings.TrimPrefix(href, smart_progress.RedirectPrefix), "%")
				if decoded, err := url.PathUnescape(href); err == nil {
					href = decoded
				}
				hrefAttr.Val = href
			}
		}
	})
	return html_util.InnerHtml(document)
}

// Returns true if the blog post is new
func (me *smartProgressImporter) savePost(post smart_progress.Post) bool {
	var goalId = gophers.GetInt64FromString(post.ObjId)
	var dateTime = me.parseDateTime(post.Date).UTC()
	var text = convertHtmlToMarkdown(me.unpackRedirects(post.Msg))
	var goalPostRow = &db_objects.GoalPostRow{
		GoalId: goalId,
		Type:   post.Type,
		Text:   text,
	}
	goalPostRow.SetDateTime(dateTime)
	var isInserted = me.database.insertGoalPost(goalPostRow)
	if !isInserted {
		var defaultLanguage = base.SupportedLanguages[0]
		me.database.setGoalPostText(goalId, dateTime, defaultLanguage, text)
	}
	return isInserted
}

func (me *smartProgressImporter) saveImages(post smart_progress.Post, imageRecords []smart_progress.ImageRecord) {
	var goalId = gophers.GetInt64FromString(post.ObjId)
	var dateEpoch = me.parseDateTime(post.Date).UTC().Unix()
	for index, image := range imageRecords {
		me.database.saveGoalPostImage(db_objects.GoalPostImageRow{
			GoalId:         goalId,
			ParentDateTime: dateEpoch,
			SequenceIndex:  int64(index),
			ContentType:    image.ContentType,
			File:           image.Data,
		})
	}
}

func (me *smartProgressImporter) readGoalInfo(goalId string) (result smart_progress.GoalRecord) {
	var url = smart_progress.Url + "/goal/" + url.PathEscape(goalId)
	var body, _ = me.httpGet("Could not load goal title", url, nil)
	var document = gophers.AssertResultError(html.Parse(bytes.NewReader(body)))

	var title = ""
	if titleNode := html_util.FindElement(document, func(node *html.Node) bool { return node.Data == "title" }); titleNode != nil {
		title = html_util.NodeText(titleNode)
	}

	var descriptionHtml = ""
	if goalDescription := html_util.FindElement(document, func(node *html.Node) bool {
		return node.Data == "div" && html_util.AttrValue(node, "id") == "goal_descr"
	}); goalDescription != nil {
		if div := html_util.FindElement(goalDescription, func(node *html.Node) bool { return node.Data == "div" }); div != nil {
			descriptionHtml = strings.TrimSpace(html_util.InnerHtml(div))
		}
	}
	var description = ""
	if descriptionHtml != "" {
		description = convertHtmlToMarkdown(descriptionHtml)
	}

	var authorName = ""
	if authorWidget := html_util.FindElement(document, func(node *html.Node) bool {
		return html_util.NodeHasClass(node, "user-widget__name")
	}); authorWidget != nil {
		if link := html_util.FindElement(authorWidget, func(node *html.Node) bool { return node.Data == "a" }); link != nil {
			authorName = strings.TrimSpace(html_util.NodeText(link))
		}
	}

	result.Image = me.readGoalImage(document)
	result.Id = gophers.GetInt64FromString(goalId)
	result.Title = title
	result.Description = description
	result.AuthorName = authorName
	return
}

func (me *smartProgressImporter) readGoalImage(document *html.Node) (result smart_progress.ImageRecord) {
	var imageUrl = ""
	if link := html_util.FindElement(document, func(node *html.Node) bool {
		return node.Data == "link" && html_util.AttrValue(node, "rel") == "image_src"
	}); link != nil {
		imageUrl = html_util.AttrValue(link, "href")
	}
	if imageUrl == "" {
		panic("Cannot find image")
	}
	var body, contentType = me.httpGet("Cannot read image", imageUrl, nil)
	result.ContentType = contentType
	result.Data = body
	return
}

func (me *smartProgressImporter) saveGoalInfo(goalRecord smart_progress.GoalRecord) {
	var goalRow = &db_objects.GoalRow{
		Id:               goalRecord.Id,
		Title:            goalRecord.Title,
		Description:      goalRecord.Description,
		AuthorName:       goalRecord.AuthorName,
		ImageData:        goalRecord.Image.Data,
		ImageContentType: goalRecord.Image.ContentType,
	}
	var isInserted = me.database.insertGoal(goalRow)
	if !isInserted {
		me.database.updateGoalSmart(goalRow)
	}
}

func (me *smartProgressImporter) readAllPosts(goalId string) (allPosts []smart_progress.Post) {
	var startId = "0"
	for {
		var posts = me.readPosts(goalId, startId)
		if len(posts.Blog) == 0 {
			break
		}
		allPosts = append(allPosts, posts.Blog...)
		startId = posts.Blog[len(posts.Blog)-1].Id
	}
	return
}

func (me *smartProgressImporter) readComments(postId string) (result []smart_progress.Comment) {
	var url = smart_progress.Url + "/blog/getComments?post_id=" + postId
	var body, _ = me.httpGet("Cannot read comments", url, map[string]string{"Accept": "application/json"})
	var responseObject smart_progress.GetCommentsResponse
	gophers.AssertError(json.Unmarshal(body, &responseObject))
	return responseObject.Comments
}

func (me *smartProgressImporter) readImages(post smart_progress.Post) (imageRecords []smart_progress.ImageRecord) {
	for _, image := range post.Images {
		var url = smart_progress.Url + image.Url
		var body, contentType = me.httpGet("Cannot read image", url, nil)
		imageRecords = append(imageRecords, smart_progress.ImageRecord{ContentType: contentType, Data: body})
	}
	return
}

func (me *smartProgressImporter) readPosts(goalId string, startId string) (result smart_progress.Posts) {
	var url = smart_progress.Url + "/blog/getPosts" +
		"?obj_id=" + goalId +
		"&sorting=old_top" +
		"&start_id=" + startId +
		"&end_id=0" +
		"&step_id=0" +
		"&only_author=0" +
		"&change_sorting=0" +
		"&obj_type=0"
	var body, _ = me.httpGet("Could not load blog posts", url, map[string]string{"Accept": "application/json"})
	gophers.AssertError(json.Unmarshal(body, &result))
	return
}

func (me *smartProgressImporter) httpGet(contextMessage string, url string, headers map[string]string) (body []byte, contentType string) {
	var request = gophers.AssertResultError(http.NewRequest(http.MethodGet, url, nil))
	request.Host = smart_progress.Host
	for name, value := range headers {
		request.Header.Set(name, value)
	}
	var response = gophers.AssertResultError(http.DefaultClient.Do(request))
	defer gophers.IoClose(response.Body)
	var readErr error
	body, readErr = io.ReadAll(response.Body)
	gophers.AssertError(readErr)
	contentType = response.Header.Get("Content-Type")
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		panic(contextMessage + ": " + response.Status + "\n" + string(body))
	}
	return
}
