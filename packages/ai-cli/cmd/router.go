package cmd

import (
	"fmt"

	"github.com/ucarion/urlpath"
)

// Router 路由管理器
type Router struct {
	routes []registeredRoute
}

// registeredRoute 已注册的路由
type registeredRoute struct {
	Method   string
	Pattern  string
	Matcher  urlpath.Path
	Callback Callback
}

// Callback 路由处理函数
// method: HTTP 方法，如 "LIST", "POST", "GET"
// pattern: 路由模式，如 "/apps/:appName/plans"
// pathname: 实际请求路径，如 "/apps/flowablex/plans"
// params: 解析出的路径参数，如 { appName: "flowablex" }
// args: 命令行传入的 flags 参数
type Callback func(method, pattern, pathname string, params map[string]string, args map[string]string) (interface{}, error)

// Verb 注册路由
// method: HTTP 方法
// pattern: 路由模式
// callback: 处理函数
func (r *Router) Verb(method, pattern string, callback Callback) {
	r.routes = append(r.routes, registeredRoute{
		Method:   method,
		Pattern:  pattern,
		Matcher:  urlpath.New(pattern),
		Callback: callback,
	})
}

// Invoke 调用路由
// method: HTTP 方法
// pathname: 请求路径
// args: 命令行传入的 flags 参数
func (r *Router) Invoke(method, pathname string, args map[string]string) (interface{}, error) {
	// fmt.Println("[DEBUG] registered routes:")
	// for i, route := range r.routes {
	// 	fmt.Printf("  [%d] %-6s %s\n", i, route.Method, route.Pattern)
	// }

	for _, route := range r.routes {
		if route.Method != method {
			continue
		}

		m, ok := route.Matcher.Match(pathname)
		if !ok {
			continue
		}

		// fmt.Printf("[DEBUG] matched route: %s %s\n", route.Method, route.Pattern)

		params := m.Params
		if params == nil {
			params = map[string]string{}
		}

		return route.Callback(method, route.Pattern, pathname, params, args)
	}

	return nil, fmt.Errorf("unsupported path: %s %s", method, pathname)
}
