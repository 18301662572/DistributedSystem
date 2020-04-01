# token 身份验证

https://blog.csdn.net/wangshubo1989/article/details/74529333

token验证是一种web常用的身份验证手段，在这里不讨论它的具体实现 <br/>
我需要在golang里实现token验证，Web框架是Gin（当然这与框架没有关系）<br/>

```text
步骤如下
1.从request获取tokenstring
2.将tokenstring转化为未解密的token对象
3.将未解密的token对象解密得到解密后的token对象
4.从解密后的token对象里取参数
```