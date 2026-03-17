#!/bin/bash
#
# check-container.sh - 判断当前是否运行在容器环境中
# 支持检测：Docker, Podman, Kubernetes, LXC 等常见容器环境
#

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检测 Docker 环境
check_docker() {
    if [ -f /.dockerenv ]; then
        return 0
    fi
    return 1
}

# 检测 Podman 环境
check_podman() {
    if [ -n "$CONTAINER_ID" ] || grep -q "container=podman" /proc/1/environ 2>/dev/null; then
        return 0
    fi
    return 1
}

# 检测 Kubernetes 环境
check_kubernetes() {
    if [ -d /var/run/secrets/kubernetes.io ] || [ -d /run/secrets/kubernetes.io ]; then
        return 0
    fi
    return 1
}

# 检测 LXC 环境
check_lxc() {
    if grep -q "container=lxc" /proc/1/environ 2>/dev/null || \
       grep -qi "lxc" /sys/class/dmi/id/product_name 2>/dev/null; then
        return 0
    fi
    return 1
}

# 检测 cgroup 特征
check_cgroup() {
    if [ -f /proc/1/cgroup ]; then
        # 检查 cgroup 是否显示容器特征
        local cgroup_content
        cgroup_content=$(cat /proc/1/cgroup 2>/dev/null)
        if echo "$cgroup_content" | grep -qE "docker|kubepods|lxc|containerd"; then
            return 0
        fi
        # 简单情况：cgroup 根路径
        if [ "$cgroup_content" = "0::/" ]; then
            return 0
        fi
    fi
    return 1
}

# 检测是否为容器环境（通用方法）
check_container() {
    # 方法 1: 检查 /.dockerenv
    if [ -f /.dockerenv ]; then
        echo "docker"
        return 0
    fi

    # 方法 2: 检查环境变量
    if [ -n "$DOCKER_CONTAINER" ] || [ -n "$KUBERNETES_SERVICE_HOST" ]; then
        echo "kubernetes"
        return 0
    fi

    # 方法 3: 检查 hostname 是否为短 ID
    local hostname
    hostname=$(hostname 2>/dev/null)
    if echo "$hostname" | grep -qE "^[a-f0-9]{12}$"; then
        echo "docker"
        return 0
    fi

    # 方法 4: 检查 cgroup
    if [ -f /proc/1/cgroup ]; then
        local cgroup_content
        cgroup_content=$(cat /proc/1/cgroup 2>/dev/null)
        if echo "$cgroup_content" | grep -qE "docker|containerd|kubepods"; then
            echo "docker"
            return 0
        fi
        # 容器内常见特征
        if [ "$cgroup_content" = "0::/" ]; then
            echo "docker"
            return 0
        fi
    fi

    # 方法 5: 检查根文件系统特征
    if grep -q "overlay\|aufs\|btrfs" /proc/mounts 2>/dev/null; then
        # 联合文件系统常见于容器
        warn "检测到联合文件系统，可能是容器环境"
        echo "unknown"
        return 0
    fi

    return 1
}

# 主函数
main() {
    local container_type
    container_type=$(check_container)
    local exit_code=$?

    if [ $exit_code -eq 0 ]; then
        case "$container_type" in
            docker)
                info "当前运行在 Docker 容器环境中"
                ;;
            kubernetes)
                info "当前运行在 Kubernetes 容器环境中"
                ;;
            podman)
                info "当前运行在 Podman 容器环境中"
                ;;
            lxc)
                info "当前运行在 LXC 容器环境中"
                ;;
            *)
                warn "检测到容器环境，但无法确定具体类型"
                ;;
        esac
        exit 0
    else
        info "当前未运行在容器环境中（宿主机）"
        exit 1
    fi
}

# 如果是直接执行脚本而非被 source
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi
