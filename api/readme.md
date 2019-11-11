v0.1.2  
名词中英对照说明：   
* 公钥 PublicKey    
* 私钥 PrivateKey   
* 助记词 Mnemonic    
* 公私钥对 Keypair    
* 加密 Encode   
* 解密 Decode   
* 地址 Address: 一个地址与一对公私钥对对应的， 一个公私钥对只有一个地址，反过来也是（有另外一种叫法帐号，但city中全部统一叫地址，不要叫别的 ）  
* 玩家 Player   
* 交易 transaction: 所有向链上写入数据，统称为交易; 简称tx   
* 生成交易 generate tx    
* 发送交易 send tx: 交易上链的一个步骤 
* 接收交易 receive tx: 这一步是上链完成   
* 确认交易 confirm tx: 也就是交易最终确认或不可撤销 
* 转帐 transfer: 是交易的一种，把token从一个地址转给到另外一个地址    
* 可用余额 usable balance ： 可以使用的余额
* 锁定余额 locked balance ： 已经使用，但还没有最后确认的， 这个余额也是不可以使用的
* 余额 balance ： 一般情况下为， 可用余额 + 锁定余额
* Nonce: 与地址对应的一个整数值，可以查询 
* 钱包 Wallet: 提供对地址的使用，如钱包提供转帐，查余额等功能  
* Hash256 : 在city中提到hash，没有特别说明的全部都是hash256   
* 订单 Purchase 
* 充值 Recharge 
* 提币 Discharge    
* 支付 Pay  
* 局 Play  
* 庄家 Banker   
* 参与者 Party   
* 赔率 Rate 
* 下注 Pour 
* 奖金池 Bonus pool  
* 注入 Inject   
* 结果 Result   
* 结算 Liquidate    
* 取回 Withdraw ，如果有任何出错在一局，庄家或参与者取回他自己的token   
* 创建721 Create721， 这里的721指的是erc721  
* 销毁721 Destroy721， 这里的721指的是erc721 

//服务端地址：    
* erc721操作地址 Erc721Address： 可以调整erc721的token。 增加当到或指定地址上；删除当前地址上的， 暂时不提供删除别人的token；锁定token，锁定后不能想互交易  
* erc20操作地址 Erc721Address：可以调整erc20的。增加当到或指定地址上；删除当前地址上的， 暂时不提供删除别人；锁定固定的token数量  

* city地址 CityAddress： 所有运行在city中的erc20与erc721都是由这个地址转给玩家， 如把token转给玩家 
* city充值地地 CityRechargeAddress:   
* ciyt提币地址 CityDischargeAddress:  

 
//接口原则  
* 一个交易对应一个唯一的订单id，且在创建交易之前产生    
* 一个玩家（或者地址）， 必须等交易发送成功（或取消）之后， 才能进行下一个交易   
* 生成交易后，一定会有处理结果（取消，成功，失败），其中的“取消”由game server保证通知链接口，成功与失败由链接口保证通知game server 
* 一但交易发送成功，保证通知到game server交易上链的结果    
* 产生重要信息时，要签名或加密（如果充值或提币） 
* 由于链接口要管理Nonce，所以把交易的生成放到链接口来实现
* 接口中的所有地址都带前缀的“0x”字符，地址有两种表示方式，一 16进制的字符串，二 byte数组，他们可以相互进化
* 所有的交易id带有“0x”前缀
* 优先使用0gas方式（如果能简单配置或修改代码实现，就使用）  
* erc721的token id, 保留100内的id(100保留)  
注： “保证通知”不管什么情况一定要通知到，如事件发生时，对方服务不在线，那么等对方服务上线后再进行通知；服务重起后，继续通知没有完成的通知 



