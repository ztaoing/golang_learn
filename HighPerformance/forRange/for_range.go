/**
* @Author:zhoutao
* @Date:2021/1/1 上午11:54
* @Desc:
 */

package forRange

//range 可以用来遍历 array,slice,map,chan

// 与for不同的是，range对迭代值都创建了一个拷贝,如果迭代的元素的内存占用很低，那么for和range的性能几乎是一样的
// 如果迭代的元素内存占用较高，例如是一个包含喝多属性的struct结构体，那么for的性能将显著高于range，有时甚至有上千倍的性能差异，这种场景建议使用for,如果使用range，建议只迭代下标，通过下标访问迭代值
// 如果想使用range同时迭代下标和值，则需要将使用指针，才能不影响性能
