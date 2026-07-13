"""fingerprint —— 全局四次 Bernstein S 曲线拟合器。

约束：S(0)=0.001、S(0.6)=4、S(0.8)=6、S(1)>9.99；输出限制在 [0,10]，
曲线必须全局、连续、单调，不允许分段；视觉上只有一个主要转折，数学上最多允许
两个曲率零点。

模型：S(x)=sum(c_i*B_i,4(x), i=0..4)，固定 c_0=0.001、c_4=10。
单调系数 c_0<=...<=c_4 保证曲线单调且不越界；四次模型最多产生两个拐点。

求解：两个中间锚点形成线性等式。使用 SVD 得到一维解空间，解析求单调可行区间，
再以最小弯曲能量选择视觉最平滑的解。只执行一次一维有界优化，不做多次非线性回归。

依赖：numpy、scipy（仅离线拟合；runtime 为纯标准库）。
"""

import numpy as np
from scipy.optimize import minimize_scalar

LOW_SCORE: float = 0.001
HIGH_SCORE: float = 10.0
ANCHOR_TOLERANCE: float = 1e-9
MONOTONIC_TOLERANCE: float = 1e-12
QUALITY_MINIMUM_AT_ONE: float = 9.99
DIAGNOSTIC_GRID_SIZE: int = 1001


def _basis4(x):
    values = np.asarray(x, dtype=float)
    one_minus = 1.0 - values
    return np.column_stack(
        [
            one_minus**4,
            4.0 * values * one_minus**3,
            6.0 * values**2 * one_minus**2,
            4.0 * values**3 * one_minus,
            values**4,
        ]
    )


def _bernstein4(x, coefficients):
    return _basis4(x) @ np.asarray(coefficients, dtype=float)


def _second_derivative(x, coefficients):
    values = np.asarray(x, dtype=float)
    second_differences = np.diff(np.asarray(coefficients, dtype=float), n=2)
    basis2 = np.column_stack(
        [
            (1.0 - values) ** 2,
            2.0 * values * (1.0 - values),
            values**2,
        ]
    )
    return 12.0 * (basis2 @ second_differences)


def _bending_energy(coefficients):
    grid = np.linspace(0.0, 1.0, DIAGNOSTIC_GRID_SIZE)
    curvature = _second_derivative(grid, coefficients)
    return float(np.mean(np.square(curvature)))


def _internal_linear_system(ratios, scores):
    basis = _basis4(ratios)
    matrix = basis[:, 1:4]
    right_hand_side = (
        np.asarray(scores, dtype=float)
        - LOW_SCORE * basis[:, 0]
        - HIGH_SCORE * basis[:, 4]
    )
    return matrix, right_hand_side


def _feasible_interval(particular, direction):
    # c1-c0 >= 0, c2-c1 >= 0, c3-c2 >= 0, c4-c3 >= 0
    offsets = np.array(
        [
            particular[0] - LOW_SCORE,
            particular[1] - particular[0],
            particular[2] - particular[1],
            HIGH_SCORE - particular[2],
        ],
        dtype=float,
    )
    slopes = np.array(
        [
            direction[0],
            direction[1] - direction[0],
            direction[2] - direction[1],
            -direction[2],
        ],
        dtype=float,
    )

    lower = -np.inf
    upper = np.inf
    for offset, slope in zip(offsets, slopes):
        if abs(slope) <= 1e-15:
            if offset < -MONOTONIC_TOLERANCE:
                raise RuntimeError("Bernstein 线性约束不存在单调可行解")
            continue

        boundary = -offset / slope
        if slope > 0.0:
            lower = max(lower, boundary)
        else:
            upper = min(upper, boundary)

    if not np.isfinite(lower) or not np.isfinite(upper) or lower > upper:
        raise RuntimeError(
            f"Bernstein 单调可行区间无效：lower={lower}, upper={upper}"
        )
    return float(lower), float(upper)


def _coefficients_at(value, particular, direction):
    internal = particular + value * direction
    return np.array(
        [LOW_SCORE, internal[0], internal[1], internal[2], HIGH_SCORE],
        dtype=float,
    )


def _count_inflections(coefficients):
    grid = np.linspace(0.0, 1.0, DIAGNOSTIC_GRID_SIZE)
    curvature = _second_derivative(grid, coefficients)
    threshold = max(1e-10, float(np.max(np.abs(curvature))) * 1e-8)
    signs = np.sign(curvature)
    signs[np.abs(curvature) <= threshold] = 0.0
    nonzero = signs[signs != 0.0]
    if len(nonzero) < 2:
        return 0
    return int(np.count_nonzero(nonzero[1:] != nonzero[:-1]))


