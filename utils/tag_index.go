package utils

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"gitlab.com/alphaticks/alphac/models"
	"regexp"
	"sort"
	"sync"
)

func HashTags(tags map[string]string) uint64 {
	var keys []string
	for k := range tags {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	str := ""
	for _, k := range keys {
		str += fmt.Sprintf("&%s=%s", k, tags[k])
	}
	hashBytes := md5.Sum([]byte(str))
	ID := binary.LittleEndian.Uint64(hashBytes[:8])
	return ID
}

type TagIndex struct {
	sync.RWMutex
	Securities map[uint64]*models.Security
	Tags       map[uint64]map[string]string
}

func NewTagIndex(securities []*models.Security) *TagIndex {
	idx := &TagIndex{
		Securities: make(map[uint64]*models.Security),
		Tags:       make(map[uint64]map[string]string),
	}
	for _, s := range securities {
		idx.Add(s)
	}
	return idx
}

func (idx *TagIndex) Add(sec *models.Security) bool {
	idx.Lock()
	defer idx.Unlock()
	tags := map[string]string{
		"ID":       fmt.Sprintf("%d", sec.SecurityID),
		"type":     sec.SecurityType,
		"base":     sec.Underlying.Symbol,
		"quote":    sec.QuoteCurrency.Symbol,
		"exchange": sec.Exchange.Name,
	}
	objectID := HashTags(tags)

	// Check if we have the tags
	if _, ok := idx.Tags[objectID]; !ok {
		idx.Tags[objectID] = tags
	}

	if _, ok := idx.Securities[objectID]; !ok {
		idx.Securities[objectID] = sec
	}

	return true
}

func (idx *TagIndex) Query(tags map[string]string) ([]*models.Security, error) {
	idx.RLock()
	defer idx.RUnlock()
	var res []*models.Security

	// Compile regex
	tagsRegex := make(map[string]*regexp.Regexp)
	for tag, regexStr := range tags {
		regex, err := regexp.Compile(regexStr)
		if err != nil {
			return nil, fmt.Errorf("error compiling regex: %v", err)
		}
		tagsRegex[tag] = regex
	}

	for k, v := range idx.Tags {
		match := true
		for tag, regex := range tagsRegex {
			if fVal, ok := v[tag]; !ok || regex.FindString(fVal) == "" {
				match = false
				break
			}
		}
		if match {
			res = append(res, idx.Securities[k])
		}
	}

	return res, nil
}
