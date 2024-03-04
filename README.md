# 参考

1. [分布式爬虫实战 | 极客时间](https://time.geekbang.org/column/intro/100124001)
2. [dreamerjackson/crawler](https://github.com/dreamerjackson/crawler)

# 专栏摘抄

1. Go 语言对协程的设计和对于 I/O 多路复用巧妙的封装，实现了同步编程的语义，但背后实则是异步 I/O 的处理模式。在减轻开发者心理负担的同时，提升了网络 I/O 的处理效率。- 出自 [16｜网络爬虫： 一次HTTP请求的魔幻旅途]
2. 空接口是实现反射的基础，因为空接口中会存储动态类型的信息，这为我们提供了复杂、意想不到的处理能力和灵活性。- 出自 [21｜采集引擎：实战接口抽象与模拟浏览器访问]
3. Rob Pike 在 2014 年的一篇博客中提到了一种优雅的处理方法叫做函数式选项模式(Functional Options)。这种模式展示了闭包函数的有趣用途，目前在很多开源库中都能看到它的身影。- 出自 [28｜调度引擎：负载均衡与调度器实战]

# 官方库 | 第三方库

- [官方处理字符集库1](golang.org/x/net/html/charset)
- [官方处理字符集库2](golang.org/x/text/encoding)
- [Xpath](https://github.com/antchfx/htmlquery)
- [CSS选择器](https://github.com/PuerkitoBio/goquery)
- [Chrome DevTools 协议库](https://github.com/chromedp/chromedp)
- [高性能协程池](https://github.com/panjf2000/ants)

# 其他

- [解决WSL2 Google Chrome找不到执行路径问题](https://github.com/oven-sh/bun/issues/5416)

# tag 标签

```sh
git commit -am "print resp"
git tag -a v0.0.1 -m "print resp"
git push origin v0.0.1
```