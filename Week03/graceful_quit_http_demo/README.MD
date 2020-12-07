## Q.基于 errgroup 实现一个 http server 的启动和关闭 ，以及 linux signal 信号的注册和处理，要保证能够 一个退出，全部注销退出。

## A. main.go
 1. 通过errgroup启动了两个http server，分别监听8080和8081端口。
 2. 通过errgruop监听 linux signal，和 两个http close的信号
 3. 可通过ctrl + c 触发SIGINT信号，触发系统关闭信号，关闭所有goroutine
 4. 可通过访问 http://127.0.0.1:8080/close 触发关闭http服务，关闭所有goroutine
 5. 可通过访问 http://127.0.0.1:8081/close 触发关闭http服务，关闭所有goroutine