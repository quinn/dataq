package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"

	"go.quinn.io/dataq/config"
	pb "go.quinn.io/dataq/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PluginProcess struct {
	Plugin *config.Plugin
	Client pb.DataQPluginClient
	Port   string
	cmd    *exec.Cmd
}

func main() {
	cfg, err := config.Get()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Start all enabled plugins
	plugins := make([]*PluginProcess, 0)
	basePort := 50051

	for _, plugin := range cfg.Plugins {
		if !plugin.Enabled {
			continue
		}

		port := fmt.Sprintf("%d", basePort)
		proc, err := startPlugin(plugin, port)
		if err != nil {
			log.Printf("Failed to start plugin %s: %v", plugin.ID, err)
			continue
		}

		plugins = append(plugins, proc)
		basePort++
	}

	// Handle shutdown gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("Shutting down...")
	
	// Stop all plugins
	for _, proc := range plugins {
		stopPlugin(proc)
	}
}

func startPlugin(plugin *config.Plugin, port string) (*PluginProcess, error) {
	// Start the plugin process
	cmd := exec.Command(plugin.BinaryPath)
	cmd.Env = append(os.Environ(), fmt.Sprintf("PORT=%s", port))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start plugin process: %w", err)
	}

	// Connect to the plugin
	conn, err := grpc.Dial(
		fmt.Sprintf("localhost:%s", port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		cmd.Process.Kill()
		return nil, fmt.Errorf("failed to connect to plugin: %w", err)
	}

	client := pb.NewDataQPluginClient(conn)

	return &PluginProcess{
		Plugin: plugin,
		Client: client,
		Port:   port,
		cmd:    cmd,
	}, nil
}

func stopPlugin(proc *PluginProcess) {
	if proc.cmd.Process != nil {
		if err := proc.cmd.Process.Signal(syscall.SIGTERM); err != nil {
			log.Printf("Failed to send SIGTERM to plugin %s: %v", proc.Plugin.ID, err)
			proc.cmd.Process.Kill()
		}
	}
}

// Example function to demonstrate plugin interaction
func extractFromPlugin(ctx context.Context, proc *PluginProcess, req *pb.ExtractRequest) (*pb.ExtractResponse, error) {
	return proc.Client.Extract(ctx, req)
}

// Example function to demonstrate plugin interaction
func transformWithPlugin(ctx context.Context, proc *PluginProcess, req *pb.TransformRequest) (*pb.TransformResponse, error) {
	return proc.Client.Transform(ctx, req)
}
