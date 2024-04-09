package cache

import (
	"CloudContent/internal/service"
	"context"
	//_ "github.com/gogf/gf/contrib/nosql/redis/v2"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/os/gcache"
	"time"
)

// sCache 缓存
type sCache struct {
	cache *gcache.Cache
}

func init() {
	service.RegisterCache(New())
}

func New() *sCache {
	c := &sCache{}
	c.cache = gcache.New()

	//使用Redis缓存
	//c.cache.SetAdapter(gcache.NewAdapterRedis(g.Redis()))

	return c
}

// IsCacheOut 有缓存直接输出，无缓存走自定义方法，并保存缓存
func (s *sCache) IsCacheOut(ctx context.Context, key interface{}, expire time.Duration, fun func() interface{}) interface{} {
	cache, _ := service.Cache().GetCache(ctx, key)
	if !cache.IsEmpty() {
		return cache.Interface()
	}
	pass := fun()
	_ = service.Cache().SetxCache(ctx, key, pass, expire)
	return pass
}

// GetCache 获取缓存
func (s *sCache) GetCache(ctx context.Context, key interface{}) (*gvar.Var, error) {
	return s.cache.Get(ctx, key)
}

// SetCache 设置缓存
func (s *sCache) SetCache(ctx context.Context, key, val interface{}) error {
	return s.cache.Set(ctx, key, val, 0)
}

// SetxCache 设置缓存,传入可key, value, 过期时间
func (s *sCache) SetxCache(ctx context.Context, key, val interface{}, expire time.Duration) error {
	return s.cache.Set(ctx, key, val, expire)
}

// GetExpire 获取获取时间
func (s *sCache) GetExpire(ctx context.Context, key interface{}) (time.Duration, error) {
	return s.cache.GetExpire(ctx, key)
}

// SetIfNotExist 当指定key的键值不存在时设置其对应的键值value并返回true，否则什么都不做并返回false。
func (s *sCache) SetIfNotExist(ctx context.Context, key, val interface{}, expire time.Duration) (bool, error) {
	return s.cache.SetIfNotExist(ctx, key, val, expire)
}

// Update 更新key的对应的键值，但不更改其过期时间，并返回旧值。如果缓存中不存在key，则返回的exist值为false。
func (s *sCache) Update(ctx context.Context, key, val interface{}) (oldValue *gvar.Var, exist bool, err error) {
	return s.cache.Update(ctx, key, val)
}

// RemoveCache 移除缓存
func (s *sCache) RemoveCache(ctx context.Context, key interface{}) error {
	_, err := s.cache.Remove(ctx, key)
	if err != nil {
		return err
	}
	return nil
}

// RemoveCacheList 移除多个缓存 g.Slice{"k1", "k2", "k3"}
func (s *sCache) RemoveCacheList(ctx context.Context, keys []interface{}) error {
	return s.cache.Removes(ctx, keys)
}

// ExistsCache 判断缓存是否存在
func (s *sCache) ExistsCache(ctx context.Context, key interface{}) (bool, error) {
	return s.cache.Contains(ctx, key)
}

// Close 关闭缓存
func (s *sCache) Close(ctx context.Context) error {
	return s.cache.Close(ctx)
}

// SetMapCache 设置Map
func (s *sCache) SetMapCache(ctx context.Context, val map[interface{}]interface{}, expire time.Duration) error {
	return s.cache.SetMap(ctx, val, expire)
}

// MustGet 读取Map
func (s *sCache) MustGet(ctx context.Context, key interface{}) *gvar.Var {
	return s.cache.MustGet(ctx, key)
}

// UpdateExpire 更新过期时间 并返回旧的过期时间值。如果缓存中不存在key，则返回-1。
func (s *sCache) UpdateExpire(ctx context.Context, key interface{}, expire time.Duration) (time.Duration, error) {
	return s.cache.UpdateExpire(ctx, key, expire)
}
