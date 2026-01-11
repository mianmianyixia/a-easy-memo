package time

import (
	"a-easy-memo/internal/api/serve/cache"
	"a-easy-memo/internal/dao"
	"log"
	"time"

	"golang.org/x/net/context"
)

// 先定时调用getcache，将过期写缓存写入数据库，还要定时将数据库内容同步到读缓存
type CacheUpdater struct {
	ticker   *time.Ticker
	name     string
	redis    dao.MemberTask
	member   dao.MemberData
	cancel   context.CancelFunc
	ctx      context.Context
	duration time.Duration
}

// @param	name:用户名
func NewCacheUpdater(name string, redis dao.MemberTask, member dao.MemberData) *CacheUpdater {
	return &CacheUpdater{
		name:   name,
		redis:  redis,
		member: member,
	}
}

func (cu *CacheUpdater) Start() {
	cu.ticker = time.NewTicker(2 * time.Hour)
	cu.ctx, cu.cancel = context.WithTimeout(context.Background(), cu.duration)
	go func() {
		for {
			select {
			case <-cu.ticker.C:
				mem := dao.Data{Name: cu.name}
				_, err := cache.GetCache(mem, cu.redis, cu.member)
				if err != nil {
					log.Printf("缓存更新失败: %v", err)
				}
			case <-cu.ctx.Done():
				log.Printf("%v", cu.duration)
				return
			}
		}
	}()
}

func (cu *CacheUpdater) Stop() {
	if cu.ticker != nil {
		cu.ticker.Stop()
		cu.ticker = nil
	}
	if cu.cancel != nil {
		cu.cancel()
	}
}
