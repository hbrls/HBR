---
name: scoring
description: "评分 scoring：将 M 中 N 通过数量映射为评分。仅支持 profile（Richards：前期缓慢、中段加速、高分拉开）、fingerprint（全局四次 Bernstein：前期较低、中段平缓、高分区提升）和 risk（全局十次 Bernstein 反 S：前期快速得分、中段放缓、尾部恢复区分度）三种模式。输出 5 分制、保留一位小数。使用 scripts/scoring.py 直接计算。"
---

# Scoring - 评分映射

> updated_by: HBR - GPT-5
> updated_at: 2026-07-14 23:07:58

## 模式

将 `total` 项判定标准中的通过数量 `hit` 映射为评分。仅支持三种模式：

- **profile**：端点校正 Richards 曲线，前期缓慢、中段加速、高分区拉开。
- **fingerprint**：全局四次 Bernstein 曲线，前期较低、中段平缓、高分区明显提升。
- **risk**：全局十次 Bernstein 反 S 曲线，前期快速得分、中段放缓、尾部恢复区分度。

三种模式均在归一化进度 `x = hit/total` 上定义，与 `total` 解耦。

| 模式 | 10 分制锚点 | 端点 |
|---|---|---|
| `profile` | `S(0.6)=4`、`S(0.8)=8` | `S(0)=0`、`S(1)=10` |
| `fingerprint` | `S(0.6)=4`、`S(0.8)=6` | `S(0)=0.001`、`S(1)=10` |
| `risk` | `S(0.2)=6`、`S(0.4)=8` | `S(0)=0`、`S(1)=10` |

## 使用

```bash
python .agents/skills/scoring/scripts/scoring.py <profile|fingerprint|risk> <hit> <total>
```

- 第一个参数：`profile`、`fingerprint` 或 `risk`。
- `hit`：通过数量，满足 `0 <= hit <= total`。
- `total`：标准总数，满足 `total >= 2`。
- 输出：5 分制 `Float`，保留一位小数。

锚点示例：

```bash
python .agents/skills/scoring/scripts/scoring.py profile 8 10
# 输出：4.0

python .agents/skills/scoring/scripts/scoring.py fingerprint 6 10
# 输出：2.0

python .agents/skills/scoring/scripts/scoring.py risk 2 10
# 输出：3.0
```

## Runtime 约束

- 直接调用 `scripts/scoring.py`，不要在运行阶段重新拟合参数。
- runtime 仅依赖 Python 标准库；NumPy、SciPy 只属于离线拟合环境。
- 三个函数内部计算 10 分制原始分，CLI 最后除以 2 并保留一位小数。
- fingerprint 和 risk 使用全局 Bernstein 多项式，不是分段曲线。
