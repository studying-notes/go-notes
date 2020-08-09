# 算法与数据结构

- [算法与数据结构](#算法与数据结构)
	- [数据结构](#数据结构)
	- [常用算法](#常用算法)
	- [数学领域](#数学领域)

## 数据结构

- 组织结构
	- 集合结构
	- 线性结构
    	- [链表 Linked List](struct/link/README.md)
		- [栈 Stack](struct/stack/README.md)
		- [队列 Queue](struct/queue/README.md)
		- 线性表(Linear List)
		- 数组(Array)
		- 串(String)
		- 广义表(General List)
		- 跳跃表
	- 树型结构
		- 树 Tree
			- [二叉树](struct/tree/README.md)
			- [Trie 树](struct/tree/README.md#Trie-树)
			- 胜者树
			- 左偏树
			- 线段树
			- 后缀树
			- 哈夫曼树
			- 静态二叉检索树
			- RMQ
		- 堆(Heap)
			二叉堆
			斜堆
	- 图型结构		
		- 图(Graph)
			- trie图
			- 有限状态自动机
- 统计结构
	- 树状数组
	- 虚二叉树
	- 线段树
	- 矩形面积并
	- 圆形面积并
- 关系结构
	- Hash表
	- 并查集
	- 路径压缩思想

## 常用算法

按问题目标的类型划分

- 排序
	- 插入排序
		- 直接插入排序(Insert Sort) O(n^2)
		- 折半插入排序(Binary Insert Sort)
		- 希尔排序(Shell Sort)

	- 交换排序
		- 冒泡排序(Bubble Sort) O(n^2)
		- 快速排序(Quick Sort)?? O(nlogn)

	- 选择排序
		- 直接选择排序(Select Sort) O(n^2)
		- 锦标赛排序(Tournament Sort) O(nlogn)
		- 堆排序(Heap Sort) O(nlogn)

	- 归并排序(Merge Sort) O(nlogn)

	- 基数排序(Radix Sort) O(d(n+radix))

	- 桶排序(Bucket Sort) O(nlogn)

- 查找/搜索

	- 二分查找(Binary Search)

	- 树型
		- 二叉搜索树(Binary Search Tree)
		- 平衡搜索树(AVL Tree)
		- 并查集(Union-Find Set)

	- 哈希(Hashing)
	
	- 最优化剪枝

	- 可行性剪枝

	- 记忆化搜索

	- 枚举搜索(Enumeration)

	- 深度优先(Depth First Search)

	- 广度优先(Breadth First Search)

	- 启发式搜索(Heuristic Search)

- 压缩
	- 哈夫曼编码

- 加密
	- 对称加密
		- AES
		- DES
		- RC4
	- 非对称加密
		- RSA
		- IDEA

- 最优解
	- 递推
	- 递归(Recursion)
	- 贪心算法(Greedy)
	- 动态规划(DynamicProgram)
	- 随机搜索
	- 爬山法
	- 模拟退火法
	- 遗传算法

- 并行
	

- 机器学习
	- 回归算法
		- 最小二乘法(Ordinary Least Square)
		- 逻辑回归(Logistic Regression)
		- 逐步式回归(Stepwise Regression)
		- 多元自适应回归样条(Multivariate Adaptive Regression Splines)
		- 本地散点平滑估计(Locally Estimated Scatterplot Smoothing)
	- 基于实例的算法
		- k-Nearest Neighbor(KNN)
		- 学习矢量量化（Learning Vector Quantization， LVQ）
		- 自组织映射算法（Self-Organizing Map ， SOM）
	- 正则化方法 
		- Ridge Regression
		- Least Absolute Shrinkage and Selection Operator（LASSO）
		- 弹性网络（Elastic Net）
	- 决策树学习
		- 分类及回归树（Classification And Regression Tree， CART）
		- ID3 (Iterative Dichotomiser 3)
		- C4.5
		- Chi-squared Automatic Interaction Detection(CHAID)
		- Decision Stump
		- 随机森林（Random Forest）
		- 多元自适应回归样条（MARS）
		- 梯度推进机（Gradient Boosting Machine， GBM）
	- 贝叶斯方法
		- 朴素贝叶斯算法
		- 平均单依赖估计（Averaged One-Dependence Estimators， AODE）
		- Bayesian Belief Network（BBN）
	- 基于核的算法
		- 支持向量机（Support Vector Machine， SVM）
		- 径向基函数（Radial Basis Function ，RBF)
		- 线性判别分析（Linear Discriminate Analysis ，LDA)
	- 聚类算法
		- k-Means算法
		- 期望最大化算法（Expectation Maximization， EM）
	- 关联规则学习
		- Apriori算法
		- Eclat算法
	- 人工神经网络
		- 感知器神经网络（Perceptron Neural Network）
		- 反向传递（Back Propagation） 
		- Hopfield网络
		- 学习矢量量化（Learning Vector Quantization， LVQ）
	- 深度学习
		- 受限波尔兹曼机（Restricted Boltzmann Machine， RBN）
		- Deep Belief Networks（DBN）
		- 卷积网络（Convolutional Network）
		- 堆栈式自动编码器（Stacked Auto-encoders）
	- 降低维度算法
		- 主成份分析（Principle Component Analysis， PCA）
		- 偏最小二乘回归（Partial Least Square Regression，PLS）
		- Sammon映射
		- 多维尺度（Multi-Dimensional Scaling, MDS）
		- 投影追踪（Projection Pursuit）
	- 集成算法
		- Boosting
		- Bootstrapped Aggregation（Bagging）
		- AdaBoost
		- 堆叠泛化（Stacked Generalization， Blending）

## 数学领域

- 图论
	- 图的深度优先遍历和广度优先遍历.
	- 最短路径算法(dijkstra,bellman-ford,floyd,heap+dijkstra)
	- 最小生成树算法(prim,kruskal)
	- 拓扑排序
	- 二分图的最大匹配 (匈牙利算法)
	- 最大流的增广路算法(KM算法)
	- 度限制最小生成树和第K最短路.
	- 最优比率生成树.
	- 次小生成树.
	- 无向图、有向图的最小环
	- 差分约束系统的建立和求解
	- 最小费用最大流
	- 双连通分量
	- 强连通分支及其缩点
	- 图的割边和割点
	- 最小割模型、网络流规约

- 动态规划
	- 四边形不等式理论
	- 函数的凸凹性
	- 规划方向
	- 旅行商问题
	- 最优二分检索树
	- 树型动态规划
	- 状态动态规划
	- 记录状态的动态规划
	- LCA(Least Common Ancestors),即最近公共祖先
	- RMQ(Range Minimum/Maximum Query),即区间最值查询
	- 最长子序列系列问题
	- 最长不下降子序列
	- 最长公共子序列
	- 最长公共不下降子序列
	- 不完全状态记录
	- 青蛙过河问题
	- 利用区间dp
	- 背包类问题
		-  0-1背包，经典问题
		- 无限背包，经典问题
		- 判定性背包问题
		- 带附属关系的背包问题
		- +-1背包问题
		- 双背包求最优值
		- 构造三角形问题
		- 带上下界限制的背包问题(012背包)
	- 线性的动态规划问题
		- 积木游戏问题
		- 决斗（判定性问题）
		- 圆的最大多边形问题
		- 统计单词个数问题
		- 棋盘分割
		- 日程安排问题
		- 最小逼近问题(求出两数之比最接近某数/两数之和等于某数等等)
		- 方块消除游戏(某区间可以连续消去求最大效益)
		- 资源分配问题
		- 数字三角形问题
		- 邮局问题与构造答案
		- 最高积木问题
		- 两段连续和最大
		- 2次幂和问题
		- N个数的最大M段子段和
		- 交叉最大数问题
	- 判定性问题的dp(如判定整除、判定可达性等)
		- 模K问题的dp
		- 特殊的模K问题，求最大(最小)模K的数
		- 变换数问题
	- 单调性优化的动态规划
		- 1-SUM问题
		- 2-SUM问题
		- 序列划分问题(单调队列优化)
	- 剖分问题(多边形剖分/石子合并/圆的剖分/乘积最大)
		- 凸多边形的三角剖分问题
		- 乘积最大问题
		- 多边形游戏(多边形边上是操作符,顶点有权值)
		- 石子合并(N^3/N^2/NLogN各种优化)
	- 贪心的动态规划
		- 最优装载问题
		- 部分背包问题
		- 乘船问题
		- 贪心策略
		- 双机调度问题Johnson算法
	- 状态dp
		- 牛仔射击问题(博弈类)
		- 哈密顿路径的状态dp
		- 两支点天平平衡问题
		- 一个有向图的最接近二部图
	- 树型dp
		- 完美服务器问题(每个节点有3种状态)
		- 小胖守皇宫问题
		- 网络收费问题
		- 树中漫游问题
		- 树上的博弈
		- 树的最大独立集问题
		- 树的最大平衡值问题
		- 构造树的最小环
