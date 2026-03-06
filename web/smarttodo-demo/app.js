const STORAGE_KEY = "smarttodo-demo:todos";
const REVIEW_KEY = "smarttodo-demo:review";
const PLAN_KEY = "smarttodo-demo:plan";
const FOCUS_KEY = "smarttodo-demo:focus";
const SESSION_API_KEY = "smarttodo-demo:api-key";
const LOCAL_API_KEY = "smarttodo-demo:api-key:persisted";

const state = {
  wasmReady: false,
  apiReady: false,
  todos: loadJSON(STORAGE_KEY, []),
  review: loadJSON(REVIEW_KEY, null),
  plan: loadJSON(PLAN_KEY, null),
  focus: loadJSON(FOCUS_KEY, null),
  filterQuery: "",
  logs: [],
};

const els = {
  wasmStatus: document.getElementById("wasm-status"),
  total: document.getElementById("metric-total"),
  open: document.getElementById("metric-open"),
  overdue: document.getElementById("metric-overdue"),
  keyInput: document.getElementById("api-key-input"),
  rememberKey: document.getElementById("remember-key"),
  keyStatus: document.getElementById("key-status"),
  connectKey: document.getElementById("connect-key"),
  forgetKey: document.getElementById("forget-key"),
  terminalOutput: document.getElementById("terminal-output"),
  commandForm: document.getElementById("command-form"),
  commandInput: document.getElementById("command-input"),
  contextInput: document.getElementById("context-input"),
  searchInput: document.getElementById("search-input"),
  reviewSummary: document.getElementById("review-summary"),
  reviewFocus: document.getElementById("review-focus"),
  reviewRisks: document.getElementById("review-risks"),
  planHeadline: document.getElementById("plan-headline"),
  planBlocks: document.getElementById("plan-blocks"),
  focusRecommendation: document.getElementById("focus-recommendation"),
  laneHot: document.getElementById("lane-hot"),
  laneReady: document.getElementById("lane-ready"),
  laneDone: document.getElementById("lane-done"),
};

boot().catch((error) => {
  logLine(`Boot failed: ${error.message}`, "error");
  render();
});

async function boot() {
  wireEvents();
  logLine("Launching control room.", "meta");
  await bootWasm();
  restoreApiKey();
  render();
}

function wireEvents() {
  els.connectKey.addEventListener("click", connectKey);
  els.forgetKey.addEventListener("click", forgetKey);
  els.commandForm.addEventListener("submit", async (event) => {
    event.preventDefault();
    const command = els.commandInput.value.trim();
    if (!command) return;
    els.commandInput.value = "";
    await runCommand(command);
  });
  els.searchInput.addEventListener("input", () => {
    state.filterQuery = els.searchInput.value.trim().toLowerCase();
    renderBoard();
  });
  document.querySelectorAll("[data-action]").forEach((button) => {
    button.addEventListener("click", async () => {
      const action = button.dataset.action;
      if (action === "prioritize") await runCommand("/prioritize");
      if (action === "focus") await runCommand("/focus");
      if (action === "review") await runCommand("/review");
      if (action === "plan") await runCommand(`/plan ${els.contextInput.value.trim()}`.trim());
    });
  });
}

async function bootWasm() {
  const go = new Go();
  const response = await fetch("smarttodo.wasm");
  const bytes = await response.arrayBuffer();
  const result = await WebAssembly.instantiate(bytes, go.importObject);
  go.run(result.instance);
  await waitForApi();
  state.wasmReady = true;
  els.wasmStatus.textContent = "WASM bridge ready";
}

async function waitForApi() {
  const started = performance.now();
  while (!window.smarttodoWasm) {
    if (performance.now() - started > 5000) {
      throw new Error("smarttodo WASM API did not initialize");
    }
    await new Promise((resolve) => setTimeout(resolve, 40));
  }
}

