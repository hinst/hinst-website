# Debugging

Looking at log:

```
	2026/10/30 21:16:06 Sync complete: goal=247488 posts=285 new=0
	2026/10/30 21:16:36 Sync complete: goal=416007 posts=132 new=0
	panic: Post "http://192.168.178.26:11434/v1/chat/completions": dial tcp 192.168.178.26:11434: connect: no route to host

	goroutine 1 [running[]:
	github.com/hinst/go-gophers.AssertError(...)
		/go/pkg/mod/github.com/hinst/go-gophers@v0.1.35/assert.go:6
	github.com/hinst/go-gophers.AssertResultError[...[](...)
		/go/pkg/mod/github.com/hinst/go-gophers@v0.1.35/assert.go:11
	github.com/hinst/hinst-website/server.(*translator).translateText(0x544381df7d80, {0x54438208a800, 0x773}, {0x0?, 0x0?, {0x0?, 0x0?}})
		/app/server/translator.go:69 +0x3ec
	github.com/hinst/hinst-website/server.(*translator).translate(0x544381ee5d80, 0x544381e52300, {0x0?, 0x0?, {0x0?, 0x0?}})
		/app/server/translator.go:48 +0x88
	github.com/hinst/hinst-website/server.(*translator).run.func1(0x544381e52300)
		/app/server/translator.go:34 +0x100
	github.com/hinst/hinst-website/server.(*database).forEachGoalPost(0x544381df31e0, 0x544381ee5d08, 0x0)
		/app/server/database.goalPosts.go:94 +0x1b8
	github.com/hinst/hinst-website/server.(*translator).run(0x0?)
		/app/server/translator.go:29 +0x4c
	github.com/hinst/hinst-website/server.(*program).updateTranslations(0x544381ffc540)
		/app/server/program.go:81 +0xa4
	github.com/hinst/hinst-website/server.(*program).update(0x544381ffc540)
		/app/server/program.go:61 +0x28
	github.com/hinst/hinst-website/server.Main()
		/app/server/main.go:32 +0x2e8
	main.main()
		/app/main.go:8 +0x1c
	stream closed: EOF for hinst-website/hinst-website-update-29889915-ksx5g (hinst-website-update)
```

Related directories:
* Go-gophers: `C:\Dev\go-gophers`
* Deployment: `C:\Dev\orange-pi-kubernetes`

Please proceed with debugging
