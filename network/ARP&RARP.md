## ARP 

### 总结 MAC地址用户内网直连设备之间直接通信，而不用通过IP、端口方式。

ARP协议位于TCP协议栈中的数据链路层，称为地址解析协议，ARP协议实现任意网络层地址到任意物理地址的转换，例如IP地址转换为MAC地址。

工作原理
主机想自己所在的网络广播一个ARP请求，该请求包含目标机器的网络地址，此网络上的其他机器都将收到这个请求，但只有被请求的目标机器会回应一个ARP应答，其中包含自己的物理地址。

ARP报文结构


·
硬件类型字段定义物理地址的类型，他的值为1表示MAC地址；

·
协议类型字段表示要映射的协议地址类型，它的值为0x800，表示IP地址；

·
硬件地址长度字段和协议地址长度字段，其单位是字节。对MAC地址来说，其长度为6，对IP（v4）来说，其长度为4；

·
操作字段指出4种操作类型：ARP请求（值为1）、ARP应答（值为2）、RARP（值为3）、RARP（值为4）。

·
最后四个字段指定通信双方的以太网地址和IP地址。发送端填充除目的端以太网地址的其他三个字段，以构建ARP请求并发送之。接收端发现该请求的目的端IP地址是自己，就把自己的以太网地址填进去，然后交换两个目的端地址和两个发送端的地址，以构建ARP应答并返回之（当然，如前所示，操作字段需要改为2）。

IP与MAC地址
我们知道每台计算机都有一个MAC地址，准确的来说每台可以连接到以太网的设备都有一个唯一的MAC地址，这个地址就是为了当别的设备向这个设备发送数据包的时候可以指定地址。当两台设备连接起来就可以使用链路层的PPP协议收发数据了，这个时候每个数据包都会直接使用两台设备的MAC地址。此时我们没有用到IP地址就实现了设备之间的数据交换。 
而当我们要实现长距离之间的设备互联就需要用到IP协议，例如我们中国的一台设备要与美国的一台设备进行数据收发，不可能再简单的拉一根网线直接把两台设备相连。此时在中国与美国之间需要众多的中转路由器，数据包要经过这些路由器才能正确的到达，所以这个时候的PPP协议就没有办法使用，因为两台设备不是直连的。

简单说就是一个在子网里面，MAC地址可以在这个子网络里面定位到不同的网络设备，IP可以在整个因特网中定位到不同的子网。

下面我来描述一下ARP协议以及IP地址和MAC地址在数据传输的过程中的作用
当你把你的电脑A连上路由器的时候，路由器会给你的电脑分配一个IP地址，你自己也会拥有一个MAC地址，当你需要给其他子网上的机器B发送数据包的时候，都会经过路由器，但是只知道路由器的IP地址不知道它的MAC地址是没有办法做到数据的传输的，因为链路层的传输协议要求知道目的端的MAC地址。

所以电脑A需要找到路由器的MAC地址，即通过路由器使用ARP协议向电脑A所在子网发送ARP广播，这个广播就是在问：请IP地址是xxxxxxx的机器告诉我你的MAC地址。

路由器收到这个请求就会给电脑A发送一个ARP响应，在里面包含自己的MAC地址，这下就可以向路由器发送数据包了，接下来路由器会把数据包传给你的ISP，方式和上面很像，而ISP和ISP之间则需要根据数据包中指定的目标IP和路由器存的路由表来确定把数据包传给接下来的哪个路由器。直到数据包被传到电脑B。
————————————————
版权声明：本文为CSDN博主「Seven17000」的原创文章，遵循 CC 4.0 BY-SA 版权协议，转载请附上原文出处链接及本声明。
原文链接：https://blog.csdn.net/MBuger/article/details/73861017