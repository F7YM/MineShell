# 应用名称
APP_NAME = mineshell

# 版本信息（从 git 获取，如果没有则使用默认值）
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME = $(shell date '+%Y-%m-%d_%H:%M:%S')
LDFLAGS = -ldflags "-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME)"

# 输出目录
BUILD_DIR = build

# 默认目标
.PHONY: all
all: clean linux windows

# 创建输出目录
$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

# 编译 Linux 版本
.PHONY: linux
linux: $(BUILD_DIR)
	@echo "编译 Linux 版本..."
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64
	@echo "✓ Linux 版本编译完成: $(BUILD_DIR)/$(APP_NAME)-linux-amd64"

# 编译 Windows 版本
.PHONY: windows
windows: $(BUILD_DIR)
	@echo "编译 Windows 版本..."
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-windows-amd64.exe
	@echo "✓ Windows 版本编译完成: $(BUILD_DIR)/$(APP_NAME)-windows-amd64.exe"

# 编译当前平台版本
.PHONY: current
current: $(BUILD_DIR)
	@echo "编译当前平台版本..."
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)
	@echo "✓ 当前平台版本编译完成: $(BUILD_DIR)/$(APP_NAME)"

# 编译所有平台
.PHONY: all-platforms
all-platforms: clean $(BUILD_DIR)
	@echo "编译所有平台版本..."
	# Linux
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-linux-arm64
	# Windows
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-windows-amd64.exe
	GOOS=windows GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-windows-arm64.exe
	# macOS
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-darwin-amd64
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-darwin-arm64
	@echo "✓ 所有平台编译完成"

# 清理编译文件
.PHONY: clean
clean:
	@echo "清理编译文件..."
	rm -rf $(BUILD_DIR)
	@echo "✓ 清理完成"

# 运行测试
.PHONY: test
test:
	go test ./...

# 安装到系统（Linux）
.PHONY: install
install: linux
	@echo "安装到 /usr/local/bin..."
	sudo cp $(BUILD_DIR)/$(APP_NAME)-linux-amd64 /usr/local/bin/$(APP_NAME)
	@echo "✓ 安装完成"

# 帮助信息
.PHONY: help
help:
	@echo "可用目标："
	@echo "  make              - 编译 Linux 和 Windows 版本（默认）"
	@echo "  make linux        - 编译 Linux 版本"
	@echo "  make windows      - 编译 Windows 版本"
	@echo "  make current      - 编译当前平台版本"
	@echo "  make all-platforms - 编译所有平台版本"
	@echo "  make clean        - 清理编译文件"
	@echo "  make test         - 运行测试"
	@echo "  make install      - 安装到系统（Linux）"
	@echo "  make help         - 显示此帮助"