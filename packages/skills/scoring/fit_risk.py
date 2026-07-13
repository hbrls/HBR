"""risk —— 全局十次 Bernstein 反 S 曲线拟合器。

约束：S(0)=0、S(0.2)=6、S(0.4)=8、S(1)=10；输出限制在 [0,10]，
曲线必须全局、连续、单调，不允许分段；前段凹、后段凸，只有一个曲率换向。

模型：S(x)=sum(c_i*B_i,10(x), i=0..10)，固定 c_0=0、c_10=10。
单调系数 c_0<=...<=c_10 保证曲线单调且不越界。

求解：两个中间锚点形成线性等式。使用 SVD 得到零空间，在每个可能的曲率换向
位置上用线性规划寻找可行点，再以最小弯曲能量选择视觉最平滑的全局解。

依赖：numpy、scipy（仅离线拟合；runtime 为纯标准库）。
"""

from math import comb

import numpy as np
from scipy.optimize import LinearConstraint, linprog, minimize

DEGREE: int = 10
LOW_SCORE: float = 0.0
HIGH_SCORE: float = 10.0
ANCHOR_TOLERANCE: float = 1e-9
SOLVER_CONSTRAINT_TOLERANCE: float = 1e-8
CURVATURE_MARGIN: float = 1e-6
DIAGNOSTIC_GRID_SIZE: int = 1001


def _basis(x, degree):
    values = np.atleast_1d(np.asarray(x, dtype=float))
    one_minus = 1.0 - values
    return np.column_stack(
        [
            comb(degree, index)
            * values**index
            * one_minus ** (degree - index)
            for index in range(degree + 1)
        ]
    )


def _bernstein10(x, coefficients):
    return _basis(x, DEGREE) @ np.asarray(coefficients, dtype=float)


def _second_derivative(x, coefficients):
    second_differences = np.diff(np.asarray(coefficients, dtype=float), n=2)
    return (
        DEGREE
        * (DEGREE - 1)
        * (_basis(x, DEGREE - 2) @ second_differences)
    )


def _bending_energy(coefficients):
    grid = np.linspace(0.0, 1.0, DIAGNOSTIC_GRID_SIZE)
    curvature = _second_derivative(grid, coefficients)
    return float(np.mean(np.square(curvature)))


def _internal_linear_system(ratios, scores):
    basis = _basis(ratios, DEGREE)
    matrix = basis[:, 1:-1]
    right_hand_side = (
        np.asarray(scores, dtype=float)
        - LOW_SCORE * basis[:, 0]
        - HIGH_SCORE * basis[:, -1]
    )
    return matrix, right_hand_side


def _coefficient_map(particular, nullspace):
    base = np.concatenate(([LOW_SCORE], particular, [HIGH_SCORE]))
    transform = np.zeros((DEGREE + 1, nullspace.shape[1]), dtype=float)
    transform[1:-1] = nullspace
    return base, transform


def _coefficients_at(free_values, base, transform):
    return base + transform @ np.asarray(free_values, dtype=float)


def _difference_operators():
    identity = np.eye(DEGREE + 1, dtype=float)
    return np.diff(identity, axis=0), np.diff(identity, n=2, axis=0)


def _profile_constraints(split, base, transform):
    first_difference, second_difference = _difference_operators()

    # 一阶差分非负保证单调；二阶差分先非正后非负保证只有一个反 S 换向。
    operators = [
        first_difference,
        -second_difference[:split],
        second_difference[split:],
        -second_difference[[0]],
        second_difference[[-1]],
    ]
    minimums = [
        np.zeros(DEGREE, dtype=float),
        np.zeros(split, dtype=float),
        np.zeros(DEGREE - 1 - split, dtype=float),
        np.array([CURVATURE_MARGIN], dtype=float),
        np.array([CURVATURE_MARGIN], dtype=float),
    ]

    operator = np.vstack(operators)
    minimum = np.concatenate(minimums)
    matrix = operator @ transform
    lower = minimum - operator @ base
    return matrix, lower


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


