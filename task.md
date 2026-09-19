# Debugging network error

## Log:

2026/09/19 13:32:46 Sync complete: goal=247488 posts=281 new=0
2026/09/19 13:33:14 Sync complete: goal=416007 posts=126 new=0
2026/09/19 13:33:14 Generated translated text for 0 of 407 posts
2026/09/19 13:33:14 Generated title for 0 of 407 posts
panic: Post "http://prettier-server.default.svc:3000?filename=index.html": EOF

goroutine 1 [running[]:
github.com/hinst/go-gophers.AssertError(...)
    /go/pkg/mod/github.com/hinst/go-gophers@v0.1.35/assert.go:6
github.com/hinst/go-gophers.AssertResultError[...[](...)
    /go/pkg/mod/github.com/hinst/go-gophers@v0.1.35/assert.go:11
github.com/hinst/hinst-website/server.formatHtml({0x687ab4a90000, 0x196f})
    /app/server/web.go:48 +0x4d8
github.com/hinst/hinst-website/server.(*webStaticGoals).formatHtml(0x687ab47dbd70, {0x687ab4a90000, 0x196f})
    /app/server/webStaticGoals.go:129 +0x40
github.com/hinst/hinst-website/server.(*webStaticGoals).generateGoalPost(0x687ab47dbd70, {0x5907?, 0x6?, {0x0?, 0x0?}}, {0x687ab48ce450, 0x23}, 0x65907, 0x6891f8aa)
    /app/server/webStaticGoals.go:86 +0xe8
github.com/hinst/hinst-website/server.(*webStaticGoals).generateGoal(0x687ab47dbd70, {0x2?, 0x0?, {0x0?, 0x0?}}, {0x687ab48ce450, 0x23}, {0x65907, {0x687ab49c0cf0, 0x3, ...}, ...})
    /app/server/webStaticGoals.go:76 +0x18c
github.com/hinst/hinst-website/server.(*webStaticGoals).generate(0x687ab47dbd70, {0x778?, 0xae?, {0x0?, 0x0?}})
    /app/server/webStaticGoals.go:60 +0x254
github.com/hinst/hinst-website/server.(*webStaticGoals).run(0x687ab47dbd70)
    /app/server/webStaticGoals.go:34 +0xbc
github.com/hinst/hinst-website/server.(*program).generateStatic(0x687ab47add40?, {0x687ab47d4120?, 0x687ab4a1fdd8?})
    /app/server/program.go:113 +0xec
github.com/hinst/hinst-website/server.(*program).update(0x687ab47add40)
    /app/server/program.go:63 +0x5c
github.com/hinst/hinst-website/server.Main()
    /app/server/main.go:32 +0x2e8
main.main()
    /app/main.go:8 +0x1c
stream closed: EOF for hinst-website/hinst-website-update-manual-tqk-6t7js (hinst-website-update)

## Task

Please proceed with guessing the likely cause of the error.
Kubernetes deployment is defined in repository `C:\Dev\orange-pi-kubernetes`
