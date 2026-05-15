---
name: 热点新闻技能
description: 当用户询问新闻、热点话题时使用此技能，获取最新资讯
---

# 热点新闻技能

## 使用方法
使用 `bash` 工具调用新闻查询：

```bash
# 查询新浪新闻热点（国内可访问）
curl -s "https://news.sina.com.cn/" | grep -o '<a[^>]*href="[^"]*"[^>]*>[^<]*</a>' | grep -i '热点\|新闻\|头条' | head -10

# 查询网易新闻
curl -s "https://news.163.com/" | grep -o '<a[^>]*href="[^"]*"[^>]*>[^<]*</a>' | head -10

# 使用新闻聚合 API
curl -s "https://api.muxiaoguo.cn/api/gupiao" | python3 -m json.tool
```

## 支持的分类
| 分类 | 说明 |
|------|------|
| general | 综合新闻 |
| tech | 科技新闻 |
| finance | 财经新闻 |
| sports | 体育新闻 |

## 注意事项
- 使用国内可访问的新闻源
- 返回标题和链接
- 可根据需求调整抓取数量