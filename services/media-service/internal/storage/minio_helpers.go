package storage

import "time"

func (s *MinioStorage) Bucket() string        { return s.bucket }
func (s *MinioStorage) Expiry() time.Duration { return s.expiry }
