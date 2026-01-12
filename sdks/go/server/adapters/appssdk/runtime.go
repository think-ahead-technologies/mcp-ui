package appssdk

// adapterRuntimeScript contains the JavaScript runtime for the Apps SDK adapter.
// This code is extracted from the TypeScript adapter-runtime.bundled.ts file.
// Generated from: sdks/typescript/server/src/adapters/appssdk/adapter-runtime.bundled.ts
// Last updated: 2026-01-12
const adapterRuntimeScript = `var __defProp = Object.defineProperty;
var __defNormalProp = (obj, key, value) => key in obj ? __defProp(obj, key, { enumerable: true, configurable: true, writable: true, value }) : obj[key] = value;
var __publicField = (obj, key, value) => {
  __defNormalProp(obj, typeof key !== "symbol" ? key + "" : key, value);
  return value;
};
class MCPUIAppsSdkAdapter {
  constructor(config = {}) {
    __publicField(this, "config");
    __publicField(this, "pendingRequests", /* @__PURE__ */ new Map());
    __publicField(this, "messageIdCounter", 0);
    __publicField(this, "originalPostMessage", null);
    this.config = {
      logger: config.logger || console,
      hostOrigin: config.hostOrigin || window.location.origin,
      timeout: config.timeout || 3e4,
      intentHandling: config.intentHandling || "prompt"
    };
  }
  install() {
    if (!window.openai) {
      this.config.logger.warn(
        "[MCPUI-Apps SDK Adapter] window.openai not detected. Adapter will not activate."
      );
      return false;
    }
    this.config.logger.log("[MCPUI-Apps SDK Adapter] Initializing adapter...");
    this.patchPostMessage();
    this.setupAppsSdkEventListeners();
    this.sendRenderData();
    this.config.logger.log("[MCPUI-Apps SDK Adapter] Adapter initialized successfully");
    return true;
  }
  uninstall() {
    for (const request of this.pendingRequests.values()) {
      clearTimeout(request.timeoutId);
      request.reject(new Error("Adapter uninstalled"));
    }
    this.pendingRequests.clear();
    if (this.originalPostMessage) {
      try {
        const parentWindow = window.parent ?? null;
        if (parentWindow) {
          parentWindow.postMessage = this.originalPostMessage;
        }
        this.config.logger.log("[MCPUI-Apps SDK Adapter] Restored original parent.postMessage");
      } catch (error) {
        this.config.logger.error(
          "[MCPUI-Apps SDK Adapter] Failed to restore original postMessage:",
          error
        );
      }
    }
    this.config.logger.log("[MCPUI-Apps SDK Adapter] Adapter uninstalled");
  }
  patchPostMessage() {
    const parentWindow = window.parent ?? null;
    this.originalPostMessage = parentWindow?.postMessage?.bind(parentWindow) ?? null;
    if (!this.originalPostMessage) {
      this.config.logger.debug(
        "[MCPUI-Apps SDK Adapter] parent.postMessage does not exist, installing shim only"
      );
    } else {
      this.config.logger.debug(
        "[MCPUI-Apps SDK Adapter] Monkey-patching parent.postMessage to intercept MCP-UI messages"
      );
    }
    const postMessageInterceptor = (message, targetOriginOrOptions, transfer) => {
      if (this.isMCPUIMessage(message)) {
        const mcpMessage = message;
        this.config.logger.debug(
          "[MCPUI-Apps SDK Adapter] Intercepted MCP-UI message:",
          mcpMessage.type
        );
        this.handleMCPUIMessage(mcpMessage);
      } else {
        if (this.originalPostMessage) {
          this.config.logger.debug(
            "[MCPUI-Apps SDK Adapter] Forwarding non-MCP-UI message to original postMessage"
          );
          if (typeof targetOriginOrOptions === "string" || targetOriginOrOptions === void 0) {
            const targetOrigin = targetOriginOrOptions ?? "*";
            this.originalPostMessage(message, targetOrigin, transfer);
          } else {
            this.originalPostMessage(message, targetOriginOrOptions);
          }
        } else {
          this.config.logger.warn(
            "[MCPUI-Apps SDK Adapter] No original postMessage to forward to, ignoring message:",
            message
          );
        }
      }
    };
    try {
      if (parentWindow) {
        parentWindow.postMessage = postMessageInterceptor;
      }
    } catch (error) {
      this.config.logger.error(
        "[MCPUI-Apps SDK Adapter] Failed to monkey-patch parent.postMessage:",
        error
      );
    }
  }
  isMCPUIMessage(message) {
    if (!message || typeof message !== "object") {
      return false;
    }
    const msg = message;
    return typeof msg.type === "string" && (msg.type.startsWith("ui-") || ["tool", "prompt", "intent", "notify", "link"].includes(msg.type));
  }
  async handleMCPUIMessage(message) {
    this.config.logger.debug("[MCPUI-Apps SDK Adapter] Received MCPUI message:", message.type);
    try {
      switch (message.type) {
        case "tool":
          await this.handleToolMessage(message);
          break;
        case "prompt":
          await this.handlePromptMessage(message);
          break;
        case "intent":
          await this.handleIntentMessage(message);
          break;
        case "notify":
          await this.handleNotifyMessage(message);
          break;
        case "link":
          await this.handleLinkMessage(message);
          break;
        case "ui-lifecycle-iframe-ready":
          this.sendRenderData();
          break;
        case "ui-request-render-data":
          this.sendRenderData(message.messageId);
          break;
        case "ui-size-change":
          this.handleSizeChange(message);
          break;
        case "ui-request-data":
          this.handleRequestData(message);
          break;
        default:
          this.config.logger.warn("[MCPUI-Apps SDK Adapter] Unknown message type:", message.type);
      }
    } catch (error) {
      this.config.logger.error("[MCPUI-Apps SDK Adapter] Error handling message:", error);
      if (message.messageId) {
        this.sendErrorResponse(message.messageId, error);
      }
    }
  }
  async handleToolMessage(message) {
    if (message.type !== "tool")
      return;
    const { toolName, params } = message.payload;
    const messageId = message.messageId || this.generateMessageId();
    this.sendAcknowledgment(messageId);
    try {
      if (!window.openai?.callTool) {
        throw new Error("Tool calling is not supported in this environment");
      }
      const result = await this.withTimeout(window.openai.callTool(toolName, params), messageId);
      this.sendSuccessResponse(messageId, result);
    } catch (error) {
      this.sendErrorResponse(messageId, error);
    }
  }
  async handlePromptMessage(message) {
    if (message.type !== "prompt")
      return;
    const prompt = message.payload.prompt;
    const messageId = message.messageId || this.generateMessageId();
    this.sendAcknowledgment(messageId);
    try {
      if (!window.openai?.sendFollowUpMessage) {
        throw new Error("Followup turns are not supported in this environment");
      }
      await this.withTimeout(window.openai.sendFollowUpMessage({ prompt }), messageId);
      this.sendSuccessResponse(messageId, { success: true });
    } catch (error) {
      this.sendErrorResponse(messageId, error);
    }
  }
  async handleIntentMessage(message) {
    if (message.type !== "intent")
      return;
    const messageId = message.messageId || this.generateMessageId();
    this.sendAcknowledgment(messageId);
    if (this.config.intentHandling === "ignore") {
      this.config.logger.log("[MCPUI-Apps SDK Adapter] Intent ignored:", message.payload.intent);
      this.sendSuccessResponse(messageId, { ignored: true });
      return;
    }
    const { intent, params } = message.payload;
    const prompt = intent + (params ? (": " + JSON.stringify(params)) : "");
    try {
      if (!window.openai?.sendFollowUpMessage) {
        throw new Error("Followup turns are not supported in this environment");
      }
      await this.withTimeout(window.openai.sendFollowUpMessage({ prompt }), messageId);
      this.sendSuccessResponse(messageId, { success: true });
    } catch (error) {
      this.sendErrorResponse(messageId, error);
    }
  }
  async handleNotifyMessage(message) {
    if (message.type !== "notify")
      return;
    const messageId = message.messageId || this.generateMessageId();
    this.config.logger.log("[MCPUI-Apps SDK Adapter] Notification:", message.payload.message);
    this.sendAcknowledgment(messageId);
    this.sendSuccessResponse(messageId, { acknowledged: true });
  }
  async handleLinkMessage(message) {
    if (message.type !== "link")
      return;
    const messageId = message.messageId || this.generateMessageId();
    this.sendAcknowledgment(messageId);
    this.sendErrorResponse(
      messageId,
      new Error("Navigation is not supported in Apps SDK environment")
    );
  }
  handleSizeChange(message) {
    this.config.logger.debug(
      "[MCPUI-Apps SDK Adapter] Size change requested (no-op in Apps SDK):",
      message.payload
    );
  }
  handleRequestData(message) {
    const messageId = message.messageId || this.generateMessageId();
    this.sendAcknowledgment(messageId);
    this.sendErrorResponse(messageId, new Error("Generic data requests not yet implemented"));
  }
  setupAppsSdkEventListeners() {
    window.addEventListener("openai:set_globals", () => {
      this.config.logger.debug("[MCPUI-Apps SDK Adapter] Globals updated");
      this.sendRenderData();
    });
  }
  sendRenderData(requestMessageId) {
    if (!window.openai)
      return;
    const renderData = {
      toolInput: window.openai.toolInput,
      toolOutput: window.openai.toolOutput,
      widgetState: window.openai.widgetState,
      locale: window.openai.locale || "en-US",
      theme: window.openai.theme || "light",
      displayMode: window.openai.displayMode || "inline",
      maxHeight: window.openai.maxHeight
    };
    this.dispatchMessageToIframe({
      type: "ui-lifecycle-iframe-render-data",
      messageId: requestMessageId,
      payload: { renderData }
    });
  }
  sendAcknowledgment(messageId) {
    this.dispatchMessageToIframe({
      type: "ui-message-received",
      payload: { messageId }
    });
  }
  sendSuccessResponse(messageId, response) {
    this.dispatchMessageToIframe({
      type: "ui-message-response",
      payload: { messageId, response }
    });
  }
  sendErrorResponse(messageId, error) {
    const errorObj = error instanceof Error ? { message: error.message, name: error.name } : { message: String(error) };
    this.dispatchMessageToIframe({
      type: "ui-message-response",
      payload: { messageId, error: errorObj }
    });
  }
  dispatchMessageToIframe(data) {
    const event = new MessageEvent("message", {
      data,
      origin: this.config.hostOrigin,
      source: null
    });
    window.dispatchEvent(event);
  }
  async withTimeout(promise, requestId) {
    return new Promise((resolve, reject) => {
      const timeoutId = setTimeout(() => {
        this.pendingRequests.delete(requestId);
        reject(new Error("Request timed out after " + this.config.timeout + "ms"));
      }, this.config.timeout);
      this.pendingRequests.set(requestId, {
        messageId: requestId,
        type: "generic",
        resolve,
        reject,
        timeoutId
      });
      promise.then((result) => {
        clearTimeout(timeoutId);
        this.pendingRequests.delete(requestId);
        resolve(result);
      }).catch((error) => {
        clearTimeout(timeoutId);
        this.pendingRequests.delete(requestId);
        reject(error);
      });
    });
  }
  generateMessageId() {
    return "adapter-" + Date.now() + "-" + (++this.messageIdCounter);
  }
}
let adapterInstance = null;
function initAdapter(config) {
  if (adapterInstance) {
    console.warn("[MCPUI-Apps SDK Adapter] Adapter already initialized");
    return true;
  }
  adapterInstance = new MCPUIAppsSdkAdapter(config);
  return adapterInstance.install();
}
function uninstallAdapter() {
  if (adapterInstance) {
    adapterInstance.uninstall();
    adapterInstance = null;
  }
}
`