def _diagnose_candidate(
    free_values,
    split,
    base,
    transform,
    diagnostic_ratios,
    diagnostic_scores,
    maximum_constraint_violation,
    candidate_source,
    optimizer_success,
    optimizer_message,
):
    coefficients = _coefficients_at(free_values, base, transform)
    predictions = _bernstein10(diagnostic_ratios, coefficients)
    residuals = predictions - diagnostic_scores
    coefficient_steps = np.diff(coefficients)
    second_differences = np.diff(coefficients, n=2)
    return {
        "shape_split": split,
        "free_values": np.asarray(free_values, dtype=float),
        "coefficients": coefficients,
        "predictions": predictions,
        "residuals": residuals,
        "max_anchor_error": float(np.max(np.abs(residuals))),
        "minimum_coefficient_step": float(np.min(coefficient_steps)),
        "start_second_difference": float(second_differences[0]),
        "end_second_difference": float(second_differences[-1]),
        "maximum_constraint_violation": maximum_constraint_violation,
        "bending_energy": _bending_energy(coefficients),
        "inflection_count": _count_inflections(coefficients),
        "candidate_source": str(candidate_source),
        "optimizer_success": bool(optimizer_success),
        "optimizer_message": str(optimizer_message),
    }


def _passes_quality(candidate):
    return (
        candidate["candidate_source"] == "bending_optimized"
        and candidate["max_anchor_error"] <= ANCHOR_TOLERANCE
        and candidate["maximum_constraint_violation"]
        <= SOLVER_CONSTRAINT_TOLERANCE
        and candidate["inflection_count"] == 1
    )


