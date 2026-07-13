"""identity —— 恒等评分映射（fit_identity）。

业务目标：
  将若干项 Pass/Fail 判定结果映射为评分。
  规则：只要通过任意一项（hit >= 1）即得满分；全不通过得 0 分。

曲线性质：
  * 单调不减：hit=0 -> 0；hit>=1 -> 满分；
  * 阶跃：本规则允许使用简单 if 判断，不要求连续
    （区别于 profile 的连续 Richards 曲线）。

端点：
  passed=0 -> 0，passed>=1 -> 10，精确成立（恒等规则天然精确，无需端点校正）。

输出：
  Float，满分 10 分制（runtime 函数内置）；最终评分与舍入由调用方/下游决定。

入参：无（恒等规则阈值固定为 1，无需控制点）。
出参：无运行时常量（runtime 为纯逻辑，不消费任何拟合参数）。

环境要求：Python 3.11+。
依赖：无（纯标准库，无需拟合）。
"""


def fit_identity() -> dict:
    return {
        # ---- Runtime 必需：无 ----
        # 恒等规则为纯逻辑，runtime 不消费任何拟合常量
        # ---- 仅离线诊断 ----
        "rms": 0.0,  # 无拟合，残差恒为 0
    }


if __name__ == "__main__":
    params = fit_identity()

    # 1) 拟合参数 + 说明
    print("# ---- 拟合参数 ----")
    print("（无）    # 恒等规则，无需拟合参数")

    # 2) 离线诊断
    print("# ---- 离线诊断 ----")
    print(f"rms = {params['rms']}    # 无拟合，恒为 0")

    # 3) 可复制即用的 runtime，纯标准库
    runtime = '''def identity(hit, total):
    if hit >= 1:
        return 10.0
    return 0.0
'''
    print("# ---- Runtime（复制即用，纯标准库）----")
    print(runtime)
