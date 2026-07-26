# Improvement Plan: Replace `StaticPath` with `LangPath`

## Summary

Replace the current `WebPath` + `StaticPath` pair in `page_data.Base` with `WebPath` + `LangPath`:

| Field | Meaning | Example value(s) |
|-------|---------|-------------------|
| **WebPath** | Full URL prefix for this site instance (the HTTP-level deployment path) | `""`, `/hinst-website`, `/blog` |
| **LangPath** | Language-specific URL segment appended after `WebPath` | `""`, `/de`, `/ru` |

**Design rule:**
- All language-dependent resources → referenced with `{{.WebPath}}{{.LangPath}}...`
- All language-independent resources (static assets) → referenced with `{{.WebPath}}...`

## Problem with Current Design

The current `StaticPath` field is overloaded and confusing:

1. **Semantic mismatch**: `StaticPath` sounds like it's for static assets, but in `header.html` it's used for *navigation links* to language switcher pages (`href="{{.StaticPath}}/de"`).
2. **Inconsistency between renderers**: The HTTP-level `webApp.webPath` is `/hinst-website`, while the template-level "WebPath" from `goalRenderer.webPath(lang)` returns only the language segment (`""` or `/de`). There's no single field that represents the actual site URL prefix.
3. **Currently broken for subdirectory deployment**: With both fields effectively always empty, deploying under a sub-path like `/blog/` would break links because neither field carries the full URL prefix information correctly.

## New Design

```go
type Base struct {
    Id         int64
    WebPath    string   // Full site URL prefix: "" or "/hinst-website" or "/blog"
    LangPath   string   // Per-language segment appended after WebPath: "" (English) or "/de", "/ru"
    ...
}
```

### Usage rules in templates

| Resource type | Template expression | Example output |
|---|---|---|
| **Language-independent** (CSS, JS, images, favicon) | `{{.WebPath}}/static/...` | `/hinst-website/static/css/style.css` or `/static/css/style.css` |
| **Language-dependent page links** | `{{.WebPath}}{{.LangPath}}/personal-goals/...` | `/blog/de/personal-goals/123.html` or `/de/personal-goals/123.html` |
| **Language switcher navigation** | `{{.WebPath}}{{.LangPath}}/de` (etc.) | `/hinst-website/de` or `/de` |
| **Favicon link to home** | `{{if .WebPath}}{{.WebPath}}{{else}}/{{end}}` (unchanged) | Same as before |

## File-by-file Change List

### 1. `server/page_data/base.go` — Struct definition

```go
type Base struct {
    Id         int64
    WebPath    string   // was: "WebPath" = language segment; now = full site prefix
    LangPath   string   // NEW: per-language URL segment ("", "/de", "/ru")
    
    SettingsSvg template.HTML
    MenuSvg     template.HTML
    InfoSvg     template.HTML
}
```

### 2. `server/goalRenderer.go` — Renderer logic

**Change:** Replace the `staticPath` parameter on all render methods, remove `me.webPath(lang)`, and introduce a unified prefix builder.

```go
func (me *goalRenderer) getBaseTemplate(webPath string, langPath string) page_data.Base {
    return page_data.Base{
        Id:          me.elementId.Add(1),
        WebPath:     webPath,      // now = full site URL prefix
        LangPath:    langPath,     // new: language segment only
        SettingsSvg: template.HTML(gophers.ReadTextFile("pages/static/images/settings.svg")),
        MenuSvg:     template.HTML(gophers.ReadTextFile("pages/static/images/menu.svg")),
        InfoSvg:     template.HTML(gophers.ReadTextFile("pages/static/images/info.svg")),
    }
}
```

Each render method signature changes from `(lang, staticPath)` to `(lang, webPath, langPath)`:

| Method | Before | After |
|--------|--------|-------|
| `renderHomePage` | `(lang, staticPath string)` | `(lang, webPath, langPath string)` |
| `renderGoalPage` | `(lang, staticPath string, goalId int64)` | `(lang, webPath, langPath string, goalId int64)` |
| `renderGoalPostPage` | `(lang, staticPath string, goalId int64, dateTime time.Time)` | `(lang, webPath, langPath string, goalId int64, dateTime time.Time)` |

The old `me.webPath(lang)` method is removed entirely. The `webPath` parameter is now provided by the caller (the HTTP handler or static generator) and represents the full site prefix.

### 3. `server/webStaticGoals.go` — Static site generator

**Change:** Pass the appropriate `webPath` and `langPath` values to render methods.

```go
func (me *webStaticGoals) getLanguagePath(tag language.Tag) string {
    if tag == language.English { return "" }
    return "/" + tag.String()  // e.g., "/de", "/ru"
}
```

- `langPath` = result of `getLanguagePath(lang)` (per-language segment, same logic as before)
- `webPath` = `""` for static file output (static files are written relative to their own directory; assets use local relative paths which map to `{{.WebPath}}/static/...` → `/static/...`)

Example calls become:
```go
me.renderer.renderHomePage(lang, "", me.getLanguagePath(lang))
me.renderer.renderGoalPage(lang, "", me.getLanguagePath(lang), goalId)
me.renderer.renderGoalPostPage(lang, "", me.getLanguagePath(lang), goalId, dateTime)
```

### 4. `server/webApp.go` — HTTP server

