# 使用官方 Go 1.22.0 镜像作为基础镜像
FROM golang:1.22.0

# 设置环境变量，避免在安装过程中出现交互提示
ENV DEBIAN_FRONTEND=noninteractive

# 更新包列表并安装必要的依赖
RUN apt-get update && \
    apt-get install -y curl gnupg && \
    rm -rf /var/lib/apt/lists/*

# 安装 Node.js 20.x（npm 10.5.0）
RUN curl -fsSL https://deb.nodesource.com/setup_20.x | bash - && \
    apt-get install -y nodejs && \
    # 验证 Node.js 和 npm 版本
    node -v && npm -v

# 设置工作目录
WORKDIR /app

# 复制 package.json 和 package-lock.json（如果存在）
COPY package.json ./
COPY package-lock.json ./

# 安装 npm 依赖
RUN npm install

# 复制项目的其余文件
COPY . .

# 暴露应用运行的端口（根据需要修改）
EXPOSE 3000

# 设置默认的启动命令
CMD ["npm", "start"]