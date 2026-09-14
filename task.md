# Backend crash during update job

panic: ERROR: null value in column "title" of relation "goalposts" violates not-null constraint (SQLSTATE 23502)

goroutine 1 [running[]:
github.com/hinst/go-gophers.AssertError(...)
    /go/pkg/mod/github.com/hinst/go-gophers@v0.1.32/assert.go:6
github.com/hinst/go-gophers.AssertResultError[...[](...)
    /go/pkg/mod/github.com/hinst/go-gophers@v0.1.32/assert.go:11
github.com/hinst/hinst-website/server.(*database).insertGoalPost(0x3932e4898d80, 0x3932e493e500)
    /app/server/database.goals.go:174 +0x194
github.com/hinst/hinst-website/server.(*smartProgressImporter).savePost(0x3932e489dd70, {{0x3932e4ac6240, 0x6}, {0x3932e4ac6238, 0x7}, {0x3932e4ac6248, 0x4}, {0x3932e4af5200, 0x412}, {0x3932e4afa048, ...}, ...})
    /app/server/smartProgressImporter.go:115 +0x178
github.com/hinst/hinst-website/server.(*smartProgressImporter).syncPosts(0x3932e489dd70, {0x3932e4892021, 0x6})
    /app/server/smartProgressImporter.go:42 +0xd4
github.com/hinst/hinst-website/server.(*smartProgressImporter).syncGoal(0x3932e4c03d70, {0x3932e4892021, 0x6})
    /app/server/smartProgressImporter.go:35 +0x74
github.com/hinst/hinst-website/server.(*smartProgressImporter).run(...)
    /app/server/smartProgressImporter.go:28
github.com/hinst/hinst-website/server.(*program).importSmartProgress(0x3932e4870940)
    /app/server/program.go:72 +0xc4
github.com/hinst/hinst-website/server.(*program).update(0x3932e4870940)
    /app/server/program.go:60 +0x20
github.com/hinst/hinst-website/server.Main()
    /app/server/main.go:32 +0x2e8
main.main()
    /app/main.go:8 +0x1c
stream closed: EOF for hinst-website/hinst-website-update-29822175-l9qkn (hinst-website-update)
