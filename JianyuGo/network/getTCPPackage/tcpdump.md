tcpdump 是一个非常好的数据抓包工具，可以帮助我们捕获和查看网络数据包

tcpdump -S -nn -vvv -i lo0 port 8000
-S 显示序列号绝对值
-vvv 为了打印更多的详细描述信息
port 将仅捕获与8000端口通信和来自8000端口的流量
-i 指定捕获换回接口 localhost

        
    tao@taodeMacBook-Pro ~ % tcpdump -S -nn -vvv -i lo0 port 8000
    tcpdump: listening on lo0, link-type NULL (BSD loopback), capture size 262144 bytes
    三次握手:
    13:29:00.473799 IP6 (flowlabel 0xf0900, hlim 64, next-header TCP (6) payload length: 44) ::1.51596 > ::1.8000: Flags [S], cksum 0x0034 (incorrect -> 0xadf2), seq 469714520, win 65535, options [mss 16324,nop,wscale 6,nop,nop,TS val 3662489990 ecr 0,sackOK,eol], length 0
    13:29:00.473881 IP6 (flowlabel 0x60200, hlim 64, next-header TCP (6) payload length: 44) ::1.8000 > ::1.51596: Flags [S.], cksum 0x0034 (incorrect -> 0x56c9), seq 1312707746, ack 469714521, win 65535, options [mss 16324,nop,wscale 6,nop,nop,TS val 227976865 ecr 3662489990,sackOK,eol], length 0
    13:29:00.473898 IP6 (flowlabel 0xf0900, hlim 64, next-header TCP (6) payload length: 32) ::1.51596 > ::1.8000: Flags [.], cksum 0x0028 (incorrect -> 0xb7c6), seq 469714521, ack 1312707747, win 6371, options [nop,nop,TS val 3662489990 ecr 227976865], length 0

    13:29:00.473910 IP6 (flowlabel 0x60200, hlim 64, next-header TCP (6) payload length: 32) ::1.8000 > ::1.51596: Flags [.], cksum 0x0028 (incorrect -> 0xb7c6), seq 1312707747, ack 469714521, win 6371, options [nop,nop,TS val 227976865 ecr 3662489990], length 0
    客户端向服务端发送数据:
    13:29:00.474484 IP6 (flowlabel 0xf0900, hlim 64, next-header TCP (6) payload length: 42) ::1.51596 > ::1.8000: Flags [P.], cksum 0x0032 (incorrect -> 0xe4ce), seq 469714521:469714531, ack 1312707747, win 6371, options [nop,nop,TS val 3662489990 ecr 227976865], length 10
    13:29:00.474516 IP6 (flowlabel 0x60200, hlim 64, next-header TCP (6) payload length: 32) ::1.8000 > ::1.51596: Flags [.], cksum 0x0028 (incorrect -> 0xb7bc), seq 1312707747, ack 469714531, win 6371, options [nop,nop,TS val 227976865 ecr 3662489990], length 0
    服务端向客户端返回数据:
    13:29:00.474897 IP6 (flowlabel 0x60200, hlim 64, next-header TCP (6) payload length: 132) ::1.8000 > ::1.51596: Flags [P.], cksum 0x008c (incorrect -> 0xe469), seq 1312707747:1312707847, ack 469714531, win 6371, options [nop,nop,TS val 227976866 ecr 3662489990], length 100
    13:29:00.474926 IP6 (flowlabel 0xf0900, hlim 64, next-header TCP (6) payload length: 32) ::1.51596 > ::1.8000: Flags [.], cksum 0x0028 (incorrect -> 0xb757), seq 469714531, ack 1312707847, win 6370, options [nop,nop,TS val 3662489991 ecr 227976866], length 0
    四次挥手:
    13:29:10.475055 IP6 (flowlabel 0xf0900, hlim 64, next-header TCP (6) payload length: 32) ::1.51596 > ::1.8000: Flags [F.], cksum 0x0028 (incorrect -> 0x9046), seq 469714531, ack 1312707847, win 6370, options [nop,nop,TS val 3662499991 ecr 227976866], length 0
    13:29:10.475106 IP6 (flowlabel 0x60200, hlim 64, next-header TCP (6) payload length: 32) ::1.8000 > ::1.51596: Flags [.], cksum 0x0028 (incorrect -> 0x6935), seq 1312707847, ack 469714532, win 6371, options [nop,nop,TS val 227986866 ecr 3662499991], length 0
    13:29:10.475332 IP6 (flowlabel 0x60200, hlim 64, next-header TCP (6) payload length: 32) ::1.8000 > ::1.51596: Flags [F.], cksum 0x0028 (incorrect -> 0x6934), seq 1312707847, ack 469714532, win 6371, options [nop,nop,TS val 227986866 ecr 3662499991], length 0
    13:29:10.475380 IP6 (flowlabel 0xf0900, hlim 64, next-header TCP (6) payload length: 32) ::1.51596 > ::1.8000: Flags [.], cksum 0x0028 (incorrect -> 0x6935), seq 469714532, ack 1312707848, win 6370, options [nop,nop,TS val 3662499991 ecr 227986866], length 0
    重点关注：Flags[] 
    [S]  代表SYN包
    [F]  代表FIN包
    [.]  代表对应的ACK包
    
    [S.]代表SYN-ACK
    [F.]代表FIN-ACK