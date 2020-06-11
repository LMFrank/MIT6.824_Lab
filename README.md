## MIT6.824

### Lab1-MapReduce

**视频资料**：

[MIT 6.824 Distributed Systems Spring 2020](https://www.bilibili.com/video/BV1x7411M7Sf)

**参考资料**：（均在Lab1目录下）

1. 文献：MapReduce: Simplified Data Processing on Large Clusters

2. 文献翻译：
   - MapReduce:在大型集群上简化数据处理（1）：https://mp.weixin.qq.com/s/sChCf07SxhTudxFIKd8pgA
   - MapReduce:在大型集群上简化数据处理（2）：https://mp.weixin.qq.com/s/h43tPiycGrKf9089pML2tw
   - MapReduce:在大型集群上简化数据处理（3）：https://mp.weixin.qq.com/s/CJrvjqpPbUIrgDvSFjeZAQ
   - MapReduce:在大型集群上简化数据处理（4）：https://mp.weixin.qq.com/s/jA4FeYBjb6fd_JyP6fzQ5g

3. 实验文档：包含中英文

MapReduce: Simplified Data Processing on Large Clusters

**Note**：

按照文档提示运行`go build -buildmode=plugin ../mrapps/wc.go`是有问题的，目前的解决方案有两种：

1. 设置`GO111MODULE=off`，将源码放在`GOPATH`下，但是看反馈，依然有报错现象，推荐第2种解决方法
2. 启动`GO111MODULE`，也就是使用`go mod`管理包
   - 在6.824根目录下打开cmd，输入`go mod init <name>`
   - 更改源码中的导包方式，将相对路径改为绝对包路径
   - 使用`go build -buildmode=plugin 6.824/src/mrapps/wc.go`编译

