一、常用的位运算：
    
    &      与 AND
    |      或OR
    ^      异或XOR
    &^     位清空 (AND NOT)
    <<     左移
    >>     右移

二、位运算的用法:
位运算都是在二进制的基础上进行运算的，所以在位运算之前要先将两个数转成二进制

    1.  &

    & 只有两个数都是1结果才为1
    
    例：var i uint8 = 20  var j uint8=15   
    求i&j
    
    i转成二进制为0001 0100,
    j转成二进制为0000 1111
    
    0001 0100 & 0000 1111 = 0000 0100
    
    0001 0100对应的十进制就是4


    2.  |

    或 两个数有一个是1 结果就是1
    
    0001 0100 | 0000 1111 = 0001 1111
    
    0001 1111转成十进制就是31
    
    故20 | 15 = 31


      3.  ^

    （1）
    
        ^可以作为二元运算符，也可以作为一元运算符
    
     ^作二元运算符就是异或，相同为0，不相同为1
    
    如1^1 =0, 0^0=0,1^0=1,0^1=1
    
        0001 0100 ^ 0000 1111 = 0001 1011
    
    故 20 ^ 15 =27

    
    （2）

    ^作一元运算符表示是按位取反

    ^0001 0100 = 1110 1011

    故结果为235

     请思想下面代码的结果：
    func main() {  
        var    i  uint8  = 20  
        fmt.Println(^i,^20)  
    }  
    结果是：235   -21

    why?

    其实原因很简单，一个是有符号的数一个是无符号的数
    
    20在编译器中默认为int类型，故最高位是符号位，符号位取反，所以得到的结果是负数
    
    串联理解：负数的二进制数怎么表示？
    
    负数的二进制数是它对应的正数按位取反得到反码，再加上1得到的补码
    
    例如：3的二进制为00000000 00000000 00000000 00000011
    
    反码：            11111111 11111111 11111111 11111100
    
    补码：反码加1：  11111111 11111111 11111111 11111101
    
    故-3的二进制为11111111 11111111 11111111 11111101
    
    
    
         所以，一个有符号位的^操作为 这个数+1的相反数 


    4. &^
    
    作用：以左边的数据为基础，将运算符 左边数据相异的位保留，相同位清零
    
    1&^1  得0
    1&^0  得1
    0&^1  得0
    0&^0  得0
    
    0001 0100 &^ 0000 1111 = 0001 0000
    
    故结果为16

    5. >>右移和 <<左移

    左移和右移算是比较常见的运算了

      左移规则：

      右边空出的位用0填补

      高位左移溢出则舍弃该高位

      右移规则：

     左边空出的位用0或者1填补。正数用0填补，负数用1填补。注：不同的环境填补方式可能不同；

     低位右移溢出则舍弃该位

    例：0001 0100 >> 1得0000 1010   转成十进制为10

    0001 0100 << 1 得0010 1000   转成十进制为40