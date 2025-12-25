package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/widrin/knights/internal/actor"
	"github.com/widrin/knights/internal/config"
	"github.com/widrin/knights/internal/logger"
	"github.com/widrin/knights/internal/service"
)

var (
	configPath = flag.String("config", "configs/server.yaml", "配置文件路径")
)

func main() {
	flag.Parse()

	logger.Info("Knights Server Starting...")

	// 加载配置
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		logger.Fatal("加载配置文件失败: %v", err)
	}

	// 验证服务类型
	serviceType := service.ServiceType(cfg.Server.ServiceType)
	if !serviceType.IsValid() {
		logger.Fatal("无效的服务类型: %s", cfg.Server.ServiceType)
	}

	logger.Info("服务类型: %s", serviceType)
	logger.Info("服务名称: %s", cfg.Server.Name)
	logger.Info("监听地址: %s:%d", cfg.Server.Address, cfg.Server.Port)

	// 创建Actor系统
	system := actor.NewActorSystem(cfg.Server.Name)
	defer system.Shutdown()

	// 创建服务管理器
	manager := service.NewManager(system, serviceType)

	// 启动对应的服务
	pid, err := manager.StartService(serviceType)
	if err != nil {
		logger.Fatal("启动服务失败: %v", err)
	}

	logger.Info("服务启动成功, PID: %s", pid.String())

	// 根据服务类型执行额外的初始化
	switch serviceType {
	case service.ServiceTypeCenter:
		logger.Info("中心服务已启动")
		// TODO: 启动HTTP管理接口

	case service.ServiceTypeLogin:
		logger.Info("登录服务已启动")
		// TODO: 连接到中心服务
		// TODO: 启动网络监听

	case service.ServiceTypeGateway:
		logger.Info("网关服务已启动")
		// TODO: 连接到中心服务
		// TODO: 启动网络监听

	case service.ServiceTypeGame:
		logger.Info("游戏服务已启动, ServerID: %s", cfg.Game.ServerID)
		// TODO: 连接到中心服务
		// TODO: 注册到中心服务
		// TODO: 启动房间管理器
		// TODO: 启动匹配服务
	}

	logger.Info("%s 服务运行中...", cfg.Server.Name)

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("收到关闭信号，正在关闭服务...")

	// 优雅关闭
	manager.StopAll()

	logger.Info("%s 服务已关闭", cfg.Server.Name)
}
