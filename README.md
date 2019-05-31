# coap
背景定义
协议背景
随着越来越多的人通过PC、手机等设备相互连接，现代互联网蓬勃发展使得人们的生活发生了翻天覆地的变化。很多人预测将会有更多其他的设备相互连接，这些设备的数量将远远超过人类的数量，到时候形成的网络将是现有网络的N个量级，这个网络带给世界的变化将是无法估量的。
不像人接入互联网的简单方便，由于物联网设备大多都是资源限制型的，有限的CPU、RAM、Flash、网络宽带等。对于这类设备来说，想要直接使用现有网络的TCP和HTTP来实现设备实现信息交换是不现实的。于是为了让这部分设备能够顺利接入网络，CoAP协议就被设计出来了。
概念定义

Constrained Application Protocol (CoAP) is a specialized Internet Application Protocol for constrained devices, as defined in RFC 7252.It enables those constrained devices called "nodes" to communicate with the wider Internet using similar protocols. CoAP is designed for use between devices on the same constrained network (e.g., low-power, lossy networks), between devices and general nodes on the Internet, and between devices on different constrained networks both joined by an internet. CoAP is also being used via other mechanisms, such as SMS on mobile communication networks.

以上是来自维基百科对CoAP的定义。简言之，CoAP是受约束设备的专用Internet应用程序协议。
协议特点

基于消息模型
请求/响应模型
双向通信
轻量、低功耗
支持可靠传输
支持IP多播
非长连接通信，支持受限设备
支持观察模式
支持异步通信

协议内容
CoAP是一个完整的二进制应用层协议，消息格式紧凑，默认运行在UDP上。
CoAP首部

【Ver】版本编号。
【T】报文类型，CoAP协议定了4种不同形式的报文，CON报文，NON报文，ACK报文和RST报文。
【TKL】CoAP标识符长度。CoAP协议中具有两种功能相似的标识符，一种为Message ID(报文编号)，一种为Token(标识符)。其中每个报文均包含消息编号，但是标识符对于报文来说是非必须的。
【Code】功能码/响应码。Code在CoAP请求报文和响应报文中具有不同的表现形式，Code占一个字节，它被分成了两部分，前3位一部分，后5位一部分，为了方便描述它被写成了c.dd结构。其中0.XX表示CoAP请求的某种方法，而2.XX、4.XX或5.XX则表示CoAP响应的某种具体表现。
【Message ID】报文编号。
【Token】标识符具体内容，通过TKL指定Token长度。
【Option】报文选项，通过报文选项可设定CoAP主机，CoAP URI，CoAP请求参数和负载媒体类型等等。
【1111 1111B】CoAP报文和具体负载之间的分隔符。

请求方法

0.01 GET：获取资源
0.02 POST：创建资源
0.03 PUT：更新资源
0.04 DELETE：删除资源

响应码
1、 Success 2.xx
这一类型的状态码，代表请求已成功被服务器接收、理解、并接受。

2.01 Created
2.02 Deleted
2.03 Valid
2.04 Changed
2.05 Content

2、 Client Error 4.xx
这类的状态码代表了客户端看起来可能发生了错误，妨碍了服务器的处理。

4.00 Bad Request
4.01 Unauthorized
4.02 Bad Option
4.03 Forbidden
4.04 Not Found
4.05 Method Not Allowed
4.06 Not Acceptable
4.12 Precondition Failed
4.13 Request Entity Too Large
4.15 Unsupported Content-Format

3、 Server Error 5.xx
这类状态码代表了服务器在处理请求的过程中有错误或者异常状态发生，也有可能是服务器的软硬件资源无法完成对请求的处理。

5.00 Internal Server Error
5.01 Not Implemented
5.02 Bad Gateway
5.03 Service Unavailable

媒体类型

【text/plain】 编号为0，表示负载为字符串形式，默认为UTF8编码。
【application/link-format】编号为40，CoAP资源发现协议中追加定义，该媒体类型为CoAP协议特有。
【application/xml】编号为41，表示负载类型为XML格式。
【application/octet-stream】编号为42，表示负载类型为二进制格式。
【application/exi】编号为47，表示负载类型为“精简XML”格式。
【applicaiton/cbor】编号为50，可以理解为二进制JSON格式。

工作模式
CoAP参考了很多HTTP的设计思路，同时也根据受限资源限制设备的具体情况改良了诸多的设计细节，增加了很多实用的功能。
消息类型

CON：需要被确认的请求，如果CON请求被发送，那么对方必须做出响应。
NON：不需要被确认的请求，如果NON请求被发送，那么对方不必做出回应。
ACK：应答消息，接受到CON消息的响应。
RST：复位消息，当接收者接收到的消息包含一个错误，接收者解析消息或者不再关心发送者发送的内容，那么复位消息将会被发送。

请求/响应模型

携带模式

分离模式

非确认模式

HTTP、CoAP、MQTT
CoAP协议的设计参考了HTTP，CoAP和MQTT都是行之有效的物联网协议，一下为它们之间的异同。
HTTP和CoAP

HTTP代表超文本传输协议，CoAP代表约束应用协议；
HTTP协议的传输层采用了TCP，CoAP协议的传输层使用UDP；
CoAP协议是HTTP协议的简化版；
CoAP协议和HTTP协议一样使用请求/响应模型，拥有相同的方法；
CoAP开销更低，并支持多播；
CoAP专为资源构成应用而设计，如：IoT/WSN/M2M等...

CoAP和MQTT

MQTT协议使用发布/订阅模型，CoAP协议使用请求/响应模型；
MQTT是长连接，CoAP协议是无连接；
MQTT通过中间代理传递消息的多对多协议，CoAP协议是Server和Client之间消息传递的单对单协议；
MQTT不支持带有类型或者其它帮助Clients理解的标签消息，CoAP内置内容协商和发现支持，允许设备彼此窥测以找到交换数据的方式。

EMQ与CoAP
EMQ的开源版提供了一个Generic CoAP的实现，读者如果有兴趣的话可以试用一下。在EMQ的商业版中，包含了对CoAP的完整支持。
