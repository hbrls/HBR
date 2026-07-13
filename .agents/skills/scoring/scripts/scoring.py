"""scoring —— 评分 runtime。

支持三种评分模式，由第一个命令行参数选择：
  profile      Richards 曲线（前期缓慢、中段加速、高分拉开）
  fingerprint  全局四次 Bernstein 曲线（前期较低、高分区提升）
  risk         全局十次 Bernstein 反 S 曲线（前期快速、中段放缓）
输出 5 分制，保留一位小数。仅依赖 Python 标准库。
"""

import math


def profile(hit, total):
    """端点校正 Richards 评分（10 分制）。"""
    q = 5.352468069154516
    k = 10.45084717683695
    m = 0.6044137118506592
    nu = 2.038298714159775
    r0 = 0.01979831546286067
    r1 = 0.9604536006148702

    x = hit / total
    r = 1.0 / (1.0 + q * math.exp(-k * (x - m))) ** (1.0 / nu)
    return 10.0 * (r - r0) / (r1 - r0)


def fingerprint(hit, total):
    """全局四次 Bernstein 评分（10 分制）。"""
    x = hit / total
    if x <= 0.0:
        return 0.001
    if x >= 1.0:
        return 10.0

    one_minus = 1.0 - x
    return (
        0.001 * one_minus**4
        + 4.0 * 3.0057589285714315 * x * one_minus**3
        + 6.0 * 3.244053571428568 * x**2 * one_minus**2
        + 4.0 * 3.244053571428568 * x**3 * one_minus
        + 10.0 * x**4
    )


def risk(hit, total):
    """全局十次 Bernstein 反 S 评分（10 分制）。"""
    x = hit / total
    if x <= 0.0:
        return 0.0
    if x >= 1.0:
        return 10.0

    values = [
        0.0,
        3.8268578751698312,
        7.6537147504583345,
        8.25690976228198,
        8.25690976224767,
        8.256909762243078,
        8.35449578065538,
        8.765871585494047,
        9.177247390323654,
        9.5886231951632,
        10.0,
    ]
    for level in range(10, 0, -1):
        for index in range(level):
            values[index] = (
                (1.0 - x) * values[index] + x * values[index + 1]
            )
    return values[0]


_MODES = {
    "profile": profile,
    "fingerprint": fingerprint,
    "risk": risk,
}


if __name__ == "__main__":
    import sys

    if len(sys.argv) != 4:
        print(
            "usage: python scoring.py <profile|fingerprint|risk> <hit> <total>",
            file=sys.stderr,
        )
        raise SystemExit(1)

    mode = sys.argv[1]
    hit = int(sys.argv[2])
    total = int(sys.argv[3])

    if mode not in _MODES:
        print(
            f"unknown mode: {mode} (profile|fingerprint|risk)",
            file=sys.stderr,
        )
        raise SystemExit(1)

    score = _MODES[mode](hit, total)
    score5 = round(score / 2.0, 1)
    print(score5)
