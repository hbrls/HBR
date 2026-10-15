import type { Plugin } from "@kilocode/plugin";

/**
 * Kilo 正常聊天请求的 system prompt 拼接顺序：
 *
 *  1. Kilo soul / Persona（可选）
 *  2. Agent 自定义 prompt，或 Provider 类型提示
 *  3. environment（包含模型身份声明）
 *  4. memory（可选）
 *  5. board context（可选）
 *  6. AGENTS / CLAUDE / CONTEXT / config 指令
 *  7. MCP instructions（可选）
 *  8. Skills（可选）
 *  9. Structured output 指令（可选）
 * 10. 当前用户消息中的 system 字段（可选）
 *
 * 上述内容在触发钩子前已用换行符拼接，
 * 一个 system 数组元素中可能包含多个段落。
 *
 * 对所有 Provider 生效：
 * 遍历所有 system 数组元素，删除以
 * "You are powered by the model named" 开头的整行。
 * 同时删除该行末尾的 LF / CRLF，其余内容原样保留。
 */
const MODEL_IDENTITY_LINE =
  /^You are powered by the model named[^\r\n]*(?:\r\n|\n|$)/gm;

const server: Plugin = async ({ client }) => {
  return {
    "experimental.chat.system.transform": async (input, output) => {
      console.error(
        "[KILO_IDENTITY_DEBUG][HOOK_ENTER] system transform 钩子已进入",
      );

      let removedLines = 0;

      for (let i = 0; i < output.system.length; i++) {
        const before = output.system[i];
        const countBefore = removedLines;
        const after = before.replace(
          MODEL_IDENTITY_LINE,
          (matchedLine, offset: number) => {
            removedLines++;
            return "";
          },
        );

        // 保留调用方持有的数组，原地修改其元素。
        output.system[i] = after;
      }
    },
  };
};

export default {
  id: "remove-model-identity",
  server,
};
