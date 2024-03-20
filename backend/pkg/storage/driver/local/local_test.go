package local

import (
	"cc.allio/fusion/pkg/storage"
	"cc.allio/fusion/pkg/util"
	"context"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"strings"
	"testing"
)

var LocalDefaultPolicy = &storage.Policy{Mode: storage.LocalMode}

func TestHandle_Upload(t *testing.T) {
	driver := NewLocalDriver(LocalDefaultPolicy)
	testCases := []*storage.FileStream{
		{
			SavePath: "TestHandle_Upload.txt",
			File:     io.NopCloser(strings.NewReader("")),
		},
		{
			SavePath: "TestHandle_Upload.txt",
			File:     io.NopCloser(strings.NewReader("")),
		},
		{
			SavePath: "inner/TestHandle_Upload.txt",
			File:     io.NopCloser(strings.NewReader("")),
		},
	}
	for _, fileCase := range testCases {
		err := driver.Upload(context.Background(), fileCase)
		assert.Empty(t, err)
	}
}

func TestHandle_Download(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	driver := NewLocalDriver(LocalDefaultPolicy)
	asserts := assert.New(t)

	file, err := os.Create(util.AbsolutePath("TestHandle_Download.txt", ""))
	asserts.NoError(err)
	asserts.NotNil(file)

	rsCloser, err := driver.Download(ctx, "TestHandle_Download.txt")
	asserts.NoError(err)
	asserts.NotNil(rsCloser)

	rsCloser, err = driver.Download(ctx, "TestHandle_Download_NoExist.txt")
	asserts.Error(err)
	asserts.Nil(rsCloser)
}

func TestHandle_List(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	driver := NewLocalDriver(LocalDefaultPolicy)
	asserts := assert.New(t)

	for _, path := range []string{
		"test/TestDriver_List/parent.txt",
		"test/TestDriver_List/parent_folder2/sub2.txt",
		"test/TestDriver_List/parent_folder1/sub_folder/sub1.txt",
		"test/TestDriver_List/parent_folder1/sub_folder/sub2.txt",
	} {
		f, _ := util.CreateNestedFile(util.AbsolutePath(path, ""))
		f.Close()
	}

	{
		res, err := driver.List(ctx, "test/TestDriver_List", false)
		asserts.NoError(err)
		asserts.Len(res, 3)
	}

	{
		res, err := driver.List(ctx, "test/TestDriver_List", true)
		asserts.NoError(err)
		asserts.Len(res, 7)
	}
}
