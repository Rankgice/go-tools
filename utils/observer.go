package utils

import (
	"fmt"
	"github.com/petermattis/goid"
	"log"
	"math/rand"
	"sync"
	"time"
)

var (
	//全局异步
	globalSyncMap = sync.Map{}

	GlobalSyncMap_Observer_Prefix = "observer_"

	GlobalSyncMap_LockFunc_Prefix = "lockFunc_"
)

// StoreInMap 存储到全局map
func StoreInMap(prefix, key string, value any) {
	globalSyncMap.Store(prefix+key, value)
}

// LoadOrStoreInMap 加载或存储数据到map，并返回数据
func LoadOrStoreInMap(prefix, key string, value any) (any, bool) {
	return globalSyncMap.LoadOrStore(prefix+key, value)
}

// LoadInMap 获取全局map中的数据
func LoadInMap(prefix, key string) (any, bool) {
	return globalSyncMap.Load(prefix + key)
}

// DeleteInMap 删除全局map中的数据
func DeleteInMap(prefix, key string) {
	globalSyncMap.Delete(prefix + key)
}

// LoadAndDeleteInMap 删除全局map中的数据，并返回被删除的数据
func LoadAndDeleteInMap(prefix, key string) (any, bool) {
	return globalSyncMap.LoadAndDelete(prefix + key)
}

// LockFunc 全局互斥锁，加锁
func LockFunc(key string) {
	mu, _ := LoadOrStoreInMap(GlobalSyncMap_LockFunc_Prefix, key, &sync.Mutex{})
	fmt.Println("Current Goroutine ID:", goid.Get(), "Try Locking", GlobalSyncMap_LockFunc_Prefix, key) //测试用
	mu.(*sync.Mutex).Lock()
	fmt.Println("Current Goroutine ID:", goid.Get(), "Locking Success", GlobalSyncMap_LockFunc_Prefix, key) //测试用
}

// UnlockFunc 全局互斥锁，解锁
func UnlockFunc(key string) {
	mu, ok := LoadInMap(GlobalSyncMap_LockFunc_Prefix, key)
	if ok {
		mutex, ok := mu.(*sync.Mutex)
		if ok {
			fmt.Println("Current Goroutine ID:", goid.Get(), "Unlocking Success", GlobalSyncMap_LockFunc_Prefix, key) //测试用
			mutex.Unlock()
		}
	}
}

// GetNoRepeatId 生成随机ID (前40位时间戳，后24位随机数)
func GetNoRepeatId() int64 {
	//获取当前时间(秒)
	curTimeSec := time.Now().Unix()
	//生成随机数
	randInt64 := rand.Int63()
	//生成随机ID (前40位时间戳，后24位随机数)
	randId := (curTimeSec << 24) | (randInt64 & 0xFFFFFF)
	return randId
}

// ObserverFunc 观察者函数
type ObserverFunc func(observerId int64, prefix, content string)

// RegisterLogObserver 注册日志观察者
func RegisterLogObserver(userId int64, observerFunc ObserverFunc) {
	randId := GetNoRepeatId()
	value, _ := LoadOrStoreInMap(GlobalSyncMap_Observer_Prefix, fmt.Sprintf("%d", userId), make(map[int64]ObserverFunc))
	if userLogObserverMap, ok := value.(map[int64]ObserverFunc); ok {
		userLogObserverMap[randId] = observerFunc
	} else {
		log.Println("RegisterLogObserver failed")
	}
}

// UnRegisterLogObserver 注销日志观察者
func UnRegisterLogObserver(userId int64, observerId int64) {
	value, _ := LoadOrStoreInMap(GlobalSyncMap_Observer_Prefix, fmt.Sprintf("%d", userId), make(map[int64]ObserverFunc))
	if userLogObserverMap, ok := value.(map[int64]ObserverFunc); ok {
		delete(userLogObserverMap, observerId)
	} else {
		log.Println("UnRegisterLogObserver failed")
	}
}

// GetLogObserverMap 获取日志观察者映射
func GetLogObserverMap(userId int64) map[int64]ObserverFunc {
	value, _ := LoadInMap(GlobalSyncMap_Observer_Prefix, fmt.Sprintf("%d", userId))
	userLogObserverMap, _ := value.(map[int64]ObserverFunc)
	return userLogObserverMap
}
