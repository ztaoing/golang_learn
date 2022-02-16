文章来源于[polarisxu ]

---
本文总结一些初学者很容易犯错的知识点

[NULL]

SQL中 NULL 是一个特殊的值，其不满足自反性，也就是NULL != NULL，如果用column = NULL试图查询值为NULL的数据注定悲剧。
所以SQL单独为其准备了IS NULL语法。

[NaN]

IEEE754规定，有且仅有NaN不满足自反性。假设NaN作为 key 存放了一个元素，那么该元素就像掉进了黑洞，再也找不回来了。针对NaN也有自己的专用方法：IsNaN。

[true != true]

这种场景需要刻意地去构造，平时是遇不到的。举一个例子：
    
    func main() {
        // b1和b2都是通过unsafe “构造” 的bool值，
        var j1 int8 = 1
        var b1 bool = *(*bool)(unsafe.Pointer(&j1))
        
        var j2 int8 = 2
        var b2 bool = *(*bool)(unsafe.Pointer(&j2))
        
        if b1 { // 编译为TESTB
            println("b1 =", b1)
        }
        if b2 == true { // true是常量，编译阶段会自行处理，不会编译为CMPB
            println("b2 = true")
        }
        if b1 == b2 { // 编译为CMPB
            println("b1 = b2")
        } else {
            println("b1 != b2")
        }
    }

[注意：]
b1和b2都是通过unsafe “构造” 的bool值，当其与 if 或者其它需要bool类型的操作符搭配时，其表现出bool类型，
当其与==操作符或者其它（可以被bool类型和int8类型同时使用的）操作符搭配时，其表现出int8类型。

根本原因在于 CPU 只认指令，它的概念里没有数据类型更不必说int8或者bool，所谓的类型是高级语言对于 CPU 指令的抽象[^0]。

譬如TESTB指令的操作数，编译器把它称为bool；CMPB对应的操作数，编译器把它称为int8或者uint8；
CPU 指令有一组有符号移位和无符号移位，编译器将其抽象为signed和unsigned。

反过来，编译器把if还原成TESTB，将==还原成CMP，于是就出现了同样的数值（同样的内存）在不同的指令下表现出了不一样的行为。

---

[nil != nil]

这个场景初学者会经常遇到，即使是大佬如果不注意也会被坑。

    package main

    type MyErr struct {
    }
    
    func (MyErr) Error() string {
        return "MyErr"
    }
    
    func main() {
        // var e error改为var e *MyErr则不会出现nil != nil的怪胎。
        var e error = GetErr()
        println(e == nil)
        // false
    }
    
    func GetErr() *MyErr {
     return nil
    }

concrete: 具体的

* 出现这种“异常”现象一定存在函数调用，在单一函数内部不存在“反常”现象。譬如在GetErr函数内部，怎么造都不会出现nil != nil的怪胎。

* 出现这种“异常”现象一定存在 interface type 引用 concrete type 的背景。譬如var e error改为var e *MyErr则不会出现nil != nil的怪胎。

* []interface type 引用 concrete type 时，隐含了convT2I[^1]。convT2I的主要作用在于用interface 的 type和concrete 的 value合成了一个新值。

* 如果是 interface type 和 nil 比较，需要 type 和 value 对应的内存同时为 0 才可认定。由于允许只有 type 没有 value[^2]，不允许只有 value 没有 type，所以可以简化为 type 是否为指定。

* 由于nil相对于编译器来说相当于字面常量，各个阶段都存在不同程度的优化。
  让我们观察一下interface type和interface type之间的比较、interface type和concrete type之间的比较，来感受一下编译器是怎么处理 type 和 value 的。


* interface type和interface type之间的比较
  
      func walkcompareInterface(n *Node, init *Nodes) *Node {
  
      n.Right = cheapexpr(n.Right, init)
      n.Left = cheapexpr(n.Left, init)
      eqtab, eqdata := eqinterface(n.Left, n.Right)
      var cmp *Node
  
      if n.Op == OEQ { // x == y
        cmp = nod(OANDAND, eqtab, eqdata) // eqtab && eqdata
      } else { // x != y
        eqtab.Op = ONE
        cmp = nod(OOROR, eqtab, nod(ONOT, eqdata, nil)) // !eqtab || !eqdata
      }
      return finishcompare(n, cmp, init)
      }
* interface type和concrete type

      if n.Left.Type.IsInterface() != n.Right.Type.IsInterface() { // 一个interface一个concrete
      l := cheapexpr(n.Left, init)
      r := cheapexpr(n.Right, init)
  
      // 如果需要，交换左右双方，保证左侧为interface右侧为concrete。
      if n.Right.Type.IsInterface() {
        l, r = r, l
      }
    
      // 如果是x == y，用&&连接；如果是x != y，用||连接。
      eq := n.Op
      andor := OOROR
      if eq == OEQ {
        andor = OANDAND
      }
      // Check for types equal.
      // For empty interface, this is:
      //   l.tab == type(r)
      // For non-empty interface, this is:
      //   l.tab != nil && l.tab._type == type(r)
      var eqtype *Node
      tab := nod(OITAB, l, nil)
      rtyp := typename(r.Type)
      // 首先比较type
      if l.Type.IsEmptyInterface() {
          tab.Type = types.NewPtr(types.Types[TUINT8])
          tab.SetTypecheck(1)
          eqtype = nod(eq, tab, rtyp)
      } else {
          nonnil := nod(brcom(eq), nodnil(), tab)
          match := nod(eq, itabType(tab), rtyp)
          eqtype = nod(andor, nonnil, match)
      }
      // 其次比较value
      eqdata := nod(eq, ifaceData(n.Pos, l, r.Type), r)
      // 使用&&或者||连接两个结果
      expr := nod(andor, eqtype, eqdata)
      n = finishcompare(n, expr, init)
      return n
      }


* [^0]: 譬如无符号长整型，C 称之为 unsigned long int，Go 称之为 uint64，Rust 称之为 u64。

* [^1]: convT2I 是一类函数调用，譬如 convT2Inoptr，convT2E，convT2Enoptr，convT64，convT32，convT16 等等，包括被编译器优化掉的函数调用。
  
* [^2]: 譬如 var err *MyErr，指明了类型，但尚未赋值。