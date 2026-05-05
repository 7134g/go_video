## ADDED Requirements

### Requirement: CA证书检查
系统 SHALL 在启动时通过 `InitCa()` 检查 CA 证书是否已安装到系统信任存储区。

#### Scenario: 证书已安装
- **WHEN** 系统启动并检查证书
- **THEN** 证书检查通过，继续启动

#### Scenario: 证书未安装
- **WHEN** 检查发现 CA 证书未安装到系统信任存储区
- **THEN** 系统 panic，提示"需要先安装证书"

### Requirement: CA证书工具
系统 SHALL 提供独立的 `cmd/proxy` 工具用于生成 CA 证书并安装到系统信任存储。

#### Scenario: 手动安装证书
- **WHEN** 用户运行 `cmd/proxy`（需管理员权限）
- **THEN** 系统生成 CA 证书（ca.crt / ca.key）并安装到当前 OS 的信任存储

#### Scenario: 证书文件已存在
- **WHEN** ca.crt 和 ca.key 文件已存在
- **THEN** 使用现有证书文件，不重新生成
