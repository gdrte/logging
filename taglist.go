package logging

import (
    "sort"
)

// A tagLevel represents a single setting for a logger.
type tagLevel struct {
    tag string
    level LogLevel
}

// A tagList holds a full set of tagLevels for a logger.
type tagList []tagLevel

// tagLists are sortable, this helps in some situations,
// in the future we may provide a "fast" version of checkTagLevel
// for presorted arrays of tags.
func (p tagList) Len() int           { return len(p) }
func (p tagList) Less(i, j int) bool { return p[i].tag < p[j].tag }
func (p tagList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func (p tagList) checkTagLevel(level LogLevel, o []string) bool {
    if p == nil {
        return false
    }
        
    lp, lo := len(p), len(o)

    for i:=0; i<lp; i++ {
        for j:=0; j<lo; j++ {
            if p[i].tag == o[j] && p[i].level <= level {
                return true
            }
        }
    }
    return false
}

func (p tagList) setTagLevel(tag string, level LogLevel) tagList {
    l := len(p)

    for i:=0;i<l;i++ {
        if p[i].tag == tag {
            p[i].level = level
            return p
        }
    }

    ret := append(p, tagLevel{
        tag: tag,
        level: level,
    })

    sort.Sort(ret)

    return ret
}

