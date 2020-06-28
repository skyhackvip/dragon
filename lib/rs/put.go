package rs
import {
	"fmt"
	"io"
}

type RSPutStream struct {
	*encoder
}

func NewRSPutStream(dataServers []string, hash string, size int64) (*RSPutStream, error) {
	if len(dataServers) != ALL_SHARDs {
		return nil, fmt.Errorf("dataservers number mismatch")
	}
	perShard := (size + DATA_SHARDS -1) / DATA_SHARDS
	writers := make([]io.Writer, ALL_SHARDS)
	var  e error
	for i := range writers {
		writers[i], e = objectstream.NewTempPutStream(dataServers[i], fmt.Sprintf("%s.%d", hash, i), perShard)
		if e != nil {
			return nil, e
		}
	}
	enc := NewEncoder(writers)
	return &RSPutStream{enc},nil
}

func (s *RSPutStream) Commit(success bool) {
	s.Flush()
	for i := range s.writers {
		s.writers[i].(*objectstream.TempPutStream).Commit(success)
	}
}