**Change:** The existing `webApp.webPath` (currently `/hinst-website`) already represents the full site URL prefix. Pass it to renderers as the new `WebPath`.

The current code:
```go
me.addFunctions(me.webPath, appGoals.init(me.db))
```

This stays the same — `me.webPath` is already the correct full-site prefix for HTTP routing. The change is only in how renderers receive it (now as a separate `webPath` argument instead of deriving language segments internally).

### 5. Template changes (`pages/html/templates/*.html`)

#### `template.html`
```html
<!-- Before -->
<link rel="icon" href="{{.StaticPath}}/static/images/favicon.png" />
<link rel="stylesheet" href="{{.StaticPath}}/static/css/minstyle.css" />
...
<script src="{{.StaticPath}}/static/js/dist/main.js"></script>

<!-- After -->
<link rel="icon" href="{{.WebPath}}/static/images/favicon.png" />
<link rel="stylesheet" href="{{.WebPath}}/static/css/minstyle.css" />
...
<script src="{{.WebPath}}/static/js/dist/main.js"></script>
```

#### `header.html`
```html
<!-- Language switcher: language-dependent links -->
<a href="{{if .WebPath}}{{.WebPath}}{{else}}/{{end}}">   <!-- favicon home link (unchanged) -->
<img src="{{.WebPath}}/static/images/favicon.png" ... />  <!-- asset (was StaticPath) -->

<!-- Before -->
<a href="{{if .StaticPath}}{{.StaticPath}}{{else}}/{{end}}">English</a>
<a class="ms-btn ms-outline" href="{{.StaticPath}}/de"> German </a>
<a class="ms-btn ms-outline" href="{{.StaticPath}}/ru"> Russian </a>

<!-- After -->
<a href="{{if .WebPath}}{{.WebPath}}{{else}}/{{end}}">English</a>
<a class="ms-btn ms-outline" href="{{.WebPath}}{{.LangPath}}/de"> German </a>
<a class="ms-btn ms-outline" href="{{.WebPath}}{{.LangPath}}/ru"> Russian </a>
```

#### `goalList.html`
```html
<!-- Before -->
<a href="{{$.WebPath}}/personal-goals/{{.Id}}.html">

<!-- After (no change — WebPath was already the right field for page links) -->
<a href="{{$.WebPath}}{{$.LangPath}}/personal-goals/{{.Id}}.html">
```

#### `goalPosts.html`
```html
<!-- Before -->
<a href="{{$.WebPath}}/personal-goals/{{$.GoalId}}/{{.DateTime}}.html">

<!-- After -->
<a href="{{$.WebPath}}{{$.LangPath}}/personal-goals/{{$.GoalId}}/{{.DateTime}}.html">
```

#### `goalPost.html` (images — language-dependent resource paths)
```html
<!-- Before -->
href="{{$.StaticPath}}/personal-goals/image/{{$.GoalId}}/{{$.DateTime}}/{{.}}.jpg"
src="{{$.StaticPath}}/personal-goals/image/{{$.GoalId}}/{{$.DateTime}}/{{.}}.jpg"

<!-- After: images are per-language (stored in language-specific directories) -->
href="{{$.WebPath}}{{$.LangPath}}/personal-goals/image/{{$.GoalId}}/{{$.DateTime}}/{{.}}.jpg"
src="{{$.WebPath}}{{$.LangPath}}/personal-goals/image/{{$.GoalId}}/{{$.DateTime}}/{{.}}.jpg"
```

### 6. `server/siteMapBuilder.go` — Sitemap generation (no change needed)

The sitemap builder uses `me.webPath` directly from the public URL (`https://hinst.github.io`) and constructs URLs by appending relative file paths. It does not use `page_data.Base`, so it needs no changes.

### 7. `server/staticFilesUpdate.go` — Static file upload (no change needed)

Uses its own `webPath` field for sitemap URL construction, independent of `page_data.Base`. No changes required.

## Correctness Improvements

1. **Clearer semantics**: `WebPath` = "where is this site deployed?" (e.g., `/`, `/hinst-website`, `/blog`). `LangPath` = "which language variant?". Each field has a single, well-defined responsibility.

2. **Subdirectory deployment support**: To deploy under `/blog/`, you only set `WebPath = "/blog"` at the HTTP handler level — all templates automatically produce correct paths for both assets (`/blog/static/...`) and pages (`/blog/de/personal-goals/...`).

3. **Correctness of image URLs**: Currently images use `StaticPath` which conflates asset path with page path semantics. With the new design, images are correctly treated as language-dependent (each language subdirectory has its own copy), referenced via `WebPath + LangPath`.

4. **Language switcher correctness**: The switcher links now explicitly compose `WebPath + LangPath`, making it clear that switching languages produces a different *page URL*, not just a different asset base path.

## Migration Checklist

- [ ] Update `server/page_data/base.go`: rename `StaticPath` → add `LangPath`
- [ ] Update `server/goalRenderer.go`: change method signatures, remove old `webPath(lang)` helper, update `getBaseTemplate`
- [ ] Update `server/webStaticGoals.go`: pass new parameters to renderer methods
- [ ] Update all 5 template files (`template.html`, `header.html`, `goalList.html`, `goalPosts.html`, `goalPost.html`)
- [ ] Verify static generation still produces correct relative paths (since static output is per-language directory with no URL prefix, `WebPath=""` works)
- [ ] Run existing tests / build to confirm no compilation errors
