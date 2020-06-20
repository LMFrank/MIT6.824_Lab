## MIT6.824

**课程主页**：http://nil.csail.mit.edu/6.824/2020/index.html

**视频资料**：

[MIT 6.824 Distributed Systems Spring 2020](https://www.bilibili.com/video/BV1x7411M7Sf)

### Lab1-MapReduce
**参考资料**：（均在Lab1目录下）

1. 文献：[MapReduce: Simplified Data Processing on Large Clusters](https://github.com/LMFrank/MIT6.824_Lab/blob/master/Lab1/mapreduce.pdf)

2. 文献翻译：
   - MapReduce:在大型集群上简化数据处理（1）：https://mp.weixin.qq.com/s/sChCf07SxhTudxFIKd8pgA
   - MapReduce:在大型集群上简化数据处理（2）：https://mp.weixin.qq.com/s/h43tPiycGrKf9089pML2tw
   - MapReduce:在大型集群上简化数据处理（3）：https://mp.weixin.qq.com/s/CJrvjqpPbUIrgDvSFjeZAQ
   - MapReduce:在大型集群上简化数据处理（4）：https://mp.weixin.qq.com/s/jA4FeYBjb6fd_JyP6fzQ5g

3. 实验文档：包含中英文
[6.824 Lab 1: MapReduce](http://nil.csail.mit.edu/6.824/2020/labs/lab-mr.html)

**Note**：

按照文档提示运行`go build -buildmode=plugin ../mrapps/wc.go`是有问题的，目前的解决方案有两种：

1. 设置`GO111MODULE=off`，将源码放在`GOPATH`下，但是看反馈，依然有报错现象，推荐第2种解决方法

2. 启动`GO111MODULE`，也就是使用`go mod`管理包
   - 在6.824根目录下打开cmd，输入`go mod init <name>`
   - 更改源码中的导包方式，将相对路径改为绝对包路径
   
3. 程序使用Go内置的rpc作为通信手段，GOB是go语言使用的二进制序列化协议，使用rpc传输的字段应该使用大写，否则不会被识别

4. 根据测试规则，我们并不需要做合并，paper中也提到经过mapreduce框架处理后的数据可能也直接用于其他场景了。

   ```shell
   # test-mr.sh会先生成正确的结果文件
   ../mrsequential ../../mrapps/wc.so ../pg*txt || exit 1
   sort mr-out-0 > mr-correct-wc.txt
   rm -f mr-out*
   //
   ...
   //
   # 然后会将复合mr-out*的文件进行排序，并合并生成mr-wc-all这个文件
   # 再使用cmp命令比对
   sort mr-out* | grep . > mr-wc-all
   if cmp mr-wc-all mr-correct-wc.txt
   then
     echo '---' wc test: PASS
   else
     echo '---' wc output is not the same as mr-correct-wc.txt
     echo '---' wc test: FAIL
     failed_any=1
   fi
   ```

   

**代码实现**：

已通过所有测试

- Master.go
  - MakeMaster：开启一个master
  - server：将进程注册到rpc中
  - Work：向worker分配任务，更新map和reduce任务的状态，使用context管理超时
  - Commit：确认task完成
- Worker.go
  - Worker：启动死循环，不断向master发出请求
  - MapWork：将master发来的文件，经过mapf处理，然后生成中间文件
  - ReduceWork：将master发来的任务，对应的中间文件，合成最终output

参考：https://www.yuque.com/abser/blog/lab1mapreduce