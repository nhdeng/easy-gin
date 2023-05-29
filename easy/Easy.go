package easy

import (
	"fmt"
	"github.com/gin-gonic/gin"
	Injector "github.com/shenyisyn/goft-ioc"
	"log"
	"reflect"
	"strings"
	"sync"
)

type Bean interface {
	Name() string
}

var innerRouter *EasyTree
var innerRouterOnce sync.Once

func getInnerRouter() *EasyTree {
	innerRouterOnce.Do(func() {
		innerRouter = NewEasyTree()
	})
	return innerRouter
}

type Easy struct {
	*gin.Engine
	g        *gin.RouterGroup
	exprData map[string]interface{}
	// 当前路由组
	currentGroup string
}

// Ignite Easy构造函数
func Ignite(ginMiddlewares ...gin.HandlerFunc) *Easy {
	g := &Easy{Engine: gin.New(), exprData: map[string]interface{}{}}
	// 强制加载异常处理中间件
	g.Use(ErrorHandle())
	for _, handler := range ginMiddlewares {
		g.Use(handler)
	}
	Injector.BeanFactory.Set(g)
	// 整个应用的配置加载进bean中
	config := InitConfig()
	Injector.BeanFactory.Set(config)
	Injector.BeanFactory.Set(NewGPAUtil())
	// 数据库实例对象加载进bean中
	db := InitGorm()
	Injector.BeanFactory.Set(db)
	return g
}

func (this *Easy) applyAll() {
	for t, v := range Injector.BeanFactory.GetBeanMapper() {
		if t.Elem().Kind() == reflect.Struct {
			Injector.BeanFactory.Apply(v.Interface())
		}
	}
}

// Launch 启动
func (this *Easy) Launch() {
	var port int32 = 8080
	if config := Injector.BeanFactory.Get((*SysConfig)(nil)); config != nil {
		port = config.(*SysConfig).Server.Port
	}
	this.applyAll()
	getCronTask().Start()
	err := this.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		log.Println(err)
	}
}

func (this *Easy) getPath(relativePath string) string {
	g := "/" + this.currentGroup
	if g == "/" {
		g = ""
	}
	g = g + relativePath
	g = strings.Replace(g, "//", "/", -1)
	return g
}

// Handle 重载gin.Handle方法
func (this *Easy) Handle(httpMethod, relativePath string, handler interface{}) *Easy {
	if h := Covert(handler); h != nil {
		methods := strings.Split(httpMethod, ",")
		for _, method := range methods {
			getInnerRouter().addRoute(method, this.getPath(relativePath), h)
			this.g.Handle(httpMethod, relativePath, h)
		}
	}
	return this
}

// Mount 挂载
func (this *Easy) Mount(group string, classes ...IClass) *Easy {
	this.g = this.Group(group)
	// 利用接口进行控制器挂载
	for _, class := range classes {
		this.currentGroup = group
		class.Build(this)
		// 将控制器也加入到bean容器中
		this.Beans(class)
	}
	// 返回自己方便链式调用
	return this
}

// Attach 添加中间件
func (this *Easy) Attach(f ...Fairing) *Easy {
	for _, f1 := range f {
		Injector.BeanFactory.Set(f1)
	}
	getFairingHandler().AddFairing(f...)
	return this
}

// Beans 设定数据库连接对象
func (this *Easy) Beans(beans ...Bean) *Easy {
	// 取出bean名称，加入到exprData里面
	for _, bean := range beans {
		this.exprData[bean.Name()] = bean
		Injector.BeanFactory.Set(bean)
	}
	return this
}

func (this *Easy) Config(cfgs ...interface{}) *Easy {
	Injector.BeanFactory.Config(cfgs...)
	return this
}

// Task 定时任务
func (this *Easy) Task(cron string, expr interface{}) *Easy {
	var err error
	if f, ok := expr.(func()); ok {
		_, err = getCronTask().AddFunc(cron, f)
	} else if exp, ok := expr.(Expr); ok {
		_, err = getCronTask().AddFunc(cron, func() {
			_, expErr := ExecExpr(exp, this.exprData)
			if expErr != nil {
				log.Println(expErr)
			}
		})
	}
	if err != nil {
		log.Println(err)
	}
	return this
}