function restoreApiKey() {
  const saved = localStorage.getItem(LOCAL_API_KEY) || sessionStorage.getItem(SESSION_API_KEY);
  if (saved) {
    els.keyInput.value = saved;
    if (localStorage.getItem(LOCAL_API_KEY)) {
      els.rememberKey.checked = true;
    }
    connectKey();
  }
}

async function connectKey() {
  const apiKey = els.keyInput.value.trim();
  if (!apiKey) {
    setKeyStatus("Enter an API key first.", true);
    return;
  }
  try {
    await callWasm("init", apiKey);
    state.apiReady = true;
    if (els.rememberKey.checked) {
      localStorage.setItem(LOCAL_API_KEY, apiKey);
      sessionStorage.removeItem(SESSION_API_KEY);
    } else {
      sessionStorage.setItem(SESSION_API_KEY, apiKey);
      localStorage.removeItem(LOCAL_API_KEY);
    }
    setKeyStatus("API key connected. Live AI features are ready.", false);
    logLine("API key connected.", "meta");
  } catch (error) {
    state.apiReady = false;
    setKeyStatus(`Key connection failed: ${error.message}`, true);
    logLine(`API init failed: ${error.message}`, "error");
  }
  render();
}

function forgetKey() {
  state.apiReady = false;
  els.keyInput.value = "";
  sessionStorage.removeItem(SESSION_API_KEY);
  localStorage.removeItem(LOCAL_API_KEY);
  setKeyStatus("Stored key removed.", false);
  render();
}

function setKeyStatus(message, isError) {
  els.keyStatus.textContent = message;
  els.keyStatus.style.color = isError ? "#fecaca" : "#94a3b8";
}

async function runCommand(command) {
  logLine(`todo> ${command}`);
  if (command.startsWith("/")) {
    const [verb, ...rest] = command.split(" ");
    const arg = rest.join(" ").trim();
    if (verb === "/clear") {
      state.logs = [];
      renderTerminal();
      return;
    }
    if (verb === "/help") {
      logLine("Commands: /prioritize, /focus, /review, /plan [context], /filter query, /clear", "meta");
      return;
    }
    if (!ensureReady()) return;
    if (verb === "/prioritize") {
      state.todos = await callWasm("prioritizeBoard", JSON.stringify(state.todos));
      persistTodos();
      logLine("Board reprioritized.", "meta");
      render();
      return;
    }
    if (verb === "/focus") {
      state.focus = await callWasm("recommendNext", JSON.stringify(state.todos));
      localStorage.setItem(FOCUS_KEY, JSON.stringify(state.focus));
      logLine(`Recommended start: ${state.focus.title}`, "meta");
      renderInsights();
      return;
    }
    if (verb === "/review") {
      state.review = await callWasm("reviewBoard", JSON.stringify(state.todos));
      localStorage.setItem(REVIEW_KEY, JSON.stringify(state.review));
      logLine("Board review updated.", "meta");
      renderInsights();
      return;
    }
    if (verb === "/plan") {
      state.plan = await callWasm("planDay", JSON.stringify(state.todos), arg || els.contextInput.value.trim());
      localStorage.setItem(PLAN_KEY, JSON.stringify(state.plan));
      logLine("Day plan generated.", "meta");
      renderInsights();
      return;
    }
    if (verb === "/filter") {
      const filtered = await callWasm("filterBoard", JSON.stringify(state.todos), arg);
      state.filterQuery = "";
      renderBoard(filtered.map((todo) => todo.id));
      logLine(`Rendered semantic filter for: ${arg}`, "meta");
      return;
    }
    logLine(`Unknown command: ${verb}`, "error");
    return;
  }

  if (!ensureReady()) return;
  const todo = await callWasm("captureTodo", command);
  state.todos.unshift(todo);
  persistTodos();
  logLine(`Captured task: ${todo.title}`, "meta");
  render();
}

function ensureReady() {
  if (!state.wasmReady) {
    logLine("WASM bridge is still loading.", "error");
    return false;
  }
  if (!state.apiReady) {
    logLine("Connect an API key before using live AI actions.", "error");
    return false;
  }
  return true;
}

