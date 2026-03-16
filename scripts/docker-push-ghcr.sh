#!/bin/bash
# 构建 Docker 镜像并推送至 GitHub Container Registry (GHCR)
#
# 使用方法:
#   ./scripts/docker-push-ghcr.sh [VERSION]
#
# 环境变量:
#   GITHUB_ACTOR    - GitHub 用户名 (用于 GHCR 登录)
#   GITHUB_TOKEN    - GitHub Personal Access Token (用于 GHCR 登录)
#   IMAGE_NAME      - 镜像名 (默认：iam)
#   VERSION         - 镜像版本标签 (默认：latest)

set -e

# ============================================================
#   配置
# ============================================================
GHCR_REGISTRY="ghcr.io"
IMAGE_NAME="${IMAGE_NAME:-iam}"
VERSION="${1:-latest}"
GITHUB_ACTOR="${GITHUB_ACTOR:-henryzhuhr}"

# 完整的镜像名称
FULL_IMAGE_NAME="${GHCR_REGISTRY}/${GITHUB_ACTOR}/${IMAGE_NAME}"

# ============================================================
#   检查必要的环境变量
# ============================================================
if [ -z "${GITHUB_TOKEN}" ]; then
    echo "错误：请设置 GITHUB_TOKEN 环境变量"
    echo "可以通过 GitHub Settings -> Developer settings -> Personal access tokens 创建 token"
    exit 1
fi

# ============================================================
#   登录 GHCR
# ============================================================
echo "🔐 登录到 ${GHCR_REGISTRY}..."
echo "${GITHUB_TOKEN}" | docker login "${GHCR_REGISTRY}" -u "${GITHUB_ACTOR}" --password-stdin

# ============================================================
#   检查本地是否已构建过镜像
# ============================================================
EXISTING_IMAGE_ID=$(docker images -q "${IMAGE_NAME}:${VERSION}" 2>/dev/null)

if [ -n "${EXISTING_IMAGE_ID}" ]; then
    echo "✅ 本地已存在镜像 ${IMAGE_NAME}:${VERSION}，跳过构建"
else
    echo "⚠️  本地未找到镜像 ${IMAGE_NAME}:${VERSION}，开始构建..."
    # ============================================================
    #   构建镜像
    # ============================================================
    export IMAGE_NAME="${IMAGE_NAME}"
    export IMAGE_TAG="${VERSION}"
    chmod +x ./scripts/docker-build.sh
    ./scripts/docker-build.sh
fi

# ============================================================
#   打标签
# ============================================================
echo "🏷️  打标签 ${FULL_IMAGE_NAME}:${VERSION}..."
docker tag "${IMAGE_NAME}:${VERSION}" "${FULL_IMAGE_NAME}:${VERSION}"

# 同时推送 latest 标签
if [ "${VERSION}" != "latest" ]; then
    echo "🏷️  打标签 ${FULL_IMAGE_NAME}:latest..."
    docker tag "${IMAGE_NAME}:${VERSION}" "${FULL_IMAGE_NAME}:latest"
fi

# ============================================================
#   推送镜像
# ============================================================
echo "📤 推送镜像到 ${GHCR_REGISTRY}..."
docker push "${FULL_IMAGE_NAME}:${VERSION}"

if [ "${VERSION}" != "latest" ]; then
    docker push "${FULL_IMAGE_NAME}:latest"
fi

# ============================================================
#   推送结果
# ============================================================
echo ""
echo "✅ 推送完成!"
echo "镜像地址："
echo "  ${FULL_IMAGE_NAME}:${VERSION}"
if [ "${VERSION}" != "latest" ]; then
    echo "  ${FULL_IMAGE_NAME}:latest"
fi
