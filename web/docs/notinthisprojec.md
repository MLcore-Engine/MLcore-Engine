```javascript
const config = {
  // `url` 是请求的服务器地址
  url: '/user',

  // `method` 是创建请求时使用的方法
  method: 'get', // 默认是 get

  // `baseURL` 将自动加在 `url` 前面，除非 `url` 是一个绝对 URL
  baseURL: 'https://api.example.com',

  // `headers` 是即将被发送的自定义请求头
  headers: {
    'X-Requested-With': 'XMLHttpRequest',
    'Content-Type': 'application/json'
  },

  // `params` 是即将与请求一起发送的 URL 参数
  params: {
    ID: 12345
  },

  // `data` 是作为请求主体被发送的数据
  // 只适用于 'PUT', 'POST', 'DELETE 和 'PATCH' 请求方法
  data: {
    firstName: 'Fred',
    lastName: 'Flintstone'
  },

  // `timeout` 指定请求超时的毫秒数
  timeout: 1000, // 默认是 `0` (永不超时)

  // `withCredentials` 表示跨域请求时是否需要使用凭证
  withCredentials: false, // 默认的

  // `responseType` 表示服务器响应的数据类型
  responseType: 'json', // 默认的

  // `validateStatus` 定义对于给定的HTTP 响应状态码是 resolve 或 reject  promise
  validateStatus: function (status) {
    return status >= 200 && status < 300; // 默认的
  },

  // `maxRedirects` 定义在 node.js 中 follow 的最大重定向数目
  maxRedirects: 5, // 默认的
};

```

- **url**: 请求的目标 URL。
- **method**: HTTP 方法，如 GET, POST, PUT, DELETE 等。
- **baseURL**: 如果指定，会被加到 url 前面（除非 url 是绝对路径）。
- **headers**: 自定义的请求头。
- **params**: URL 参数，会被添加到 URL 的查询字符串中。
- **data**: 请求体数据，通常用于 POST, PUT 等方法。
- **timeout**: 请求超时时间（毫秒）。
- **withCredentials**: 是否携带跨域请求凭证。
- **responseType**: 期望的响应数据类型。
- **validateStatus**: 自定义哪些 HTTP 状态码应该被视为成功。
- **maxRedirects**: 最大重定向次数（仅在 Node.js 中有效）。