async function callWasm(method, ...args) {
  if (!window.smarttodoWasm || typeof window.smarttodoWasm[method] !== "function") {
    throw new Error(`WASM method not available: ${method}`);
  }
  const raw = await window.smarttodoWasm[method](...args);
  return JSON.parse(raw);
}

function persistTodos() {
  localStorage.setItem(STORAGE_KEY, JSON.stringify(state.todos));
}

function render(filteredIds = null) {
  renderMetrics();
  renderTerminal();
  renderBoard(filteredIds);
  renderInsights();
}

function renderMetrics() {
  const total = state.todos.length;
  const open = state.todos.filter((todo) => !todo.completed).length;
  const overdue = state.todos.filter((todo) => todo.deadline && !todo.completed && new Date(todo.deadline).getTime() < Date.now()).length;
  els.total.textContent = total;
  els.open.textContent = open;
  els.overdue.textContent = overdue;
}

function renderTerminal() {
  els.terminalOutput.innerHTML = "";
  const lines = state.logs.length ? state.logs : [{ text: "Awaiting command.", kind: "meta" }];
  for (const line of lines) {
    const row = document.createElement("div");
    row.className = `terminal-line ${line.kind || ""}`.trim();
    row.textContent = line.text;
    els.terminalOutput.appendChild(row);
  }
  els.terminalOutput.scrollTop = els.terminalOutput.scrollHeight;
}

function renderBoard(filteredIds = null) {
  const idFilter = filteredIds ? new Set(filteredIds) : null;
  const query = state.filterQuery;
  const visible = state.todos.filter((todo) => {
    if (idFilter && !idFilter.has(todo.id)) return false;
    if (!query) return true;
    return [todo.title, todo.description, todo.category, todo.location, todo.context].join(" ").toLowerCase().includes(query);
  });

  const hot = [];
  const ready = [];
  const done = [];

  visible.forEach((todo) => {
    if (todo.completed) {
      done.push(todo);
      return;
    }
    if ((todo.priority || "").toLowerCase() === "high" || isOverdue(todo)) {
      hot.push(todo);
      return;
    }
    ready.push(todo);
  });

  fillLane(els.laneHot, hot, "No urgent or overdue work.");
  fillLane(els.laneReady, ready, "Nothing queued in the ready lane.");
  fillLane(els.laneDone, done, "Nothing completed yet.");
}

function fillLane(container, todos, emptyMessage) {
  container.innerHTML = "";
  if (!todos.length) {
    const empty = document.createElement("div");
    empty.className = "empty-state";
    empty.textContent = emptyMessage;
    container.appendChild(empty);
    return;
  }
  todos.forEach((todo) => container.appendChild(renderCard(todo)));
}

