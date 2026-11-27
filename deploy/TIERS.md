# 三档分组运营手册

本站提供三个分组：满血、标准、低价。分组、倍率、模型白名单、定价配置和订阅套餐都在后台维护，不需要改代码。

## 一、档位定义

| 档位 | rate_multiplier 起点 | allowed_models | model_mapping | 计费来源 | rpm / 并发 |
| --- | --- | --- | --- | --- | --- |
| 满血 | 1.50 | 全部官方模型 | 空 | requested | 不限或高 |
| 标准 | 1.00 | 主力模型 | 空 | group_mapped | 300 |
| 低价 | 0.60 | 仅列出的便宜模型（DeepSeek / Qwen / 豆包等） | 空 | upstream | 60，低并发 |

倍率是用户侧价格系数，不是成本。上线前按实际上游成本重算，低价档至少覆盖 `上游成本 × 1.10`。

## 二、创建分组

后台 → 分组管理：

1. 新建三个分组，填写中英文名称和描述。
2. 设置 `display_brand`：满血写「官方模型直连」，标准写「标准池」，低价写实际模型名 + 「低并发」。
3. `allowed_protocols` 只勾选所接 provider 真正支持的协议。
4. `models_list_config` 保持默认，不要用它隐藏或改名模型。
5. 低价组的三个 fallback 字段（`fallback_group_id`、`fallback_group_id_on_invalid_request`、`unavailable_fallback_group_id`）全部留空，禁止回落到满血组。
6. `max_reasoning_effort` 留空；必须限制时用 `deny`，不要用 `downgrade`，让用户看到明确报错而不是静默降级。

## 三、路由策略

低价组：

```json
{
  "enabled": true,
  "model_mapping": {},
  "restrict_models": true,
  "restriction_model_source": "upstream",
  "allowed_models": ["<精确上游模型 ID 1>", "<精确上游模型 ID 2>"],
  "features": "",
  "features_config": {}
}
```

- `allowed_models` 必须是上游真实模型 ID，含版本/日期后缀。
- 不要用别名或通配符，除非该写法本身就透明且在定价页写明。
- `model_mapping` 必须为空。

满血/标准组可以使用友好别名，但要在分组描述和定价页说明别名指向的实际模型。

## 四、Provider 约束

低价组接入的每个 provider：

- `credentials.model_mapping` 与 `compact_model_mapping` 必须为空。
- provider 级模型白名单只能包含低价组 `allowed_models` 里的模型。
- 不允许 provider 级别名、重定向或强制改写。
- provider 必须支持该组的协议和精确模型 ID。

provider 内部需要别名时，不要接入低价组，直到别名移除或目标模型被明确展示。

## 五、定价配置

后台 → 定价：

1. 建「满血 + 标准」配置：`group_ids` 选两档，`billing_model_source=group_mapped`，价卡按模型配置。
2. 建「低价直连」配置：`group_ids` 只选低价组，`billing_model_source=upstream`。
3. 低价组白名单里的每个模型都必须有价卡；缺价会直接拒单，不会静默按别的模型计费。
4. 一个分组最多挂一个定价配置。

## 六、订阅套餐

后台 → 套餐：

- 每个档位建一个套餐，`group_ids` 明确选择对应分组；空列表等于全部分组，分层产品不要留空。
- 设置日/周/月限额，低价组和标准组必须有上限。
- 套餐里的 `group_rate_multipliers` 会覆盖组倍率，促销时使用。
- `for_sale` 验证通过后再打开。

## 七、验收

```sql
-- 组与定价配置对应
SELECT g.name, pc.billing_model_source, pc.status
FROM groups g
LEFT JOIN pricing_config_groups pcg ON pcg.group_id = g.id
LEFT JOIN pricing_configs pc ON pc.id = pcg.pricing_config_id
WHERE g.deleted_at IS NULL AND g.name IN ('满血', '标准', '低价');

-- 低价组 provider 不得有映射
SELECT p.name, p.credentials->'model_mapping' AS model_mapping
FROM providers p
JOIN provider_groups pg ON pg.provider_id = p.id
JOIN groups g ON g.id = pg.group_id
WHERE g.name = '低价';

-- 用量透明：model = requested_model = upstream_model，mapping_chain 为空
SELECT created_at, model, requested_model, upstream_model,
       upstream_model_mismatch, model_mapping_chain, actual_cost
FROM usage_logs
WHERE group_id = (SELECT id FROM groups WHERE name = '低价')
ORDER BY created_at DESC
LIMIT 50;
```

再用低价组的测试 Key 各发一次流式和非流式请求，确认：

- `/v1/models` 只返回低价组明示的模型。
- 扣费倍率符合定价配置。
- 后台用量记录里的上游模型与请求模型一致，`model_mapping_chain` 为空。
- 出现 `upstream_model_mismatch` 视为事故：定位 provider，移除映射或下线该 provider，并在公告里说明。

## 八、回滚

- 低价组 `status` 设为 inactive，停止新用户进入。
- 套餐 `for_sale` 设为 false。
- 必要时解绑低价组定价配置、移除 provider。
- 不删除历史用量和账单；配置错误影响计费时，用数据库备份恢复分组与定价配置。
