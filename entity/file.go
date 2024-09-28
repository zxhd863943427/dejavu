// DejaVu - Data snapshot and sync.
// Copyright (c) 2022-present, b3log.org
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package entity

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/siyuan-note/dejavu/util"
)

// File 描述了文件。
type File struct {
	ID      string   `json:"id"`      // Hash
	Path    string   `json:"path"`    // 文件路径
	Size    int64    `json:"size"`    // 文件大小
	Updated int64    `json:"updated"` // 最后更新时间
	Chunks  []string `json:"chunks"`  // 文件分块列表
}

const (
	TimestampIndex = 22
	ChunkHashIndex = 32
	IDLen          = 40
)

func NewFile(path string, size int64, updated int64) (ret *File) {
	ret = &File{
		Path:    path,
		Size:    size,
		Updated: updated,
	}

	ret.ID = genCheckID(ret)
	return
}

func (f *File) SecUpdated() int64 {
	return f.Updated / 1000
}

func genCheckID(f *File) string {
	buf := bytes.Buffer{}
	buf.WriteString(f.Path)
	pathHash := util.Hash(buf.Bytes())

	TimestampLen := ChunkHashIndex - TimestampIndex
	updated := fmt.Sprintf("%0[1]*[2]x", TimestampLen, f.Updated/1000)

	var builder strings.Builder
	ChunkHashLen := IDLen - ChunkHashIndex
	builder.WriteString(pathHash[:TimestampIndex])
	builder.WriteString(updated[:TimestampLen])
	builder.WriteString(fmt.Sprintf("%0[1]*[2]x", ChunkHashLen, 0))
	return builder.String()
}

func (f *File) GenTrueID() {

	prefixLen := ChunkHashIndex
	ChunkHashLen := IDLen - ChunkHashIndex

	buf3 := bytes.Buffer{}
	// slices.SortFunc(f.Chunks, func(a, b string) int {
	// 	return strings.Compare(a, b)
	// })

	for _, chunk := range f.Chunks {
		buf3.WriteString(chunk)
	}
	chunkHash := util.Hash(buf3.Bytes())

	var builder strings.Builder
	builder.WriteString(f.ID[:prefixLen])
	builder.WriteString(chunkHash[:ChunkHashLen])
	f.ID = builder.String()
}
