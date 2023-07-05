// Copyright © 2017 Microsoft <wastore@microsoft.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"context"
	"github.com/Azure/azure-storage-azcopy/v10/common"
	"github.com/Azure/azure-storage-azcopy/v10/ste"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsSourceDirWithStub(t *testing.T) {
	a := assert.New(t)
	bsu := getBSU()

	// Generate source container and blobs
	containerURL, containerName := createNewContainer(a, bsu)
	defer deleteContainer(a, containerURL)
	a.NotNil(containerURL)

	dirName := "source_dir"
	createNewDirectoryStub(a, containerURL, dirName)
	// set up to create blob traverser
	ctx := context.WithValue(context.TODO(), ste.ServiceAPIVersionOverride, ste.DefaultServiceApiVersion)
	p := azblob.NewPipeline(azblob.NewAnonymousCredential(), azblob.PipelineOptions{})

	// List
	rawBlobURLWithSAS := scenarioHelper{}.getRawBlobURLWithSAS(a, containerName, dirName)
	blobTraverser := newBlobTraverser(&rawBlobURLWithSAS, p, ctx, true, true, func(common.EntityType) {}, false, common.CpkOptions{}, false, false, false, common.EPreservePermissionsOption.None())

	isDir, err := blobTraverser.IsDirectory(true)
	a.True(isDir)
	a.Nil(err)
}

func TestIsSourceDirWithNoStub(t *testing.T) {
	a := assert.New(t)
	bsu := getBSU()

	// Generate source container and blobs
	containerURL, containerName := createNewContainer(a, bsu)
	defer deleteContainer(a, containerURL)
	a.NotNil(containerURL)

	dirName := "source_dir/"
	ctx := context.WithValue(context.TODO(), ste.ServiceAPIVersionOverride, ste.DefaultServiceApiVersion)
	p := azblob.NewPipeline(azblob.NewAnonymousCredential(), azblob.PipelineOptions{})

	// List
	rawBlobURLWithSAS := scenarioHelper{}.getRawBlobURLWithSAS(a, containerName, dirName)
	blobTraverser := newBlobTraverser(&rawBlobURLWithSAS, p, ctx, true, true, func(common.EntityType) {}, false, common.CpkOptions{}, false, false, false, common.EPreservePermissionsOption.None())

	isDir, err := blobTraverser.IsDirectory(true)
	a.True(isDir)
	a.Nil(err)
}

func TestIsSourceFileExists(t *testing.T) {
	a := assert.New(t)
	bsu := getBSU()

	// Generate source container and blobs
	containerURL, containerName := createNewContainer(a, bsu)
	defer deleteContainer(a, containerURL)
	a.NotNil(containerURL)

	fileName := "source_file"
	_, fileName = createNewBlockBlob(a, containerURL, fileName)

	ctx := context.WithValue(context.TODO(), ste.ServiceAPIVersionOverride, ste.DefaultServiceApiVersion)
	p := azblob.NewPipeline(azblob.NewAnonymousCredential(), azblob.PipelineOptions{})

	// List
	rawBlobURLWithSAS := scenarioHelper{}.getRawBlobURLWithSAS(a, containerName, fileName)
	blobTraverser := newBlobTraverser(&rawBlobURLWithSAS, p, ctx, true, true, func(common.EntityType) {}, false, common.CpkOptions{}, false, false, false, common.EPreservePermissionsOption.None())

	isDir, err := blobTraverser.IsDirectory(true)
	a.False(isDir)
	a.Nil(err)
}

func TestIsSourceFileDoesNotExist(t *testing.T) {
	a := assert.New(t)
	bsu := getBSU()

	// Generate source container and blobs
	containerURL, containerName := createNewContainer(a, bsu)
	defer deleteContainer(a, containerURL)
	a.NotNil(containerURL)

	fileName := "file_does_not_exist"
	ctx := context.WithValue(context.TODO(), ste.ServiceAPIVersionOverride, ste.DefaultServiceApiVersion)
	p := azblob.NewPipeline(azblob.NewAnonymousCredential(), azblob.PipelineOptions{})

	// List
	rawBlobURLWithSAS := scenarioHelper{}.getRawBlobURLWithSAS(a, containerName, fileName)
	blobTraverser := newBlobTraverser(&rawBlobURLWithSAS, p, ctx, true, true, func(common.EntityType) {}, false, common.CpkOptions{}, false, false, false, common.EPreservePermissionsOption.None())

	isDir, err := blobTraverser.IsDirectory(true)
	a.False(isDir)
	a.Equal(common.FILE_NOT_FOUND, err.Error())
}