package main

import (
	"context"
	"io/fs"
	"path/filepath"

	"go.quinn.io/dataq/plugin"
	pb "go.quinn.io/dataq/proto"
)

type FileScanPlugin struct {
	rootPath string
}

func New() *FileScanPlugin {
	return &FileScanPlugin{}
}

func (p *FileScanPlugin) ID() string {
	return "filescan"
}

func (p *FileScanPlugin) Name() string {
	return "File System Scanner"
}

func (p *FileScanPlugin) Description() string {
	return "Scans specified directories for files and returns their metadata"
}

func (p *FileScanPlugin) Configure(config map[string]string) error {
	if path, ok := config["root_path"]; ok {
		p.rootPath = path
	} else {
		p.rootPath = "."
	}
	return nil
}

func (p *FileScanPlugin) Extract(ctx context.Context, action *pb.Action, api *plugin.PluginAPI) error {
	err := filepath.WalkDir(p.rootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if d.IsDir() {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		item := &pb.DataItem{
			Meta: &pb.DataItemMetadata{
				PluginId:    p.ID(),
				Id:          path,
				Kind:        "file",
				Timestamp:   info.ModTime().Unix(),
				ContentType: "application/octet-stream",
			},
			RawData: []byte{},
		}

		api.WriteItem(item)

		return nil
	})

	if err != nil {
		// In a real implementation, we'd want to handle this error better
		return err
	}

	return nil
}

func (p *FileScanPlugin) Transform(ctx context.Context, item *pb.DataItem, api *plugin.PluginAPI) error {
	return nil
}
