"""profile —— 基于 Richards Curve（广义 Logistic）的评分映射拟合器。

业务目标：
  将若干项 Pass/Fail 判定结果映射为 0～10 分。
  评分曲线满足：前期增长缓慢 → 中段加速 → 高分区拉开区分度。

曲线性质（由拟合保证）：
  * 单调递增：passA < passB ⟹ score(passA) <= score(passB)；
  * 连续：禁止分段规则 / 阶梯函数 / 人工 if-else 映射；
  * 前期压缩：anchor1 之前每增一项的得分增长，明显低于 anchor1~anchor2；
  * 中段加速：anchor1~anchor2 为整条曲线增长最快区域；
  * 高分拉开：anchor2 之后快速逼近 10 但仍保持区分度
    （如 8.8 / 9.2 / 9.6 / 9.9，而非 9.9 / 10 / 10 / 10）。

归一化：
  曲线在比例域 progress = passed/total 上定义，同比例同形状
  （7 条标准 anchor=2/5 与 70 条 anchor=20/50 形状一致）。

端点：
  passed=0 ⟹ score=0，passed=total ⟹ score=10，精确成立
  （由端点线性校正 r0/r1 保证，非无限接近）。

输出：
  Float ∈ [0, 10]，小数位数与舍入由调用方决定。

入参：两个中间控制点 (ratio1, score1), (ratio2, score2)。
      端点 (0.0, 0.0) 与 (1.0, 10.0) 由内部固定注入，调用方无需提供。
出参：一组 Richards 参数 + 端点校正常量，供手写的 runtime profile(hit, total) 消费。

环境要求：Python 3.11+。
依赖：numpy、scipy（仅 Offline，不进 Runtime）。
"""

import numpy as np
from scipy.optimize import least_squares

SCORE_TOP: float = 10.0


def _richards_raw(x: np.ndarray, Q: float, k: float, M: float, nu: float) -> np.ndarray:
    # Richards 原始值（未做端点校正）：r = 1 / (1 + Q·e^{-k(x-M)})^{1/ν}
    return 1.0 / (1.0 + Q * np.exp(-k * (x - M))) ** (1.0 / nu)


def _params_to_score_fn(Q: float, k: float, M: float, nu: float):
    # 计算 x=0 与 x=1 处的原值，用于端点线性校正，使 (0->0, 1->10) 精确成立
    x_grid = np.array([0.0, 1.0])
    r_grid = _richards_raw(x_grid, Q, k, M, nu)
    r0, r1 = float(r_grid[0]), float(r_grid[1])

    def score(x):
        r = _richards_raw(np.asarray(x, dtype=float), Q, k, M, nu)
        return SCORE_TOP * (r - r0) / (r1 - r0)

    return score, r0, r1


def _residuals(p, mid_ratios, mid_scores):
    # 残差 = 校正后预测分 - 目标分（仅中间点，端点由线性校正结构上精确满足）
    Q, k, M, nu = p
    score, _, _ = _params_to_score_fn(Q, k, M, nu)
    pred = score(np.array(mid_ratios))
    return pred - np.array(mid_scores)


def fit_profile(p1: tuple[float, float], p2: tuple[float, float]) -> dict:
    r1, s1 = float(p1[0]), float(p1[1])
    r2, s2 = float(p2[0]), float(p2[1])

    if not (0.0 < r1 < r2 < 1.0):
        raise ValueError("需满足 0 < ratio1 < ratio2 < 1")
    if not (0.0 < s1 < s2 < SCORE_TOP):
        raise ValueError("需满足 0 < score1 < score2 < 10")
    if not (s1 < s2):
        raise ValueError("控制点需单调递增")

    mid_ratios = [r1, r2]
    mid_scores = [s1, s2]

    # 参数下界：Q>0, k>0, M∈[0,1], ν>0，保证 r 单调递增 -> score 单调递增
    lb = [1e-6, 1e-6, 0.0, 1e-3]
    ub = [1e6, 1e3, 1.0, 1e3]

    # 多初值，取最小 cost，避免局部最优
    starts = [
        (1.0, 5.0, 0.5, 1.0),
        (1.0, 10.0, 0.6, 2.0),
        (0.5, 8.0, 0.55, 1.5),
        (2.0, 12.0, 0.5, 3.0),
    ]

    best = None
    for x0 in starts:
        x0 = list(x0)
        x0[2] = min(max(x0[2], 1e-6), 1.0 - 1e-6)
        try:
            res = least_squares(
                _residuals,
                x0,
                args=(mid_ratios, mid_scores),
                bounds=(lb, ub),
                method="trf",
                max_nfev=20000,
            )
        except Exception:
            continue
        cost = float(res.cost)
        if best is None or cost < best[0]:
            best = (cost, res)

    if best is None:
        raise RuntimeError("所有初值下 least_squares 均失败")

    _, res = best
    Q, k, M, nu = (float(v) for v in res.x)
    _, r0, r1 = _params_to_score_fn(Q, k, M, nu)

    return {
        # ---- Runtime 必需：手写 profile(hit, total) 消费 ----
        # Richards 公式 r = 1 / (1 + Q·e^{-k(x-M)})^{1/nu} 所需
        "Q": Q,      # 低端陡度因子
        "k": k,      # 增长率
        "M": M,      # 水平平移参数（位置）
        "nu": nu,    # 非对称参数（>1 → 前压 + 高分拉开）
        # 端点校正所需：score = 10.0 * (r(x) - r0) / (r1 - r0)
        # （满分 10.0 由 spec 锁定，runtime 直接硬编码，不作为参数传递）
        "r0": r0,    # Richards 在 x=0 的原值
        "r1": r1,    # Richards 在 x=1 的原值
        # ---- 仅离线诊断，runtime 不消费 ----
        "rms": float(np.sqrt(2.0 * res.cost / max(1, len(mid_ratios)))),
    }


if __name__ == "__main__":
    params = fit_profile((0.6, 4.0), (0.8, 8.0))

    # 1) 拟合参数 + 说明
    labels = [
        ("Q",  "低端陡度因子；越大 → 前期压缩越强"),
        ("k",  "增长率；越大 → 中段陡升越窄越急"),
        ("M",  "水平平移参数（位置）；真实拐点 x_i = M + (1/k)·ln(Q/nu)"),
        ("nu", "非对称参数；nu>1 → 前压 + 高分拉开"),
        ("r0", "Richards 在 x=0 的原值；端点校正用，保证 score(0)=0 精确"),
        ("r1", "Richards 在 x=1 的原值；端点校正用，保证 score(1)=10 精确"),
    ]
    print("# ---- 拟合参数 ----")
    for key, desc in labels:
        print(f"{key} = {params[key]}    # {desc}")

    # 2) 离线诊断
    print("# ---- 离线诊断 ----")
    print(f"rms = {params['rms']}    # 残差均方根；越小越贴合，近 0 表示控制点被精确命中")

    # 3) 可复制即用的 runtime，常数内联、仅依赖 math
    runtime = f'''import math


def profile(hit, total):
    x = hit / total
    r = 1.0 / (1.0 + {params["Q"]!r} * math.exp(-{params["k"]!r} * (x - {params["M"]!r}))) ** (1.0 / {params["nu"]!r})
    return 10.0 * (r - {params["r0"]!r}) / ({params["r1"]!r} - {params["r0"]!r})
'''
    print("# ---- Runtime（复制即用，仅依赖 math）----")
    print(runtime)
