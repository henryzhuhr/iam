#!/bin/bash
# 预先构建项目镜像的脚本，加快 docker compose up 的速度


# ============================================================
#   构建的镜像配置 (可通过环境变量覆盖)
# ============================================================
IMAGE_NAME=${IMAGE_NAME:-"iam"}
IMAGE_TAG=${IMAGE_TAG:-"latest"}

# ============================================================
#   基础镜像版本配置 (可通过环境变量覆盖)
# ============================================================
NODE_TAG=${NODE_TAG:-"25"}
GO_TAG=${GO_TAG:-"1.26"}
UV_TAG=${UV_TAG:-"0.10.12"}

# ============================================================
#   镜像源配置 (本地使用国内镜像，CI 使用官方源)
# ============================================================
MIRRORS_URL="mirrors.ustc.edu.cn"
NPM_CONFIG_REGISTRY="https://registry.npmjs.org"
UV_DEFAULT_INDEX="https://pypi.tuna.tsinghua.edu.cn/simple"

# 镜像列表（格式：镜像名:标签）
IMAGES=(
  "ubuntu:26.04"
  "node:${NODE_TAG}"
  "golang:${GO_TAG}"
  "ghcr.io/astral-sh/uv:${UV_TAG}"
  "mysql:9.6"
  "redis:8.6"
  "apache/kafka:4.1.2"
)

for IMAGE in "${IMAGES[@]}"; do
  NAME=$(echo "${IMAGE}" | cut -d: -f1)
  TAG=$(echo "${IMAGE}" | cut -d: -f2-)
  if ! docker images | grep -q "^${NAME}[[:space:]]\+${TAG}[[:space:]]"; then
    echo "pull image: ${IMAGE}"
    docker pull "${IMAGE}" || {
      echo "failed to pull image ${IMAGE}, aborting!";
      exit 1;
    }
  else
    echo "found ${IMAGE}, skip docker pull."
  fi
done


BUILD_ARGS=(
  "--build-arg" "UV_TAG=${UV_TAG}"
  "--build-arg" "GO_TAG=${GO_TAG}"
  "--build-arg" "NODE_TAG=${NODE_TAG}"
  "--build-arg" "MIRRORS_URL=${MIRRORS_URL}"
  "--build-arg" "NPM_CONFIG_REGISTRY=${NPM_CONFIG_REGISTRY}"
  "--build-arg" "UV_DEFAULT_INDEX=${UV_DEFAULT_INDEX}"
)

docker build --no-cache \
  -t "${IMAGE_NAME}:${IMAGE_TAG}" \
  -f dockerfiles/Dockerfile \
  "${BUILD_ARGS[@]}" \
  .

# ============================================================
#   打印构建结果
# ============================================================
echo ""
echo "Built images:"
docker images --format "table {{.Repository}}:{{.Tag}}\t{{.Size}}\t{{.CreatedSince}}" | grep "${IMAGE_NAME}" || true
