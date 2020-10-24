# 规则引擎合约开发工程

## 基本功能定义

 
### 基本形式:

#### Rule

```
 {} success_command failed_command
 {[()]}
```

#### Factor

Factor: {}, OP只能是&&, ||, 所以Factor的返回值为bool
```
eg.:
{ factorA OP factorB }, factor可以仍然是factor
```

#### Atom

 Atom: [], 只有一种构成方式, OP只能是 >, <, >=, <=, ==, !=, 所以Atom的返回值为 bool
 Atom的嵌套层数最多只会有一层
 ```
 [ Proton OP Proton]
 eg.:
 [ (p1 + (p2 + (p3))) ] >= (p3 + (p4 + (p5))) ]
```

#### Proton

 Proton () 只涉及 +, -, *, /, %, 返回(uint256, uint256)，第一个代表Key, 第二个代表value, 假定 devID设定的值都比较小，在有限范围内所有Id相加不会溢出
  dev可以是设备，也可以是预定的环境变量。
```
 ( dev1 curStatus OP dev2 curStatus )
```

## 规则示例
### 一条简单地完整的规则示例
 某时刻, 如果空调A的设定温度大于空调B的设定温度, 那么5min后将空调B的温度设定为A的温度, 否则关闭设备A及设备B
```
 { [ ( devA 0 - devB 0 ) > 0 ] } "Set(devB, 0, 5min)" "Set(devA, off), Set(devB, off)"
```
 0 用于填充参数, 或者占位

### 一条复杂的规则示例
 某时刻, 如果{房间A中的三盏灯全亮，并且房间B的电视工作}, 或者房间C的三台设备功率大于500W, 那么开启节能模式
```
 { { [ [ ( ( A.light1 0 + A.light2 0 ) + A.light3 0 ) == 3 ] } && { [ ( B.TV 0 0 0 0 ) == 1 ] } || { [ ( C.dev1 0 + ( C.dev2 0 + C.dev3 0 ) ) > 500 ] } } "Enter into energy saving mode" 0
```

