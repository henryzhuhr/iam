// Main entry point for the Hello GoZero application.
package main

import (
	"context"
	"flag"
	"fmt"
	_ "net/http/pprof" // 导入 pprof
	"os"
	"os/signal"
	"syscall"
	"time"

	"iam/internal/component"
	"iam/internal/config"
	"iam/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
)

// 全局配置文件路径
var configFile = flag.String("f", "etc/dev.yaml", "the config file")

func main() {
	flag.Parse() // 加载配置文件

	var c config.Config
	conf.MustLoad(*configFile, &c)

	// 创建服务上下文
	svcCtx, err := svc.NewServiceContext(c)
	if err != nil {
		fmt.Printf("failed to create service context: %v\n", err)
		return
	}
	defer svcCtx.Close()

	// ========== 使用组件管理器统一启动所有组件 ==========

	// 创建组件管理器（30秒超时）
	componentManager := component.NewManager(30 * time.Second)

	// 按顺序注册组件（先注册的先启动）
	componentManager.Register(component.NewPprofComponent(c.Pprof))

	// 统一启动所有组件
	if err := componentManager.StartAll(context.Background()); err != nil {
		fmt.Printf("❌ Failed to start components: %v\n", err)
		return
	}

	// 等待退出信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("\n🛑 Shutting down gracefully...")

	// 优雅关闭：统一停止所有组件
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	// 停止所有组件（按逆序）
	if err := componentManager.StopAll(shutdownCtx); err != nil {
		fmt.Printf("⚠️  Some components failed to stop: %v\n", err)
	}

	select {
	case <-shutdownCtx.Done():
		fmt.Println("⚠️ Shutdown timeout exceeded")
	default:
		fmt.Println("✅ Server stopped successfully")
	}
}
