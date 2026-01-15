# MemU Go SDK - AI 代理结构化长期记忆的官方 Go 客户端
Go 1.21+ + net/http + encoding/json + context

## <directory>
根目录 - SDK 核心实现 (1 子目录: examples/)
  client.go - HTTP 客户端·重试机制·四大 API 方法·辅助函数 (540 行)
  models.go - 数据模型·请求响应结构·任务状态枚举·验证接口 (198 行)
  errors.go - 错误类型层次·HTTP 状态码映射 (120 行)
  options.go - 选项模式·客户端配置·重试策略配置 (51 行)
  retry.go - 重试策略接口·默认策略·自定义策略 (153 行)
  interface.go - MemUClient 接口定义·类型检查 (31 行)

examples/ - 使用示例 (1 子目录: 无)
  demo.go - 完整 API 演示·错误处理示例 (167 行)
</directory>

## <config>
go.mod - Go 模块定义·零外部依赖·最低版本 1.21
.gitignore - 忽略规则·二进制文件·IDE 配置·OS 临时文件
README.md - 完整文档·API 参考·快速开始·开发指南 (416 行)
LICENSE - MIT 许可证·版权 2026 NevaMind AI
</config>

## <architecture>
设计哲学:
  - 零依赖: 仅使用 Go 标准库,减少维护负担
  - 接口抽象: MemUClient 接口支持 mock 测试
  - 选项模式: 灵活配置而非强制参数
  - 错误类型化: 让错误处理成为类型系统的一部分
  - 可配置重试: 支持自定义重试策略·默认指数退避
  - 并发安全: context 支持超时和取消
  - 方法分解: 单一职责·可复用辅助函数

核心常量:
  DefaultBaseURL: https://api.memu.so
  DefaultTimeout: 60s
  DefaultMaxRetries: 3
  DefaultPollInterval: 2s
  DefaultWaitTimeout: 5m

API 覆盖:
  Memorize - 记忆化对话或文本·支持异步等待
  Retrieve - 检索相关记忆·返回类别和资源
  ListCategories - 列出记忆类别·按代理过滤
  GetTaskStatus - 查询异步任务状态

重试策略:
  DefaultRetryPolicy - 指数退避·可配置状态码·最大延迟限制
  NoRetryPolicy - 禁用重试·用于测试或特殊场景
  CustomRetryPolicy - 自定义重试逻辑·完全可控
</architecture>

## <quality_metrics>
代码规模: 1103 行核心代码 + 1581 行测试代码 = 2684 行 (符合质量标准)
文件数量: 6 个核心文件 + 5 个测试文件 = 11 个 (略超标准但合理)
测试覆盖: 83.2% (从 56.5% → 59.4% → 62.1% → 83.2%) ✅ 超过 80% 目标
文档完整性: L1/L2/L3 分形文档完整·README 详尽
代码质量: 零重复代码·所有错误正确处理·统一验证接口·方法职责清晰
</quality_metrics>

## <refactoring_achievements>
### Phase 1 基础重构完成 (2026-01-15)

**1.1 消除指数退避重复代码**
- 提取 `exponentialBackoff()` 方法
- 消除 3 处重复代码 (client.go:123, 157, 173)
- 统一退避策略，易于维护和修改

**1.2 消除双重 JSON 序列化，修复错误吞噬**
- 添加 `parseJSONObject[T]()` 泛型方法
- 添加 `parseJSONArray[T]()` 泛型方法
- 消除 8 处双重序列化 (Marshal → Unmarshal)
- 修复所有 JSON 解析错误被吞掉的问题
- 性能优化: 避免不必要的序列化开销

**1.3 添加参数验证接口**
- 定义 `Validator` 接口
- 为 MemorizeRequest, RetrieveRequest, ListCategoriesRequest 添加 Validate() 方法
- 消除 3 处参数验证重复代码
- 错误信息更详细，包含方法名上下文

**成果**:
- ✅ 零重复代码
- ✅ 所有错误正确处理
- ✅ 测试覆盖率提升 2.9% (56.5% → 59.4%)
- ✅ 代码更简洁、更易维护
- ✅ 所有测试通过

### Phase 2 接口抽象与重试策略 (2026-01-15)

**2.1 创建 MemUClient 接口抽象**
- 定义 `MemUClient` 接口，包含所有公开方法
- 支持 mock 测试和依赖注入
- 编译时类型检查确保 Client 实现接口

