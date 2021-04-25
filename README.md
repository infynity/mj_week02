# mj_week02

## 不应该，

## 第一点原因  dao层属于高度复用层，不知道啥情况需要wrap ，应该只返回root error（sentinel error)，

## 第二点原因 如果使用了wrap   上层如果需要wrap的话 会出现冗余