function renderCard(todo) {
  const card = document.createElement("article");
  card.className = `card ${todo.completed ? "done" : ""} ${isOverdue(todo) ? "hot" : ""}`.trim();
  const subtasks = Array.isArray(todo.tasks) ? todo.tasks : [];
  const completedSubtasks = subtasks.filter((task) => task.completed).length;
  const description = [todo.description, todo.context].filter(Boolean).join(" | ");
  card.innerHTML = `
    <div class="card-top">
      <h4 class="card-title">${escapeHtml(todo.title || "Untitled task")}</h4>
      <span class="pill ${priorityClass(todo.priority)}">${escapeHtml((todo.priority || "open").toUpperCase())}</span>
    </div>
    <div class="card-meta">
      <span class="pill subtle">${escapeHtml(todo.category || "general")}</span>
      <span class="pill subtle">${escapeHtml(todo.location || "anywhere")}</span>
      <span class="pill subtle">${escapeHtml(todo.effort || "medium")}</span>
      ${todo.deadline ? `<span class="pill subtle">${escapeHtml(formatDeadline(todo.deadline))}</span>` : ""}
    </div>
    <p class="card-copy">${escapeHtml(description || "No additional context.")}</p>
    ${subtasks.length ? `<p class="card-copy">Subtasks ${completedSubtasks}/${subtasks.length}</p>` : ""}
    <div class="card-actions">
      <button data-action="toggle">${todo.completed ? "Reopen" : "Complete"}</button>
      <button data-action="focus">Focus</button>
      <button data-action="edit">Revise</button>
      <button data-action="delete">Delete</button>
    </div>
  `;

  card.querySelector('[data-action="toggle"]').addEventListener("click", () => {
    todo.completed = !todo.completed;
    persistTodos();
    render();
  });
  card.querySelector('[data-action="delete"]').addEventListener("click", () => {
    state.todos = state.todos.filter((item) => item.id !== todo.id);
    persistTodos();
    render();
  });
  card.querySelector('[data-action="focus"]').addEventListener("click", () => {
    state.focus = todo;
    localStorage.setItem(FOCUS_KEY, JSON.stringify(todo));
    renderInsights();
  });
  card.querySelector('[data-action="edit"]').addEventListener("click", async () => {
    if (!ensureReady()) return;
    const instruction = window.prompt("How should this task change?", "make the next step clearer");
    if (!instruction) return;
    const revised = await callWasm("reviseTodo", JSON.stringify(todo), instruction);
    state.todos = state.todos.map((item) => (item.id === todo.id ? { ...item, ...revised } : item));
    persistTodos();
    logLine(`Revised task: ${revised.title}`, "meta");
    render();
  });
  return card;
}

function renderInsights() {
  els.reviewSummary.innerHTML = state.review ? escapeHtml(state.review.summary || "No summary available.") : 'Run <code>/review</code> to synthesize a board brief.';
  renderList(els.reviewFocus, state.review?.focus_areas || []);
  renderList(els.reviewRisks, state.review?.risks || []);

  if (state.plan) {
    els.planHeadline.textContent = state.plan.headline || "Day plan ready.";
    els.planBlocks.innerHTML = "";
    (state.plan.blocks || []).forEach((block) => {
      const item = document.createElement("div");
      item.className = "plan-block";
      item.innerHTML = `<strong>${escapeHtml(block.label || "Block")}</strong><div>${escapeHtml(block.window || "Window TBD")}</div><p>${escapeHtml(block.goal || "No goal specified.")}</p>`;
      els.planBlocks.appendChild(item);
    });
    if (!state.plan.blocks?.length) {
      const empty = document.createElement("div");
      empty.className = "empty-state";
      empty.textContent = "No blocks returned.";
      els.planBlocks.appendChild(empty);
    }
  } else {
    els.planHeadline.innerHTML = 'Run <code>/plan</code> to project the board into time blocks.';
    els.planBlocks.innerHTML = "";
  }

  els.focusRecommendation.textContent = state.focus?.title || state.focus?.recommended_start || "Run /focus to choose the next task.";
}

function renderList(container, items) {
  container.innerHTML = "";
  if (!items.length) return;
  items.forEach((item) => {
    const li = document.createElement("li");
    li.textContent = item;
    container.appendChild(li);
  });
}

function logLine(text, kind = "") {
  state.logs.push({ text, kind });
  if (state.logs.length > 120) {
    state.logs = state.logs.slice(-120);
  }
  renderTerminal();
}

function loadJSON(key, fallback) {
  try {
    const raw = localStorage.getItem(key);
    return raw ? JSON.parse(raw) : fallback;
  } catch {
    return fallback;
  }
}

function priorityClass(priority) {
  const value = (priority || "").toLowerCase();
  if (value === "high") return "high";
  if (value === "low") return "low";
  return "medium";
}

function isOverdue(todo) {
  return Boolean(todo.deadline && !todo.completed && new Date(todo.deadline).getTime() < Date.now());
}

function formatDeadline(deadline) {
  const date = new Date(deadline);
  return Number.isNaN(date.getTime()) ? deadline : date.toLocaleString([], { month: "short", day: "numeric", hour: "numeric", minute: "2-digit" });
}

function escapeHtml(value) {
  return String(value)
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;")
    .replaceAll("'", "&#39;");
}