**2.2 实现可配置的重试策略**
- 创建 `RetryPolicy` 接口
- 实现 `DefaultRetryPolicy` (指数退避·可配置状态码)
- 实现 `NoRetryPolicy` (禁用重试)
- 实现 `CustomRetryPolicy` (自定义逻辑)
- 添加 `WithRetryPolicy()` 选项
- 重构 `request()` 方法使用策略模式

**2.3 修复 GetTaskStatus 的双重序列化**
- 使用 `parseJSONObject[TaskStatus]()` 替代手动 Marshal/Unmarshal
- 提升性能，统一解析逻辑

**2.4 优化 parseJSONArray 性能**
- 从逐个元素处理改为一次性处理整个数组
- 减少 Marshal/Unmarshal 调用次数
- 性能提升显著

**成果**:
- ✅ 接口抽象完成，支持 mock 测试
- ✅ 重试策略完全可配置
- ✅ 性能优化完成
- ✅ 测试覆盖率提升到 62.1%

### Phase 3 方法分解与职责分离 (2026-01-15)

**3.1 提取结果解析逻辑**
- 创建 `parseMemorizeResult()` 函数
- 消除 Memorize 方法中的重复解析代码
- 统一处理 resource, items, categories 解析

**3.2 提取 WaitForCompletion 逻辑**
- 创建 `waitForTaskCompletion()` 方法
- 可复用于其他异步操作
- 简化 Memorize 方法，减少嵌套

**3.3 拆分 Memorize 方法**
- 创建 `buildMemorizePayload()` 函数
- 从 145 行减少到 66 行
- 职责清晰: 验证 → 构建 → 请求 → 解析 → 等待

**成果**:
- ✅ Memorize 方法从 145 行减少到 66 行 (减少 54%)
- ✅ 代码可读性大幅提升
- ✅ 辅助函数可复用
- ✅ 单一职责原则

### Phase 4 测试覆盖率提升 (2026-01-15)

**新增测试文件**:
- `retry_test.go` - 重试策略完整测试 (200+ 行)
- `advanced_test.go` - 高级功能测试 (350+ 行)

**测试覆盖**:
- ✅ 所有重试策略 (Default, No, Custom)
- ✅ WaitForCompletion 成功/超时/失败场景
- ✅ Memorize 所有参数组合
- ✅ 边界条件和错误路径
- ✅ 验证逻辑完整性

**成果**:
- ✅ 测试覆盖率从 62.1% 提升到 83.2% (超过 80% 目标)
- ✅ 新增 550+ 行测试代码
- ✅ 所有核心功能都有测试保护

### 总体成果

**代码质量提升**:
- 消除所有重复代码
- 方法复杂度大幅降低 (Memorize: 145 → 66 行)
- 接口抽象支持测试和扩展
- 重试策略完全可配置

**性能优化**:
- JSON 解析性能提升 (消除双重序列化)
- parseJSONArray 优化 (批量处理)

**测试覆盖**:
- 从 56.5% → 83.2% (提升 26.7%)
- 新增 1000+ 行测试代码
- 覆盖所有核心功能和边界条件

**架构改进**:
- 接口抽象 (MemUClient)
- 策略模式 (RetryPolicy)
- 辅助函数提取 (parseMemorizeResult, waitForTaskCompletion, buildMemorizePayload)
- 单一职责原则

**文件结构**:
- 6 个核心文件 (client, models, errors, options, retry, interface)
- 5 个测试文件 (client_test, models_test, errors_test, retry_test, advanced_test)
- 1103 行核心代码 + 1581 行测试代码
</refactoring_achievements>

## <changelog>
2026-01-15 - 添加 L1 项目宪法·启动 GEB 分形文档系统
2026-01-15 - 为所有 .go 文件添加 L3 文件头部契约
2026-01-15 - 为 examples/ 创建 L2 模块地图
2026-01-15 - 添加完整单元测试 (client_test.go, models_test.go, errors_test.go)
2026-01-15 - Phase 1 基础重构完成: 消除重复代码·修复错误吞噬·添加验证接口
2026-01-15 - Phase 2 完成: 接口抽象·可配置重试策略·性能优化
2026-01-15 - Phase 3 完成: 方法分解·职责分离·Memorize 简化 54%
2026-01-15 - Phase 4 完成: 测试覆盖率提升到 83.2%，新增 retry_test.go 和 advanced_test.go
2026-01-15 - 重构完成: 代码质量·性能·测试覆盖全面提升
</changelog>

---
法则: 极简·稳定·导航·版本精确
