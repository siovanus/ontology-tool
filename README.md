# ontology-tool

tools for ontology governance

## Instruction

### 1. Clone repo

```shell
git clone https://github.com/ontio/ontology-tool.git
```

### 2. Build or download latest release

bulid:

```shell
go build
```

or download latest release:

https://github.com/ontio/ontology-tool/releases

### 3. Update config file

```shell
vim config.json
```

content of config.json：

```json
{
  "JsonRpcAddress":"http://dappnode2.ont.io:20336",
  "GasPrice":2500,
  "GasLimit":20000
}
```

`JsonRpcAddress`：rpc of ontology nodes

for mainnet: 
`"http://dappnode1.ont.io:20336","http://dappnode2.ont.io:20336","http://dappnode3.ont.io:20336","http://dappnode4.ont.io:20336"`

for polaris testnet: 
`"http://polaris1.ont.io:20336"，"http://polaris2.ont.io:20336"，"http://polaris3.ont.io:20336"，"http://polaris4.ont.io:20336"`

### 4. Run command line

list of supported command line: 
config file is under params directory.

| command line                                    | config file                                | function                                                   |
| ----------------------------------------------- | ------------------------------------------ | ---------------------------------------------------------- |
| `./main -t RegisterCandidate`                   | `RegisterCandidate.json`                   | 注册成为候选节点                                           |
| `./main -t ChangeMaxAuthorization`              | `ChangeMaxAuthorization.json`              | 修改节点最大接受的质押数                                   |
| `./main -t SetFeePercentage`                    | `SetFeePercentage.json`                    | 修改节点收益的分配比例，独占的initpos部分和独占的stake部分 |
| `./main -t AddInitPos`                          | `AddInitPos.json`                          | 增加初始质押                                               |
| `./main -t ReduceInitPos`                       | `ReduceInitPos.json`                       | 减少初始质押                                               |
| `./main -t AuthorizeForPeer`                    | `AuthorizeForPeer.json`                    | 向节点投票质押                                             |
| `./main -t UnAuthorizeForPeer`                  | `UnAuthorizeForPeer.json`                  | 取消向节点投票质押                                         |
| `./main -t Withdraw`                            | `Withdraw.json`                            | 提取质押的ont                                              |
| `./main -t QuitNode`                            | `QuitNode.json`                            | 退出节点                                                   |
| `./main -t BlackNode`                           | `BlackNode.json`                           | 拉黑节点                                                   |
| `./main -t WhiteNode`                           | `WhiteNode.json`                           | 取消拉黑节点                                               |
| `./main -t CommitDpos`                          | `CommitDpos.json`                          | 强行切换共识周期                                           |
| `./main -t UpdateConfig`                        | `UpdateConfig.json`                        | 更改共识配置                                               |
| `./main -t UpdateGlobalParam`                   | `UpdateGlobalParam.json`                   | 更改全局参数                                               |
| `./main -t UpdateGlobalParam2`                  | `UpdateGlobalParam2.json`                  | 更改全局参数2                                              |
| `./main -t UpdateSplitCurve`                    | `UpdateSplitCurve.json`                    | 更改分润曲线                                               |
| `./main -t TransferPenalty`                     | `TransferPenalty.json`                     | 提取拉黑罚没的ont                                          |
| `./main -t SetPromisePos`                       | `SetPromisePos.json`                       | 设置节点的承诺质押                                         |
| `./main -t GetVbftConfig`                       | 无                                         | 查询当前共识配置                                           |
| `./main -t GetPreConfig`                        | 无                                         | 查询下轮生效的共识配置                                     |
| `./main -t GetGlobalParam`                      | 无                                         | 查询全局参数                                               |
| `./main -t GetGlobalParam2`                     | 无                                         | 查询全局参数2                                              |
| `./main -t GetSplitCurve`                       | 无                                         | 查询分润曲线                                               |
| `./main -t GetGovernanceView`                   | 无                                         | 查询当前周期信息                                           |
| `./main -t GetPeerPoolItem`                     | `GetPeerPoolItem.json`                     | 查询某个节点信息                                           |
| `./main -t GetPeerPoolMap`                      | 无                                         | 查询所有节点信息                                           |
| `./main -t GetAuthorizeInfo`                    | `GetAuthorizeInfo.json`                    | 查询某个地址对某个节点的质押信息                           |
| `./main -t GetTotalStake`                       | `GetTotalStake.json`                       | 查询地址的总质押                                           |
| `./main -t GetPenaltyStake`                     | `GetPenaltyStake.json`                     | 查询罚没的ont信息                                          |
| `./main -t GetAttributes`                       | `GetAttributes.json`                       | 查询节点的属性信息                                         |
| `./main -t GetSplitFee`                         | 无                                         | 查询总的已经分出还未提取的ong                              |
| `./main -t GetSplitFeeAddress`                  | `GetSplitFeeAddress.json`                  | 查询某个地址已经分出还未提取的ong                          |
| `./main -t GetPromisePos`                       | `GetPromisePos.json`                       | 查询节点的承诺质押                                         |
| `./main -t InBlackList`                         | `InBlackList.json`                         | 查询节点是否在黑名单                                       |
| `./main -t WithdrawOng`                         | `WithdrawOng.json`                         | 提取ong收益                                                |
| `./main -t Vrf`                                 | 无                                         | 查询vrf信息                                                |
| `./main -t TransferOntMultiSign`                | `TransferOntMultiSign.json`                | 多签转账ont，目的地是account                               |
| `./main -t TransferOngMultiSign`                | `TransferOngMultiSign.json`                | 多签转账ong，目的地是account                               |
| `./main -t TransferFromOngMultiSign`            | `TransferFromOngMultiSign.json`            | 多签转账ong, ，目的地是account，transferfrom方法           |
| `./main -t TransferOntMultiSignAddress`         | `TransferOntMultiSignAddress.json`         | 多签转账ont，目的地是地址                                  |
| `./main -t TransferOngMultiSignAddress`         | `TransferOngMultiSignAddress.json`         | 多签转账ong，目的地是地址                                  |
| `./main -t TransferFromOngMultiSignAddress`     | `TransferFromOngMultiSignAddress.json`     | 多签转账ong, ，目的地是地址，transferfrom方法              |
| `./main -t GetAddressMultiSign`                 | `GetAddressMultiSign.json`                 | 根据公钥算出多签地址                                       |
| `./main -t TransferOntMultiSignToMultiSign`     | `TransferOntMultiSignToMultiSign.json`     | 多签对多签转ont                                            |
| `./main -t TransferOngMultiSignToMultiSign`     | `TransferOngMultiSignToMultiSign.json`     | 多签对多签转ong                                            |
| `./main -t TransferFromOngMultiSignToMultiSign` | `TransferFromOngMultiSignToMultiSign.json` | 多签对多签transferfrom ong                                 |
| `./main -t GetVbftInfo`                         | `GetVbftInfo.json`                         | 查询vbftInfo                                               |

And now you can run your command and input your password if needed.