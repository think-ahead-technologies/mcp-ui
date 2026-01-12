package mcpapps

// adapterRuntimeScript contains the JavaScript runtime for the MCP Apps adapter.
// This code is extracted from the TypeScript adapter-runtime.bundled.ts file.
// Generated from: sdks/typescript/server/src/adapters/mcp-apps/adapter-runtime.bundled.ts
// Last updated: 2026-01-12
const adapterRuntimeScript = `var __defProp = Object.defineProperty;
var __defNormalProp = (obj, key, value) => key in obj ? __defProp(obj, key, { enumerable: true, configurable: true, writable: true, value }) : obj[key] = value;
var __publicField = (obj, key, value) => {
  __defNormalProp(obj, typeof key !== "symbol" ? key + "" : key, value);
  return value;
};
const LATEST_PROTOCOL_VERSION = "2025-11-21";
const METHODS = {
  INITIALIZE: "ui/initialize",
  INITIALIZED: "ui/notifications/initialized",
  TOOL_INPUT: "ui/notifications/tool-input",
  TOOL_INPUT_PARTIAL: "ui/notifications/tool-input-partial",
  TOOL_RESULT: "ui/notifications/tool-result",
  TOOL_CANCELLED: "ui/notifications/tool-cancelled",
  HOST_CONTEXT_CHANGED: "ui/notifications/host-context-changed",
  SIZE_CHANGED: "ui/notifications/size-changed",
  RESOURCE_TEARDOWN: "ui/resource-teardown",
  TOOLS_CALL: "tools/call",
  NOTIFICATIONS_MESSAGE: "notifications/message",
  OPEN_LINK: "ui/open-link",
  MESSAGE: "ui/message"
};
class McpAppsAdapter {
  constructor(config = {}) {
    __publicField(this, "config");
    __publicField(this, "pendingRequests", /* @__PURE__ */ new Map());
    __publicField(this, "messageIdCounter", 0);
    __publicField(this, "originalPostMessage", null);
    __publicField(this, "parentWindow", null);
    __publicField(this, "hostCapabilities", null);
    __publicField(this, "hostContext", null);
    __publicField(this, "initialized", false);
    __publicField(this, "currentRenderData", {});
    this.config = {
      logger: config.logger || console,
      timeout: config.timeout || 3e4
    };
  }
  install() {
    this.parentWindow = window.parent;
    this.config.logger.log("[MCP Apps Adapter] Checking parent window...");
    this.config.logger.log("[MCP Apps Adapter] window.parent exists:", !!this.parentWindow);
    this.config.logger.log(
      "[MCP Apps Adapter] window.parent === window:",
      this.parentWindow === window
    );
    if (!this.parentWindow || this.parentWindow === window) {
      this.config.logger.warn(
        "[MCP Apps Adapter] No parent window detected. Adapter will not activate."
      );
      return false;
    }
    this.config.logger.log("[MCP Apps Adapter] Initializing adapter...");
    this.patchPostMessage();
    window.addEventListener("message", this.handleHostMessage.bind(this));
    this.performInitialization();
    this.config.logger.log("[MCP Apps Adapter] Adapter initialized successfully");
    return true;
  }
  async performInitialization() {
    const jsonRpcId = this.generateJsonRpcId();
    const initPromise = new Promise((resolve, reject) => {
      this.pendingRequests.set(String(jsonRpcId), {
        messageId: "init",
        type: "init",
        resolve: (result) => {
          const res = result;
          this.hostCapabilities = res?.hostCapabilities ?? null;
          this.hostContext = res?.hostContext ?? null;
          this.initialized = true;
          this.sendJsonRpcNotification(METHODS.INITIALIZED, {});
          if (this.hostContext) {
            if (this.hostContext.theme)
              this.currentRenderData.theme = this.hostContext.theme;
            if (this.hostContext.displayMode)
              this.currentRenderData.displayMode = this.hostContext.displayMode;
            if (this.hostContext.locale)
              this.currentRenderData.locale = this.hostContext.locale;
            if (this.hostContext.viewport?.maxHeight)
              this.currentRenderData.maxHeight = this.hostContext.viewport.maxHeight;
          }
          this.sendRenderData();
          this.dispatchMessageToIframe({
            type: "ui-lifecycle-iframe-ready"
          });
          resolve();
        },
        reject: (error) => {
          this.config.logger.error("[MCP Apps Adapter] Initialization failed:", error);
          reject(error);
        },
        timeoutId: setTimeout(() => {
          this.pendingRequests.delete(String(jsonRpcId));
          this.config.logger.warn("[MCP Apps Adapter] Initialization timed out, proceeding anyway");
          this.dispatchMessageToIframe({
            type: "ui-lifecycle-iframe-ready"
          });
          resolve();
        }, this.config.timeout)
      });
    });
    this.config.logger.log("[MCP Apps Adapter] Sending ui/initialize request with id:", jsonRpcId);
    this.sendJsonRpcRequest(jsonRpcId, METHODS.INITIALIZE, {
      appInfo: {
        name: "mcp-ui-adapter",
        version: "1.0.0"
      },
      appCapabilities: {},
      protocolVersion: LATEST_PROTOCOL_VERSION
    });
    this.config.logger.log("[MCP Apps Adapter] ui/initialize request sent");
    try {
      await initPromise;
    } catch (_error) {
      this.config.logger.warn("[MCP Apps Adapter] Continuing despite initialization error");
    }
  }
  uninstall() {
    for (const request of this.pendingRequests.values()) {
      clearTimeout(request.timeoutId);
      request.reject(new Error("Adapter uninstalled"));
    }
    this.pendingRequests.clear();
    if (this.originalPostMessage && this.parentWindow) {
      try {
        this.parentWindow.postMessage = this.originalPostMessage;
        this.config.logger.log("[MCP Apps Adapter] Restored original parent.postMessage");
      } catch (error) {
        this.config.logger.error(
          "[MCP Apps Adapter] Failed to restore original postMessage:",
          error
        );
      }
    }
    window.removeEventListener("message", this.handleHostMessage.bind(this));
    this.config.logger.log("[MCP Apps Adapter] Adapter uninstalled");
  }
  patchPostMessage() {
    this.originalPostMessage = this.parentWindow?.postMessage.bind(this.parentWindow) ?? null;
    const postMessageInterceptor = (message, targetOriginOrOptions, transfer) => {
      if (this.isMCPUIMessage(message)) {
        const mcpMessage = message;
        this.config.logger.debug("[MCP Apps Adapter] Intercepted MCP-UI message:", mcpMessage.type);
        this.handleMCPUIMessage(mcpMessage);
      } else {
        if (this.originalPostMessage) {
          if (typeof targetOriginOrOptions === "string" || targetOriginOrOptions === void 0) {
            const targetOrigin = targetOriginOrOptions ?? "*";
            this.originalPostMessage(message, targetOrigin, transfer);
          } else {
            this.originalPostMessage(message, targetOriginOrOptions);
          }
        }
      }
    };
    try {
      if (this.parentWindow) {
        this.parentWindow.postMessage = postMessageInterceptor;
      }
    } catch (error) {
      this.config.logger.error(
        "[MCP Apps Adapter] Failed to monkey-patch parent.postMessage:",
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
  handleHostMessage(event) {
    const data = event.data;
    if (!data || typeof data !== "object" || !data.jsonrpc) {
      return;
    }
    this.config.logger.debug("[MCP Apps Adapter] Received JSON-RPC message:", data);
    if (data.method) {
      switch (data.method) {
        case METHODS.TOOL_INPUT:
          this.currentRenderData.toolInput = data.params?.arguments;
          this.sendRenderData();
          break;
        case METHODS.TOOL_INPUT_PARTIAL:
          this.currentRenderData.toolInput = data.params?.arguments;
          this.sendRenderData();
          break;
        case METHODS.TOOL_RESULT:
          this.currentRenderData.toolOutput = data.params;
          this.sendRenderData();
          break;
        case METHODS.HOST_CONTEXT_CHANGED:
          if (data.params?.theme)
            this.currentRenderData.theme = data.params.theme;
          if (data.params?.displayMode)
            this.currentRenderData.displayMode = data.params.displayMode;
          if (data.params?.locale)
            this.currentRenderData.locale = data.params.locale;
          if (data.params?.viewport?.maxHeight)
            this.currentRenderData.maxHeight = data.params.viewport.maxHeight;
          this.sendRenderData();
          break;
        case METHODS.SIZE_CHANGED:
          if (data.params?.height)
            this.currentRenderData.maxHeight = data.params.height;
          this.sendRenderData();
          break;
        case METHODS.TOOL_CANCELLED:
          this.dispatchMessageToIframe({
            type: "ui-lifecycle-tool-cancelled",
            payload: {
              reason: data.params?.reason
            }
          });
          break;
        case METHODS.RESOURCE_TEARDOWN:
          this.dispatchMessageToIframe({
            type: "ui-lifecycle-teardown",
            payload: {
              reason: data.params?.reason
            }
          });
          if (data.id) {
            this.sendJsonRpcResponse(data.id, {});
          }
          break;
      }
    } else if (data.id) {
      const pendingRequest = this.pendingRequests.get(String(data.id));
      if (pendingRequest) {
        if (data.error) {
          pendingRequest.reject(new Error(data.error.message));
        } else {
          pendingRequest.resolve(data.result);
        }
        this.pendingRequests.delete(String(data.id));
        clearTimeout(pendingRequest.timeoutId);
        this.dispatchMessageToIframe({
          type: "ui-message-response",
          messageId: pendingRequest.messageId,
          payload: {
            messageId: pendingRequest.messageId,
            response: data.result,
            error: data.error
          }
        });
      }
    }
  }
  async handleMCPUIMessage(message) {
    const messageId = message.messageId || this.generateMessageId();
    this.dispatchMessageToIframe({
      type: "ui-message-received",
      payload: { messageId }
    });
    try {
      switch (message.type) {
        case "tool": {
          const { toolName, params } = message.payload;
          const jsonRpcId = this.generateJsonRpcId();
          this.pendingRequests.set(String(jsonRpcId), {
            messageId,
            type: "tool",
            resolve: () => {
            },
            reject: () => {
            },
            timeoutId: setTimeout(() => {
              this.pendingRequests.delete(String(jsonRpcId));
              this.dispatchMessageToIframe({
                type: "ui-message-response",
                messageId,
                payload: { messageId, error: "Timeout" }
              });
            }, this.config.timeout)
          });
          this.sendJsonRpcRequest(jsonRpcId, METHODS.TOOLS_CALL, {
            name: toolName,
            arguments: params
          });
          break;
        }
        case "ui-size-change": {
          const { width, height } = message.payload;
          this.sendJsonRpcNotification(METHODS.SIZE_CHANGED, { width, height });
          break;
        }
        case "notify": {
          const { message: msg } = message.payload;
          this.sendJsonRpcNotification(METHODS.NOTIFICATIONS_MESSAGE, {
            level: "info",
            data: msg
          });
          break;
        }
        case "link": {
          const { url } = message.payload;
          const jsonRpcId = this.generateJsonRpcId();
          this.pendingRequests.set(String(jsonRpcId), {
            messageId,
            type: "link",
            resolve: () => {
            },
            reject: () => {
            },
            timeoutId: setTimeout(() => {
              this.pendingRequests.delete(String(jsonRpcId));
              this.dispatchMessageToIframe({
                type: "ui-message-response",
                messageId,
                payload: { messageId, error: "Timeout" }
              });
            }, this.config.timeout)
          });
          this.sendJsonRpcRequest(jsonRpcId, METHODS.OPEN_LINK, { url });
          break;
        }
        case "prompt": {
          const { prompt } = message.payload;
          const jsonRpcId = this.generateJsonRpcId();
          this.pendingRequests.set(String(jsonRpcId), {
            messageId,
            type: "prompt",
            resolve: () => {
            },
            reject: () => {
            },
            timeoutId: setTimeout(() => {
              this.pendingRequests.delete(String(jsonRpcId));
              this.dispatchMessageToIframe({
                type: "ui-message-response",
                messageId,
                payload: { messageId, error: "Timeout" }
              });
            }, this.config.timeout)
          });
          this.sendJsonRpcRequest(jsonRpcId, METHODS.MESSAGE, {
            role: "user",
            content: [{ type: "text", text: prompt }]
          });
          break;
        }
        case "ui-lifecycle-iframe-ready": {
          this.sendJsonRpcNotification(METHODS.INITIALIZED, {});
          this.sendRenderData();
          break;
        }
        case "ui-request-render-data": {
          this.sendRenderData(messageId);
          break;
        }
        case "intent": {
          const { intent, params } = message.payload;
          const jsonRpcId = this.generateJsonRpcId();
          this.pendingRequests.set(String(jsonRpcId), {
            messageId,
            type: "intent",
            resolve: () => {
            },
            reject: () => {
            },
            timeoutId: setTimeout(() => {
              this.pendingRequests.delete(String(jsonRpcId));
              this.dispatchMessageToIframe({
                type: "ui-message-response",
                messageId,
                payload: { messageId, error: "Timeout" }
              });
            }, this.config.timeout)
          });
          this.sendJsonRpcRequest(jsonRpcId, METHODS.MESSAGE, {
            role: "user",
            content: [
              { type: "text", text: "Intent: " + intent + ". Parameters: " + JSON.stringify(params) }
            ]
          });
          break;
        }
      }
    } catch (error) {
      this.config.logger.error("[MCP Apps Adapter] Error handling message:", error);
      this.dispatchMessageToIframe({
        type: "ui-message-response",
        messageId,
        payload: { messageId, error }
      });
    }
  }
  sendRenderData(requestMessageId) {
    this.dispatchMessageToIframe({
      type: "ui-lifecycle-iframe-render-data",
      messageId: requestMessageId,
      payload: {
        renderData: {
          toolInput: this.currentRenderData.toolInput,
          toolOutput: this.currentRenderData.toolOutput,
          widgetState: this.currentRenderData.widgetState,
          locale: this.currentRenderData.locale,
          theme: this.currentRenderData.theme,
          displayMode: this.currentRenderData.displayMode,
          maxHeight: this.currentRenderData.maxHeight
        }
      }
    });
  }
  sendJsonRpcRequest(id, method, params) {
    this.originalPostMessage?.(
      {
        jsonrpc: "2.0",
        id,
        method,
        params
      },
      "*"
    );
  }
  sendJsonRpcResponse(id, result) {
    this.originalPostMessage?.(
      {
        jsonrpc: "2.0",
        id,
        result
      },
      "*"
    );
  }
  sendJsonRpcNotification(method, params) {
    this.originalPostMessage?.(
      {
        jsonrpc: "2.0",
        method,
        params
      },
      "*"
    );
  }
  dispatchMessageToIframe(data) {
    const event = new MessageEvent("message", {
      data,
      origin: window.location.origin,
      source: window
    });
    window.dispatchEvent(event);
  }
  generateMessageId() {
    return "adapter-" + Date.now() + "-" + (++this.messageIdCounter);
  }
  generateJsonRpcId() {
    return ++this.messageIdCounter;
  }
}
let adapterInstance = null;
function initAdapter(config) {
  if (adapterInstance) {
    console.warn("[MCP Apps Adapter] Adapter already initialized");
    return true;
  }
  adapterInstance = new McpAppsAdapter(config);
  return adapterInstance.install();
}
function uninstallAdapter() {
  if (adapterInstance) {
    adapterInstance.uninstall();
    adapterInstance = null;
  }
}
`
