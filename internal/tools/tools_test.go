package tools

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestReadFile(t *testing.T) {
	workDir := t.TempDir()
	tool := NewReadFileTool(workDir)

	testFile := filepath.Join(workDir, "test.txt")
	testContent := "Hello, World!"
	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		path    string
		want    string
		wantErr bool
	}{
		{
			name:    "read existing file",
			path:    "test.txt",
			want:    testContent,
			wantErr: false,
		},
		{
			name:    "read non-existent file",
			path:    "nonexistent.txt",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args, err := json.Marshal(map[string]string{"path": tt.path})
			if err != nil {
				t.Fatal(err)
			}
			got, err := tool.Execute(context.Background(), args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("Execute() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestWriteFile(t *testing.T) {
	workDir := t.TempDir()
	tool := NewWriteFileTool(workDir)

	tests := []struct {
		name    string
		path    string
		content string
		want    string
		wantErr bool
	}{
		{
			name:    "write new file",
			path:    "new.txt",
			content: "test content",
			want:    "成功将内容写入到文件: new.txt",
			wantErr: false,
		},
		{
			name:    "overwrite existing file",
			path:    "overwrite.txt",
			content: "overwritten content",
			want:    "成功将内容写入到文件: overwrite.txt",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args, err := json.Marshal(map[string]string{
				"path":    tt.path,
				"content": tt.content,
			})
			if err != nil {
				t.Fatal(err)
			}
			got, err := tool.Execute(context.Background(), args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("Execute() = %q, want %q", got, tt.want)
			}

			if !tt.wantErr {
				content, err := os.ReadFile(filepath.Join(workDir, tt.path))
				if err != nil {
					t.Fatal(err)
				}
				if string(content) != tt.content {
					t.Errorf("file content = %q, want %q", string(content), tt.content)
				}
			}
		})
	}
}

func TestEditFile(t *testing.T) {
	workDir := t.TempDir()
	tool := NewEditFileTool(workDir)

	testFile := filepath.Join(workDir, "edit.txt")
	if err := os.WriteFile(testFile, []byte("old content"), 0644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		path    string
		oldStr  string
		newStr  string
		want    string
		wantErr bool
	}{
		{
			name:    "edit existing content",
			path:    "edit.txt",
			oldStr:  "old",
			newStr:  "new",
			want:    "成功编辑文件 edit.txt",
			wantErr: false,
		},
		{
			name:    "edit non-existent content",
			path:    "edit.txt",
			oldStr:  "nonexistent",
			newStr:  "new",
			want:    "",
			wantErr: true,
		},
		{
			name:    "edit non-existent file",
			path:    "nonexistent.txt",
			oldStr:  "old",
			newStr:  "new",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args, err := json.Marshal(map[string]string{
				"path":    tt.path,
				"old_str": tt.oldStr,
				"new_str": tt.newStr,
			})
			if err != nil {
				t.Fatal(err)
			}
			got, err := tool.Execute(context.Background(), args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("Execute() = %q, want %q", got, tt.want)
			}
		})
	}
}