def fit_fingerprint(p1: tuple[float, float], p2: tuple[float, float]) -> dict:
    r1, s1 = float(p1[0]), float(p1[1])
    r2, s2 = float(p2[0]), float(p2[1])

    if not (0.0 < r1 < r2 < 1.0):
        raise ValueError("需满足 0 < ratio1 < ratio2 < 1")
    if not (LOW_SCORE < s1 < s2 < HIGH_SCORE):
        raise ValueError("需满足 0.001 < score1 < score2 < 10")

    fit_ratios = np.array([r1, r2], dtype=float)
    fit_scores = np.array([s1, s2], dtype=float)
    matrix, right_hand_side = _internal_linear_system(fit_ratios, fit_scores)

    particular = np.linalg.lstsq(matrix, right_hand_side, rcond=None)[0]
    _, singular_values, vh = np.linalg.svd(matrix, full_matrices=True)
    direction = vh[-1]
    lower, upper = _feasible_interval(particular, direction)

    result = minimize_scalar(
        lambda value: _bending_energy(
            _coefficients_at(value, particular, direction)
        ),
        bounds=(lower, upper),
        method="bounded",
        options={"xatol": 1e-14, "maxiter": 1000},
    )

    trial_values = [lower, upper]
    if np.isfinite(result.x):
        trial_values.append(float(result.x))

    candidates = []
    diagnostic_ratios = np.array([0.0, r1, r2, 1.0], dtype=float)
    diagnostic_scores = np.array([LOW_SCORE, s1, s2, HIGH_SCORE], dtype=float)

    for value in trial_values:
        coefficients = _coefficients_at(value, particular, direction)
        predictions = _bernstein4(diagnostic_ratios, coefficients)
        residuals = predictions - diagnostic_scores
        differences = np.diff(coefficients)
        candidates.append(
            {
                "free_value": value,
                "coefficients": coefficients,
                "predictions": predictions,
                "residuals": residuals,
                "max_anchor_error": float(np.max(np.abs(residuals))),
                "minimum_coefficient_step": float(np.min(differences)),
                "bending_energy": _bending_energy(coefficients),
                "inflection_count": _count_inflections(coefficients),
            }
        )

    candidates.sort(
        key=lambda candidate: (
            candidate["max_anchor_error"],
            candidate["bending_energy"],
        )
    )
    best = candidates[0]

    if best["max_anchor_error"] > ANCHOR_TOLERANCE:
        raise RuntimeError(
            f"Bernstein 锚点误差超限：{best['residuals'].tolist()}"
        )
    if best["minimum_coefficient_step"] < -MONOTONIC_TOLERANCE:
        raise RuntimeError(
            f"Bernstein 系数不单调：{best['coefficients'].tolist()}"
        )
    if best["predictions"][-1] <= QUALITY_MINIMUM_AT_ONE:
        raise RuntimeError(
            f"质量门禁失败：S(1)={best['predictions'][-1]}"
        )
    if best["inflection_count"] > 2:
        raise RuntimeError(
            f"视觉形状失败：检测到 {best['inflection_count']} 个曲率零点"
        )

    return {
        # ---- Runtime 必需 ----
        "model": "global_bernstein_4",
        "coefficients": [float(value) for value in best["coefficients"]],
        # ---- 仅离线诊断 ----
        "predictions": [float(value) for value in best["predictions"]],
        "residuals": [float(value) for value in best["residuals"]],
        "max_anchor_error": best["max_anchor_error"],
        "minimum_coefficient_step": best["minimum_coefficient_step"],
        "bending_energy": best["bending_energy"],
        "inflection_count": best["inflection_count"],
        "free_value": best["free_value"],
        "singular_values": [float(value) for value in singular_values],
    }


if __name__ == "__main__":
    params = fit_fingerprint((0.6, 4.0), (0.8, 6.0))

    print("# ---- 拟合参数 ----")
    print(f"model = {params['model']}")
    print(f"coefficients = {params['coefficients']}")

    print("# ---- 离线诊断 ----")
    for key in (
        "predictions",
        "residuals",
        "max_anchor_error",
        "minimum_coefficient_step",
        "bending_energy",
        "inflection_count",
        "free_value",
        "singular_values",
    ):
        print(f"{key} = {params[key]}")

    coefficients = params["coefficients"]
    runtime = f'''def fingerprint(progress):
    if progress <= 0.0:
        return {coefficients[0]!r}
    if progress >= 1.0:
        return {coefficients[4]!r}

    x = progress
    one_minus = 1.0 - x
    return (
        {coefficients[0]!r} * one_minus ** 4
        + 4.0 * {coefficients[1]!r} * x * one_minus ** 3
        + 6.0 * {coefficients[2]!r} * x ** 2 * one_minus ** 2
        + 4.0 * {coefficients[3]!r} * x ** 3 * one_minus
        + {coefficients[4]!r} * x ** 4
    )
'''
    print("# ---- Runtime（复制即用，纯标准库）----")
    print(runtime)
