# Debugging

2026/09/09 20:05:12 Sync complete: goal=247488 posts=280 new=0
2026/09/09 20:05:41 Sync complete: goal=416007 posts=125 new=0
panic: Post "http://192.168.0.30:11434/v1/chat/completions": http: ContentLength=2691 with Body length 0

goroutine 1 [running[]:
github.com/hinst/go-gophers.AssertError(...)
    /go/pkg/mod/github.com/hinst/go-gophers@v0.1.30/assert.go:6
github.com/hinst/go-gophers.AssertResultError[...[](...)
    /go/pkg/mod/github.com/hinst/go-gophers@v0.1.30/assert.go:11
github.com/hinst/hinst-website/server.(*translator).translateText(0x502401aedd80, {0x502401c7b500, 0x95d}, {0xe0?, 0x0?, {0x0?, 0x0?}})
    /app/server/translator.go:67 +0x4fc
github.com/hinst/hinst-website/server.(*translator).translate(0x502401be7d80, 0x502401e14370, {0x9b68?, 0x1e9?, {0x0?, 0x0?}})
    /app/server/translator.go:49 +0x54
github.com/hinst/hinst-website/server.(*translator).run.func1(0x502401e14370?)
    /app/server/translator.go:36 +0xc4
github.com/hinst/hinst-website/server.(*database).forEachGoalPost(0x502401ba3340, 0x502401be7d08, {0x502401c306c0?, 0x0?}, 0x0)
    /app/server/database.goals.go:79 +0x178
github.com/hinst/hinst-website/server.(*translator).run(0x502401e99d80)
    /app/server/translator.go:29 +0x9c
github.com/hinst/hinst-website/server.(*program).updateTranslations(0x502401c94540)
    /app/server/program.go:81 +0xa4
github.com/hinst/hinst-website/server.(*program).update(0x502401c94540)
    /app/server/program.go:61 +0x28
github.com/hinst/hinst-website/server.Main()
    /app/server/main.go:32 +0x2e8
main.main()
    /app/main.go:8 +0x1c
stream closed: EOF for hinst-website/hinst-website-update-manual-nrz-mb4vm (hinst-website-update)

