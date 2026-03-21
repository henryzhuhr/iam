# Get the base image
ARG NODE_TAG=latest
ARG GO_TAG=latest
ARG UV_TAG=latest

FROM node:${NODE_TAG} AS node-builder
FROM golang:${GO_TAG} AS golang-builder
FROM ghcr.io/astral-sh/uv:${UV_TAG} AS uv-builder

FROM ubuntu:26.04

WORKDIR /root
USER root

# ============================================================
#   系统软件安装列表
# ============================================================
ARG INSTALL_LIST=""
# ARG INSTALL_LIST="$INSTALL_LIST ubuntu-server"
ARG INSTALL_LIST="$INSTALL_LIST ca-certificates locales"
ARG INSTALL_LIST="$INSTALL_LIST vim nano tree curl pciutils iputils-ping"
ARG INSTALL_LIST="$INSTALL_LIST zsh"
ARG INSTALL_LIST="$INSTALL_LIST git openssh-client"
ARG INSTALL_LIST="$INSTALL_LIST build-essential binutils binutils-gold"
ARG INSTALL_LIST="$INSTALL_LIST python3-dev"
ARG INSTALL_LIST="$INSTALL_LIST postgresql-client"
ARG INSTALL_LIST="$INSTALL_LIST redis-tools"
ARG INSTALL_LIST="$INSTALL_LIST mysql-client"

# 自定义镜像源 (本地构建可传入，如 MIRRORS_URL=mirrors.ustc.edu.cn)
# CI 环境不传此参数，使用官方源
ARG MIRRORS_URL=""

# ============================================================
#   系统环境配置
# ============================================================
# 安装软件
ENV DEBIAN_FRONTEND=noninteractive
ENV LANG=zh_CN.UTF-8
ENV LANGUAGE=zh_CN:zh
ENV LC_ALL=zh_CN.UTF-8
RUN set -x && \
    # Enable bash color prompt
    sed -i 's/#force_color_prompt=yes/force_color_prompt=yes/g' ${HOME}/.bashrc && \
    if [ -n "${MIRRORS_URL}" ]; then \
      # 24.04+ DEB822 format
      sed -i -e "s#//.*archive.ubuntu.com#//${MIRRORS_URL}#g" -e "s#//ports.ubuntu.com#//${MIRRORS_URL}#g" /etc/apt/sources.list.d/ubuntu.sources; \
    fi && \
    apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y --no-install-recommends ${INSTALL_LIST} && \
    # Generate locale
    locale-gen zh_CN.UTF-8 && \
    # Cleanup
    apt-get autoremove -y && \
    apt-get clean all && \
    rm -rf /tmp/* /var/tmp/* /var/lib/apt/lists/*

# 系统时间错误可能导致: `tls: failed to verify certificate: x509: certificate signed by unknown authority`
# 修改时区 (date -R 检查当前时间)
# 设置时区
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

# 启用 bash 补全
# RUN sed -i.bak '/^#if ! shopt -oq posix; then$/,/^#fi$/s/^#//' /etc/bash.bashrc

# 设置 zsh 为默认 shell
RUN set -x && \
    # ========================
    # Step: Set zsh as the default shell
    # ========================
    chsh -s /usr/bin/zsh root && \
    # ====== install oh-my-zsh ======
    git clone https://github.com/ohmyzsh/ohmyzsh.git --depth 1 /tmp/ohmyzsh && \
    sh /tmp/ohmyzsh/tools/install.sh --unattended && \
    rm -rf /tmp/ohmyzsh && \
    # ====== install zsh-autosuggestions ======
    ZSH_CUSTOM="/root/.oh-my-zsh/custom" && \
    mkdir -p "$ZSH_CUSTOM/plugins" && \
    git clone https://github.com/zsh-users/zsh-autosuggestions --depth 1 "$ZSH_CUSTOM/plugins/zsh-autosuggestions" && \
    # ==============================
    # set zsh theme and plugins
    echo "ZSH_THEME='candy'" >> ~/.zshrc && \
    echo "plugins=(git cp python zsh-autosuggestions)" >> ~/.zshrc && \
    echo "ZSH_AUTOSUGGEST_HIGHLIGHT_STYLE='fg=#00BFFF,bg=gray,bold,underline'" >> ~/.zshrc && \
    echo "ZSH_AUTOSUGGEST_STRATEGY=(history completion)" >> ~/.zshrc && \
    echo "source /root/.oh-my-zsh/oh-my-zsh.sh" >> ~/.zshrc && \
    echo "alias ll='ls -alh'" >> ~/.zshrc && \
    echo "alias la='ls -A'" >> ~/.zshrc && \
    echo "alias l='ls -CF'" >> ~/.zshrc && \
    echo "alias h='history'" >> ~/.zshrc && \
    echo "alias mkdir='mkdir -p'" >> ~/.zshrc && \
    echo "alias grep='grep --color=auto'" >> ~/.zshrc && \
    echo "alias egrep='egrep --color=auto'" >> ~/.zshrc && \
    echo "alias fgrep='fgrep --color=auto'" >> ~/.zshrc


# ============================================================
#   Golang 环境
# ============================================================
COPY --from=golang-builder /usr/local/go /usr/local/go
ENV PATH="/usr/local/go/bin:/root/go/bin:$PATH"
ENV GOPROXY=https://goproxy.cn,direct


# ============================================================
#   Python 环境
# ============================================================
COPY --from=uv-builder /uv /uvx /bin/
ARG UV_DEFAULT_INDEX="https://pypi.org/simple"
ENV UV_DEFAULT_INDEX=${UV_DEFAULT_INDEX}
ENV VIRTUAL_ENV_DISABLE_PROMPT=1


# ============================================================
#   Node.js 环境
# ============================================================
COPY --from=node-builder /usr/local/bin/node          /usr/local/bin/
COPY --from=node-builder /usr/local/lib/node_modules  /usr/local/lib/node_modules
ENV PATH="/usr/local/bin:$PATH"
ARG NPM_CONFIG_REGISTRY="https://registry.npmjs.org"
ENV NPM_CONFIG_REGISTRY=${NPM_CONFIG_REGISTRY}

RUN set -x && \
  # Create symlinks for npm and npx
  ln -s ../lib/node_modules/npm/bin/npm-cli.js /usr/local/bin/npm && \
  ln -s ../lib/node_modules/npm/bin/npx-cli.js /usr/local/bin/npx && \
  # Install Global npm packages
  npm install -g \
    @openai/codex \
    @anthropic-ai/claude-code && \
  # Clean up npm cache and temporary files
  npm cache clean --force && \
  rm -rf /tmp/* /var/tmp/* /root/.npm

CMD ["/usr/bin/zsh"]