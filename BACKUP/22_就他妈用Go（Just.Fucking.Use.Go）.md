# [就他妈用Go（Just Fucking Use Go）](https://github.com/IjichiNijika99/go-gitblog/issues/22)

> 原文：[Just Fucking Use Go - Blain Smith May 08 2026](https://blainsmith.com/articles/just-fucking-use-go/)
## 就他妈用Go
嘿，傻缺。你知道什么东西能在两秒内编译完成，以单文件二进制的形式部署，并且不会在凌晨 3 点因为 npm 上某个间接依赖被下架就当场崩成狗吗？是 Go。就像 HTML 从该死的互联网诞生之初就一直坐在那里，等着你停止把前端搞得那么复杂一样，Go 也在那里坐了十多年，等着你停止把后端搞得那么复杂。
但你偏不。你非要在那儿把 15 个 Node 包、3 个 TypeScript 构建工具和一个 Kubernetes 集群缝合在一起，就为了提供一个他妈的表单服务。你雇了一个平台团队来给你那堆 Rails 单体屎山当保姆。你甚至说服了你的 CTO，说写一个每秒顶多 40 个请求的 CRUD（增删改查）应用必须得上 Rust。恭喜你，混蛋，你把自己玩进去了。

### 语言是故意这么无聊的
你知道为什么 Go 让人觉得无聊吗？因为它是真无聊，而这他妈的正是它的意义所在。没有装饰器（decorators），没有元类（metaclasses），没有宏（macros），没有 trait、monad，也没有这周 Haskell 圈子里又在沉迷的任何该死的抽象概念。它只有结构体（structs）、函数（functions）、接口（interfaces）、协程 (goroutines) 和 通道 (channels)，就这些。你可以在午休时间读完它的语言规范，下午就能直接开始干活出产出。
无聊意味着你上个月刚招的初级开发能看懂主程两年前写的代码。只有一种格式化代码的方式，而且 gofmt 早就帮你做好了。你那些“自作聪明”的同事没法在代码库里塞进 17 层抽象，因为语言根本不让他这么干。当没人对着自己的小聪明流口水时，这才是交付该有的样子。

### 标准库就是框架
别再找什么框架了，你这个纯种榆木脑袋。
标准库就是框架。
```
package main

import (
    "embed"
    "html/template"
    "net/http"
)

//go:embed templates/*.html
var files embed.FS

var tmpl = template.Must(template.ParseFS(files, "templates/*.html"))

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        tmpl.ExecuteTemplate(w, "index.html", map[string]string{
            "Name": "asshole",
        })
    })

    http.ListenAndServe(":8080", nil)
}
```
这就是一个能跑的 Web 应用，HTML 模板直接编译进了二进制文件里。没有 webpack，没有 Vite，没有“开发服务器”，没有体积大得像辆他妈的大众汽车一样的 node_modules。你敲下 go build，然后你就只需交付这一个文件，把它扔到服务器上，完事。
你想要数据库？用 database/sql；要 JSON？用 encoding/json；想调用其他服务？net/http 同时也是个客户端；想同时干五件事？在前面甩个 go 关键字；测试？用 go test；基准测试？用 go test -bench；性能分析？pprof 早就准备好在那儿嘲笑你用 console.log 进行调试了。

### 标准库也他妈的极度硬核
io.Reader 和 io.Writer 是两个各只有一个方法的接口，但它们也是为什么你能用三行代码，不用过脑子就能把一个 HTTP 响应体管道传输给一个 gzip writer，再写入到磁盘文件中的原因。整个生态系统里所有正儿八经的包都在说这门“语言”。一旦你顿悟了这一点，你就会发现 Go 一半的“魔法”不过是这两个接口在到处串场罢了。
context.Context 是你用来取消操作的玩意。用户关了浏览器标签页，请求上下文被取消，数据库查询被取消，下游的 HTTP 调用也被取消。一路取消到底。没有 goroutine 泄漏，没有吞噬你连接池的僵尸查询。你把它作为第一个参数传进去，并且尊重它的规则，整个 API 就这么简单。
encoding/json、encoding/xml、encoding/csv、encoding/binary，全都在标准库里。相同的结构体标签模式，相同的“解码到指针”的人体工学设计，学会了一个，你基本上就全懂了。

### 不会让你痛哭流涕的并发
Goroutine 不是线程。它们是由运行时（runtime）复用并映射到 OS 线程上的、带有自身栈的执行体，启动一个大约只需消耗 2KB 的内存。你可以在一台笔记本上轻松拉起十万个 goroutine。拿你的 Node 事件循环试试看，看着它当场拉胯。
Channel 是 goroutine 之间带类型的管道。你在一头发送，在另一头接收，运行时会处理好所有同步问题。如果你需要的是共享状态，sync.Mutex 就在那儿，而且竞态检测器（race detector）会在你搞砸的时候立刻提醒你。
```
results := make(chan string, len(urls))
for _, url := range urls {
    go func(u string) {
        resp, _ := http.Get(u)
        results <- resp.Status
    }(url)
}
for range urls {
    fmt.Println(<-results)
}
```
这就是一个并行的 HTTP 抓取器。没有依赖库，没有框架，没有 async/await 那套繁文缛节。语言本身就搞定了这一切。

### 一个真实例子，而不是什么 hello-world
这里有一个从 Postgres 读取数据并渲染 HTML 的 CRUD 路由，全部代码都在这。
```
//go:embed templates/*.html
var tmplFS embed.FS

var tmpl = template.Must(template.ParseFS(tmplFS, "templates/*.html"))

type Post struct {
    ID    int
    Title string
    Body  string
}

func postsHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        rows, err := db.QueryContext(r.Context(),
            "SELECT id, title, body FROM posts ORDER BY id DESC LIMIT 50")
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        defer rows.Close()

        var posts []Post
        for rows.Next() {
            var p Post
            if err := rows.Scan(&p.ID, &p.Title, &p.Body); err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
            posts = append(posts, p)
        }

        tmpl.ExecuteTemplate(w, "posts.html", posts)
    }
}
```
数据库、模板和一个 HTTP 处理器，一屏就能放下。请求上下文（Context）被连通到了查询中，所以关闭的连接能直接取消 SQL 查询。没有 ORM，没有依赖注入（DI）容器，没有 service 层，没有那种装了 17 个抽象基类的 controllers/ 目录。你从上读到下，就能清清楚楚地知道它到底干了什么。

### 不会毁了你休闲时间的依赖管理
go mod init，搞定。你的依赖活在 go.mod 和 go.sum 里。sum 文件是一份包含了你实际获取内容的密码学记录，所以当有人给你搞一出 left-pad 下架的戏码时，你立刻就能发现。这里没有 node_modules 目录。在开发环境和 CI 之间不存在锁文件漂移（lockfile drift）。这里没有同级依赖（peer dependencies）、没有可选依赖（optional dependencies）、没有 devDependencies，也没有 peerDependenciesMeta。只有一个文件列出你用了什么，以及另一个文件证明你拿到了你期望的东西。
想要离线构建？go mod vendor 会把所有东西扔进 vendor/ 目录，工具链会自动去使用它。整个项目，包含所有依赖项在内，一个 tar 压缩包就能装下。你们的安全团队会感激涕零的。

### 工具链自带于编译器中
gofmt 负责格式化你的代码，没有任何争论的余地，不会有什么关于 .prettierrc 的圣战。格式定下来就是定下来了，大家都这么用。你的 diff 代码变更会保持得很小，因为没人会去瞎调空格和换行。
go vet 揪出明显的低级错误；go test 跑你的测试；go test -race 带着竞态检测器跑测试，找出那些你以为不存在的数据竞争；go test -bench 跑基准测试；go test -cover 告诉你遗漏了哪些测试覆盖；go tool pprof 则通过你用两行代码接进去的 HTTP 端点，直接为你生成运行中生产服务的 CPU 和内存火焰图。
所有这些都不是第三方的，不是什么插件，也不需要你去维护一个配置文件。它开箱即用。

### 部署就是一条复制命令
这是让 Rails 和 Node 开发者在生理上感到愤怒的部分。你构建一个 Go 二进制文件，把它复制到服务器上，运行它。
```
GOOS=linux GOARCH=amd64 go build -o myapp ./cmd/myapp
scp myapp user@server:/usr/local/bin/
ssh user@server 'systemctl restart myapp'
```
三条命令，完事。没有 Dockerfile，没有多阶段构建，不会在每个星期二收到基础镜像的 CVE 漏洞警告，没有 Kubernetes manifest（清单），没有 Helm chart，没有 ArgoCD，没有 Service Mesh，没有 Sidecar。
一个 12MB 的静态链接二进制文件和一个 20 行的 systemd unit 文件，这就是一个生产级别的部署，它的寿命会比你的职业生涯还长。唯一让你去碰 Docker 的理由，是你们的运维团队在合同里被要求必须用它，即便如此，你也可以直接把二进制文件塞进一个 FROM scratch 的空镜像里，然后打卡下班。

### “那 Rails / Django / Express / Next 呢？”
它们怎么了？ Rails 需要一个包含 Capistrano、三个配置文件以及献祭一只山羊的部署仪式。Django 想让你学习它的 ORM、它的 admin 界面、它的中间件系统，以及它对世界上所有事情的武断看法。Express 是靠着 npm audit 警告和祈祷拼凑起来的。Next.js 每六个月改一次路由约定，然后对你搞煤气灯操纵（PUA），死不承认。
你的 Go 二进制文件根本不在乎。它编译过了，它就能跑，而且五年后它依然能在现在还不存在的硬件上跑。你用的框架呢？到圣诞节就被弃用了，然后维护者还会跑去 Medium 上发篇小作文哭诉自己有多倦怠。

### “但是我们要搞微服务！”
不。
写你他妈的单体应用（Monolith）。一个 Go 二进制文件， 一个 Postgres 数据库，如果非要不可的话再加一个 Redis，用提供 JSON API 的同一个端口来伺服 HTML。把它跑在一个月费连你平时买燕麦奶的钱都不如的单台 VPS 上，毫不费力地把并发扩展到每秒一万次请求而面不改色，因为 Go 生来就是干这个的，并且 goroutine 便宜得要命。
当你真的需要把它拆分的时候——虽然你大概率用不着——拆分一个 Go 单体应用也不过是把对应的包移到它们自己的代码库里。接口早就定义好了，你甚至都没刻意去设计就已经为拆分做好了准备，因为语言本身就会逼你这么做。

### “但是泛型！但是错误处理！但是没有异常机制！”
if err != nil 是特性（feature），不是 Bug。它强迫你去审视每一个可能出错的地方，并决定到底怎么处理它。你那层层嵌套的 try/catch 地狱并没有让错误消失，它只是把错误藏了起来，直到凌晨 2 点在生产环境里炸开。
泛型在 1.18 版本就落地了，它们挺好的，需要用的时候就用，别再抱怨了。

### 就他妈用 Go
别再装作你需要一个框架了。你不需要微服务，不需要用 Rust 重写，也不需要上周二刚发布的、号称能拯救你的某个 JavaScript 元框架——毕竟前面那六个框架都没救得了你。
打开你的编辑器，运行 go mod init，写一个 main.go，嵌入你的模板，然后编译，把这该死的东西发布上线吧。
最无聊的选择就是最正确的选择，历来如此。






