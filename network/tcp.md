## tcp 通信
* 与udp对比特点。 面向连接，稳定，数据可靠的协议

tcp 三次握手

    client      ||            ||  server
                ||     syn    || 
    syn_send    ||  ------>   ||  syn_rcvd
                || ack,syn    || 
    established ||  <------   || 
                ||     ack    || 
                ||  ------>   ||  established


tcp 四次挥手 

    client      ||            ||  server
                ||     fin    || 
    fin_wait_1  ||  ------>   ||  close_wait
                ||     ack    || 
    fin_wait_2  ||  <------   || 
                ||     fin    || 
    time_wait   ||  <------   ||  last_ack
                ||     ack    || 
                ||  ------>   ||  closed


                