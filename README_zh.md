# sgv - Go 版本管理器

`sgv`（Simple Go Version）是一个轻量级的命令行工具，用于在本地系统上管理多个 Go 版本。支持 Go 1.13 及以上，适用于 macOS 和 Linux。

## 功能特性

- **安装 Go 版本**：下载并安装任意受支持的 Go 版本。
- **切换 Go 版本**：一键切换到已安装的 Go 版本。
- **自动切换**：根据当前项目的 `go.mod` 自动切换到所需 Go 版本。
- **获取最新版**：一条命令安装并切换到最新 Go 版本。
- **环境变量管理**：为每个 Go 版本管理项目特定的环境变量。
- **列出已安装版本**：按主版本分组查看所有已安装的 Go 版本。
- **列出可用补丁版本**：列出指定主版本下所有可用补丁版本，并标记已安装。
- **卸载 Go 版本**：卸载任意已安装的 Go 版本（当前激活版本除外）。
- **显示 sgv 版本**：显示 sgv 的构建版本和 commit hash。

---

## 安装

### macOS / Linux

```bash
curl -sSL https://raw.githubusercontent.com/fun7257/sgv/main/install.sh | bash
```

- 安装到 `/usr/local/bin/sgv`
- 自动配置 `GOROOT` 和 `PATH` 到 `~/.bashrc` 或 `~/.zshrc`
- 安装后请重启终端或执行 `source ~/.bashrc` 或 `source ~/.zshrc`

---

## 使用说明

### 切换或安装并切换 Go 版本

```bash
sgv <version>
```
- 例：`sgv 1.22.1` 或 `sgv go1.21.0`
- 若未安装则自动下载安装并切换
- 若当前目录为 Go 项目且请求版本低于 `go.mod` 要求，则会报错并中止

### 仅安装（不切换）

```bash
sgv install <version>
```
- 例：`sgv install 1.22.1`
- 仅下载安装，不切换当前版本

### 自动切换（基于 go.mod）

```bash
sgv auto
```
- 检测 `go.mod` 所需 Go 版本（优先 `toolchain`，若存在且更高）
- 若未安装则提示下载安装
- 若已是当前激活版本则无操作
- 若非 Go 项目则提示并无操作

### 获取并切换到最新版

```bash
sgv latest
```
- 若未安装则下载安装最新版，并切换为当前版本

### 列出已安装 Go 版本

```bash
sgv list
```
- 按主版本分组显示所有已安装版本
- 当前激活版本标记为 `<- current`

### 列出主版本下所有补丁版本

```bash
sgv sub <major_version>
```
- 例：`sgv sub 1.22`
- 列出所有可用的 Go 1.22.x 版本，已安装的标记为 `(installed)`
- 仅支持 Go 1.13 及以上

### 卸载 Go 版本

```bash
sgv uninstall <version>
```
- 例：`sgv uninstall 1.22.1`
- 不能卸载当前激活版本

### 显示 sgv 版本

```bash
sgv version
```
- 显示 sgv 构建时的 Go 版本和 commit hash

### 管理环境变量

```bash
sgv env                           # 列出当前 Go 版本的环境变量
sgv env -w KEY=VALUE             # 为当前 Go 版本设置环境变量
sgv env -u KEY                   # 删除当前 Go 版本的环境变量
sgv env --shell                  # 以 shell 格式输出环境变量
```

**使用示例：**
```bash
sgv env -w GOWORK=auto           # 启用 Go workspace 模式
sgv env -w GODEBUG=gctrace=1     # 启用 GC 追踪调试
sgv env -u GODEBUG               # 删除 GODEBUG 设置
```

**主要特性：**
- **版本隔离**：每个 Go 版本拥有独立的环境变量配置
- **无感加载**：切换版本或修改环境变量时自动应用到当前 shell
- **保护机制**：关键 Go 变量（GOROOT、GOPATH 等）不可修改
- **持久存储**：变量按版本保存，切换时自动恢复

---

## 无感体验

sgv 提供无感的自动环境加载体验：

- **自动切换**：执行 `sgv go1.21.0` 时，环境变量自动加载到当前 shell
- **环境管理**：使用 `sgv env -w` 或 `sgv env -u` 的更改立即应用到当前 shell
- **自动命令**：`sgv auto` 和 `sgv latest` 在版本切换后自动加载环境变量
- **无需手动操作**：无需运行 `eval` 命令或重启终端

这通过安装脚本自动添加到 shell 配置中的包装函数实现。

---

## 环境变量

- `SGV_DOWNLOAD_URL_PREFIX`  
  更改 Go 下载源（如中国大陆用户可设为 `https://golang.google.cn/dl/`）

```sh
export SGV_DOWNLOAD_URL_PREFIX=https://golang.google.cn/dl/
```
- 可在运行 sgv 前设置，或加入 shell 配置文件实现持久化

---

## 其他说明

- 所有 Go 版本均安装在 `~/.sgv/versions/`，当前激活版本通过软链接 `~/.sgv/current` 实现
- 切换版本后请确保 `GOROOT` 和 `PATH` 已正确配置（安装脚本会自动处理）

---

## 贡献

欢迎提交 issue 或 PR！

## 许可证

MIT License，详见 [LICENSE](./LICENSE)