def fit_risk(p1: tuple[float, float], p2: tuple[float, float]) -> dict:
    r1, s1 = float(p1[0]), float(p1[1])
    r2, s2 = float(p2[0]), float(p2[1])

    if not (0.0 < r1 < r2 < 1.0):
        raise ValueError("需满足 0 < ratio1 < ratio2 < 1")
    if not (LOW_SCORE < s1 < s2 < HIGH_SCORE):
        raise ValueError("需满足 0 < score1 < score2 < 10")

    fit_ratios = np.array([r1, r2], dtype=float)
    fit_scores = np.array([s1, s2], dtype=float)
    matrix, right_hand_side = _internal_linear_system(fit_ratios, fit_scores)

    particular = np.linalg.lstsq(matrix, right_hand_side, rcond=None)[0]
    _, singular_values, vh = np.linalg.svd(matrix, full_matrices=True)
    rank = int(np.linalg.matrix_rank(matrix))
    if rank != len(fit_ratios):
        raise RuntimeError(
            f"Bernstein 锚点线性系统秩不足：rank={rank}"
        )

    nullspace = vh[rank:].T
    base, transform = _coefficient_map(particular, nullspace)
    diagnostic_ratios = np.array([0.0, r1, r2, 1.0], dtype=float)
    diagnostic_scores = np.array([LOW_SCORE, s1, s2, HIGH_SCORE], dtype=float)
    _, second_difference = _difference_operators()
    curvature_operator = (
        DEGREE
        * (DEGREE - 1)
        * (
            _basis(
                np.linspace(0.0, 1.0, DIAGNOSTIC_GRID_SIZE),
                DEGREE - 2,
            )
            @ second_difference
        )
    )
    base_curvature = curvature_operator @ base
    free_curvature = curvature_operator @ transform

    def bending_objective(free_values):
        curvature = base_curvature + free_curvature @ free_values
        return float(np.mean(np.square(curvature)))

    def bending_gradient(free_values):
        curvature = base_curvature + free_curvature @ free_values
        return (
            2.0
            * (free_curvature.T @ curvature)
            / float(len(curvature))
        )

    candidates = []

    # 二阶差分共有 9 项；split 表示首个非负项的位置。
    for split in range(1, DEGREE - 1):
        constraint_matrix, constraint_lower = _profile_constraints(
            split, base, transform
        )
        feasible = linprog(
            np.zeros(nullspace.shape[1], dtype=float),
            A_ub=-constraint_matrix,
            b_ub=-constraint_lower,
            bounds=[(None, None)] * nullspace.shape[1],
            method="highs",
        )
        if not feasible.success:
            continue

        linear_constraint = LinearConstraint(
            constraint_matrix,
            constraint_lower,
            np.inf,
        )
        try:
            result = minimize(
                bending_objective,
                feasible.x,
                jac=bending_gradient,
                method="SLSQP",
                constraints=[linear_constraint],
                options={"ftol": 1e-12, "maxiter": 1000},
            )
        except Exception as error:
            result = None
            optimizer_error = str(error)
        else:
            optimizer_error = ""

        trials = [
            (
                feasible.x,
                "linear_feasible",
                False,
                "仅找到线性可行点；"
                f"linprog: {feasible.message}; "
                f"optimizer_error: {optimizer_error or 'N/A'}",
            )
        ]
        if result is not None and np.all(np.isfinite(result.x)):
            trials.append(
                (
                    result.x,
                    "bending_optimized",
                    result.success,
                    result.message,
                )
            )

        for free_values, source, success, message in trials:
            maximum_constraint_violation = max(
                0.0,
                float(
                    np.max(
                        constraint_lower
                        - constraint_matrix @ free_values
                    )
                ),
            )
            candidates.append(
                _diagnose_candidate(
                    free_values,
                    split,
                    base,
                    transform,
                    diagnostic_ratios,
                    diagnostic_scores,
                    maximum_constraint_violation,
                    source,
                    success,
                    message,
                )
            )

    if not candidates:
        raise RuntimeError("所有反 S 曲率窗口均不存在单调可行解")

    candidates.sort(
        key=lambda candidate: (
            not _passes_quality(candidate),
            candidate["bending_energy"],
            candidate["max_anchor_error"],
        )
    )
    best = candidates[0]

    if not _passes_quality(best):
        raise RuntimeError(
            "所有反 S 曲率窗口均未通过质量门禁；"
            f"最佳窗口={best['shape_split']}, "
            f"predictions={best['predictions'].tolist()}, "
            f"residuals={best['residuals'].tolist()}, "
            f"minimum_coefficient_step={best['minimum_coefficient_step']}, "
            f"start_second_difference={best['start_second_difference']}, "
            f"end_second_difference={best['end_second_difference']}, "
            f"maximum_constraint_violation="
            f"{best['maximum_constraint_violation']}, "
            f"inflection_count={best['inflection_count']}, "
            f"candidate_source={best['candidate_source']}, "
            f"optimizer_success={best['optimizer_success']}, "
            f"message={best['optimizer_message']}"
        )

    return {
        # ---- Runtime 必需 ----
        "model": "global_bernstein_10",
        "coefficients": [float(value) for value in best["coefficients"]],
        # ---- 仅离线诊断 ----
        "predictions": [float(value) for value in best["predictions"]],
        "residuals": [float(value) for value in best["residuals"]],
        "max_anchor_error": best["max_anchor_error"],
        "minimum_coefficient_step": best["minimum_coefficient_step"],
        "start_second_difference": best["start_second_difference"],
        "end_second_difference": best["end_second_difference"],
        "maximum_constraint_violation": best["maximum_constraint_violation"],
        "bending_energy": best["bending_energy"],
        "inflection_count": best["inflection_count"],
        "shape_split": best["shape_split"],
        "candidate_source": best["candidate_source"],
        "free_values": [float(value) for value in best["free_values"]],
        "singular_values": [float(value) for value in singular_values],
        "optimizer_success": best["optimizer_success"],
        "optimizer_message": best["optimizer_message"],
    }


if __name__ == "__main__":
    params = fit_risk((0.2, 6.0), (0.4, 8.0))

    print("# ---- 拟合参数 ----")
    print(f"model = {params['model']}")
    print(f"coefficients = {params['coefficients']}")

    print("# ---- 离线诊断 ----")
    for key in (
        "predictions",
        "residuals",
        "max_anchor_error",
        "minimum_coefficient_step",
        "start_second_difference",
        "end_second_difference",
        "maximum_constraint_violation",
        "bending_energy",
        "inflection_count",
        "shape_split",
        "candidate_source",
        "free_values",
        "singular_values",
        "optimizer_success",
        "optimizer_message",
    ):
        print(f"{key} = {params[key]}")

    coefficients = params["coefficients"]
    runtime = f'''def risk(hit, total):
    x = hit / total
    if x <= 0.0:
        return {coefficients[0]!r}
    if x >= 1.0:
        return {coefficients[-1]!r}

    values = {coefficients!r}
    for level in range(10, 0, -1):
        for index in range(level):
            values[index] = (
                (1.0 - x) * values[index] + x * values[index + 1]
            )
    return values[0]
'''
    print("# ---- Runtime（复制即用，纯标准库）----")
    print(runtime)
