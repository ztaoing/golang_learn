/**
* @Author:zhoutao
* @Date:2022/3/1 09:12
* @Desc:
 */

package _7compression

//状态压缩：将2维的降维到1维数组
//dp[i][j]只和 dp[i][j-1]、dp[i+1][j-1]、dp[i+][j]有关
//dp[i][j-1]、dp[i-1][j-1]会存在覆盖
//压缩后的一维数组就是dp[i][...]这一行
//dp[j]原来是dp[i+1][j]：是外层for循环上一次算出来的值
//dp[j-1]原来是dp[i][j-1]:是内层for循环上一次算出来的值
//dp[i+1][j]也是内层for循环计算出来的值
