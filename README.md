## DHT爬虫服务
[![Build Status](https://drone.io/github.com/btlike/crawl/status.png)](https://drone.io/github.com/btlike/crawl/latest)


基于底层爬虫库，进行infohash去重并存储。如果是announce_peer请求，尝试从DHT网络(BEP0009)获取资源并存储。



## 特性

- 分表查询，高效去重
- 更新资源热度
- 实现从DHT获取metadata(BEP0009)，基于[github.com/shiyanhui/dht](https://github.com/shiyanhui/dht)


## 安装
`go get github.com/btlike/crawl`



## 配置

```
{
    "database":"root:password@tcp(127.0.0.1:3306)/torrent?charset=utf8", //数据库地址
    "elastic":"http://127.0.0.1:9200", //搜索引擎地址
    "spider_number":10, //爬虫数量
    "spider_listen_port":20000, //爬虫监听起始端口,此处将监听20000~20009
    "enable_metadata":true //是否从DHT网络获取资源(会略微降低爬虫速度)
}

```



## 示例

配置从DHT网络获取资源，这部分每天就可获取数十上百万资源，通过[storage](http://github.com/btlike/storage)从种子库还能获取等数资源。

 ![](http://77g42f.com1.z0.glb.clouddn.com/crawl.png)




## 常见问题
终于运行起了爬虫，但运行没几分钟，各种linux问题出现了，最开始应该是ulimit问题，这个问题很好解决，参考[这个文章](http://www.stutostu.com/?p=1322)。然后会出现开始大量报出：`nf_conntrack: table full, dropping packet`。这个问题参考[这个文章](http://jaseywang.me/2012/08/16/%E8%A7%A3%E5%86%B3-nf_conntrack-table-full-dropping-packet-%E7%9A%84%E5%87%A0%E7%A7%8D%E6%80%9D%E8%B7%AF/)。原因就是，

```
nf_conntrack/ip_conntrack 跟 nat 有关，用来跟踪连接条目，它会使用一个哈希表来记录 established 的记录。nf_conntrack 在 2.6.15 被引入，而 ip_conntrack 在 2.6.22 被移除，如果该哈希表满了，就会出现：nf_conntrack: table full, dropping packet。
```

解决办法很简单，我们让某些端口的流量不要被记录即可。假如我们运行100个节点，而节点监听的端口是20000到20099，我们只需要执行以下命令即可。

```
iptables -A INPUT -m state --state UNTRACKED -j ACCEPT
iptables -t raw -A PREROUTING -p udp -m udp --dport 20000 -j NOTRACK
...... //从端口20000一直到20099，每个端口一行
iptables -t raw -A PREROUTING -p udp -m udp --dport 20099 -j NOTRACK
```
