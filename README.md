# haochain
基于超级账本开发一个区块链应用--票据背书



## 先决条件

- Docker version 17.03.0-ce or greater is require.
- Docker-compose version 1.8 or greater is required.
- Go version 1.9.x or greater is required.  


## 架构设计


### 1.1应用层
用户登陆，发布票据，发起背书，拒绝背书，我的票据，待签收票据，票据详情

### 1.2业务层
用户管理，票据管理，其他

### 1.3智能合约(链码)
票据发布，票据背书，票据背书签收，票据背书拒签，查询持票人票据，查询待签收票据，根据票据号码查询票据详情

### 1.4区块链底层平台
CA，区块链，账本


### 每层的主要功能如下：

- 区块链底层平台: 提供分布式共享账本的维护、状态数据库维护、智能合约的全 生命周期管理等区块链功能，实现数据的不可篡改和智能合约的业务逻辑。根 据第11章的内容搭建区块链网络以后，默认就提供了这部分功能。另外，通过 fabric-ca提供成员注册和注销等功能。

- 智能合约: 智能合约通过链码来实现，包括票据发布、票据背书、票据背书签 收、票据背书拒绝等链码调用功能，链码查询包括查询持票人票据、查询待签 收票据、根据链码号码查询票据信息等。

- 业务层: 业务层是应用程序的后端服务，给Web应用提供RESTful的接口，处理 前端的业务请求。后端服务的基本功能包括用户管理和票据管理，通过 Hyperledger Fabri提供的Go SDK和区块链网络进行通信。业务层也可以和其 他的业务系统进行交互。

- 应用层，Web应用采用jQuery+HTML+CSS 的前端架构编写页面，提供用户交 互的界面操作，包括用户操作的功能

说明：业务操作的功能。用户是内置的，只提供用户登录和用户退出操作。业务操作 包括发布 查询持票人持有的票据、发起票据背书、查询待签收票据、签收票据 背书、拒绝票据背书等功能。 各个层之间采用不同的接口，业务层的Go SDK、智能合约和区块链底层平台之间用gRPC的接口。


### 安装向导

参考<br/>
[Tutorial Hyperledger Fabric: How to build your first network](https://chainhero.io/2018/04/tutorial-hyperledger-fabric-how-to-build-your-first-network/)<br/>
[Tutorial Hyperledger Fabric SDK Go: How to build your first app](https://chainhero.io/2018/06/tutorial-build-blockchain-app-v1-1-0/)






 

